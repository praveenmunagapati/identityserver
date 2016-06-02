package siteservice

import (
	"net/http"
	"strings"
	"net/url"
)

//ShowAuthorizeForm shows the scopes an application requested and asks a user for confirmation
func (service *Service) ShowAuthorizeForm(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.RequestURI
	parameters := make(url.Values)
	//Don't parse the redirect url, can only give errors while we don't gain much
	if !strings.Contains(redirectURI, "#") {
		redirectURI += "#"
	}
	redirectURI += parameters.Encode()
	http.Redirect(w, r, redirectURI, http.StatusFound)
}
