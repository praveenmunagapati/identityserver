package siteservice

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"

	"github.com/gorilla/context"
)

//SessionType is used to define the type of session
type SessionType int

const (
	//SessionForRegistration is the short anynymous session used during registration
	SessionForRegistration SessionType = iota
	//SessionInteractive is the session of an authenticated user on the itsyou.online website
	SessionInteractive SessionType = iota
	//SessionLogin is the session during the login flow
	SessionLogin SessionType = iota
)

//initializeSessionStore creates a cookieStore
// mageAge is the maximum age in seconds
func initializeSessionStore(cookieSecret string, maxAge int) (sessionStore *sessions.CookieStore) {
	sessionStore = sessions.NewCookieStore([]byte(cookieSecret))
	sessionStore.Options.HttpOnly = true

	sessionStore.Options.Secure = true
	sessionStore.Options.MaxAge = maxAge
	return
}

func (service *Service) initializeSessions(cookieSecret string) {
	service.Sessions = make(map[SessionType]*sessions.CookieStore)

	service.Sessions[SessionForRegistration] = initializeSessionStore(cookieSecret, 10*60)
	service.Sessions[SessionInteractive] = initializeSessionStore(cookieSecret, 10*60)
	service.Sessions[SessionLogin] = initializeSessionStore(cookieSecret, 5*60)

}

//GetSession returns the a session of the specified kind and a spefic name
func (service *Service) GetSession(request *http.Request, kind SessionType, name string) (*sessions.Session, error) {
	return service.Sessions[kind].Get(request, name)
}

//SetLoggedInUser creates a session for an authenticated user
func (service *Service) SetLoggedInUser(w http.ResponseWriter, request *http.Request, username string) (err error) {
	authenticatedSession, err := service.GetSession(request, SessionInteractive, "authenticatedsession")
	if err != nil {
		log.Error(err)
		return
	}
	authenticatedSession.Values["username"] = username

	//TODO: rework this, is not really secure I think
	// Set user cookie after successful login
	cookie := &http.Cookie{
		Name:  "itsyou.online.user",
		Path:  "/",
		Value: username,
	}
	http.SetCookie(w, cookie)

	return
}

//SetAPIAccessToken sets the api access token in a cookie
//TODO: is not safe to do. Now there are also two ways of passing tokens to the client
func (service *Service) SetAPIAccessToken(w http.ResponseWriter, token string) (err error) {
	// Set token cookie
	cookie := &http.Cookie{
		Name:  "itsyou.online.apitoken",
		Path:  "/",
		Value: token,
	}
	http.SetCookie(w, cookie)

	return
}

//GetLoggedInUser returns an authenticated user, or an empty string if there is none
func (service *Service) GetLoggedInUser(request *http.Request, w http.ResponseWriter) (username string, err error) {
	authenticatedSession, err := service.GetSession(request, SessionInteractive, "authenticatedsession")
	if err != nil {
		log.Error(err)
		return
	}
	err = authenticatedSession.Save(request, w)
	if err != nil {
		log.Error(err)
		return
	}
	savedusername := authenticatedSession.Values["username"]
	if savedusername != nil {
		username, _ = savedusername.(string)
	}
	return
}

//SetAuthenticatedUserMiddleWare puthe the authenticated user on the context
func (service *Service) SetAuthenticatedUserMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if username, err := service.GetLoggedInUser(request, w); err == nil {
			context.Set(request, "authenticateduser", username)
		}

		next.ServeHTTP(w, request)
	})
}
