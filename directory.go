package stormpath

import "golang.org/x/net/context"

//Directory represents a Stormpath directory object
//
//See: http://docs.stormpath.com/rest/product-guide/#directories
type Directory struct {
	accountStoreResource
	Name                  string                 `json:"name,omitempty"`
	Description           string                 `json:"description,omitempty"`
	Status                string                 `json:"status,omitempty"`
	Groups                *Groups                `json:"groups,omitempty"`
	Tenant                *Tenant                `json:"tenant,omitempty"`
	AccountCreationPolicy *AccountCreationPolicy `json:"accountCreationPolicy,omitempty"`
}

//Directories represnets a paged result of directories
type Directories struct {
	collectionResource
	Items []Directory `json:"items"`
}

//NewDirectory creates a new directory with the given name
func NewDirectory(name string) *Directory {
	return &Directory{Name: name}
}

//GetDirectory loads a directory by href and criteria
func GetDirectory(ctx context.Context, href string, criteria Criteria) (*Directory, error) {
	directory := &Directory{}

	err := getClient(ctx).get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		directory,
	)

	if err != nil {
		return nil, err
	}

	return directory, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (dir *Directory) Refresh(ctx context.Context) error {
	return getClient(ctx).get(dir.Href, emptyPayload(), dir)
}

//Update updates the given resource, by doing a POST to the resource Href
func (dir *Directory) Update(ctx context.Context) error {
	return getClient(ctx).post(dir.Href, dir, dir)
}

//GetAccountCreationPolicy loads the directory account creation policy
func (dir *Directory) GetAccountCreationPolicy(ctx context.Context) (*AccountCreationPolicy, error) {
	err := getClient(ctx).get(buildAbsoluteURL(dir.AccountCreationPolicy.Href), emptyPayload(), dir.AccountCreationPolicy)

	if err != nil {
		return nil, err
	}

	return dir.AccountCreationPolicy, nil
}

//GetGroups returns all the groups from a directory
func (dir *Directory) GetGroups(ctx context.Context, criteria Criteria) (*Groups, error) {
	err := getClient(ctx).get(
		buildAbsoluteURL(dir.Groups.Href, criteria.ToQueryString()),
		emptyPayload(),
		dir.Groups,
	)

	if err != nil {
		return nil, err
	}

	return dir.Groups, nil
}

//CreateGroup creates a new group in the directory
func (dir *Directory) CreateGroup(ctx context.Context, group *Group) error {
	return getClient(ctx).post(dir.Groups.Href, group, group)
}

//RegisterAccount registers a new account into the directory
//
//See: http://docs.stormpath.com/rest/product-guide/#directory-accounts
func (dir *Directory) RegisterAccount(ctx context.Context, account *Account) error {
	return getClient(ctx).post(dir.Accounts.Href, account, account)
}

//RegisterSocialAccount registers a new account into the application using an external provider Google, Facebook
//
//See: http://docs.stormpath.com/rest/product-guide/#accessing-accounts-with-google-authorization-codes-or-an-access-tokens
func (dir *Directory) RegisterSocialAccount(ctx context.Context, socialAccount *SocialAccount) (*Account, error) {
	account := &Account{}

	err := getClient(ctx).post(dir.Accounts.Href, socialAccount, account)

	if err != nil {
		return nil, err
	}

	return account, nil
}
