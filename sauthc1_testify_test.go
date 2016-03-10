// +build enc
package stormpath_test

import (
	"net/http"
	"time"

	"github.com/sappenin/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"strings"
)

// The original StormPath client (https://github.com/jarias/stormpath-sdk-go) has two problems when SAuthc1 Digests are computed inside of AppEngine.  First,
// the original client uses the Go http Client and errantly specifies an empty set of quotes as the payload for GET requests (see https://github.com/jarias/stormpath-sdk-go/issues/23).
// Second, the urlFetch service appends extra information into the user-agent header of all requests, which makes signing the user-agent field prone to runtime
// errors depending on the environment (i.e., dev server, appengine production, etc).  This set of tests validates each scenario.

func init() {
	var err error
	cred, err = stormpath.NewDefaultCredentials()
	if err != nil {
		panic(err)
	}
}

func TestKnownSAuthc1ForUrlFetch_StormPath(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	//log.Printf("Assemling request...")
	ds := "20160308T230740Z"

	req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/tenants/current", nil)
	//log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"api.stormpath.com"},
		"User-Agent": []string{"sappenin-sp-client"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},

	}
	//fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "03d9cf9c-352e-4bf5-5308-40c98fc54958"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2PXPDIBC8NPB500J73ZNGJEVB/20160308/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=7ff61dc6d245071fc1be636f2cd973fdfb91e82a8dd638c523e8fca391cd17fa"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}

// Tests the SAuthc1 for a valid request bin.  We see that, while using the urlFetch service, the payload for the GET request is empty and the User-Agent is ignored when computing the
// Auth Digest.
func TestKnownSAuthc1ForUrlFetch_RequestBin(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	//log.Printf("Assemling request...")
	ds := "20160308T221346Z"

	req, _ := http.NewRequest("GET", "https://requestb.in/1kvzjhg1", nil)
	//log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"requestb.in"},
		"User-Agent": []string{"Foo Bar This value is ignored"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	//fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "8fd1415e-6e0b-45b9-7513-7fd6e9887795"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2PXPDIBC8NPB500J73ZNGJEVB/20160308/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=3c6e15b0d7d3ce68f337952f5306d9ea5c69a576ea6c062f44918afdb8888a6d"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}

func TestKnownSAuthc1ForUrlFetch_StormPathAccounts_GET(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	//log.Printf("Assemling request...")
	ds := "20160309T002222Z"

	req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/applications/4jIEdHsNp17DWMQd8HTKQY/accountStoreMappings/?limit=25&offset=0", nil)
	//log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"api.stormpath.com"},
		"User-Agent": []string{"sappenin/stormpath-sdk-go/0.1.0-beta.12 (darwin; amd64)"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	//fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "3abf1257-2b65-41bf-6773-9b5f607850bc"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2PXPDIBC8NPB500J73ZNGJEVB/20160309/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=85d60a8fa70553f021100206d7cc927632be7d265edaad8c99760d4f4db7661f"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}

func TestKnownSAuthc1ForUrlFetch_Delete(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	//log.Printf("Assemling request...")
	ds := "20160309T002224Z"

	req, _ := http.NewRequest("DELETE", "https://api.stormpath.com/v1/directories/4jJ4VvOa3wljEdUFZr7L3W", nil)
	//log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"api.stormpath.com"},
		"User-Agent": []string{"sappenin/stormpath-sdk-go/0.1.0-beta.12 (darwin; amd64)"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	//fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "459c56d5-a391-49d2-7fc3-a48e7b4b5539"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2PXPDIBC8NPB500J73ZNGJEVB/20160309/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=808d09dbaf549486f4f27ec1e4f4100c43f4d53ac777b2f86ed1ddf2dbd67bcf"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}