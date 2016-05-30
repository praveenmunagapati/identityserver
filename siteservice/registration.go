package siteservice

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"
	"github.com/itsyouonline/identityserver/validation"
)

const (
	mongoRegistrationCollectionName = "registrationsessions"
)

//initLoginModels initialize models in mongo
func (service *Service) initRegistrationModels() {
	index := mgo.Index{
		Key:      []string{"sessionkey"},
		Unique:   true,
		DropDups: false,
	}

	db.EnsureIndex(mongoRegistrationCollectionName, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 10, //10 minutes
		Background:  true,
	}
	db.EnsureIndex(mongoRegistrationCollectionName, automaticExpiration)

}

type registrationSessionInformation struct {
	SessionKey           string
	SMSCode              string
	Confirmed            bool
	ConfirmationAttempts uint
	CreatedAt            time.Time
}

func newRegistrationSessionInformation() (sessionInformation *registrationSessionInformation, err error) {
	sessionInformation = &registrationSessionInformation{CreatedAt: time.Now()}
	sessionInformation.SessionKey, err = generateRandomString()
	if err != nil {
		return
	}
	numbercode, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return
	}
	sessionInformation.SMSCode = fmt.Sprintf("%06d", numbercode)
	return
}

const (
	registrationFileName                        = "registration.html"
	registrationPhonenumberconfirmationFileName = "registrationsmsform.html"
	registrationResendSMSFileName               = "registrationresendsms.html"
)

func (service *Service) renderRegistrationFrom(w http.ResponseWriter, request *http.Request, validationErrors []string, totpsecret string) {
	htmlData, err := html.Asset(registrationFileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Don't use go templates since angular uses "{{ ... }}" syntax as well and this way the standalone page also works
	htmlData = bytes.Replace(htmlData, []byte("secret=1234123412341234"), []byte("secret="+totpsecret), 2)

	errorMap := make(map[string]bool)
	for _, errorkey := range validationErrors {
		errorMap[errorkey] = true
	}
	jsonErrors, err := json.Marshal(errorMap)

	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	htmlData = bytes.Replace(htmlData, []byte(`{"invalidsomething": true}`), jsonErrors, 1)

	sessions.Save(request, w)
	w.Write(htmlData)
}

//CheckRegistrationSMSConfirmation is called by the sms code form to check if the sms is already confirmed on the mobile phone
func (service *Service) CheckRegistrationSMSConfirmation(w http.ResponseWriter, request *http.Request) {
	registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := map[string]bool{}

	if registrationSession.IsNew {
		response["confirmed"] = true //This way the form will be submitted, let the form handler deal with redirect to login
	} else {
		validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)

		confirmed, err := service.phonenumberValidationService.IsConfirmed(request, validationkey)
		if err == validation.ErrInvalidOrExpiredKey {
			confirmed = true //This way the form will be submitted, let the form handler deal with redirect to login
			return
		}
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response["confirmed"] = confirmed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//ShowRegistrationForm shows the user registration page
func (service *Service) ShowRegistrationForm(w http.ResponseWriter, request *http.Request) {
	validationErrors := make([]string, 0, 0)

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

	service.renderRegistrationFrom(w, request, validationErrors, token.Secret)
}

//ShowPhonenumberConfirmationForm shows the user a form to enter the code sent by sms
func (service *Service) ShowPhonenumberConfirmationForm(w http.ResponseWriter, request *http.Request) {
	service.renderForm(w, request, registrationPhonenumberconfirmationFileName, []string{})
}

func (service *Service) renderForm(w http.ResponseWriter, request *http.Request, filename string, validationErrors []string) {
	htmlData, err := html.Asset(filename)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Don't use go templates since angular uses "{{ ... }}" syntax as well and this way the standalone page also works

	errorMap := make(map[string]bool)
	for _, errorkey := range validationErrors {
		errorMap[errorkey] = true
	}
	jsonErrors, err := json.Marshal(errorMap)

	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	htmlData = bytes.Replace(htmlData, []byte(`{"invalidsomething": true}`), jsonErrors, 1)

	sessions.Save(request, w)
	w.Write(htmlData)
}

//ProcessPhonenumberConfirmationForm processes the Phone number confirmation form
func (service *Service) ProcessPhonenumberConfirmationForm(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		redirectToDifferentPage(w, request, true, "registersmsconfirmation", "login")
		return
	}

	username, _ := registrationSession.Values["username"].(string)
	validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)

	if isConfirmed, _ := service.phonenumberValidationService.IsConfirmed(request, validationkey); isConfirmed {
		service.loginUser(w, request, username)
		return
	}

	smscode := request.Form.Get("smscode")
	if err != nil || smscode == "" {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = service.phonenumberValidationService.ConfirmValidation(request, validationkey, smscode)
	if err == validation.ErrInvalidCode {
		service.renderForm(w, request, registrationPhonenumberconfirmationFileName, []string{"invalidcredentials"})
		return
	}
	if err == validation.ErrInvalidOrExpiredKey {
		redirectToDifferentPage(w, request, true, "registersmsconfirmation", "login")
		return
	}
	service.loginUser(w, request, username)
}

//ShowResendPhonenumberConfirmation renders the Resend phonenumberconfirmation form
func (service *Service) ShowResendPhonenumberConfirmation(w http.ResponseWriter, request *http.Request) {
	service.renderForm(w, request, registrationResendSMSFileName, []string{})
}

//ResendPhonenumberConfirmation resend the phonenumberconfirmation to a possbily new phonenumber
func (service *Service) ResendPhonenumberConfirmation(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		redirectToDifferentPage(w, request, true, "registersmsconfirmation", "login")
		return
	}

	username, _ := registrationSession.Values["username"].(string)

	//Invalidate the previous validation request, ignore a possible error
	validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)
	_ = service.phonenumberValidationService.ExpireValidation(request, validationkey)

	phonenumber := user.Phonenumber(request.FormValue("phonenumber"))
	if !phonenumber.IsValid() {
		log.Debug("Invalid phone number")
		service.renderForm(w, request, registrationPhonenumberconfirmationFileName, []string{"invalidphonenumber"})
		return
	}

	uMgr := user.NewManager(request)
	err = uMgr.SavePhone(username, "main", phonenumber)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	validationkey, err = service.phonenumberValidationService.RequestValidation(request, username, phonenumber, fmt.Sprintf("https://%s/phonevalidation", request.Host))
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	registrationSession.Values["phonenumbervalidationkey"] = validationkey

	redirectToDifferentPage(w, request, true, "registerresendsms", "registersmsconfirmation")
}

