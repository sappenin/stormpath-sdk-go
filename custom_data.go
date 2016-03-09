package stormpath

import "golang.org/x/net/context"

type customDataAwareResource struct {
	resource
	CustomData *CustomData `json:"customData,omitempty"`
}

//CustomData represents Stormpath's custom data resouce
type CustomData map[string]interface{}

//GetCustomData returns the given resource custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (r *customDataAwareResource) GetCustomData(ctx context.Context) (CustomData, error) {
	customData := make(CustomData)

	err := getClient(ctx).get(buildAbsoluteURL(r.Href, "customData"), emptyPayload(), &customData)

	if err != nil {
		return nil, err
	}

	return customData, nil
}

//UpdateCustomData sets or updates the given resource custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (r *customDataAwareResource) UpdateCustomData(ctx context.Context, customData CustomData) (CustomData, error) {
	customData = cleanCustomData(customData)

	err := getClient(ctx).post(buildAbsoluteURL(r.Href, "customData"), customData, &customData)

	if err != nil {
		return nil, err
	}

	return customData, nil
}

//DeleteCustomData deletes all the resource custom data
//
//See: http://docs.stormpath.com/rest/product-guide/#custom-data
func (r *customDataAwareResource) DeleteCustomData(ctx context.Context) error {
	return getClient(ctx).delete(buildAbsoluteURL(r.Href, "customData"), emptyPayload())
}
