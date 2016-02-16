package oauthservice

import (
	"github.com/gorilla/mux"
)

//Service is the oauthserver http service
type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (service *Service) AddRoutes(router *mux.Router) {
}
