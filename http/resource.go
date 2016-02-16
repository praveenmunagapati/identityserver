package http

import (
	"encoding/json"
	"net/http"
	"path"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type ResourceHandler interface {
	GetList(w http.ResponseWriter, req *http.Request)
	PostList(w http.ResponseWriter, req *http.Request)
	PutList(w http.ResponseWriter, req *http.Request)
	DeleteList(w http.ResponseWriter, req *http.Request)

	GetDetail(w http.ResponseWriter, req *http.Request)
	PostDetail(w http.ResponseWriter, req *http.Request)
	PutDetail(w http.ResponseWriter, req *http.Request)
	DeleteDetail(w http.ResponseWriter, req *http.Request)

	GetRoutes() Routes

	Respond(w http.ResponseWriter, response interface{})
}

type ResourceDispatcher interface {
	DispatchList(w http.ResponseWriter, req *http.Request)
	DispatchDetail(w http.ResponseWriter, req *http.Request)
}

type Resource struct {
	ResourceHandler
}

func (r *Resource) DispatchList(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case GET:
		r.ResourceHandler.GetList(w, req)
	case POST:
		r.ResourceHandler.PostList(w, req)
	case PUT:
		r.ResourceHandler.PutList(w, req)
	case DELETE:
		r.ResourceHandler.DeleteList(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Resource) DispatchDetail(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case GET:
		r.ResourceHandler.GetDetail(w, req)
	case POST:
		r.ResourceHandler.PostDetail(w, req)
	case PUT:
		r.ResourceHandler.PutDetail(w, req)
	case DELETE:
		r.ResourceHandler.DeleteDetail(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Resource) GetRoutes() Routes {
	return Routes{}
}

func (r *Resource) GetList(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) PostList(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) PutList(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) DeleteList(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) GetDetail(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) PostDetail(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) PutDetail(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) DeleteDetail(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented", http.StatusInternalServerError)
	return
}

func (r *Resource) Respond(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (r *Resource) BuildUri(components ...string) string {
	return path.Join(components...) + "/"
}
