package stormpath

import (
	"encoding/base64"
	"errors"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
)

//Application represents a Stormpath application object
//
//See: http://docs.stormpath.com/rest/product-guide/#applications
type Application struct {
	accountStoreResource
	Name                       string                `json:"name,omitempty"`
	Description                string                `json:"description,omitempty"`
	Status                     string                `json:"status,omitempty"`
	Groups                     *Groups               `json:"groups,omitempty"`
	Tenant                     *Tenant               `json:"tenant,omitempty"`
	PasswordResetTokens        *resource             `json:"passwordResetTokens,omitempty"`
	AccountStoreMappings       *AccountStoreMappings `json:"accountStoreMappings,omitempty"`
	DefaultAccountStoreMapping *AccountStoreMapping  `json:"defaultAccountStoreMapping,omitempty"`
	DefaultGroupStoreMapping   *AccountStoreMapping  `json:"defaultGroupStoreMapping,omitempty"`
}

//Applications represents a paged result or applications
type Applications struct {
	collectionResource
	Items []Application `json:"items"`
}

//IDSiteCallbackResult holds the ID Site callback parsed JWT token information + the acccount if one was given
type IDSiteCallbackResult struct {
	Account *Account
	State   string
	IsNew   bool
	Status  string
}

//OAuthResponse represents an OAuth2 response from StormPath
type OAuthResponse struct {
	AccessToken              string `json:"access_token"`
	RefreshToken             string `json:"refresh_token"`
	TokenType                string `json:"token_type"`
	ExpiresIn                int    `json:"expires_in"`
	StormpathAccessTokenHref string `json:"stormpath_access_token_href"`
}

type AccessToken struct {
	resource
	Account     *Account               `json:"account,omitempty"`
	Tenant      *Tenant                `json:"tenant,omitempty"`
	Application *Application           `json:"application,omitempty"`
	JWT         string                 `json:"jwt"`
	ExpandedJWT map[string]interface{} `json:"expandedJwt"`
}

//NewApplication creates a new application
func NewApplication(name string) *Application {
	return &Application{Name: name}
}

//GetApplication loads an application by href and criteria
func GetApplication(ctx context.Context, href string, criteria Criteria) (*Application, error) {
	application := &Application{}

	err := getClient(ctx).get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		application,
	)

	return application, err
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (app *Application) Refresh(ctx context.Context) error {
	return getClient(ctx).get(app.Href, emptyPayload(), app)
}

//Update updates the given resource, by doing a POST to the resource Href
func (app *Application) Update(ctx context.Context) error {
	return getClient(ctx).post(app.Href, app, app)
}

//Purge deletes all the account stores before deleting the application
//
//See: http://docs.stormpath.com/rest/product-guide/#delete-an-application
func (app *Application) Purge(ctx context.Context) error {
	accountStoreMappings, err := app.GetAccountStoreMappings(ctx, MakeAccountStoreMappingCriteria().Offset(0).Limit(25))
	if err != nil {
		return err
	}
	for _, m := range accountStoreMappings.Items {
		getClient(ctx).delete(m.AccountStore.Href, emptyPayload())
	}

	return app.Delete(ctx)
}

//GetAccountStoreMappings returns all the applications account store mappings
//
//See: http://docs.stormpath.com/rest/product-guide/#application-account-store-mappings
func (app *Application) GetAccountStoreMappings(ctx context.Context, criteria Criteria) (*AccountStoreMappings, error) {
	accountStoreMappings := &AccountStoreMappings{}

	err := getClient(ctx).get(
		buildAbsoluteURL(app.AccountStoreMappings.Href, criteria.ToQueryString()),
		emptyPayload(),
		accountStoreMappings,
	)

	if err != nil {
		return nil, err
	}

	return accountStoreMappings, nil
}

//RegisterAccount registers a new account into the application
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts
func (app *Application) RegisterAccount(ctx context.Context, account *Account) error {
	err := getClient(ctx).post(app.Accounts.Href, account, account)
	if err == nil {
		//Password should be cleanup so we don't keep an unhash password in memory
		account.Password = ""
	}
	return err
}