//ProcessRegistrationForm processes the user registration form
func (service *Service) ProcessRegistrationForm(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing registration form:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	validationErrors := make([]string, 0, 0)

	values := request.Form

	twoFAMethod := values.Get(".twoFAMethod")
	if twoFAMethod != "sms" && twoFAMethod != "totp" {
		log.Info("Invalid 2fa method during registration: ", twoFAMethod)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	totpsession, err := service.GetSession(request, SessionForRegistration, "totp")
	if err != nil {
		log.Error("ERROR while getting the totp registration session", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if totpsession.IsNew {
		//TODO: indicate expired registration session
		log.Debug("New registration session while processing the registration form")
		service.ShowRegistrationForm(w, request)
		return
	}
	totpsecret, ok := totpsession.Values["secret"].(string)
	if !ok {
		log.Error("Unable to convert the stored session totp secret to a string")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newuser := &user.User{
		Username:    values.Get("login"),
		Email:       map[string]string{"main": values.Get("email")},
		TwoFAMethod: twoFAMethod,
	}
	//TODO: validate newuser

	//validate the username is not taken yet
	userMgr := user.NewManager(request)
	//we now just depend on mongo unique index to avoid duplicates when concurrent requests are made
	userExists, err := userMgr.Exists(newuser.Username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		validationErrors = append(validationErrors, "duplicateusername")
		log.Debug("USER ", newuser.Username, " already registered")
		service.renderRegistrationFrom(w, request, validationErrors, totpsecret)
		return
	}

	if twoFAMethod == "sms" {
		phonenumber := user.Phonenumber(values.Get("phonenumber"))
		if !phonenumber.IsValid() {
			log.Debug("Invalid phone number")
			validationErrors = append(validationErrors, "invalidphonenumber")
			service.renderRegistrationFrom(w, request, validationErrors, totpsecret)
			return
		}
		newuser.Phone = map[string]user.Phonenumber{"main": phonenumber}
	} else {
		totpcode := values.Get("totpcode")

		token := totp.TokenFromSecret(totpsecret)
		if !token.Validate(totpcode) {
			log.Debug("Invalid totp code")
			validationErrors = append(validationErrors, "invalidtotpcode")
			service.renderRegistrationFrom(w, request, validationErrors, totpsecret)
			return
		}
	}

	//TODO: this should only be a temporary user registration (until the email/phone validation is completed)
	userMgr.Save(newuser)
	passwdMgr := password.NewManager(request)
	passwdMgr.Save(newuser.Username, values.Get("password"))

	if twoFAMethod == "sms" {
		validationkey, err := service.phonenumberValidationService.RequestValidation(request, newuser.Username, newuser.Phone["main"], fmt.Sprintf("https://%s/phonevalidation", request.Host))
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
		registrationSession.Values["username"] = newuser.Username
		registrationSession.Values["phonenumbervalidationkey"] = validationkey

		redirectToDifferentPage(w, request, true, "register", "registersmsconfirmation")
		return
	}

	totpMgr := totp.NewManager(request)
	totpMgr.Save(newuser.Username, totpsecret)

	log.Debugf("Registered %s", newuser.Username)
	service.loginUser(w, request, newuser.Username)
}
