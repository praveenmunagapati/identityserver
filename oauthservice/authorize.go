package oauthservice

import (
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	//AuthorizationGrantCodeType is requested response_type for an 'authorization code' oauth2 flow
	AuthorizationGrantCodeType = "code"
)

func (service *Service) validateRedirectURI(redirectURI, clientID string) {
	//TODO:
}

//Authorize is the handler of the /login/oauth/authorize endpoint
func (service *Service) Authorize(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Check if the requested authorization grant type is supported
	if r.Form.Get("response_type") != AuthorizationGrantCodeType {
		log.Debug("Invalid authorization grant type requested")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Check if the user is already authenticated, if not, redirect to the login page before returning here
	username, err := service.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if username == "" {
		queryvalues := r.URL.Query()
		queryvalues.Add("endpoint", r.URL.EscapedPath())
		//TODO: redirect according the the received http method
		http.Redirect(w, r, "/login?"+queryvalues.Encode(), http.StatusFound)
	}

	//Validate client and redirect_uri
	redirectURI, err := url.QueryUnescape(r.Form.Get("redirect_uri"))
	if err != nil {
		log.Debug("Unparsable redirect_uri")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	clientID := r.Form.Get("client_id")
	service.validateRedirectURI(redirectURI, clientID)

	//requestedScopes := r.Form.Get("scope")
	//TODO: check if the client still has a valid authorization for the requested scope, if not ask the user

	clientState := r.Form.Get("state")
	//TODO: generate and store the Authorization Code
	authorizationCode := "FIXTHIS"

	//TODO: store the state to pass it when redirecting
	parameters := make(url.Values)
	parameters.Add("code", authorizationCode)
	parameters.Add("state", clientState)
	//Don't parse the redirect url, can only give errors while we don't gain much
	if !strings.Contains(redirectURI, "?") {
		redirectURI += "?"
	} else {
		if !strings.HasSuffix(redirectURI, "&") {
			redirectURI += "&"
		}
	}
	redirectURI += parameters.Encode()

	http.Redirect(w, r, redirectURI, http.StatusFound)

}