//RegisterSocialAccount registers a new account into the application using an external provider Google, Facebook
//
//See: http://docs.stormpath.com/rest/product-guide/#accessing-accounts-with-google-authorization-codes-or-an-access-tokens
func (app *Application) RegisterSocialAccount(ctx context.Context, socialAccount *SocialAccount) (*Account, error) {
	account := &Account{}

	err := getClient(ctx).post(app.Accounts.Href, socialAccount, account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//AuthenticateAccount authenticates an account against the application
//
//See: http://docs.stormpath.com/rest/product-guide/#authenticate-an-account
func (app *Application) AuthenticateAccount(ctx context.Context, username string, password string) (*Account, error) {
	accountRef := &accountRef{Account: &Account{}}

	loginAttemptPayload := make(map[string]string)
	loginAttemptPayload["type"] = "basic"
	loginAttemptPayload["value"] = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	err := getClient(ctx).post(buildAbsoluteURL(app.Href, "loginAttempts"), loginAttemptPayload, accountRef)

	if err != nil {
		return nil, err
	}

	account := accountRef.Account
	err = account.Refresh(ctx)
	if err != nil {
		return nil, err
	}

	return account, nil
}

//GetOAuthToken creates a OAuth2 token response for a given user credentials
func (app *Application) GetOAuthToken(ctx context.Context, username string, password string) (*OAuthResponse, error) {
	response := &OAuthResponse{}

	values := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	body := canonicalizeQueryString(values)

	err := getClient(ctx).postURLEncodedForm(
		buildAbsoluteURL(app.Href, "oauth/token"),
		body,
		response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//Validate Token against Application
func (app *Application) ValidateToken(ctx context.Context, token string) (*AccessToken, error) {
	response := &AccessToken{}

	err := getClient(ctx).get(
		buildAbsoluteURL(app.Href, "authTokens", token),
		emptyPayload(),
		response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

//SendPasswordResetEmail sends a password reset email to the given user
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) SendPasswordResetEmail(ctx context.Context, email string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	passwordResetPayload := make(map[string]string)
	passwordResetPayload["email"] = email

	err := getClient(ctx).post(buildAbsoluteURL(app.Href, "passwordResetTokens"), passwordResetPayload, passwordResetToken)

	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

//ValidatePasswordResetToken validates a password reset token
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) ValidatePasswordResetToken(ctx context.Context, token string) (*AccountPasswordResetToken, error) {
	passwordResetToken := &AccountPasswordResetToken{}

	err := getClient(ctx).get(buildAbsoluteURL(app.Href, "passwordResetTokens", token), emptyPayload(), passwordResetToken)

	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

//ResetPassword resets a user password based on the reset token
//
//See: http://docs.stormpath.com/rest/product-guide/#reset-an-accounts-password
func (app *Application) ResetPassword(ctx context.Context, token string, newPassword string) (*Account, error) {
	accountRef := &accountRef{}
	account := &Account{}

	resetPasswordPayload := make(map[string]string)
	resetPasswordPayload["password"] = newPassword

	err := getClient(ctx).post(buildAbsoluteURL(app.Href, "passwordResetTokens", token), resetPasswordPayload, accountRef)

	if err != nil {
		return nil, err
	}
	account.Href = accountRef.Account.Href

	return account, nil
}

//CreateGroup creates a new group in the application
//
//See: http://docs.stormpath.com/rest/product-guide/#application-groups
func (app *Application) CreateGroup(ctx context.Context, group *Group) error {
	return getClient(ctx).post(app.Groups.Href, group, group)
}

//GetGroups returns all the application groups
//
//See: http://docs.stormpath.com/rest/product-guide/#application-groups
func (app *Application) GetGroups(ctx context.Context, criteria Criteria) (*Groups, error) {
	groups := &Groups{}

	err := getClient(ctx).get(
		buildAbsoluteURL(app.Groups.Href, criteria.ToQueryString()),
		emptyPayload(),
		groups,
	)

	if err != nil {
		return nil, err
	}

	return groups, nil
}

//CreateIDSiteURL creates the IDSite URL for the application
func (app *Application) CreateIDSiteURL(ctx context.Context, options map[string]string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	nonce, _ := uuid.NewV4()

	if options["path"] == "" {
		options["path"] = "/"
	}

	token.Claims["jti"] = nonce.String()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["iss"] = getClient(ctx).Credentials.ID
	token.Claims["sub"] = app.Href
	token.Claims["state"] = options["state"]
	token.Claims["path"] = options["path"]
	token.Claims["cb_uri"] = options["callbackURI"]

	tokenString, err := token.SignedString([]byte(getClient(ctx).Credentials.Secret))
	if err != nil {
		return "", err
	}

	p, _ := url.Parse(app.Href)
	ssoURL := p.Scheme + "://" + p.Host + "/sso"

	if options["logout"] == "true" {
		ssoURL = ssoURL + "/logout" + "?jwtRequest=" + tokenString
	} else {
		ssoURL = ssoURL + "?jwtRequest=" + tokenString
	}

	return ssoURL, nil
}

//HandleIDSiteCallback handles the URL from an ID Site callback it parses the JWT token
//validates it and return an IDSiteCallbackResult with the token info + the Account if the sub was given
func (app *Application) HandleIDSiteCallback(ctx context.Context, URL string) (*IDSiteCallbackResult, error) {
	result := &IDSiteCallbackResult{}

	cbURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	jwtResponse := cbURL.Query().Get("jwtResponse")

	token, err := jwt.Parse(jwtResponse, func(token *jwt.Token) (interface{}, error) {
		return []byte(getClient(ctx).Credentials.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if token.Claims["aud"].(string) != getClient(ctx).Credentials.ID {
		return nil, errors.New("ID Site invalid aud")
	}

	if time.Now().Unix() > int64(token.Claims["exp"].(float64)) {
		return nil, errors.New("ID Site JWT has expired")
	}

	if token.Claims["sub"] != nil {
		account, err := GetAccount(ctx, token.Claims["sub"].(string), MakeAccountCriteria())
		if err != nil {
			return nil, err
		}
		result.Account = account
	}
	if token.Claims["state"] != nil {
		result.State = token.Claims["state"].(string)
	}
	result.Status = token.Claims["status"].(string)

	return result, nil
}
