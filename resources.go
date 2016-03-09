package stormpath

import (
	"strings"
	"time"
	"golang.org/x/net/context"
)

//collectionResource represent the basic attributes of collection of resources (Application, Group, Account, etc.)
type collectionResource struct {
	Href       string     `json:"href,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	ModifiedAt *time.Time `json:"modifiedAt,omitempty"`
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
}

func (r collectionResource) IsCacheable() bool {
	return false
}

//resource resprents the basic attributes of any resource (Application, Group, Account, etc.)
type resource struct {
	Href       string     `json:"href,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	ModifiedAt *time.Time `json:"modifiedAt,omitempty"`
}

func (r resource) IsCacheable() bool {
	return true
}

//Delete deletes the given account, it wont modify the calling account
func (r *resource) Delete(ctx context.Context) error {
	return getClient(ctx).delete(r.Href, emptyPayload())
}

type accountStoreResource struct {
	customDataAwareResource
	Accounts *Accounts `json:"accounts,omitempty"`
}

//GetAccounts returns all the accounts of the application
//
//See: http://docs.stormpath.com/rest/product-guide/#application-accounts
func (r *accountStoreResource) GetAccounts(ctx context.Context, criteria Criteria) (*Accounts, error) {
	accounts := &Accounts{}

	err := getClient(ctx).get(
		buildAbsoluteURL(r.Accounts.Href, criteria.ToQueryString()),
		emptyPayload(),
		accounts,
	)

	return accounts, err
}

func GetToken(href string) string {
	return href[strings.LastIndex(href, "/") + 1:]
}
