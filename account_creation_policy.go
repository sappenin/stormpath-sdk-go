package stormpath

import (
	"golang.org/x/net/context"
)

//AccountCreationPolicy represents a directory account creation policy object
//
//See: http://docs.stormpath.com/rest/product-guide/#directory-account-creation-policy
type AccountCreationPolicy struct {
	resource
	VerificationEmailStatus           string          `json:"verificationEmailStatus,omitempty"`
	VerificationEmailTemplates        *EmailTemplates `json:"verificationEmailTemplates,omitempty"`
	VerificationSuccessEmailStatus    string          `json:"verificationSuccessEmailStatus,omitempty"`
	VerificationSuccessEmailTemplates *EmailTemplates `json:"verificationSuccessEmailTemplates,omitempty"`
	WelcomeEmailStatus                string          `json:"welcomeEmailStatus,omitempty"`
	WelcomeEmailTemplates             *EmailTemplates `json:"welcomeEmailTemplates,omitempty"`
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (policy *AccountCreationPolicy) Refresh(ctx context.Context) error {
	return getClient(ctx).get(policy.Href, emptyPayload(), policy)
}

//Update updates the given resource, by doing a POST to the resource Href
func (policy *AccountCreationPolicy) Update(ctx context.Context) error {
	return getClient(ctx).post(policy.Href, policy, policy)
}

//GetVerificationEmailTemplates loads the policy VerificationEmailTemplates collection and returns it
func (policy *AccountCreationPolicy) GetVerificationEmailTemplates(ctx context.Context) (*EmailTemplates, error) {
	err := getClient(ctx).get(policy.VerificationEmailTemplates.Href, emptyPayload(), policy.VerificationEmailTemplates)

	if err != nil {
		return nil, err
	}

	return policy.VerificationEmailTemplates, nil
}

//GetVerificationSuccessEmailTemplates loads the policy VerificationSuccessEmailTemplates collection and returns it
func (policy *AccountCreationPolicy) GetVerificationSuccessEmailTemplates(ctx context.Context) (*EmailTemplates, error) {
	err := getClient(ctx).get(policy.VerificationSuccessEmailTemplates.Href, emptyPayload(), policy.VerificationSuccessEmailTemplates)

	if err != nil {
		return nil, err
	}

	return policy.VerificationSuccessEmailTemplates, nil
}

//GetWelcomeEmailTemplates loads the policy WelcomeEmailTemplates collection and returns it
func (policy *AccountCreationPolicy) GetWelcomeEmailTemplates(ctx context.Context) (*EmailTemplates, error) {
	err := getClient(ctx).get(policy.WelcomeEmailTemplates.Href, emptyPayload(), policy.WelcomeEmailTemplates)

	if err != nil {
		return nil, err
	}

	return policy.WelcomeEmailTemplates, nil
}
