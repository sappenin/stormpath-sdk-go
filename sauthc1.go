package stormpath

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"golang.org/x/net/context"
	"fmt"
	"google.golang.org/appengine"
)

//SAuthc1 algorithm constants
const (
	IDTerminator = "sauthc1_request"
	AuthenticationScheme = "SAuthc1"
	NL = "\n"
	HostHeader = "Host"
	AuthorizationHeader = "Authorization"
	StormpathDateHeader = "X-Stormpath-Date"
	Algorithm = "HMAC-SHA-256"
	SAUTHC1Id = "sauthc1Id"
	SAUTHC1SignedHeaders = "sauthc1SignedHeaders"
	SAUTHC1Signature = "sauthc1Signature"
	DateFormat = "20060102"
	TimestampFormat = "20060102T150405Z0700"
)

//Authenticate generates the proper authentication header for the SAuthc1 algorithm use by Stormpath
func Authenticate(ctx context.Context, req *http.Request, payload []byte, date time.Time, credentials Credentials, nonce string) {

	//fmt.Printf("%vPAYLOAD: %#v", NL, string(payload))

	//fmt.Printf("%vStart AUTHENTICATE", NL)

	timestamp := date.Format(TimestampFormat)
	//fmt.Printf("%vAUTH(TIMESTAMP): %v", NL, timestamp)

	dateStamp := date.Format(DateFormat)
	//fmt.Printf("%vAUTH(DATESTAMP): %v", NL, dateStamp)

	req.Header.Set(HostHeader, req.URL.Host)
	//fmt.Printf("%vAUTH(HOSTHEADER): %v", NL, req.URL.Host)

	req.Header.Set(StormpathDateHeader, timestamp)
	//fmt.Printf("%vAUTH(STORMPATHDATEHEADER): %v", NL, timestamp)

	signedHeadersString := signedHeadersStringWithoutUserAgent(req.Header)
	//fmt.Printf("%vAUTH(SIGNEDHEADERSSTRING): %v", NL, signedHeadersString)

	//fmt.Printf("%vAUTH(PATH): %v", NL, req.URL.Path)
	//fmt.Printf("%vAUTH(QUERY): %v", NL, req.URL.Query())

	canonicalRequest :=
	req.Method +
	NL +
	canonicalizeresourcePath(req.URL.Path) +
	NL +
	canonicalizeQueryString(req.URL.Query()) +
	NL +
	canonicalizeHeadersStringWithoutUserAgent(ctx, req.Header) +
	NL +
	signedHeadersString +
	NL +
	hex.EncodeToString(hash(payload))
	//fmt.Printf("%vAUTH(CANONICALREQUEST): %v", NL, canonicalRequest)

	id := credentials.ID + "/" + dateStamp + "/" + nonce + "/" + IDTerminator
	//fmt.Printf("%vAUTH(ID): %v", NL, id)

	canonicalRequestHashHex := hex.EncodeToString(hash([]byte(canonicalRequest)))
	//fmt.Printf("%vAUTH(CANONICALREQUESTHASHHEX): %v", NL, canonicalRequestHashHex)

	stringToSign :=
	Algorithm +
	NL +
	timestamp +
	NL +
	id +
	NL +
	canonicalRequestHashHex
	//fmt.Printf("%vAUTH(STRINGTOSIGN): %v", NL, stringToSign)

	secret := []byte(AuthenticationScheme + credentials.Secret)
	//fmt.Printf("%vAUTH(SECRET): %v", NL, secret)

	singDate := sing(dateStamp, secret)
	//fmt.Printf("%vAUTH(SINGDATE): %v", NL, singDate)

	singNonce := sing(nonce, singDate)
	//fmt.Printf("%vAUTH(SINGNONCE): %v", NL, singNonce)

	signing := sing(IDTerminator, singNonce)
	//fmt.Printf("%vAUTH(SIGNING): %v", NL, signing)

	signature := sing(stringToSign, signing)
	//fmt.Printf("%vAUTH(SIGNATURE): %v", NL, signature)

	signatureHex := hex.EncodeToString(signature)
	//fmt.Printf("%vAUTH(SIGNATUREHEX): %v", NL, signatureHex)

	authorizationHeader :=
	AuthenticationScheme + " " +
	createNameValuePair(SAUTHC1Id, id) + ", " +
	createNameValuePair(SAUTHC1SignedHeaders, signedHeadersString) + ", " +
	createNameValuePair(SAUTHC1Signature, signatureHex)
	//log.Printf("FOO: AuthorizationHeader: %v", authorizationHeader)
	req.Header.Set(AuthorizationHeader, authorizationHeader)
}

