package stormpath

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"log"
	"errors"
	"google.golang.org/appengine/urlfetch"
	"runtime"
	ae "google.golang.org/appengine/log"
)

var BaseURL = "https://api.stormpath.com/v1/"

//Version is the current SDK Version
const version = "0.1.0-beta.12"

const (
	Enabled = "ENABLED"
	Disabled = "DISABLED"
	ApplicationJson = "application/json"
	ApplicationFormURLencoded = "application/x-www-form-urlencoded"
)

var _credentials Credentials
var _cache       *CacheableCache

//ClientProperties is low level REST client for any Stormpath request,
//it holds the credentials, an the actual http client, and the cache.
//The Cache can be initialize in nil and the client would simply ignore it
//and don't cache any response.
type Client struct {
	Credentials Credentials
	Cache       *CacheableCache
	ctx         context.Context
	httpClient  *http.Client
}

// Init initializes the underlying client that communicates with Stormpath.
// This cache will be shared by all requests, but this is ok for the purposes of this SDK because all requests
// will be operating upon the same StormPath account, even across multiple calling threads.
func Init(credentials Credentials, cache *cache.Cache) {
	_credentials = credentials
	if cache != nil {
		_cache = &CacheableCache{Cache: cache}
	}
}

func getClient(ctx context.Context) *Client {
	trans := urlfetch.Transport{Context:ctx, AllowInvalidServerCertificate:false}
	httpClient := &http.Client{Transport: &trans}
	client := &Client{Credentials: _credentials, Cache: _cache, ctx: ctx, httpClient: httpClient}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return checkRedirect(client, req, via)
	}

	log.Printf("Credentials: %#v", _credentials)
	return client
}

func (client *Client) postURLEncodedForm(urlStr string, body string, result interface{}) error {
	return client.execute("POST", urlStr, []byte(body), result, ApplicationFormURLencoded)
}

func (client *Client) post(urlStr string, body interface{}, result interface{}) error {
	return client.execute("POST", urlStr, body, result, ApplicationJson)
}

func (client *Client) get(urlStr string, body interface{}, result interface{}) error {
	return client.execute("GET", urlStr, body, result, ApplicationJson)
}

func (client *Client) delete(urlStr string, body interface{}) error {
	return client.do(client.newRequest("DELETE", urlStr, body, ApplicationJson))
}

func (client *Client) execute(method string, urlStr string, body interface{}, result interface{}, contentType string) error {
	return client.doWithResult(client.newRequest(method, urlStr, body, contentType), result)
}

func buildRelativeURL(parts ...string) string {
	buffer := bytes.NewBufferString(BaseURL)

	for i, part := range parts {
		buffer.WriteString(part)
		if !strings.HasSuffix(part, "/") && i + 1 < len(parts) {
			buffer.WriteString("/")
		}
	}

	return buffer.String()
}

func buildAbsoluteURL(parts ...string) string {
	buffer := bytes.NewBufferString("")

	for i, part := range parts {
		buffer.WriteString(part)
		if !strings.HasSuffix(part, "/") && i + 1 < len(parts) {
			buffer.WriteString("/")
		}
	}

	return buffer.String()
}

func (client *Client) newRequest(method string, urlStr string, body interface{}, contentType string) *http.Request {
	fmt.Sprintln("Building new request for: %v", urlStr)

	var encodedBody []byte
	if strings.ToLower(method) == strings.ToLower("GET") || strings.ToLower(method) == strings.ToLower("DELETE") {
		encodedBody = body.([]byte)
	} else if contentType == ApplicationJson {
		encodedBody, _ = json.Marshal(body)
	} else {
		//If content type is not application/json then it is application/x-www-form-urlencoded in which case the body should be the encoded params as a []byte
		encodedBody = body.([]byte)
	}
	fmt.Printf("%vENCODED_BODY: %v", NL, encodedBody)
	fmt.Printf("%vCONTENT_TYPE=%v", NL, contentType)

	req, _ := http.NewRequest(method, urlStr, bytes.NewReader(encodedBody))
	req.Header.Set("User-Agent", fmt.Sprintf("sappenin/stormpath-sdk-go/%s (%s; %s)", version, runtime.GOOS, runtime.GOARCH))
	//req.Header.Set("User-Agent", fmt.Sprintf("sappenin-sp-client"))
	req.Header.Set("Accept", ApplicationJson)
	req.Header.Set("Content-Type", contentType)

	uuid, _ := uuid.NewV4()
	nonce := uuid.String()

	Authenticate(client.ctx, req, encodedBody, time.Now().In(time.UTC), client.Credentials, nonce)

	fmt.Sprintln("Returning Request: %#v", req)
	return req
}

