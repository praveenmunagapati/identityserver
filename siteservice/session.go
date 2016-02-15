package siteservice

import (
	"net/http"

	"github.com/gorilla/sessions"
)

//SessionType is used to define the type of session
type SessionType int

const (
	//SessionForRegistration is the short anynymous session used during registration
	SessionForRegistration SessionType = iota
)

func (service *Service) initializeSessions() {
	service.Sessions = make(map[SessionType]*sessions.CookieStore)

	//TODO: https://github.com/itsyouonline/identityserver/issues/6
	cookieStoreSecret := "TODO: ISSUE #6"

	registrationSessionStore := sessions.NewCookieStore([]byte(cookieStoreSecret))
	registrationSessionStore.Options.HttpOnly = true
	registrationSessionStore.Options.Secure = true
	registrationSessionStore.Options.MaxAge = 10 * 60 //10 minutes

	service.Sessions[SessionForRegistration] = registrationSessionStore
}

//GetSession returns the a session of the specified kind and a spefic name
func (service *Service) GetSession(request *http.Request, kind SessionType, name string) (*sessions.Session, error) {
	return service.Sessions[kind].Get(request, name)
}
