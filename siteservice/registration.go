package siteservice

import (
	"bytes"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/totp"
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
	token, err := totp.NewToken()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totpsession, err := service.GetSession(request, SessionForRegistration, "totp")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totpsession.Values["secret"] = token.Secret

	//Don't use go templates since angular uses "{{ ... }}" syntax as well and this way it the standalone page also works
	htmlData = bytes.Replace(htmlData, []byte("secret=1234123412341234"), []byte("secret="+token.Secret), 2)

	sessions.Save(request, w)
	w.Write(htmlData)
}

//ProcessRegistrationForm processes the user registration form
func (service *Service) ProcessRegistrationForm(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing registration form:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	values := request.Form

	totpcode := values.Get("totpcode")

	totpsession, err := service.GetSession(request, SessionForRegistration, "totp")
	if err != nil {
		log.Error("EROR while getting the totp registration session", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if totpsession.IsNew {
		//TODO: indicate expired registration session
		log.Debug("New registration session while processing the registration form")
		service.ShowRegistrationForm(w, request)
		return
	}
	secret, ok := totpsession.Values["secret"].(string)
	if !ok {
		log.Error("Unable to convert the stored session totp secret to a string")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	token := totp.TokenFromSecret(secret)
	if !token.Validate(totpcode) {
		//TODO: indidate invalid totpcode
		service.ShowRegistrationForm(w, request)
	}
	log.Debug(values)

	http.Redirect(w, request, "", http.StatusFound)
}
