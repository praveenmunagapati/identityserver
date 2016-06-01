package user

//This file is auto-generated by go-raml
//Do not edit this file by hand since it will be overwritten during the next generation

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

// UsersInterface is interface for /users root endpoint
type UsersInterface interface { // Post is the handler for POST /users
	// Create a new user
	Post(http.ResponseWriter, *http.Request)
	// usernamevalidateGet is the handler for GET /users/{username}/validate
	usernamevalidateGet(http.ResponseWriter, *http.Request)
	// UpdatePassword is the handler for PUT /users/{username}/password
	UpdatePassword(http.ResponseWriter, *http.Request)
	// usernamephonenumbersGet is the handler for GET /users/{username}/phonenumbers
	usernamephonenumbersGet(http.ResponseWriter, *http.Request)
	// RegisterNewPhonenumber is the handler for POST /users/{username}/phonenumbers
	// Register a new phonenumber
	RegisterNewPhonenumber(http.ResponseWriter, *http.Request)
	// usernamephonenumberslabelGet is the handler for GET /users/{username}/phonenumbers/{label}
	usernamephonenumberslabelGet(http.ResponseWriter, *http.Request)
	// UpdatePhonenumber is the handler for PUT /users/{username}/phonenumbers/{label}
	// Update the label and/or value of an existing phonenumber.
	UpdatePhonenumber(http.ResponseWriter, *http.Request)
	// DeletePhonenumber is the handler for DELETE /users/{username}/phonenumbers/{label}
	// Removes a phonenumber
	DeletePhonenumber(http.ResponseWriter, *http.Request)
	// usernamebanksGet is the handler for GET /users/{username}/banks
	usernamebanksGet(http.ResponseWriter, *http.Request)
	// usernamebanksPost is the handler for POST /users/{username}/banks
	// Create new bank account
	usernamebanksPost(http.ResponseWriter, *http.Request)
	// usernamenotificationsGet is the handler for GET /users/{username}/notifications
	// Get the list of notifications, these are pending invitations or approvals
	usernamenotificationsGet(http.ResponseWriter, *http.Request)
	// usernameGet is the handler for GET /users/{username}
	usernameGet(http.ResponseWriter, *http.Request)
	// DeleteFacebookAccount is the handler for DELETE /users/{username}/facebook
	// Delete the associated facebook account
	DeleteFacebookAccount(http.ResponseWriter, *http.Request)
	// RegisterNewEmailAddress is the handler for POST /users/{username}/emailaddresses
	// Register a new email address
	RegisterNewEmailAddress(http.ResponseWriter, *http.Request)
	// UpdateEmailAddress is the handler for PUT /users/{username}/emailaddresses/{label}
	// Updates the label and/or value of an email address
	UpdateEmailAddress(http.ResponseWriter, *http.Request)
	// DeleteEmailAddress is the handler for DELETE /users/{username}/emailaddresses/{label}
	// Removes an email address
	DeleteEmailAddress(http.ResponseWriter, *http.Request)
	// DeleteGithubAccount is the handler for DELETE /users/{username}/github
	// Unlink Github Account
	DeleteGithubAccount(http.ResponseWriter, *http.Request)
	// GetUserInformation is the handler for GET /users/{username}/info
	GetUserInformation(http.ResponseWriter, *http.Request)
	// usernameaddressesGet is the handler for GET /users/{username}/addresses
	usernameaddressesGet(http.ResponseWriter, *http.Request)
	// RegisterNewAddress is the handler for POST /users/{username}/addresses
	// Register a new address
	RegisterNewAddress(http.ResponseWriter, *http.Request)
	// usernameaddresseslabelGet is the handler for GET /users/{username}/addresses/{label}
	usernameaddresseslabelGet(http.ResponseWriter, *http.Request)
	// UpdateAddress is the handler for PUT /users/{username}/addresses/{label}
	// Update the label and/or value of an existing address.
	UpdateAddress(http.ResponseWriter, *http.Request)
	// DeleteAddress is the handler for DELETE /users/{username}/addresses/{label}
	// Removes an address
	DeleteAddress(http.ResponseWriter, *http.Request)
	// usernamebankslabelGet is the handler for GET /users/{username}/banks/{label}
	usernamebankslabelGet(http.ResponseWriter, *http.Request)
	// usernamebankslabelPut is the handler for PUT /users/{username}/banks/{label}
	// Update an existing bankaccount and label.
	usernamebankslabelPut(http.ResponseWriter, *http.Request)
	// usernamebankslabelDelete is the handler for DELETE /users/{username}/banks/{label}
	// Delete a BankAccount
	usernamebankslabelDelete(http.ResponseWriter, *http.Request)
	// usernamecontractsGet is the handler for GET /users/{username}/contracts
	// Get the contracts where the user is 1 of the parties. Order descending by date.
	usernamecontractsGet(http.ResponseWriter, *http.Request)
	// GetAllAuthorizations is the handler for GET /users/{username}/authorizations
	// Get the list of authorizations.
	GetAllAuthorizations(http.ResponseWriter, *http.Request)
	// GetAuthorization is the handler for GET /users/{username}/authorizations/{grantedTo}
	// Get the authorization for a specific organization.
	GetAuthorization(http.ResponseWriter, *http.Request)
	// UpdateAuthorization is the handler for PUT /users/{username}/authorizations/{grantedTo}
	// Modify which information an organization is able to see.
	UpdateAuthorization(http.ResponseWriter, *http.Request)
	// DeleteAuthorization is the handler for DELETE /users/{username}/authorizations/{grantedTo}
	// Remove the authorization for an organization, the granted organization will no longer
	// have access the user's information.
	DeleteAuthorization(http.ResponseWriter, *http.Request)
	// Add API Key
	AddAPIKey(http.ResponseWriter, *http.Request)
	GetAPIKey(http.ResponseWriter, *http.Request)
	UpdateAPIKey(http.ResponseWriter, *http.Request)
	DeleteAPIKey(http.ResponseWriter, *http.Request)
	ListAPIKeys(http.ResponseWriter, *http.Request)

}