func createNameValuePair(name string, value string) string {
	return name + "=" + value
}

func encodeURL(value string, path bool, canonical bool) string {
	if value == "" {
		return ""
	}

	encoded := url.QueryEscape(value)

	if canonical {
		encoded = strings.Replace(encoded, "+", "%20", -1)
		encoded = strings.Replace(encoded, "*", "%2A", -1)
		encoded = strings.Replace(encoded, "%7E", "~", -1)

		if path {
			encoded = strings.Replace(encoded, "%2F", "/", -1)
		}
	}

	return encoded
}

func canonicalizeQueryString(queryValues url.Values) string {
	stringBuffer := bytes.NewBufferString("")

	keys := sortedMapKeys(queryValues)

	for _, k := range keys {
		key := encodeURL(k, false, true)
		v := queryValues[k]
		for _, vv := range v {
			value := encodeURL(vv, false, true)

			if stringBuffer.Len() > 0 {
				stringBuffer.WriteString("&")
			}

			stringBuffer.WriteString(key + "=" + value)
		}
	}
	return stringBuffer.String()
}

func canonicalizeresourcePath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	return encodeURL(path, true, true)
}

func canonicalizeHeadersString(ctx context.Context, headers http.Header) string {
	stringBuffer := bytes.NewBufferString("")

	keys := sortedMapKeys(headers)
	//log.Printf("KEYS: %#v", keys)

	for _, k := range keys {
		stringBuffer.WriteString(strings.ToLower(k))
		stringBuffer.WriteString(":")

		first := true

		for _, v := range headers[k] {
			if !first {
				stringBuffer.WriteString(",")
			}

			// For URL Fetch, append the Google User Agent.
			if strings.ToLower(k) == strings.ToLower("User-Agent") {

				appName := appengine.AppID(ctx)
				//var devText string
				//if appengine.IsDevAppServer() {
				devText := "dev~"
				//}

				v = v + fmt.Sprintf(" AppEngine-Google; (+http://code.google.com/appengine; appid: %v%v)", devText, appName)
				//fmt.Printf("%vNew User-Agent: %v", NL, v)
			}

			stringBuffer.WriteString(v)
			first = false
		}
		stringBuffer.WriteString(NL)
	}

	//log.Printf("HEADERS: %v", stringBuffer.String())
	return stringBuffer.String()
}

func canonicalizeHeadersStringWithoutUserAgent(ctx context.Context, headers http.Header) string {
	stringBuffer := bytes.NewBufferString("")

	keys := sortedMapKeys(headers)
	//log.Printf("KEYS: %#v", keys)

	for _, k := range keys {
		if strings.ToLower(k) != strings.ToLower("User-Agent") {
			stringBuffer.WriteString(strings.ToLower(k))
			stringBuffer.WriteString(":")

			first := true

			for _, v := range headers[k] {

				//if k == strings.ToLower("Content-Type") && v == "" {
				// Skip an empty content-type.
				//}else {
				if !first {
					stringBuffer.WriteString(",")
				}

				stringBuffer.WriteString(v)
				first = false
				//}
			}
			stringBuffer.WriteString(NL)
		}
	}

	//fmt.Printf("HEADERS without UA: %v", stringBuffer.String())
	return stringBuffer.String()
}

//func isContentTypeEmpty(headers http.Header) bool {
//	keys := sortedMapKeys(headers)
//	for _, k := range keys {
//		if strings.ToLower(k) == strings.ToLower("Content-Type") {
//			for _, v := range headers[k] {
//				if v != "" {
//					// If we get here, it means there's at least one non-empty Content-Type header.
//					return false
//				}
//			}
//		}
//	}
//	return true
//}

func signedHeadersStringWithoutUserAgent(headers http.Header) string {
	stringBuffer := bytes.NewBufferString("")

	keys := sortedMapKeys(headers)

	for _, k := range keys {
		if strings.ToLower(k) != strings.ToLower("User-Agent") {
			//if skipContentType && strings.ToLower(k) == strings.ToLower("Content-Type") {
			// Skip an empty content-type.
			//	} else {
			if stringBuffer.Len() > 0 {
				stringBuffer.WriteString(";")
			}
			stringBuffer.WriteString(strings.ToLower(k))
			//	}
		}
	}

	return stringBuffer.String()
}

func sortedMapKeys(m map[string][]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func sing(data string, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}
