package identityservice

import (
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/identityservice/company"
	"github.com/itsyouonline/identityserver/identityservice/user"
)

//Service is the identityserver http service
type Service struct {
}

func NewService() *Service {
	return &Service{}
}

//AddRoutes registers the http routes with the router.
func (service *Service) AddRoutes(router *mux.Router) {
	// User API
	user.UsersInterfaceRoutes(router, user.UsersAPI{})
	user.InitModels()

	// Company API
	company.CompaniesInterfaceRoutes(router, company.CompaniesAPI{})
	company.InitModels()
}
