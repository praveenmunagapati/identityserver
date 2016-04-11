package siteservice

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"

	log "github.com/Sirupsen/logrus"
)

const authorizeFileName = "authorize.html"

//renderAuthorizeForm renders the html page for the authorize form
func (service *Service) renderAuthorizeForm(w http.ResponseWriter, request *http.Request, postbackURL string) {
	htmlData, err := html.Asset(authorizeFileName)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	htmlData = bytes.Replace(htmlData, []byte(`action="authorize"`), []byte(fmt.Sprintf("action=\"%s\"", postbackURL)), 1)
	w.Write(htmlData)
}

//ShowAuthorizeForm shows the scopes an application requested and asks a user for confirmation
func (service *Service) ShowAuthorizeForm(w http.ResponseWriter, r *http.Request) {
	service.renderAuthorizeForm(w, r, r.RequestURI)
}

//ProcessAuthorizeForm saves or cancels a requested authorization
//TODO: is this staying like this or de talk to the api from the client?
func (service *Service) ProcessAuthorizeForm(w http.ResponseWriter, request *http.Request) {
	//TODO: validate csrf token
	//TODO: limit the number of failed/concurrent requests
}
