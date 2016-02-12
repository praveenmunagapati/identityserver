package siteservice

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/website/packaged/assets"
	"github.com/itsyouonline/website/packaged/html"
)

//Service is the identityserver http service
type Service struct {
}

//AddRoutes registers the http routes with the router
func (service *Service) AddRoutes(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(service.HomePage)
	router.Methods("GET").Path("/register").HandlerFunc(service.ShowRegistrationForm)
	router.Methods("POST").Path("/register").HandlerFunc(service.ProcessRegistrationForm)
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(
		&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assets.AssetInfo})))

}

const homepageFileName = "index.html"

//HomePage shows the homepage
func (service *Service) HomePage(w http.ResponseWriter, request *http.Request) {
	htmlData, err := html.Asset(homepageFileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(htmlData)
}
