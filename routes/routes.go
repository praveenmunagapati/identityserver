package routes

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice"
	"github.com/itsyouonline/identityserver/oauthservice"
	"github.com/itsyouonline/identityserver/siteservice"
)

//GetRouter contructs the router hierarchy and registers all handlers and middleware
func GetRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	cookieSecret := identityservice.GetCookieSecret()

	sc := siteservice.NewService(cookieSecret)
	sc.AddRoutes(r)

	apiRouter := r.PathPrefix("/api").Subrouter()
	is := identityservice.NewService()
	is.AddRoutes(apiRouter)

	oauthservice.NewService(sc, is).AddRoutes(r)

	// Add middlewares
	router := NewRouter(r)

	dbmw := db.DBMiddleware()
	recovery := handlers.RecoveryHandler()

	router.Use(recovery, LoggingMiddleware, dbmw)

	return router.Handler()
}
