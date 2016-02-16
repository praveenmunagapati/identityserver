package db

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
)

type DBHandler struct {
	handler http.Handler
	session *mgo.Session
}

func (d *DBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.session = SetDBSession(r)
	if d.session == nil {
		panic(errors.New("Failed to retrieve a DB session!"))
	}

	defer d.closeSession()

	d.handler.ServeHTTP(w, r)
}

func (d *DBHandler) closeSession() {
	d.session.Close()
}

func DBMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &DBHandler{
			handler: h,
		}
	}
}