// UsersInterfaceRoutes is routing for /users root endpoint
func UsersInterfaceRoutes(r *mux.Router, i UsersInterface) {
	r.Handle("/users", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.Post))).Methods("POST")
	r.Handle("/users/{username}/validate", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.usernamevalidateGet))).Methods("GET")
	r.Handle("/users/{username}/password", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.UpdatePassword))).Methods("PUT")
	r.Handle("/users/{username}/phonenumbers", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamephonenumbersGet))).Methods("GET")
	r.Handle("/users/{username}/phonenumbers", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewPhonenumber))).Methods("POST")
	r.Handle("/users/{username}/phonenumbers/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamephonenumberslabelGet))).Methods("GET")
	r.Handle("/users/{username}/phonenumbers/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdatePhonenumber))).Methods("PUT")
	r.Handle("/users/{username}/phonenumbers/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeletePhonenumber))).Methods("DELETE")
	r.Handle("/users/{username}/banks", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamebanksGet))).Methods("GET")
	r.Handle("/users/{username}/banks", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamebanksPost))).Methods("POST")
	r.Handle("/users/{username}/notifications", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamenotificationsGet))).Methods("GET")
	r.Handle("/users/{username}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernameGet))).Methods("GET")
	r.Handle("/users/{username}/apikeys", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.ListAPIKeys))).Methods("GET")
	r.Handle("/users/{username}/apikeys", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.AddAPIKey))).Methods("POST")
	r.Handle("/users/{username}/apikeys/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetAPIKey))).Methods("GET")
	r.Handle("/users/{username}/apikeys/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAPIKey))).Methods("PUT")
	r.Handle("/users/{username}/apikeys/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteAPIKey))).Methods("DELETE")
	r.Handle("/users/{username}/facebook", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteFacebookAccount))).Methods("DELETE")
	r.Handle("/users/{username}/emailaddresses", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewEmailAddress))).Methods("POST")
	r.Handle("/users/{username}/emailaddresses/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateEmailAddress))).Methods("PUT")
	r.Handle("/users/{username}/emailaddresses/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteEmailAddress))).Methods("DELETE")
	r.Handle("/users/{username}/github", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.DeleteGithubAccount))).Methods("DELETE")
	r.Handle("/users/{username}/info", alice.New(newOauth2oauth_2_0Middleware([]string{"user:info", "user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserInformation))).Methods("GET")
	r.Handle("/users/{username}/addresses", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernameaddressesGet))).Methods("GET")
	r.Handle("/users/{username}/addresses", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewAddress))).Methods("POST")
	r.Handle("/users/{username}/addresses/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernameaddresseslabelGet))).Methods("GET")
	r.Handle("/users/{username}/addresses/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAddress))).Methods("PUT")
	r.Handle("/users/{username}/addresses/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteAddress))).Methods("DELETE")
	r.Handle("/users/{username}/banks/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamebankslabelGet))).Methods("GET")
	r.Handle("/users/{username}/banks/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamebankslabelPut))).Methods("PUT")
	r.Handle("/users/{username}/banks/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamebankslabelDelete))).Methods("DELETE")
	r.Handle("/users/{username}/contracts", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.usernamecontractsGet))).Methods("GET")
	r.Handle("/users/{username}/authorizations", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetAllAuthorizations))).Methods("GET")
	r.Handle("/users/{username}/authorizations/{grantedTo}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetAuthorization))).Methods("GET")
	r.Handle("/users/{username}/authorizations/{grantedTo}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAuthorization))).Methods("PUT")
	r.Handle("/users/{username}/authorizations/{grantedTo}", alice.New(newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteAuthorization))).Methods("DELETE")
}
