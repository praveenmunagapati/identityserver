package siteservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/user"
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

const (
	registrationFileName = "registration.html"
)

func (service *Service) renderRegistrationFrom(w http.ResponseWriter, request *http.Request) {
	htmlData, err := html.Asset(registrationFileName)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

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
		// todo: registrationSession is new with SMS, something must be wrong
		log.Warn("Registration is new")
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
	service.renderRegistrationFrom(w, request)
}

//ProcessPhonenumberConfirmationForm processes the Phone number confirmation form
func (service *Service) ProcessPhonenumberConfirmationForm(w http.ResponseWriter, request *http.Request) {
	values := struct {
		Smscode string `json:"smscode"`
	}{}

	response := struct {
		RedirectUrL string `json:"redirecturl"`
		Error       string `json:"error"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ProcessPhonenumberConfirmation request:", err)
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
		sessions.Save(request, w)
		response.RedirectUrL = fmt.Sprintf("https://%s/register/#/smsconfirmation", request.Host)
		json.NewEncoder(w).Encode(&response)
		return
	}

	username, _ := registrationSession.Values["username"].(string)
	validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)

	if isConfirmed, _ := service.phonenumberValidationService.IsConfirmed(request, validationkey); isConfirmed {
		service.loginUser(w, request, username)
		return
	}

	smscode := values.Smscode
	if err != nil || smscode == "" {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = service.phonenumberValidationService.ConfirmValidation(request, validationkey, smscode)
	if err == validation.ErrInvalidCode {
		w.WriteHeader(422)
		response.Error = "invalidsmscode"
		json.NewEncoder(w).Encode(&response)
		return
	}
	if err == validation.ErrInvalidOrExpiredKey {
		sessions.Save(request, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&response)
		return
	}
	service.loginUser(w, request, username)
}

//ResendPhonenumberConfirmation resend the phonenumberconfirmation to a possbily new phonenumber
func (service *Service) ResendPhonenumberConfirmation(w http.ResponseWriter, request *http.Request) {
	values := struct {
		PhoneNumber string `json:"phonenumber"`
	}{}

	response := struct {
		RedirectUrL string `json:"redirecturl"`
		Error       string `json:"error"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ResendPhonenumberConfirmation request:", err)
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
		sessions.Save(request, w)
		log.Debug("Registration session expired")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username, _ := registrationSession.Values["username"].(string)

	//Invalidate the previous validation request, ignore a possible error
	validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)
	_ = service.phonenumberValidationService.ExpireValidation(request, validationkey)

	phonenumber := user.Phonenumber(values.PhoneNumber)
	if !phonenumber.IsValid() {
		log.Debug("Invalid phone number")
		w.WriteHeader(422)
		response.Error = "invalidphonenumber"
		json.NewEncoder(w).Encode(&response)
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

	sessions.Save(request, w)
	response.RedirectUrL = fmt.Sprintf("https://%s/register/#smsconfirmation", request.Host)
	json.NewEncoder(w).Encode(&response)
}

//ProcessRegistrationForm processes the user registration form
func (service *Service) ProcessRegistrationForm(w http.ResponseWriter, request *http.Request) {
	response := struct {
		Redirecturl string `json:"redirecturl"`
		Error       string `json:"error"`
	}{}
	values := struct {
		TwoFAMethod string `json:"twofamethod"`
		Login       string `json:"login"`
		Email       string `json:"email"`
		Phonenumber string `json:"phonenumber"`
		TotpCode    string `json:"totpcode"`
		Password    string `json:"password"`
	}{}
	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the registration request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	twoFAMethod := values.TwoFAMethod
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
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	totpsecret, ok := totpsession.Values["secret"].(string)
	if !ok {
		log.Error("Unable to convert the stored session totp secret to a string")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newuser := &user.User{
		Username:    values.Login,
		Email:       map[string]string{"main": values.Email},
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
		log.Debug("USER ", newuser.Username, " already registered")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if twoFAMethod == "sms" {
		phonenumber := user.Phonenumber(values.Phonenumber)
		if !phonenumber.IsValid() {
			log.Debug("Invalid phone number")
			w.WriteHeader(422)
			response.Error = "invalidphonenumber";
			json.NewEncoder(w).Encode(&response)
			return
		}
		newuser.Phone = map[string]user.Phonenumber{"main": phonenumber}
	} else {
		token := totp.TokenFromSecret(totpsecret)
		if !token.Validate(values.TotpCode) {
			log.Debug("Invalid totp code")
			w.WriteHeader(422)
			response.Error = "invalidtotpcode";
			json.NewEncoder(w).Encode(&response)
			return
		}
	}

	//TODO: this should only be a temporary user registration (until the email/phone validation is completed)
	userMgr.Save(newuser)
	passwdMgr := password.NewManager(request)
	err = passwdMgr.Save(newuser.Username, values.Password)
	if err != nil {
		log.Error(err)
		if (err.Error() != "internal_error") {
			w.WriteHeader(422)
			response.Error = "invalidpassword";
			json.NewEncoder(w).Encode(&response)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

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

		sessions.Save(request, w)
		response.Redirecturl = fmt.Sprintf("https://%s/register#/smsconfirmation?%s", request.Host, request.URL.Query().Encode());
		json.NewEncoder(w).Encode(&response)
		return
	}

	totpMgr := totp.NewManager(request)
	totpMgr.Save(newuser.Username, totpsecret)

	log.Debugf("Registered %s", newuser.Username)
	service.loginUser(w, request, newuser.Username)
}

//ValidateUsername checks if a username is already taken or not
func (service *Service) ValidateUsername(w http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("username")
	response := struct {
		Valid bool `json:"valid"`
		Error string `json:"error"`
	}{
		true,
		"",
	}
	userMgr := user.NewManager(request)
	userExists, err := userMgr.Exists(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		log.Debug("username ", username, " already taken")
		response.Error = "duplicateusername"
		response.Valid = false
	}
	// TODO: validate username
	json.NewEncoder(w).Encode(&response)
	return
}

