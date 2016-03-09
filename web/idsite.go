package stormpathweb

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/sappenin/stormpath-sdk-go"
	"google.golang.org/appengine/log"
	"fmt"
	"net/url"
	"google.golang.org/appengine"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	gorilla_context "github.com/gorilla/context"
)


//Our appengine.Context http handler.  Essentially, this is a http.Handler that adds an appending context.
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// A container that holds the Handler for GAE.
type ContextHandler struct {
	Real ContextHandlerFunc
}

const REQUEST_ID = "request-id"

// Makes ContextHandler conform to http.Handler
func (f ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// For each appengine request, we create a UUID to uniquely identify the request.
	if rid, err := uuid.NewV4(); err != nil {
		panic(err)
	} else {
		gorilla_context.Set(r, REQUEST_ID, rid)
		ctx := appengine.NewContext(r)
		f.Real(ctx, w, r)
	}
}

//IDSiteLoginHandler is an http.Handler for Strompath's IDSite login
type IDSiteLoginHandler struct {
	Options map[string]string
}

//IDSiteLogoutHandler is an http.Handler for Strompath's IDSite logout
type IDSiteLogoutHandler struct {
	Options map[string]string
}

//ServeHTTP implements the http.Handler interface for IDSiteLoginHandler type and ContextHandlerFunc to support App Engine Contexts
func (h IDSiteLoginHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idSiteURLHandler(ctx, w, r, ensureOption("logout", "", h.Options))
}

// Implement ContextHandlerFunc for IDSiteLogoutHandler
func (h IDSiteLogoutHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idSiteURLHandler(ctx, w, r, ensureOption("logout", "true", h.Options))
}

func idSiteURLHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, options map[string]string) {
	if options["callbackURI"][0] == '/' {
		u, _ := url.Parse(r.Header.Get("Referer"))
		callbackURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, options["callbackURI"])
		options["callbackURI"] = callbackURL
	}

	idSiteURL, _ := GetApplication(r).CreateIDSiteURL(ctx, options)
	http.Redirect(w, r, idSiteURL, http.StatusFound)
}

func ensureOption(key string, value string, options map[string]string) map[string]string {
	if options == nil {
		options = make(map[string]string)
	}
	options[key] = value
	return options
}

//IDSiteAuthCallbackHandler is an http.Handler for the ID Site callback
type IDSiteAuthCallbackHandler struct {
	SessionStore      sessions.Store
	SessionName       string
	LoginRedirectURI  string
	LogoutRedirectURI string
	ErrorHandler      http.Handler
}

//ServeHTTP implements the http.Handler interface for the IDSiteAuthCallbackHandler type and ContextHandlerFunc to support App Engine Contexts
// Implement  for IDSiteAuthCallbackHandler
func (h IDSiteAuthCallbackHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	app := GetApplication(r)
	if app == nil {
		log.Debugf(ctx, "No StormPath Application found in Context.  Clearing Account Session!")
		h.clearAccountInSession(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	log.Debugf(ctx, "StormPath Application was found in Context: %#v", app)
	result, err := app.HandleIDSiteCallback(ctx, r.URL.String())

	// TODO: Make these nicer.  If there's a JWT error of some sort, we should redirect the User to the Oops page.
	if err != nil {
		log.Debugf(ctx, "IDSite %s", err)
		h.ErrorHandler.ServeHTTP(w, r)
		return
	}

	if result.Status == "AUTHENTICATED" {
		//Login succesful
		log.Debugf(ctx, "Login Successful!  Storing Account in Session.")
		h.storeAccountInSession(result.Account, w, r)
		http.Redirect(w, r, h.LoginRedirectURI, http.StatusFound)
	} else {
		//Logout
		log.Debugf(ctx, "Logout Successful!  Clearing Account from Session.")
		h.clearAccountInSession(w, r)
		http.Redirect(w, r, h.LogoutRedirectURI, http.StatusFound)
	}
}

//StoreAccountInSession stores a given account in the session as the current account
func (h IDSiteAuthCallbackHandler) storeAccountInSession(account *stormpath.Account, w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionStore.Get(r, h.SessionName)

	jsonBody, _ := json.Marshal(account)

	session.Values[AccountKey] = jsonBody
	session.Save(r, w)
}

//ClearAccountInSession removes the current account form the session
func (h IDSiteAuthCallbackHandler) clearAccountInSession(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionStore.Get(r, h.SessionName)

	session.Values[AccountKey] = nil
	session.Save(r, w)
}
