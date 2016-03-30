package userorganization

import "github.com/itsyouonline/identityserver/identityservice/user"

// Oauth2oauth_2_0Middleware is oauth2 middleware for oauth_2_0
type Oauth2oauth_2_0Middleware struct {
	user.Oauth2oauth_2_0Middleware
}

// newOauth2oauth_2_0Middlewarecreate new Oauth2oauth_2_0Middleware struct
func newOauth2oauth_2_0Middleware(scopes []string) *Oauth2oauth_2_0Middleware {
	om := &Oauth2oauth_2_0Middleware{}
	om.Scopes = scopes

	om.DescribedBy = "headers"
	om.Field = "Authorization"

	return om
}
