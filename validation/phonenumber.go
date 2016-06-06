package validation

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/validation"
)

//SMSService is the interface an sms communicaction channel should have to be used by the IYOPhonenumberValidationService
type SMSService interface {
	Send(phonenumber string, message string) (err error)
}

//IYOPhonenumberValidationService is the itsyou.online implementation of a PhonenumberValidationService
type IYOPhonenumberValidationService struct {
	SMSService SMSService
}


//RequestValidation validates the phonenumber by sending an SMS
func (service *IYOPhonenumberValidationService) RequestValidation(request *http.Request, username string, phonenumber user.Phonenumber, confirmationurl string) (key string, err error) {
	valMngr := validation.NewManager(request)
	info, err := valMngr.NewPhonenumberValidationInformation(username, phonenumber)
	if err != nil {
		return
	}
	err = valMngr.SavePhonenumberValidationInformation(info)
	if err != nil {
		return
	}
	smsmessage := fmt.Sprintf("To verify your phonenumber on itsyou.online enter the code %s in the form or use this link: %s?c=%s&k=%s", info.SMSCode, confirmationurl, info.SMSCode, url.QueryEscape(info.Key))

	go service.SMSService.Send(string(phonenumber), smsmessage)
	key = info.Key
	return
}

//ExpireValidation removes a pending validation
func (service *IYOPhonenumberValidationService) ExpireValidation(request *http.Request, key string) (err error) {
	if key == "" {
		return
	}
	valMngr := validation.NewManager(request)
	err = valMngr.RemovePhonenumberValidationInformation(key)
	return
}

var (
	//ErrInvalidCode denotes that the supplied code is invalid
	ErrInvalidCode = errors.New("Invalid code")
	//ErrInvalidOrExpiredKey denotes that the key is not found, it can be invalid or expired
	ErrInvalidOrExpiredKey = errors.New("Invalid key")
)

func (service *IYOPhonenumberValidationService) getPhonenumberValidationInformation(request *http.Request, key string) (info *validation.PhonenumberValidationInformation, err error) {
	if key == "" {
		return
	}
	valMngr := validation.NewManager(request)
	info, err = valMngr.GetByKeyPhonenumberValidationInformation(key)
	return
}

//IsConfirmed checks wether a validation request is already confirmed
func (service *IYOPhonenumberValidationService) IsConfirmed(request *http.Request, key string) (confirmed bool, err error) {
	info, err := service.getPhonenumberValidationInformation(request, key)
	if err != nil {
		return
	}
	if info == nil {
		err = ErrInvalidOrExpiredKey
		return
	}
	confirmed = info.Confirmed
	return
}

//ConfirmValidation checks if the supplied code matches the username and key
func (service *IYOPhonenumberValidationService) ConfirmValidation(request *http.Request, key, code string) (err error) {
	info, err := service.getPhonenumberValidationInformation(request, key)
	if err != nil {
		return
	}
	if info == nil {
		err = ErrInvalidOrExpiredKey
		return
	}
	if info.SMSCode != code {
		err = ErrInvalidCode
		return
	}
	valMngr := validation.NewManager(request)
	p := valMngr.NewValidatedPhonenumber(info.Username, info.Phonenumber)
	err = valMngr.SaveValidatedPhonenumber(p)
	if err != nil {
		return
	}
	err = valMngr.UpdatePhonenumberValidationInformation(key, true)
	if err != nil {
		return
	}
	return
}

