package oauthservice

import (
	"crypto/ecdsa"
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
	//FilterAuthorizedScopes filters the requested scopes to the ones that are authorizated, if no authorization exists, authorizedScops is nil
	FilterAuthorizedScopes(r *http.Request, username string, grantedTo string, requestedscopes []string) (authorizedScopes []string, err error)
	//FilterPossibleScopes filters the requestedScopes to the relevant ones that are possible
	// For example, a `user:memberof:orgid1` is not possible if the user is not a member the `orgid1` organization
	FilterPossibleScopes(r *http.Request, username string, clientID string, requestedScopes []string) (possibleScopes []string, err error)
}

//Service is the oauthserver http service
type Service struct {
	sessionService  SessionService
	identityService IdentityService
	router          *mux.Router
	jwtSigningKey   *ecdsa.PrivateKey
}

//NewService creates and initializes a Service
func NewService(sessionService SessionService, identityService IdentityService, ecdsaKey *ecdsa.PrivateKey) (service *Service, err error) {
	service = &Service{sessionService: sessionService, identityService: identityService, jwtSigningKey: ecdsaKey}
	return
}

const (
	//AuthorizationGrantCodeType is the requested response_type for an 'authorization code' oauth2 flow
	AuthorizationGrantCodeType = "code"
	//ClientCredentialsGrantCodeType is the requested grant_type for a 'client credentials' oauth2 flow
	ClientCredentialsGrantCodeType = "client_credentials"
)

//GetWebuser returns the authenticated user if any or an empty string if not
func (service *Service) GetWebuser(r *http.Request) (username string, err error) {
	username, err = service.sessionService.GetLoggedInUser(r)
	return
}

//AddRoutes adds the routes and handlerfunctions to the router
func (service *Service) AddRoutes(router *mux.Router) {
	service.router = router
	router.HandleFunc("/v1/oauth/authorize", service.AuthorizeHandler).Methods("GET")
	router.HandleFunc("/v1/oauth/access_token", service.AccessTokenHandler).Methods("POST")
	router.HandleFunc("/v1/oauth/jwt", service.JWTHandler).Methods("POST", "GET")
	InitModels()
}
