package oauthservice

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/siteservice"
)

//Service is the oauthserver http service
type Service struct {
	siteService *siteservice.Service
	router      *mux.Router
}

//NewService creates and initializes a Service
func NewService(siteService *siteservice.Service) *Service {
	return &Service{siteService: siteService}
}

//GetAuthenticatedUser returns the authenticated user if any or an empty string if not
func (service *Service) GetAuthenticatedUser(r *http.Request) (username string, err error) {
	username, err = service.siteService.GetLoggedInUser(r)
	return
}

//AddRoutes adds the routes and handlerfunctions to the router
func (service *Service) AddRoutes(router *mux.Router) {
	service.router = router
	router.HandleFunc("/v1/oauth/authorize", service.Authorize).Methods("GET")
}
