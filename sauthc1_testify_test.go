// +build enc
package stormpath_test

import (
	"net/http"
	"time"

	"github.com/sappenin/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
	"log"
	"strings"
	"fmt"
)

// The original StormPath client (https://github.com/jarias/stormpath-sdk-go) has two problems when SAuthc1 Digests are computed inside of AppEngine.  First,
// the original client uses the Go http Client and errantly specifies an empty set of quotes as the payload for GET requests (see https://github.com/jarias/stormpath-sdk-go/issues/23).
// Second, the urlFetch service appends extra information into the user-agent header of all requests, which makes signing the user-agent field prone to runtime
// errors depending on the environment (i.e., dev server, appengine production, etc).  This set of tests validates each scenario.

func init() {
	var STORMPATH_API_KEY_ID = os.Getenv("STORMPATH_API_KEY_ID")
	log.Printf("STORMPATH_API_KEY_ID: %v", STORMPATH_API_KEY_ID)

	var STORMPATH_API_KEY_SECRET = os.Getenv("STORMPATH_API_KEY_SECRET")
	log.Printf("STORMPATH_API_KEY_SECRET: %v", STORMPATH_API_KEY_SECRET)

	cred = stormpath.Credentials{ID: STORMPATH_API_KEY_ID, Secret: STORMPATH_API_KEY_SECRET}
}

func TestKnownSAuthc1ForUrlFetch_StormPath(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	log.Printf("Assemling request...")
	ds := "20160308T230740Z"

	req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/tenants/current", nil)
	log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"api.stormpath.com"},
		"User-Agent": []string{"sappenin-sp-client"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},

	}
	fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "03d9cf9c-352e-4bf5-5308-40c98fc54958"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2SF81PCVA776S8QA9SZ7PCREX/20160308/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=3bb749a5cbcd6e4a465d7933986e8397317b503767e6a35aba631dfe1ae29a17"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}

// Tests the SAuthc1 for a valid request bin.  We see that, while using the urlFetch service, the payload for the GET request is empty and the User-Agent is ignored when computing the
// Auth Digest.
func TestKnownSAuthc1ForUrlFetch_RequestBin(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	log.Printf("Assemling request...")
	ds := "20160308T221346Z"

	req, _ := http.NewRequest("GET", "https://requestb.in/1kvzjhg1", nil)
	log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"requestb.in"},
		"User-Agent": []string{"Foo Bar This value is ignored"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "8fd1415e-6e0b-45b9-7513-7fd6e9887795"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2SF81PCVA776S8QA9SZ7PCREX/20160308/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=ddf6b594851f069d54962f5c5720cfc6c69b11ab3c2f955f16554ed2b2ab0a84"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}

func TestKnownSAuthc1ForUrlFetch_StormPathAccounts_GET(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	log.Printf("Assemling request...")
	ds := "20160309T002222Z"

	req, _ := http.NewRequest("GET", "https://api.stormpath.com/v1/applications/4jIEdHsNp17DWMQd8HTKQY/accountStoreMappings/?limit=25&offset=0", nil)
	log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"api.stormpath.com"},
		"User-Agent": []string{"sappenin/stormpath-sdk-go/0.1.0-beta.12 (darwin; amd64)"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "3abf1257-2b65-41bf-6773-9b5f607850bc"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2SF81PCVA776S8QA9SZ7PCREX/20160309/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=8944a27a4b0fae2d1af5bd4ddc90ae493551f19b8f89d119ac46c0387b3665d6"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}

func TestKnownSAuthc1ForUrlFetch_Delete(t *testing.T) {
	assert.NotNil(t, cred.ID)
	assert.NotNil(t, cred.Secret)

	log.Printf("Assemling request...")
	ds := "20160309T002224Z"

	req, _ := http.NewRequest("DELETE", "https://api.stormpath.com/v1/directories/4jJ4VvOa3wljEdUFZr7L3W", nil)
	log.Printf("Request assembled!")

	req.Header = map[string][]string{
		"Host": []string{"api.stormpath.com"},
		"User-Agent": []string{"sappenin/stormpath-sdk-go/0.1.0-beta.12 (darwin; amd64)"},
		"Accept": []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	fmt.Printf("Request was: %#v", req)

	now, err := time.Parse(stormpath.TimestampFormat, ds)
	if err != nil {
		panic(err)
	}

	var nonce = "459c56d5-a391-49d2-7fc3-a48e7b4b5539"
	stormpath.Authenticate(ctx, req, []byte{}, now.In(time.UTC), cred, nonce)

	expected := strings.Join([]string{"SAuthc1 sauthc1Id=2SF81PCVA776S8QA9SZ7PCREX/20160309/",
		nonce,
		"/sauthc1_request, sauthc1SignedHeaders=accept;content-type;host;x-stormpath-date, sauthc1Signature=47bb11e2e92fb0cae81368c4127a10f53e13755bf4c1a2413191de1b45d835c1"}, "")
	actual := req.Header.Get("Authorization")
	assert.Equal(t, expected, actual)
}