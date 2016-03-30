package user

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/oauthservice"
)

// Oauth2oauth_2_0Middleware is oauth2 middleware for oauth_2_0
type Oauth2oauth_2_0Middleware struct {
	DescribedBy string
	Field       string
	Scopes      []string
}

// newOauth2oauth_2_0Middlewarecreate new Oauth2oauth_2_0Middleware struct
func newOauth2oauth_2_0Middleware(scopes []string) *Oauth2oauth_2_0Middleware {
	om := Oauth2oauth_2_0Middleware{
		Scopes: scopes,
	}

	om.DescribedBy = "headers"
	om.Field = "Authorization"

	return &om
}

// CheckScopes checks whether user has needed scopes
func (om *Oauth2oauth_2_0Middleware) CheckScopes(scopes []string) bool {
	if len(om.Scopes) == 0 {
		return true
	}

	for _, allowed := range om.Scopes {
		for _, scope := range scopes {
			if scope == allowed {
				return true
			}
		}
	}
	return false
}

// Handler return HTTP handler representation of this middleware
func (om *Oauth2oauth_2_0Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessToken string

		// access token checking
		if om.DescribedBy == "queryParameters" {
			accessToken = r.URL.Query().Get(om.Field)
		} else if om.DescribedBy == "headers" {
			accessToken = r.Header.Get(om.Field)
		}
		//Get the actual token out of the header (accept 'token ABCD' as well as just 'ABCD' and ignore some possible whitespace)
		accessToken = strings.TrimSpace(strings.TrimPrefix(accessToken, "token"))
		if accessToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Debug("Access Token: ", accessToken)
		var scopes []string
		//TODO: cache
		oauthMgr := oauthservice.NewManager(r)
		at, err := oauthMgr.GetAccessToken(accessToken)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if at == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		protectedUsername := mux.Vars(r)["username"]

		if protectedUsername == at.Username && at.ClientID == "itsyouonline" && at.Scope == "admin" {
			scopes = []string{"user:admin"}
		}
		// TODO: "user:info" scope

		log.Debug("Available scopes: ", scopes)

		// check scopes
		if !om.CheckScopes(scopes) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
