package stormpathweb

import (
	"github.com/sappenin/stormpath-sdk-go"
	"net/http"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine"
	gorillia_context "github.com/gorilla/context"
)

//ApplicationKey is the key of the current application in the context
const ApplicationKey = "application"

//AccountKey is the key of the current account in the context and session
const AccountKey = "account"

//GetApplication returns the application from the context previously set by the ApplicationMiddleware
func GetApplication(r *http.Request) *stormpath.Application {
	ctx := appengine.NewContext(r)
	app := gorillia_context.Get(r, ApplicationKey)
	if app == nil {
		log.Debugf(ctx, "Application did NOT exist in Context with Key '%v'", ApplicationKey)
		return nil
	} else {
		application := app.(stormpath.Application)
		log.Debugf(ctx, "Application DID exist in Context with Key '%v': %#v", application, ApplicationKey)
		return &application
	}
}

//GetCurrentAccount retrieves the current account if any from the request context
func GetCurrentAccount(r *http.Request) *stormpath.Account {
	ctx := appengine.NewContext(r)
	acc := gorillia_context.Get(r, AccountKey)
	if acc == nil {
		log.Debugf(ctx, "Account did NOT exist in Context with Key '%v'.", AccountKey)
		return nil
	} else {
		account := acc.(stormpath.Account)
		log.Debugf(ctx, "Account DID exist in Context with Key '%v':  %#v", AccountKey, account)
		return &account
	}
}

//IsAuthenticated checks if there is an authenticated user
func IsAuthenticated(r *http.Request) bool {
	return GetCurrentAccount(r) != nil
}
