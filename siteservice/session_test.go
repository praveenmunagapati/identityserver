package siteservice

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAvailableSessions(t *testing.T) {

	router := mux.NewRouter()

	siteService := &Service{}
	siteService.AddRoutes(router)
	request := &http.Request{}

	session, err := siteService.GetSession(request, SessionForRegistration, "akey")
	assert.NoError(t, err)
	assert.NotNil(t, session)

}
