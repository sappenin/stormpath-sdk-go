package stormpathweb

import (
	"encoding/json"
	"net/http"

	gorillia_context "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/sappenin/stormpath-sdk-go"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

//ApplicationMiddleware is an http.Handler that stores a given account in the request context
//to be use by any other handlers in the chain.
type ApplicationMiddleware struct {
	ApplicationHref string
}

//ServeHTTP implements the http.Handler interface for the ApplicationMiddleware type
func (m ApplicationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Check if it the current app already exists
	ctx := appengine.NewContext(r)
	app := GetApplication(r)
	if app == nil {
		app, err := stormpath.GetApplication(ctx, m.ApplicationHref, stormpath.MakeApplicationCriteria())
		if err == nil {
			gorillia_context.Set(r, ApplicationKey, *app)
			log.Debugf(ctx, "ApplicationMiddleware.ServeHTTP(): Successfully set Application into Context with Key '%v'", ApplicationKey)
		} else {
			log.Debugf(ctx, "ApplicationMiddleware.ServeHTTP(): Unable to fetch Application from StormPath: %v", err)
			panic(err)
		}
	} else {
		log.Debugf(ctx, "ApplicationMiddleware.ServeHTTP(): Application was: %#v", app)
	}
}

//AccountMiddleware is an http.Handler that unmarshals the current account store in the session
//and stores it in the request context to be use by any other handler in the chain
type AccountMiddleware struct {
	SessionStore sessions.Store
	SessionName  string
}

//ServeHTTP implements the http.Handler interface for the AccountMiddleware type
func (m AccountMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "AccountMiddleware.ServeHTTP(): Operating on Session with Name: '%v'", m.SessionName)
	session, _ := m.SessionStore.Get(r, m.SessionName)
	if session.Values[AccountKey] != nil {
		account := stormpath.Account{}

		json.Unmarshal([]byte(session.Values[AccountKey].([]uint8)), &account)
		gorillia_context.Set(r, AccountKey, account)
		log.Debugf(ctx, "AccountMiddleware.ServeHTTP(): Successfully set Account into Context with Key '%v'", AccountKey)
	} else {
		log.Debugf(ctx, "AccountMiddleware.ServeHTTP(): No Account existed in session for Key '%v'", AccountKey)
	}
}

//AuthenticationMiddleware handles authentication for a web application, it should only be apply to http.Handlers
//that require authentication it checks the session for current account if exists it calls handler else it applies
//the UnauthorizedHandler
type AuthenticationMiddleware struct {
	Next                http.Handler
	SessionStore        sessions.Store
	SessionName         string
	UnauthorizedHandler http.Handler
}

//ServeHTTP implements the http.Handler interface for the AuthenticationMiddleware type
func (m AuthenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, _ := m.SessionStore.Get(r, m.SessionName)

	if session.Values[AccountKey] == nil {
		//No account in session
		m.UnauthorizedHandler.ServeHTTP(w, r)
		return
	}

	//We are good move along
	if m.Next != nil {
		m.Next.ServeHTTP(w, r)
	}
}
