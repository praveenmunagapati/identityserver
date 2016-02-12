package siteservice

import (
	"net/http"

	"github.com/itsyouonline/website/packaged/html"
)

const registrationFileName = "registration.html"

//ShowRegistrationForm shows the user registration page
func (service *Service) ShowRegistrationForm(w http.ResponseWriter, request *http.Request) {
	htmlData, err := html.Asset(registrationFileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(htmlData)
}

//ProcessRegistrationForm processes the user registration form
func (service *Service) ProcessRegistrationForm(response http.ResponseWriter, request *http.Request) {

}
