package stormpath

import "golang.org/x/net/context"

//Group represents a Stormpath Group
//
//See: http://docs.stormpath.com/rest/product-guide/#groups
type Group struct {
	accountStoreResource
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	Tenant      *Tenant    `json:"tenant,omitempty"`
	Directory   *Directory `json:"directory,omitempty"`
}

//Groups represent a paged result of groups
type Groups struct {
	collectionResource
	Items []Group `json:"items"`
}

//NewGroup creates a new Group with the given name
func NewGroup(name string) *Group {
	return &Group{Name: name}
}

//GetGroup loads a group by href and criteria
func GetGroup(ctx context.Context, href string, criteria Criteria) (*Group, error) {
	group := &Group{}

	err := getClient(ctx).get(
		buildAbsoluteURL(href, criteria.ToQueryString()),
		emptyPayload(),
		group,
	)

	if err != nil {
		return nil, err
	}

	return group, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (group *Group) Refresh(ctx context.Context) error {
	return getClient(ctx).get(group.Href, emptyPayload(), group)
}

//Update updates the given resource, by doing a POST to the resource Href
func (group *Group) Update(ctx context.Context) error {
	return getClient(ctx).post(group.Href, group, group)
}

//GetGroupMemberships loads the given group memeberships
func (group *Group) GetGroupMemberships(ctx context.Context, criteria Criteria) (*GroupMemberships, error) {
	groupMemberships := &GroupMemberships{}

	err := getClient(ctx).get(
		buildAbsoluteURL(group.Href, "accountMemberships", criteria.ToQueryString()),
		emptyPayload(),
		groupMemberships,
	)

	if err != nil {
		return nil, err
	}

	return groupMemberships, nil
}
