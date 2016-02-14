package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	resource "github.com/itsyouonline/identityserver/http"
	userModel "github.com/itsyouonline/identityserver/identityservice/user/models"
)

type UserDetails struct {
	userModel.User
	Uri          string `json:"uri"`
	CompaniesUri string `json:"companiesUri"`
}

type UserResource struct {
	resource.Resource
}

func NewUserResource() *UserResource {
	u := &UserResource{}
	u.ResourceHandler = u

	return u
}

func (u *UserResource) GetRoutes() resource.Routes {
	routes := resource.Routes{
		resource.Route{
			Name: "UserList",
			Methods: resource.RouteMethods{
				resource.POST,
			},
			Path:        "/users/",
			HandlerFunc: u.DispatchList,
		},
		resource.Route{
			Name: "UserDetail",
			Methods: resource.RouteMethods{
				resource.GET,
			},
			Path:        "/users/{username}/",
			HandlerFunc: u.DispatchDetail,
		},
	}

	return routes
}

func (u *UserResource) PostList(w http.ResponseWriter, r *http.Request) {
	user, err := u.deserialize(r)
	if err != nil {
		log.Debug(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = user.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := u.serialize(user)

	u.Respond(w, response)
}

func (u *UserResource) GetDetail(w http.ResponseWriter, r *http.Request) {
	userMgr := userModel.NewUserManager(r)

	username := mux.Vars(r)["username"]

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	response := u.serialize(user)

	u.Respond(w, response)
}

func (u *UserResource) serialize(user *userModel.User) *UserDetails {
	uri := u.resourceUri(user)
	companiesUri := u.BuildUri(uri, "companies")

	return &UserDetails{
		*user,
		uri,
		companiesUri,
	}
}

func (u *UserResource) deserialize(r *http.Request) (*userModel.User, error) {
	user := userModel.NewUser(r)

	data, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserResource) resourceUri(user *userModel.User) string {
	return u.BuildUri("/users/", user.Username)
}