//buildExpandParam coverts a slice of expand attributes to a url.Values with
//only one value "expand=attr1,attr2,etc"
func buildExpandParam(expandAttributes []string) url.Values {
	stringBuffer := bytes.NewBufferString("")

	first := true
	for _, expandAttribute := range expandAttributes {
		if !first {
			stringBuffer.WriteString(",")
		}
		stringBuffer.WriteString(expandAttribute)
		first = false
	}

	values := url.Values{}
	expandValue := stringBuffer.String()
	//Should not include the expand query param if the value is empty
	if expandValue != "" {
		values.Add("expand", expandValue)
	}

	return values
}

func requestParams(values ...url.Values) string {
	params := url.Values{}

	for _, v := range values {
		params = appendParams(params, v)
	}

	encodedParams := params.Encode()
	if encodedParams != "" {
		return "?" + encodedParams
	}
	return ""
}

func appendParams(params url.Values, toAppend url.Values) url.Values {
	for k, v := range toAppend {
		params[k] = v
	}
	return params
}

func emptyPayload() []byte {
	return []byte{}
}

//doWithResult executes the given StormpathRequest and serialize the response body into the given expected result,
//it returns an error if any occurred while executing the request or serializing the response
func (client *Client) doWithResult(request *http.Request, result interface{}) error {
	var err error
	var response *http.Response

	key := request.URL.String()

	if client.Cache != nil && request.Method == "GET" && client.Cache.Exists(key) {
		err = client.Cache.Get(key, result)
	} else {
		response, err = client.execRequest(request)
		if err != nil {
			return err
		}
		err = json.NewDecoder(response.Body).Decode(result)
	}

	if client.Cache != nil && err == nil {
		switch request.Method {
		case "POST", "DELETE", "PUT":
			client.Cache.Del(key)
			break
		case "GET":
			cacheResource(key, result, client.Cache)
		}
	}

	return err
}

//do executes the StormpathRequest without expecting a response body as a result,
//it returns an error if any occurred while executing the request
func (client *Client) do(request *http.Request) error {
	_, err := client.execRequest(request)
	return err
}

//execRequest executes a request, it would return a byte slice with the raw resoponse data and an error if any occurred
func (client *Client) execRequest(req *http.Request) (*http.Response, error) {

	var dump []byte
	dump, _ = httputil.DumpRequest(req, true)
	ae.Debugf(client.ctx, "Stormpath request\n%s", dump)

	resp, err := client.httpClient.Do(req)

	dump, _ = httputil.DumpResponse(resp, true)
	ae.Debugf(client.ctx, "Stormpath response\n%s", dump)

	return resp, client.handleResponseError(resp, err)
}

func cleanCustomData(customData map[string]interface{}) map[string]interface{} {
	// delete illegal keys from data
	// http://docs.stormpath.com/rest/product-guide/#custom-data
	keys := []string{
		"href", "createdAt", "modifiedAt", "meta",
		"spMeta", "spmeta", "ionmeta", "ionMeta",
	}

	for i := range keys {
		delete(customData, keys[i])
	}

	return customData
}

func checkRedirect(client *Client, req *http.Request, via []*http.Request) error {
	//Go client default behavior is to bail after 10 redirects
	if len(via) > 10 {
		return errors.New("stopped after 10 redirects")
	}
	//No redirect do nothing
	if len(via) == 0 {
		// No redirects
		return nil
	}
	// Re-Authenticate the redirect request
	uuid, _ := uuid.NewV4()
	nonce := uuid.String()

	//We can use an empty payload cause the only redirect is for the current tenant
	//this could change in the future
	Authenticate(client.ctx, req, emptyPayload(), time.Now().In(time.UTC), client.Credentials, nonce)

	return nil
}