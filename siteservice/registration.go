package siteservice

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"
)

const registrationFileName = "registration.html"

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
	totpsecret, ok := totpsession.Values["secret"].(string)
	if !ok {
		log.Error("Unable to convert the stored session totp secret to a string")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newuser := &user.User{
		Username: values.Get("login"),
		Email:    map[string]string{"registration": values.Get("email")},
	}
	//TODO: validate newuser

	//validate the username is not taken yet
	userMgr := user.NewManager(request)
	//TODO: distributed lock
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

	totpcode := values.Get("totpcode")

	token := totp.TokenFromSecret(totpsecret)
	if !token.Validate(totpcode) {
		log.Debug("Invalid totp code")
		validationErrors = append(validationErrors, "invalidtotpcode")
		service.renderRegistrationFrom(w, request, validationErrors, totpsecret)
		return
	}

	userMgr.Save(newuser)
	passwdMgr := password.NewManager(request)
	passwdMgr.Save(newuser.Username, values.Get("password"))
	totpMgr := totp.NewManager(request)
	totpMgr.Save(newuser.Username, totpsecret)

	log.Debugf("Registered %s", newuser.Username)
	service.SetLoggedInUser(request, newuser.Username)
	http.Redirect(w, request, "", http.StatusFound)
}
