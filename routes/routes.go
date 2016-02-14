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

func GetRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	identityservice.NewService().AddRoutes(r)
	oauthservice.NewService().AddRoutes(r)
	siteservice.NewService().AddRoutes(r)

	// Add middlewares
	router := NewRouter(r)

	dbmw := db.DBMiddleware()
	recovery := handlers.RecoveryHandler()

	router.Use(recovery, LoggingMiddleware, dbmw)

	return router.Handler()
}
