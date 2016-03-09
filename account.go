package stormpath

import "golang.org/x/net/context"

//Account represents an Stormpath account object
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts
type Account struct {
	customDataAwareResource
	Username               string            `json:"username"`
	Email                  string            `json:"email"`
	Password               string            `json:"password,omitempty"`
	FullName               string            `json:"fullName,omitempty"`
	GivenName              string            `json:"givenName"`
	MiddleName             string            `json:"middleName,omitempty"`
	Surname                string            `json:"surname"`
	Status                 string            `json:"status,omitempty"`
	Groups                 *Groups           `json:"groups,omitempty"`
	GroupMemberships       *GroupMemberships `json:"groupMemberships,omitempty"`
	Directory              *Directory        `json:"directory,omitempty"`
	Tenant                 *Tenant           `json:"tenant,omitempty"`
	EmailVerificationToken *resource         `json:"emailVerificationToken,omitempty"`
}

//Accounts represents a paged result of Account objects
//
//See: http://docs.stormpath.com/rest/product-guide/#accounts-collectionResource
type Accounts struct {
	collectionResource
	Items []Account `json:"items"`
}

//AccountPasswordResetToken represents an password reset token for a given account
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts (Reset An Account’s Password)
type AccountPasswordResetToken struct {
	Href    string
	Email   string
	Account Account
}

type accountRef struct {
	Account *Account `json:"account"`
}

//SocialAccount represents the JSON payload use to create an account for a social backend directory
//(Google, Facebook, Github, etc)
type SocialAccount struct {
	Data ProviderData `json:"providerData"`
}

//ProviderData represents the especific information needed by the social provider (Google, Github, Faceboo, etc)
type ProviderData struct {
	ProviderID  string `json:"providerId"`
	AccessToken string `json:"accessToken,omitempty"`
}

//NewAccount returns a pointer to an Account with the minimum data required
func NewAccount(username, password, email, givenName, surname string) *Account {
	return &Account{Username: username, Password: password, Email: email, GivenName: givenName, Surname: surname}
}

//GetAccount fetches an account by href and criteria
func GetAccount(ctx context.Context, href string, criteria Criteria) (*Account, error) {
	account := &Account{}

	err := getClient(ctx).get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		account,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//GetAccount fetches an account by href and criteria
func (client *Account) GetAccount(ctx context.Context, href string, criteria Criteria) (*Account, error) {
	account := &Account{}

	err := getClient(ctx).get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		account,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (account *Account) Refresh(ctx context.Context) error {
	return getClient(ctx).get(account.Href, emptyPayload(), account)
}

//Update updates the given resource, by doing a POST to the resource Href
func (account *Account) Update(ctx context.Context) error {
	return getClient(ctx).post(account.Href, account, account)
}

//AddToGroup adds the given account to a given group and returns the respective GroupMembership
func (account *Account) AddToGroup(ctx context.Context, group *Group) (*GroupMembership, error) {
	groupMembership := NewGroupMembership(account.Href, group.Href)

	err := getClient(ctx).post(buildRelativeURL("groupMemberships"), groupMembership, groupMembership)

	if err != nil {
		return nil, err
	}

	return groupMembership, nil
}

//RemoveFromGroup removes the given account from the given group by searching the account groupmemberships,
//and deleting the corresponding one
func (account *Account) RemoveFromGroup(ctx context.Context, group *Group) error {
	groupMemberships, err := account.GetGroupMemberships(ctx, MakeGroupMemershipCriteria().Offset(0).Limit(25),
	)

	if err != nil {
		return err
	}

	for i := 1; len(groupMemberships.Items) > 0; i++ {
		for _, gm := range groupMemberships.Items {
			if gm.Group.Href == group.Href {
				return gm.Delete(ctx)
			}
		}
		groupMemberships, err = account.GetGroupMemberships(ctx, MakeGroupMemershipCriteria().Offset(i * 25).Limit(25),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

//GetGroupMemberships returns a paged result of the group memeberships of the given account
func (account *Account) GetGroupMemberships(ctx context.Context, criteria Criteria) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := getClient(ctx).get(
		buildAbsoluteURL(
			account.GroupMemberships.Href,
			criteria.ToQueryString(),
		),
		emptyPayload(),
		groupMemberships,
	)

	if err != nil {
		return nil, err
	}

	return groupMemberships, nil
}

//VerifyEmailToken verifies an email verification token associated with an account
//
//See: http://docs.stormpath.com/rest/product-guide/#account-verify-email
func VerifyEmailToken(ctx context.Context, token string) (*Account, error) {
	account := &Account{}
	client := getClient(ctx)
	err := client.post(buildAbsoluteURL(BaseURL, "accounts/emailVerificationTokens", token), emptyPayload(), account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
