package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sappenin/stormpath-sdk-go/web"

	"github.com/codegangsta/negroni"
	"github.com/sappenin/stormpath-sdk-go"
	"github.com/julienschmidt/httprouter"

	"github.com/gorilla/sessions"
	"google.golang.org/appengine/aetest"
)

var indexHTML = `
<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title></title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">
    </head>
    <body>
    	Hello! <a href="/login">Login</a>
    </body>
</html>
`

var appHTML = `
<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title></title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">
    </head>
    <body>
    	Cool your in! <a href="/logout">Logout</a>
    </body>
</html>
`

const sessionName = "go-sdk-demo"

var store = sessions.NewCookieStore([]byte("go-sdk-demo"))

var ctx, done, ctx_err = aetest.NewContext()

// TODO: Wrap the Negroni to supply a context.

func main() {

	credentials, _ := stormpath.NewDefaultCredentials()
	stormpath.Init(credentials, nil)

	n := negroni.Classic()

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if stormpathweb.IsAuthenticated(r) {
			http.Redirect(w, r, "/app", http.StatusFound)
			return
		}
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprint(w, indexHTML)
	})

	router.Handler("GET", "/login", stormpathweb.ContextHandler{Real: loginHandler()})
	router.Handler("GET", "/logout", stormpathweb.ContextHandler{Real: logoutHandler()})
	router.Handler("GET", "/callback", stormpathweb.ContextHandler{Real: callbackHandler()})

	authRouter := httprouter.New()

	authRouter.GET("/app", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprint(w, appHTML)
	})

	router.Handler("GET", "/app", negroni.New(
		negroni.HandlerFunc(authenticationMiddleware),
		negroni.Wrap(authRouter),
	))

	n.UseHandler(applicationMiddleware())
	n.UseHandler(accountMiddleware())
	n.UseHandler(router)

	n.Run(":9999")

	defer done()
}

func authenticationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m := stormpathweb.AuthenticationMiddleware{
		Next:                next,
		SessionStore:        store,
		SessionName:         sessionName,
		UnauthorizedHandler: http.HandlerFunc(unauthorizedHandler),
	}
	m.ServeHTTP(rw, r)
}

func unauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?unauthorize", http.StatusFound)
}

func accountMiddleware() stormpathweb.AccountMiddleware {
	return stormpathweb.AccountMiddleware{
		SessionStore: store,
		SessionName:  sessionName,
	}
}

func applicationMiddleware() stormpathweb.ApplicationMiddleware {
	return stormpathweb.ApplicationMiddleware{
		ApplicationHref: os.Getenv("APPLICATION_HREF"),
	}
}

func loginHandler() stormpathweb.ContextHandlerFunc {
	return stormpathweb.IDSiteLoginHandler{
		Options: map[string]string{"callbackURI": "/callback"},
	}.ServeHTTP
}

func logoutHandler() stormpathweb.ContextHandlerFunc {
	return stormpathweb.IDSiteLogoutHandler{
		Options: map[string]string{"callbackURI": "/callback"},
	}.ServeHTTP
}

func callbackHandler() stormpathweb.ContextHandlerFunc {
	return stormpathweb.IDSiteAuthCallbackHandler{
		SessionStore:      store,
		SessionName:       sessionName,
		LoginRedirectURI:  "/app",
		LogoutRedirectURI: "/",
		ErrorHandler:      http.HandlerFunc(idSiteErrorHandler),
	}.ServeHTTP
}

func idSiteErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Ooops", http.StatusInternalServerError)
}
