package siteservice

import (
	"net/http"

	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"

	log "github.com/Sirupsen/logrus"
)

const authorizeFileName = "authorize.html"

//ShowAuthorizeForm shows the scopes an application requested and asks a user for confirmation
func (service *Service) ShowAuthorizeForm(w http.ResponseWriter, r *http.Request) {
	htmlData, err := html.Asset(authorizeFileName)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(htmlData)
}
