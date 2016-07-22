package organization

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/oauthservice"
)

// Oauth2oauth_2_0Middleware is oauth2 middleware for oauth_2_0
type Oauth2oauth_2_0Middleware struct {
	describedBy string
	field       string
	scopes      []string
}

// newOauth2oauth_2_0Middlewarecreate new Oauth2oauth_2_0Middleware struct
func newOauth2oauth_2_0Middleware(scopes []string) *Oauth2oauth_2_0Middleware {
	om := Oauth2oauth_2_0Middleware{
		scopes: scopes,
	}

	om.describedBy = "headers"
	om.field = "Authorization"

	return &om
}

// CheckScopes checks whether user has needed scopes
func (om *Oauth2oauth_2_0Middleware) CheckScopes(scopes []string) bool {
	if len(om.scopes) == 0 {
		return true
	}

	for _, allowed := range om.scopes {
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
		if om.describedBy == "queryParameters" {
			accessToken = r.URL.Query().Get(om.field)
		} else if om.describedBy == "headers" {
			accessToken = r.Header.Get(om.field)
		}

		var scopes []string
		protectedOrganization := mux.Vars(r)["globalid"]
		var atscopestring string
		var username string
		var clientID string
		var globalID string

		if strings.HasPrefix(accessToken, "token") {
			accessToken = strings.TrimSpace(strings.TrimPrefix(accessToken, "token"))
			if accessToken == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			log.Debug("Access Token: ", accessToken)
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
			globalID = at.GlobalID
			username = at.Username
			atscopestring = at.Scope
			clientID = at.ClientID
		} else {
			if webuser, ok := context.GetOk(r, "webuser"); ok {
				if parsedusername, ok := webuser.(string); ok && parsedusername != "" {
					username = parsedusername
					atscopestring = "admin"
					clientID = "itsyouonline"
				}
			}
		}
		if (username == "" && globalID == "") || clientID == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		context.Set(r, "authenticateduser", username)
		if globalID == protectedOrganization {
			scopes = []string{atscopestring}
		} else {
			orgMgr := organization.NewManager(r)
			isOwner, err := orgMgr.IsOwner(protectedOrganization, username)
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if isOwner && clientID == "itsyouonline" && atscopestring == "admin" {
				scopes = []string{"organization:owner"}
			} else {
				isMember, err := orgMgr.IsMember(protectedOrganization, username)
				if err != nil {
					log.Error(err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if isMember && clientID == "itsyouonline" && atscopestring == "admin" {
					scopes = []string{"organization:member"}
				}
			}
		}

		//TODO: scopes "organization:info", "organization:contracts:read"

		log.Debug("Available scopes: ", scopes)

		// check scopes
		if !om.CheckScopes(scopes) {
			w.WriteHeader(403)
			return
		}

		next.ServeHTTP(w, r)
	})
}
