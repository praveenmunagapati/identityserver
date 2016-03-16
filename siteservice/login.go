package siteservice

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"

	log "github.com/Sirupsen/logrus"
)

const loginFileName = "login.html"

//renderLoginForm shows the user login page
func (service *Service) renderLoginForm(w http.ResponseWriter, request *http.Request, indicateError bool, postbackURL string) {
	htmlData, err := html.Asset(loginFileName)
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
	service.renderLoginForm(w, request, false, request.RequestURI)

}

//ProcessLoginForm logs a user in if the credentials are valid
func (service *Service) ProcessLoginForm(w http.ResponseWriter, request *http.Request) {
	//TODO: validate csrf token
	//TODO: limit the number of failed/concurrent requests

	log.Debug(request.RequestURI)

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

	var validtotpcode bool
	totpMgr := totp.NewManager(request)
	if validtotpcode, err = totpMgr.Validate(username, values.Get("totpcode")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	validcredentials := userexists && validpassword && validtotpcode
	if !validcredentials {
		service.renderLoginForm(w, request, true, request.RequestURI)
		return
	}

	if err := service.SetLoggedInUser(request, username); err != nil {
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
	}

	// Set user cookie after successful login
	cookie := &http.Cookie{
		Name:  "itsyou.online.user",
		Path:  "/",
		Value: username,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, request, redirectURL, http.StatusFound)

}
