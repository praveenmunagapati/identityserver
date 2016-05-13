package siteservice

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"

	log "github.com/Sirupsen/logrus"
)

const (
	mongoLoginCollectionName = "loginsessions"
)

//initLoginModels initialize models in mongo
func (service *Service) initLoginModels() {
	index := mgo.Index{
		Key:      []string{"sessionkey"},
		Unique:   true,
		DropDups: false,
	}

	db.EnsureIndex(mongoLoginCollectionName, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 10,
		Background:  true,
	}
	db.EnsureIndex(mongoLoginCollectionName, automaticExpiration)

}

type loginSessionInformation struct {
	sessionKey string
	smsCode    string
	createdAt  time.Time
}

func newLoginSessionInformation() (sessionInformation *loginSessionInformation, err error) {
	sessionInformation = &loginSessionInformation{createdAt: time.Now()}
	sessionInformation.sessionKey, err = generateRandomString()
	if err != nil {
		return
	}
	numbercode, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return
	}
	sessionInformation.smsCode = fmt.Sprintf("%06d", numbercode)
	return
}

const loginFileName = "login.html"
const totpFileName = "logintotpform.html"
const smsFormFileName = "loginsmsform.html"

//renderForm shows the user login page
func (service *Service) renderLoginForm(w http.ResponseWriter, request *http.Request, pageFileName string, indicateError bool, postbackURL string) {
	htmlData, err := html.Asset(pageFileName)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if indicateError {
		htmlData = bytes.Replace(htmlData, []byte(`{"invalidsomething": true}`), []byte(`{"invalidcredentials": true}`), 1)
	}
	htmlData = bytes.Replace(htmlData, []byte(`action="login"`), []byte(fmt.Sprintf("action=\"%s\"", postbackURL)), 1)
	sessions.Save(request, w)
	w.Write(htmlData)
}

//ShowLoginForm shows the user login page on the initial request
func (service *Service) ShowLoginForm(w http.ResponseWriter, request *http.Request) {
	service.renderLoginForm(w, request, loginFileName, false, request.RequestURI)

}

//ProcessLoginForm logs a user in if the credentials are valid
func (service *Service) ProcessLoginForm(w http.ResponseWriter, request *http.Request) {
	//TODO: validate csrf token
	//TODO: limit the number of failed/concurrent requests

	err := request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing registration form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	values := request.Form

	username := values.Get("login")

	//validate the username exists
	var userexists bool
	userMgr := user.NewManager(request)

	if userexists, err = userMgr.Exists(username); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var validpassword bool
	passwdMgr := password.NewManager(request)
	if validpassword, err = passwdMgr.Validate(username, values.Get("password")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	validcredentials := userexists && validpassword
	if !validcredentials {
		service.renderLoginForm(w, request, loginFileName, true, request.RequestURI)
		return
	}
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		return
	}
	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loginSession.Values["username"] = username
	if u.TwoFAMethod == "sms" {
		sessionInfo, err := newLoginSessionInformation()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		loginSession.Values["sessionkey"] = sessionInfo.sessionKey
		mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)
		mgoCollection.Insert(sessionInfo)
		smsmessage := fmt.Sprintf("https://%s/smsvalidation?c=%s or enter the code %s in the form", request.Host, sessionInfo.smsCode, sessionInfo.smsCode)
		//TODO: check which phonenumber to use
		phonenumber := u.Phone["main"]
		go service.smsService.Send(string(phonenumber), smsmessage)
		redirectToDifferentLoginPage(w, request, "login", "loginsmsconfirmation")
	} else {
		redirectToDifferentLoginPage(w, request, "login", "logintotpconfirmation")
	}
}

func generateRandomString() (randomString string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	randomString = base64.StdEncoding.EncodeToString(b)
	return
}

//getUserLoggingIn returns an user trying to log in, or an empty string if there is none
func (service *Service) getUserLoggingIn(request *http.Request) (username string, err error) {
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		return
	}
	savedusername := loginSession.Values["username"]
	if savedusername != nil {
		username, _ = savedusername.(string)
	}
	return
}

