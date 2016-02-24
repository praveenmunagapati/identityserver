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

	userMgr := NewManager(r)
	if err := userMgr.Save(&u); err != nil {
		log.Error("ERROR while saving user:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&u)
	w.Header().Set("Content-type", "application/json")
}

// Update existing user. Updating ``username`` is not allowed.
// It is handler for PUT /users/{username}
func (api UsersAPI) usernamePut(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	var u User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := NewManager(r)

	oldUser, uerr := userMgr.GetByName(username)
	if uerr != nil {
		log.Debug(uerr)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if u.Username != username || u.getID() != oldUser.getID() {
		http.Error(w, "Changing username or id is Forbidden!", http.StatusForbidden)
		return
	}

	if err := userMgr.Save(&u); err != nil {
		log.Error("ERROR while saving user:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&u)
	w.Header().Set("Content-type", "application/json")
}

// It is handler for GET /users/{username}/info
func (api UsersAPI) usernameinfoGet(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := NewManager(r)

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// TODO: apply authorization limits.
	addresses := []Address{}
	emails := []string{}
	phones := []Phonenumber{}

	for _, address := range user.Address {
		addresses = append(addresses, address)
	}

	for _, email := range user.Email {
		emails = append(emails, email)
	}

	for _, phone := range user.Phone {
		phones = append(phones, phone)
	}

	respBody := &Userview{
		Address:  addresses,
		Email:    emails,
		Phone:    phones,
		Username: user.Username,
	}

	json.NewEncoder(w).Encode(respBody)
	w.Header().Set("Content-type", "application/json")
}

// It is handler for GET /users/{username}/validate
func (api UsersAPI) usernamevalidateGet(w http.ResponseWriter, r *http.Request) {

	// token := req.FormValue("token")

	// uncomment below line to add header
	// w.Header.Set("key","value")
}
