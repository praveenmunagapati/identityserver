package oauthservice

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/user/apikey"
	"gopkg.in/mgo.v2/bson"
)

//AccessTokenExpiration is the time in seconds an access token expires
var AccessTokenExpiration = time.Second * 3600 * 24 //Tokens expire after 1 day

//AccessToken is an oauth2 accesstoken together with the access information it stands for
type AccessToken struct {
	ID          bson.ObjectId `json:"-" bson:"_id,omitempty"`
	AccessToken string
	Type        string
	Username    string
	GlobalID    string //The organization that granted the token (in case of a client credentials flow)
	Scope       string
	ClientID    string //The client_id of the organization that was granted the token
	CreatedAt   time.Time
}

//IsExpiredAt checks if the token is expired at a specific time
func (at *AccessToken) IsExpiredAt(testtime time.Time) bool {
	return testtime.After(at.ExpirationTime())
}

//IsExpired is a convenience method for IsExpired(time.Now())
func (at *AccessToken) IsExpired() bool {
	return at.IsExpiredAt(time.Now())
}

//ExpirationTime return the time at which this token expires
func (at *AccessToken) ExpirationTime() time.Time {
	return at.CreatedAt.Add(AccessTokenExpiration)
}

func newAccessToken(username, globalID, clientID, scope string) *AccessToken {
	var at AccessToken

	randombytes := make([]byte, 21) //Multiple of 3 to make sure no padding is added
	rand.Read(randombytes)
	at.AccessToken = base64.URLEncoding.EncodeToString(randombytes)
	at.CreatedAt = time.Now()
	at.Username = username
	at.GlobalID = globalID
	at.ClientID = clientID
	at.Scope = scope
	at.Type = "bearer"

	return &at
}

