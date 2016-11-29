package oauthservice

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

var errUnauthorized = errors.New("Unauthorized")

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

	requestedScopeParameter := r.FormValue("scope")

	audiences := strings.TrimSpace(r.FormValue("aud"))
	tokenString, err := service.convertAccessTokenToJWT(r, at, requestedScopeParameter, audiences)
	if err == errUnauthorized {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/jwt")
	w.Write([]byte(tokenString))
}

func (service *Service) convertAccessTokenToJWT(r *http.Request, at *AccessToken, requestedScopeString, audiences string) (tokenString string, err error) {

	requestedScopes := splitScopeString(requestedScopeString)
	acquiredScopes := splitScopeString(at.Scope)

	if !jwtScopesAreAllowed(acquiredScopes, requestedScopes) {
		err = errUnauthorized
		return
	}

	token := jwt.New(jwt.SigningMethodES384)

	if at.Username != "" {
		token.Claims["username"] = at.Username
		possibleScopes, e := service.filterPossibleScopes(r, at.Username, requestedScopes, false)
		if e != nil {
			err = e
			return
		}
		token.Claims["scope"] = strings.Join(possibleScopes, ",")
	}
	if at.GlobalID != "" {
		token.Claims["globalid"] = at.GlobalID
		token.Claims["scope"] = requestedScopes
	}

	audiencesArr := strings.Split(audiences, ",")
	if len(audiencesArr) > 0 {
		token.Claims["aud"] = audiencesArr

		// azp claim is only needed when the ID Token has a single
		// audience value and that audience is different than the authorized
		// party
		if len(audiencesArr) == 1 && audiences != at.ClientID {
			token.Claims["azp"] = at.ClientID
		}
	}
	token.Claims["exp"] = at.ExpirationTime().Unix()
	token.Claims["iss"] = "itsyouonline"

	tokenString, err = token.SignedString(service.jwtSigningKey)
	return
}

func jwtScopesAreAllowed(grantedScopes []string, requestedScopes []string) (valid bool) {
	valid = true
	for _, rs := range requestedScopes {
		log.Debug(fmt.Sprintf("Checking if '%s' is allowed", rs))
		valid = valid && checkIfScopeInList(grantedScopes, rs)
	}

	return
}

func checkIfScopeInList(grantedScopes []string, scope string) (valid bool) {
	for _, as := range grantedScopes {
		//Allow all user scopes if the 'user:admin' scope is part of the autorized scopes
		if as == "user:admin" {
			if strings.HasPrefix(scope, "user:") {
				valid = true
				return
			}
		}
		if strings.HasPrefix(scope, as) {
			valid = true
			return
		}
	}
	return
}
