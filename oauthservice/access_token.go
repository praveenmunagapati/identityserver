package oauthservice

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

//AccessTokenExpiration is the time in seconds an access token expires
var AccessTokenExpiration = time.Second * 3600

type accessToken struct {
	AccessToken string
	Type        string
	Username    string
	Scope       string
	ClientID    string
	CreatedAt   time.Time
}

func (at *accessToken) IsExpiredAt(testtime time.Time) bool {
	return testtime.After(at.CreatedAt.Add(AccessTokenExpiration))
}

func newAccessToken(username, clientID, scope string) *accessToken {
	var at accessToken

	randombytes := make([]byte, 21) //Multiple of 3 to make sure no padding is added
	rand.Read(randombytes)
	at.AccessToken = base64.URLEncoding.EncodeToString(randombytes)
	at.CreatedAt = time.Now()
	at.Username = username
	at.ClientID = clientID
	at.Scope = scope
	at.Type = "bearer"

	return &at
}

//AccessTokenHandler is the handler of the /login/oauth/access_token endpoint
func (service *Service) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing form: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	clientSecret := r.FormValue("client_secret")
	clientID := r.FormValue("client_id")

	if code == "" || clientSecret == "" || clientID == "" {
		log.Debug("Required parameter missing in the request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")

	mgr := NewManager(r)
	ar, err := mgr.Get(code)
	if err != nil {
		log.Error("ERROR getting the original authorization request:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if ar == nil {
		log.Debug("No original authorization request found with this authorization code")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if ar.ClientID != clientID || ar.State != state {
		log.Info("Bad client or hacking attempt, state or client_id is different from the original authorization request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if ar.IsExpiredAt(time.Now()) {
		log.Info("Token request for an expired authorizationrequest")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//TODO: check redirecturl
	//TODO: check clientID/clientSecret
	at := newAccessToken(ar.Username, ar.ClientID, ar.Scope)
	mgr.SaveAccessToken(at)

	response := struct {
		AccessToken string      `json:"access_token"`
		TokenType   string      `json:"token_type"`
		Scope       string      `json:"scope"`
		Info        interface{} `json:"info"`
	}{
		AccessToken: at.AccessToken,
		TokenType:   at.Type,
		Scope:       at.Scope,
		Info: struct {
			Username string `json:"username"`
		}{
			Username: at.Username,
		},
	}

	json.NewEncoder(w).Encode(&response)
	w.Header().Set("Content-type", "application/json")
}