//AccessTokenHandler is the handler of the /v1/oauth/access_token endpoint
func (service *Service) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing form: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var clientID, clientSecret string
	code := r.FormValue("code")
	grantType := r.FormValue("grant_type")
	clientSecret = r.FormValue("client_secret")
	clientID = r.FormValue("client_id")

	//If clientSecret if missing from form data check if its available as basicauth
	//See https://tools.ietf.org/html/rfc6749#section-2.3.1
	if clientSecret == "" {
		var ok bool
		clientID, clientSecret, ok = r.BasicAuth()
		if !ok {
			log.Debug("clientSecret not found in form data nor basicauth")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	//Also accept some alternatives
	if grantType == "authorization_code" {
		grantType = ""
	}

	if clientSecret == "" || clientID == "" || (grantType == "" && code == "") {
		log.Debug("Required parameter missing in the request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var at *AccessToken
	httpStatusCode := http.StatusOK

	mgr := NewManager(r)
	if grantType != "" {
		if grantType == ClientCredentialsGrantCodeType {
			at, httpStatusCode = clientCredentialsTokenHandler(clientID, clientSecret, mgr, r)
		} else {
			log.Debug("Invalid grant_type")
			httpStatusCode = http.StatusBadRequest
		}
	} else {
		redirectURI := r.FormValue("redirect_uri")
		at, httpStatusCode = convertCodeToAccessTokenHandler(code, clientID, clientSecret, redirectURI, mgr)
	}

	if httpStatusCode != http.StatusOK {
		http.Error(w, http.StatusText(httpStatusCode), httpStatusCode)
		return
	}

	// It is also possible to immediately get a JWT by specifying 'id_token' as the response type
	// In this case, the scope parameter needs to be given to prevent consumers to accidentally handing out too powerful tokens to third party services
	// It is also possible to specify additional audiences
	responseType := r.FormValue("response_type")

	if responseType == "id_token" {
		requestedScopeParameter := r.FormValue("scope")
		extraAudiences := r.FormValue("aud")

		validity := parseValidity(r)

		var tokenString string
		tokenString, err = service.convertAccessTokenToJWT(r, at, requestedScopeParameter, extraAudiences, validity)
		if err == errUnauthorized {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// if client could accept JSON we give the token as JSON string
		// if not, in plain text with the application/jwt mime-type
		if strings.Index(r.Header.Get("Accept"), "application/json") >= 0 {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"access_token": tokenString})
		} else {
			w.Header().Set("Content-type", "application/jwt")
			w.Write([]byte(tokenString))
		}
		return
	}
	mgr.saveAccessToken(at)

	orgMgr := organization.NewManager(r)
	scope, err := verifyScopes(at.Scope, at.Username, at.ClientID, orgMgr)
	if err != nil {
		log.Error("Failed to verify token scopes: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		AccessToken string      `json:"access_token"`
		TokenType   string      `json:"token_type"`
		Scope       string      `json:"scope"`
		ExpiresIn   int64       `json:"expires_in"`
		Info        interface{} `json:"info"`
	}{
		AccessToken: at.AccessToken,
		TokenType:   at.Type,
		Scope:       scope,
		ExpiresIn:   int64(at.ExpirationTime().Sub(time.Now()).Seconds() - 600),

		Info: struct {
			Username string `json:"username"`
		}{
			Username: at.Username,
		},
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

func clientCredentialsTokenHandler(clientID string, secret string, mgr *Manager, r *http.Request) (at *AccessToken, httpStatusCode int) {
	httpStatusCode = http.StatusOK
	var scopes string
	username := ""
	organization := ""

	client, err := mgr.getClientByCredentials(clientID, secret)
	if err != nil {
		log.Error("Error getting the oauth client: ", err)
		httpStatusCode = http.StatusInternalServerError
		return
	}
	if client == nil || !client.ClientCredentialsGrantType {
		log.Info("Checking user api")
		apikeyMgr := apikey.NewManager(r)
		apikey, err := apikeyMgr.GetByApplicationAndSecret(clientID, secret)
		if err != nil || apikey.ApiKey != secret {
			log.Debug("Failed to get the user api key: ", err)
			httpStatusCode = http.StatusBadRequest
			return
		}
		if apikey.ApiKey != secret {
			log.Debug("Invalid credentials")
			httpStatusCode = http.StatusBadRequest
			return
		}
		log.Info("apikey", apikey)
		scopes = strings.Join(apikey.Scopes, " ")
		log.Info("scopes ", scopes)
		username = apikey.Username
	} else {
		organization = clientID
		scopes = "organization:owner"
	}

	at = newAccessToken(username, organization, clientID, scopes)
	return
}

func convertCodeToAccessTokenHandler(code string, clientID string, secret string, redirectURI string, mgr *Manager) (at *AccessToken, httpStatusCode int) {
	httpStatusCode = http.StatusOK

	ar, err := mgr.getAuthorizationRequest(code)
	if err != nil {
		log.Error("ERROR getting the original authorization request:", err)
		httpStatusCode = http.StatusInternalServerError
		return
	}
	if ar == nil {
		log.Debug("No original authorization request found with this authorization code")
		httpStatusCode = http.StatusBadRequest
		return
	}

	if ar.ClientID != clientID || ar.RedirectURL != redirectURI {
		log.Debugf("Client id:%s - Expected client id:%s", clientID, ar.ClientID)
		log.Debugf("Redirect url:%s - Expected redirect url:%s", redirectURI, ar.RedirectURL)
		log.Info("Bad client or hacking attempt, client_id or redirect_uri is different from the original authorization request")
		httpStatusCode = http.StatusBadRequest
		return
	}

	if ar.IsExpiredAt(time.Now()) {
		log.Info("Token request for an expired authorizationrequest")
		httpStatusCode = http.StatusBadRequest
		return
	}

	client, err := mgr.getClientByCredentials(clientID, secret)
	if err != nil {
		log.Error("Error getting the oauth client: ", err)
		httpStatusCode = http.StatusInternalServerError
		return
	}
	if client == nil {
		log.Info("(client_id - secret) combination not found")
		httpStatusCode = http.StatusBadRequest
		return
	}

	if !strings.HasPrefix(redirectURI, client.CallbackURL) {
		log.Debug("return_uri does not match the callback uri")
		httpStatusCode = http.StatusBadRequest
		return
	}

	at = newAccessToken(ar.Username, "", ar.ClientID, ar.Scope)
	return
}

func (service *Service) createItsYouOnlineAdminToken(username string, r *http.Request) (token string, err error) {
	at := newAccessToken(username, "", "itsyouonline", "admin")

	mgr := NewManager(r)
	err = mgr.saveAccessToken(at)
	if err == nil {
		token = at.AccessToken
	}
	return
}

// verifyScopes checks for a user:memberof:clientid scope. If present, check if the user
// really is a member or owner in this organization. Should this not be the case,
// remove the scope and add the correct memberof:childorg scopes
func verifyScopes(scopeString string, username string, clientID string, orgMgr *organization.Manager) (string, error) {
	scopes := strings.Split(scopeString, ",")
	for i, scope := range scopes {
		// If the scope is memberof:clientid, where clientid is the globalid of the
		// organization requesting authorization, verify if we really have this scope
		// or if we have the scope on one or more children. In case of the latter, replace
		// the scopes by the appropriate ones
		if scope == "user:memberof:"+clientID {
			isMember, err := orgMgr.IsMember(clientID, username)
			if err != nil {
				return scopeString, err
			}
			if isMember {
				return scopeString, nil
			}
			isOwner, err := orgMgr.IsOwner(clientID, username)
			if err != nil {
				return scopeString, err
			}
			if isOwner {
				return scopeString, nil
			}
			log.Debugf("Access token contains scope %v, but user is not a member of %v - parse suborganizations",
				scope, clientID)
			// Remove the memberof:parentOrg scope
			scopes = append(scopes[:i], scopes[i+1:]...)
			var subOrgScopes []string
			subOrgScopes, err = findSuborgScopes(clientID, username, orgMgr)
			if err != nil {
				return scope, err
			}
			// Add the correct memberof scopes
			scopes = append(scopes, subOrgScopes...)
			scopeString = strings.Join(scopes, ",")
		}
	}
	return scopeString, nil
}

func findSuborgScopes(parentID string, username string, orgMgr *organization.Manager) ([]string, error) {
	subOrgs, err := orgMgr.GetSubOrganizations(parentID)
	if err != nil {
		return nil, err
	}
	orgIds := make([]string, len(subOrgs))
	for i, subOrg := range subOrgs {
		orgIds[i] = subOrg.Globalid
	}
	sort.Strings(orgIds)

	var foundScopes []string

	var foundID string
	for _, orgID := range orgIds {
		if foundID != "" && strings.HasPrefix(orgID, foundID) {
			// This is a suborg of an organization we already know about so skip it
			continue
		}
		isMember, err := orgMgr.IsMember(orgID, username)
		if err != nil {
			return nil, err
		}
		if isMember {
			foundScopes = append(foundScopes, "user:memberof:"+orgID)
			foundID = orgID
		}
		isOwner, err := orgMgr.IsOwner(orgID, username)
		if err != nil {
			return nil, err
		}
		if isOwner {
			foundScopes = append(foundScopes, "user:memberof:"+orgID)
			foundID = orgID
		}
	}
	return foundScopes, nil
}
