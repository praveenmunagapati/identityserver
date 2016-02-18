package user

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type UsersAPI struct {
}

// It is handler for POST /users
func (api UsersAPI) Post(w http.ResponseWriter, r *http.Request) {

	var u User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := NewUserManager(r)
	if err := userMgr.Save(&u); err != nil {
		log.Error("ERROR while saving user:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&u)
}

// Update existing user. Updating ``username`` i s not allowed.
// It is handler for PUT /users/{username}
func (api UsersAPI) usernamePut(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	var u User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := NewUserManager(r)

	oldUser, uerr := userMgr.GetByName(username)
	if uerr != nil {
		log.Debug(uerr)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if u.Username != username || u.GetId() != oldUser.GetId() {
		http.Error(w, "Changing username or id is Forbidden!", http.StatusForbidden)
		return
	}

	if err := userMgr.Save(&u); err != nil {
		log.Error("ERROR while saving user:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&u)
}

// It is handler for GET /users/{username}/info
func (api UsersAPI) usernameinfoGet(w http.ResponseWriter, r *http.Request) {

	var respBody Userview
	json.NewEncoder(w).Encode(&respBody)

	// uncomment below line to add header
	// w.Header.Set("key","value")
}

// It is handler for GET /users/{username}/validate
func (api UsersAPI) usernamevalidateGet(w http.ResponseWriter, r *http.Request) {

	// token := req.FormValue("token")

	// uncomment below line to add header
	// w.Header.Set("key","value")
}
