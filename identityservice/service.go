package identityservice

import (
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/identityservice/company"
	"github.com/itsyouonline/identityserver/identityservice/organization"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/identityservice/userorganization"
)

//Service is the identityserver http service
type Service struct {
}

//NewService creates and initializes a Service
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

	// Organization API
	organization.OrganizationsInterfaceRoutes(router, organization.OrganizationsAPI{})
	userorganization.UserorganizationsInterfaceRoutes(router, userorganization.UsersusernameorganizationsAPI{})
	organization.InitModels()

}
