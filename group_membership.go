package stormpath

import "golang.org/x/net/context"

type GroupMembership struct {
	resource
	Account Account `json:"account"`
	Group   Group   `json:"group"`
}

type GroupMemberships struct {
	collectionResource
	Items []GroupMembership `json:"items"`
}

func NewGroupMembership(accountHref string, groupHref string) *GroupMembership {
	account := Account{}
	account.Href = accountHref
	group := Group{}
	group.Href = groupHref
	return &GroupMembership{
		Account: account,
		Group:   group,
	}
}

func (groupmembership *GroupMembership) GetAccount(ctx context.Context, criteria Criteria) (*Account, error) {
	account := &Account{}

	err := getClient(ctx).get(
		buildAbsoluteURL(groupmembership.Account.Href, criteria.ToQueryString()),
		emptyPayload(),
		account,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (groupmembership *GroupMembership) GetGroup(ctx context.Context, criteria Criteria) (*Group, error) {
	group := &Group{}

	err := getClient(ctx).get(
		buildAbsoluteURL(groupmembership.Group.Href, criteria.ToQueryString()),
		emptyPayload(),
		group,
	)

	if err != nil {
		return nil, err
	}

	return group, nil
}
