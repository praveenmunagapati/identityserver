package company

import (
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/identityservice/company/api"
	"github.com/itsyouonline/identityserver/identityservice/company/models"
)

func AddRoutes(r *mux.Router) {
	companyApiRoutes := api.NewCompanyResource().GetRoutes()

	for _, route := range companyApiRoutes {
		r.
			Methods(route.Methods...).
			Name(route.Name).
			Path(route.Path).
			Handler(route.HandlerFunc)
	}

	// Create indices ...
	go models.InitModels()
}
