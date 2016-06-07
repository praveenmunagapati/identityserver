package siteservice

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/validation"
)

//PhonenumberValidation is the page that is linked to in the SMS for phonenumbervalidation and is thus accessed on the mobile phone
func (service *Service) PhonenumberValidation(w http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := request.Form
	key := values.Get("k")
	smscode := values.Get("c")

	err = service.phonenumberValidationService.ConfirmValidation(request, key, smscode)
	if err == validation.ErrInvalidCode || err == validation.ErrInvalidOrExpiredKey {
		service.renderSMSConfirmationPage(w, request, "Invalid or expired link")
		return
	}
	if err != nil {
		log.Error(err)
		service.renderSMSConfirmationPage(w, request, "An unexpected error occurred, please try again later")
		return
	}

	service.renderSMSConfirmationPage(w, request, "Your phonenumber is confirmed")
}


func (service *Service) EmailValidation(w http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := request.Form
	key := values.Get("k")
	smscode := values.Get("c")

	err = service.emailaddressValidationService.ConfirmValidation(request, key, smscode)
	if err == validation.ErrInvalidCode || err == validation.ErrInvalidOrExpiredKey {
		service.renderEmailConfirmationPage(w, request, "Invalid or expired link")
		return
	}
	if err != nil {
		log.Error(err)
		service.renderEmailConfirmationPage(w, request, "An unexpected error occurred, please try again later")
		return
	}

	service.renderEmailConfirmationPage(w, request, "Your Email Address is confirmed")
}
