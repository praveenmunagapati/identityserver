package oauthservice

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

//JWTHandler returns a JWT with claims that are a subset of the scopes available to the authorizing token
func (service *Service) JWTHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Debug("Error parsing form: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	accessToken := r.Header.Get("Authorization")

	//Get the actual token out of the header (accept 'token ABCD' as well as just 'ABCD' and ignore some possible whitespace)
	accessToken = strings.TrimSpace(strings.TrimPrefix(accessToken, "token"))
	if accessToken == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	oauthMgr := NewManager(r)
	at, err := oauthMgr.GetAccessToken(accessToken)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if at == nil || at.IsExpired() {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	requestedScopes := r.FormValue("scope")
	extraAudiences := strings.TrimSpace(r.FormValue("aud"))

	if !jwtScopesAreAllowed(at.Scope, requestedScopes) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodES384)
	if at.Username != "" {
		token.Claims["username"] = at.Username
	}
	if at.GlobalID != "" {
		token.Claims["globalid"] = at.GlobalID
	}

	audiences := []string{at.ClientID}
	if extraAudiences != "" {
		audiences = append(audiences, strings.Split(extraAudiences, ",")...)
	}

	token.Claims["aud"] = audiences
	token.Claims["exp"] = at.ExpirationTime().Unix()
	token.Claims["iss"] = "itsyouonline"
	token.Claims["scope"] = requestedScopes

	tokenString, err := token.SignedString(service.jwtSigningKey)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}

func jwtScopesAreAllowed(allowedScopes string, requestedScopes string) (valid bool) {
	if strings.TrimSpace(requestedScopes) == "" {
		valid = true
		return
	}

	//Split and clean the scope string in to seperate scopes
	var allowedScopesList []string

	for _, value := range strings.Split(allowedScopes, ",") {
		scope := strings.TrimSpace(value)
		if scope != "" {
			allowedScopesList = append(allowedScopesList, scope)
		}
	}
	requestedScopesList := strings.Split(requestedScopes, ",")
	for i, value := range requestedScopesList {
		requestedScopesList[i] = strings.TrimSpace(value)
	}

	valid = true
	for _, rs := range requestedScopesList {
		log.Info(fmt.Sprintf("Checking if '%s' is allowed", rs))
		valid = valid && checkIfScopeInList(allowedScopesList, rs)
	}

	return
}

func checkIfScopeInList(allowedScopes []string, scope string) (valid bool) {
	for _, as := range allowedScopes {
		if strings.HasPrefix(scope, as) {
			valid = true
			return
		}
	}
	return
}
