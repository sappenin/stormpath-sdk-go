package stormpath

import "net/url"

type GroupCriteria struct {
	baseCriteria
}

func MakeGroupCriteria() GroupCriteria {
	return GroupCriteria{baseCriteria{filter: url.Values{}}}
}

func MakeGroupsCriteria() GroupCriteria {
	return GroupCriteria{baseCriteria{limit: 25, filter: url.Values{}}}
}

//Filter related functions

//Possible filters:
//* name
//* description
//* status

func (c GroupCriteria) NameEq(name string) GroupCriteria {
	c.filter.Add("name", name)
	return c
}

func (c GroupCriteria) DescriptionEq(description string) GroupCriteria {
	c.filter.Add("description", description)
	return c
}

func (c GroupCriteria) StatusEq(status string) GroupCriteria {
	c.filter.Add("statu", status)
	return c
}

//Expansion related functions

func (c GroupCriteria) WithCustomData() GroupCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "customData")
	return c
}

func (c GroupCriteria) WithAccounts(pageRequest PageRequest) GroupCriteria {
	c.expandedAttributes = append(c.expandedAttributes, pageRequest.toExpansion("accounts"))
	return c
}

func (c GroupCriteria) WithTenant() GroupCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "tenant")
	return c
}

func (c GroupCriteria) WithDirectory() GroupCriteria {
	c.expandedAttributes = append(c.expandedAttributes, "directory")
	return c
}
