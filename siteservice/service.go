package siteservice

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/website/packaged/assets"
	"github.com/itsyouonline/website/packaged/html"
	"github.com/itsyouonline/website/packaged/thirdpartyassets"
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
	router.PathPrefix("/thirdpartyassets/").Handler(http.StripPrefix("/thirdpartyassets/", http.FileServer(
		&assetfs.AssetFS{Asset: thirdpartyassets.Asset, AssetDir: thirdpartyassets.AssetDir, AssetInfo: thirdpartyassets.AssetInfo})))

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
