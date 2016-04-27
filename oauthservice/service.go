package oauthservice

import (
	"net/http"

	"github.com/gorilla/mux"
)

//SessionService declares a context where you can have a logged in user
type SessionService interface {
	//GetLoggedInUser returns an authenticated user, or an empty string if there is none
	GetLoggedInUser(request *http.Request) (username string, err error)
	//SetAPIAccessToken sets the api access token for this session
	SetAPIAccessToken(w http.ResponseWriter, token string) (err error)
}

//IdentityService provides some basic knowledge about authorizations required for the oauthservice
type IdentityService interface {
	//ValidAuthorizationForScopes checks if there is a valid authorization for the requested scopes
	ValidAuthorizationForScopes(r *http.Request, username string, grantedTo string, scopes string) (valid bool, err error)
}

//Service is the oauthserver http service
type Service struct {
	sessionService  SessionService
	identityService IdentityService
	router          *mux.Router
}

//NewService creates and initializes a Service
func NewService(sessionService SessionService, identityService IdentityService) *Service {
	return &Service{sessionService: sessionService, identityService: identityService}
}

const (
	//AuthorizationGrantCodeType is the requested response_type for an 'authorization code' oauth2 flow
	AuthorizationGrantCodeType = "code"
	//ImplicitGrantCodeType is the requested response_type for an 'implicit' oauth2 flow
	ImplicitGrantCodeType = "token"
	//ClientCredentialsGrantCodeType is the requested grant_type for a 'client credentials' oauth2 flow
	ClientCredentialsGrantCodeType = "client_credentials"
)

//GetAuthenticatedUser returns the authenticated user if any or an empty string if not
func (service *Service) GetAuthenticatedUser(r *http.Request) (username string, err error) {
	username, err = service.sessionService.GetLoggedInUser(r)
	return
}

//AddRoutes adds the routes and handlerfunctions to the router
func (service *Service) AddRoutes(router *mux.Router) {
	service.router = router
	router.HandleFunc("/v1/oauth/authorize", service.AuthorizeHandler).Methods("GET")
	router.HandleFunc("/v1/oauth/access_token", service.AccessTokenHandler).Methods("POST")

	InitModels()
}
