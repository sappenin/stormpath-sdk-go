package stormpath

import "golang.org/x/net/context"

//AccountStoreMapping represents an Stormpath account store mapping
//
//See: http://docs.stormpath.com/rest/product-guide/#account-store-mappings
type AccountStoreMapping struct {
	resource
	ListIndex             *int        `json:"collectionResourceIndex,omitempty"`
	IsDefaultAccountStore *bool       `json:"isDefaultAccountStore,omitempty"`
	IsDefaultGroupStore   *bool       `json:"isDefaultGroupStore,omitempty"`
	Application           Application `json:"application"`
	AccountStore          resource    `json:"accountStore"`
}

//AccountStoreMappings represents a pages result of account store mappings
//
//See: http://docs.stormpath.com/rest/product-guide/#collectionResource-account-store-mappings
type AccountStoreMappings struct {
	collectionResource
	Items []AccountStoreMapping `json:"items"`
}

//NewAccountStoreMapping creates a new account store mappings
func NewAccountStoreMapping(applicationHref string, accountStoreHref string) *AccountStoreMapping {
	app := Application{}
	app.Href = applicationHref
	return &AccountStoreMapping{
		Application:  app,
		AccountStore: resource{Href: accountStoreHref},
	}
}

//Save saves the given account store mapping
func (mapping *AccountStoreMapping) Save(ctx context.Context) error {
	url := buildRelativeURL("accountStoreMappings")
	if mapping.Href != "" {
		url = mapping.Href
	}

	return getClient(ctx).post(url, mapping, mapping)
}
