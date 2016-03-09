package stormpath

import "golang.org/x/net/context"

const (
//TextPlain "text/plain" mime type
	TextPlain = "text/plain"
//TextHTML "text/html" mime type
	TextHTML = "text/html"
)

//EmailTemplate represents an account creation policy email template
type EmailTemplate struct {
	resource
	FromEmailAddress string            `json:"fromEmailAddress"`
	FromName         string            `json:"fromName"`
	Subject          string            `json:"subject"`
	HTMLBody         string            `json:"htmlBody"`
	TextBody         string            `json:"textBody"`
	MimeType         string            `json:"mimeType"`
	DefaultModel     map[string]string `json:"defaultModel"`
}

//EmailTemplates represents a collection of EmailTemplate
type EmailTemplates struct {
	collectionResource
	Items []EmailTemplate `json:"items"`
}

//GetEmailTemplate loads an email template by href
func GetEmailTemplate(ctx context.Context, href string) (*EmailTemplate, error) {
	emailTemplate := &EmailTemplate{}

	err := getClient(ctx).get(
		href,
		emptyPayload(),
		emailTemplate,
	)

	if err != nil {
		return nil, err
	}

	return emailTemplate, nil
}

//Refresh refreshes the resource by doing a GET to the resource href endpoint
func (template *EmailTemplate) Refresh(ctx context.Context) error {
	return getClient(ctx).get(template.Href, emptyPayload(), template)
}

//Update updates the given resource, by doing a POST to the resource Href
func (template *EmailTemplate) Update(ctx context.Context) error {
	return getClient(ctx).post(template.Href, template, template)
}
