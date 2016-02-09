package site

import (
	"github.com/gorilla/mux"
)

//Service is the identityserver http service
type Service struct {
}

//AddRoutes registers the http routes with the router
func (service *Service) AddRoutes(router *mux.Router) {
	router.Methods("GET").Path("register").HandlerFunc(service.ShowRegistrationForm)
	router.Methods("POST").Path("register").HandlerFunc(service.ProcessRegistrationForm)
}
