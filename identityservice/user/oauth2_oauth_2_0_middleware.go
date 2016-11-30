package user

import (
	"net/http"
	"strings"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/identityservice/security"
	"github.com/itsyouonline/identityserver/oauthservice"
	"regexp"
)

// Oauth2oauth_2_0Middleware is oauth2 middleware for oauth_2_0
type Oauth2oauth_2_0Middleware struct {
	security.OAuth2Middleware
}

// newOauth2oauth_2_0Middlewarecreate new Oauth2oauth_2_0Middleware struct
func newOauth2oauth_2_0Middleware(scopes []string) *Oauth2oauth_2_0Middleware {
	om := Oauth2oauth_2_0Middleware{}
	om.Scopes = scopes
	return &om
}

// CheckScopes checks whether user has needed scopes
func checkScopes(possibleScopes []string, authorizedScopes []string) bool {
	if len(possibleScopes) == 0 {
		return true
	}

	for _, allowed := range possibleScopes {
		for _, scope := range authorizedScopes {
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
		var atscopestring string
		var username string
		var clientID string
		authorizedScopes := []string{}

		jwtstring := om.GetJWT(r)
		accessToken := om.GetAccessToken(r)

		if jwtstring != "" {
			token, err := jwt.Parse(jwtstring, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if token.Method != jwt.SigningMethodES384 {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return &security.JWTPublicKey, nil
			})
			if err != nil || !token.Valid {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			username = token.Claims["username"].(string)
			clientID = token.Claims["aud"].(string)
			atscopestring = token.Claims["scope"].(string)

		} else if accessToken != "" {
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
		if username == "" || clientID == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		protectedUsername := mux.Vars(r)["username"]

		if protectedUsername == username && clientID == "itsyouonline" && atscopestring == "admin" {
			authorizedScopes = append(authorizedScopes, "user:admin")
		}
		if strings.HasPrefix(atscopestring, "user:") {
			authorizedScopes = append(authorizedScopes, "user:info")
		}
		for _, scope := range strings.Split(atscopestring, ",") {
			scope = strings.Trim(scope, " ")
			authorizedScopes = append(authorizedScopes, scope)
		}

		context.Set(r, "client_id", clientID)
		context.Set(r, "availablescopes", atscopestring)

		possibleScopes := []string{}
		// Replace {variables} in the scopes with the real vars present in the url.
		for _, scope := range om.Scopes {
			regex, _ := regexp.Compile(`\{(.*?)\}`)
			vars := regex.FindStringSubmatch(scope)
			if vars != nil {
				// e.g. vars -> ["{label}", "label"]
				re, _ := regexp.Compile(vars[0])
				urlValue := mux.Vars(r)[vars[1]]
				scope = re.ReplaceAllString(scope, urlValue)
			}
			possibleScopes = append(possibleScopes, scope)
		}
		if clientID != "itsyouonline" {
			// todo: cache
			userMgr := user.NewManager(r)
			authorization, err := userMgr.GetAuthorization(protectedUsername, clientID)
			if err != nil {
				log.Error("Error while getting authorization: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if authorization == nil {
				log.Debugf("Authorization for client id %s not found for user %s", clientID, protectedUsername)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			authorizedScopes = authorization.FilterAuthorizedScopes(authorizedScopes)

		}
		// check scopes
		log.Debug("Authorized scopes: ", authorizedScopes)
		log.Debug("Needed possible scopes: ", possibleScopes)
		if !checkScopes(possibleScopes, authorizedScopes) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
