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

// Update an existing user. Updating ``username`` is not allowed. The labelled lists
// can not be updated this way, the normal properties can (like github and facebook account).
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

	respBody := &userview{
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

// It is handler for GET /users/{username}/phonenumbers/{label}
func (api UsersAPI) usernamephonenumberslabelGet(w http.ResponseWriter, r *http.Request) {

}

// Update or create an existing phonenumber.
// It is handler for PUT /users/{username}/phonenumbers/{label}
func (api UsersAPI) usernamephonenumberslabelPut(w http.ResponseWriter, r *http.Request) {

}

// Delete a phonenumber
// It is handler for DELETE /users/{username}/phonenumbers/{label}
func (api UsersAPI) usernamephonenumberslabelDelete(w http.ResponseWriter, r *http.Request) {

}

// It is handler for GET /users/{username}/banks/{label}
func (api UsersAPI) usernamebankslabelGet(w http.ResponseWriter, r *http.Request) {

}

// Update or create an existing bankaccount.
// It is handler for PUT /users/{username}/banks/{label}
func (api UsersAPI) usernamebankslabelPut(w http.ResponseWriter, r *http.Request) {

}

// Delete a BankAccount
// It is handler for DELETE /users/{username}/banks/{label}
func (api UsersAPI) usernamebankslabelDelete(w http.ResponseWriter, r *http.Request) {

}

// It is handler for GET /users/{username}
func (api UsersAPI) usernameGet(w http.ResponseWriter, r *http.Request) {

}

// It is handler for GET /users/{username}/addresses
func (api UsersAPI) usernameaddressesGet(w http.ResponseWriter, r *http.Request) {

}

// It is handler for GET /users/{username}/addresses/{label}
func (api UsersAPI) usernameaddresseslabelGet(w http.ResponseWriter, r *http.Request) {

}

// Update or create an existing address.
// It is handler for PUT /users/{username}/addresses/{label}
func (api UsersAPI) usernameaddresseslabelPut(w http.ResponseWriter, r *http.Request) {

}

// Delete an address
// It is handler for DELETE /users/{username}/addresses/{label}
func (api UsersAPI) usernameaddresseslabelDelete(w http.ResponseWriter, r *http.Request) {

}

// It is handler for GET /users/{username}/phonenumbers
func (api UsersAPI) usernamephonenumbersGet(w http.ResponseWriter, r *http.Request) {

}

// It is handler for GET /users/{username}/banks
func (api UsersAPI) usernamebanksGet(w http.ResponseWriter, r *http.Request) {

}