//getSessionKey returns an the login session key, or an empty string if there is none
func (service *Service) getSessionKey(request *http.Request) (sessionKey string, err error) {
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		return
	}
	savedSessionKey := loginSession.Values["sessionkey"]
	if savedSessionKey != nil {
		sessionKey, _ = savedSessionKey.(string)
	}
	return
}

func redirectToDifferentLoginPage(w http.ResponseWriter, request *http.Request, from string, to string) {
	log.Debugf("Redirecting from %s to %s", from, to)
	sessions.Save(request, w)
	u, _ := url.Parse(request.RequestURI)
	path := strings.TrimSuffix(u.Path, from)
	path += to
	u.Path = path
	http.Redirect(w, request, u.RequestURI(), http.StatusFound)
}

//ShowTOTPConfirmationForm shows the user login page on the initial request
func (service *Service) ShowTOTPConfirmationForm(w http.ResponseWriter, request *http.Request) {
	service.renderLoginForm(w, request, totpFileName, false, request.RequestURI)
}

//ProcessTOTPConfirmation checks the totp 2 factor authentication code
func (service *Service) ProcessTOTPConfirmation(w http.ResponseWriter, request *http.Request) {
	username, err := service.getUserLoggingIn(request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if username == "" {
		redirectToDifferentLoginPage(w, request, "logintotpconfirmation", "login")
		return
	}

	err = request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing totp form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	values := request.Form
	var validtotpcode bool
	totpMgr := totp.NewManager(request)
	if validtotpcode, err = totpMgr.Validate(username, values.Get("totpcode")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !validtotpcode { //TODO: limit to 3 failed attempts
		service.renderLoginForm(w, request, totpFileName, true, request.RequestURI)
		return
	}
	service.loginUser(w, request, username)
}

//ShowSMSConfirmationForm shows the user login page on the initial request
func (service *Service) ShowSMSConfirmationForm(w http.ResponseWriter, request *http.Request) {
	service.renderLoginForm(w, request, smsFormFileName, false, request.RequestURI)
}

//ProcessSMSConfirmation checks the totp 2 factor authentication code
func (service *Service) ProcessSMSConfirmation(w http.ResponseWriter, request *http.Request) {
	username, err := service.getUserLoggingIn(request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if username == "" {
		redirectToDifferentLoginPage(w, request, "loginsmsconfirmation", "login")
		return
	}

	err = request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing sms confirmation form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	values := request.Form
	smscode := values.Get("smscode")

	sessionKey, err := service.getSessionKey(request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if sessionKey == "" {
		redirectToDifferentLoginPage(w, request, "loginsmsconfirmation", "login")
		return
	}
	var validsmscode bool

	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)
	sessionInfo := &loginSessionInformation{}
	err = mgoCollection.Find(bson.M{"sessionkey": sessionKey}).One(sessionInfo)
	if err == mgo.ErrNotFound {
		redirectToDifferentLoginPage(w, request, "loginsmsconfirmation", "login")
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	validsmscode = (smscode != sessionInfo.sessionKey)

	if !validsmscode { //TODO: limit to 3 failed attempts
		service.renderLoginForm(w, request, smsFormFileName, true, request.RequestURI)
		return
	}
	service.loginUser(w, request, username)
}

func (service *Service) loginUser(w http.ResponseWriter, request *http.Request, username string) {
	//TODO: Clear login session
	if err := service.SetLoggedInUser(w, request, username); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sessions.Save(request, w)

	log.Debugf("Successfull login by '%s'", username)

	redirectURL := ""
	queryValues := request.URL.Query()
	endpoint := queryValues.Get("endpoint")
	if endpoint != "" {
		queryValues.Del("endpoint")
		redirectURL = endpoint + "?" + queryValues.Encode()
	} else {
		parameters := make(url.Values)
		parameters.Add("client_id", "itsyouonline")
		parameters.Add("response_type", "token")
		redirectURL = "v1/oauth/authorize?" + parameters.Encode()
	}

	http.Redirect(w, request, redirectURL, http.StatusFound)
}
