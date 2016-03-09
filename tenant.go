package stormpath

import (
	"net/url"
	"golang.org/x/net/context"
)

//Tenant represents a Stormpath tenant see http://docs.stormpath.com/rest/product-guide/#tenants
type Tenant struct {
	customDataAwareResource
	Name         string       `json:"name"`
	Key          string       `json:"key"`
	Applications Applications `json:"applications"`
	Directories  Directories  `json:"directories"`
}

//CurrentTenant returns the current tenant see http://docs.stormpath.com/rest/product-guide/#retrieve-the-current-tenant
func CurrentTenant(ctx context.Context) (*Tenant, error) {
	tenant := &Tenant{}

	err := getClient(ctx).doWithResult(
		getClient(ctx).newRequest(
			"GET",
			buildRelativeURL("tenants", "current"),
			emptyPayload(),
			ApplicationJson,
		), tenant)

	return tenant, err
}

//CreateApplication creates a new application for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-applications
func (tenant *Tenant) CreateApplication(ctx context.Context, app *Application) error {
	var extraParams = url.Values{}
	extraParams.Add("createDirectory", "true")

	return getClient(ctx).post(buildRelativeURL("applications", requestParams(extraParams)), app, app)
}

//CreateDirectory creates a new directory for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-directories
func (tenant *Tenant) CreateDirectory(ctx context.Context, dir *Directory) error {
	return getClient(ctx).post(buildRelativeURL("directories"), dir, dir)
}

//GetApplications returns all the applications for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-applications
func (tenant *Tenant) GetApplications(ctx context.Context, criteria Criteria) (*Applications, error) {
	apps := &Applications{}

	err := getClient(ctx).get(buildAbsoluteURL(tenant.Applications.Href, criteria.ToQueryString()), emptyPayload(), apps)

	return apps, err
}

//GetDirectories returns all the directories for the given tenant
//
//See: http://docs.stormpath.com/rest/product-guide/#tenant-directories
func (tenant *Tenant) GetDirectories(ctx context.Context, criteria Criteria) (*Directories, error) {
	directories := &Directories{}

	err := getClient(ctx).get(buildAbsoluteURL(tenant.Directories.Href, criteria.ToQueryString()), emptyPayload(), directories)

	return directories, err
}
