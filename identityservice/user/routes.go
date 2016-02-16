package user

import (
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/identityservice/user/api"
)

func AddRoutes(r *mux.Router) {
	userApiRoutes := api.NewUserResource().GetRoutes()

	for _, route := range userApiRoutes {
		r.
			Methods(route.Methods...).
			Name(route.Name).
			Path(route.Path).
			Handler(route.HandlerFunc)
	}
}
