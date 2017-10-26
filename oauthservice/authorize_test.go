package oauthservice

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationRequestExpiration(t *testing.T) {
	ar := &authorizationRequest{CreatedAt: time.Now()}

	assert.True(t, ar.IsExpiredAt(ar.CreatedAt.Add(time.Second*11)))
	assert.False(t, ar.IsExpiredAt(ar.CreatedAt.Add(time.Second*10)))
}

func TestNewAuthorizationRequest(t *testing.T) {
	ar := newAuthorizationRequest("user1", "client1", "scope3", "https://localhost")
	assert.NotEmpty(t, ar.AuthorizationCode)
	assert.False(t, strings.HasSuffix(ar.AuthorizationCode, "="))
	assert.NotEqual(t, time.Time{}, ar.CreatedAt)
	assert.Equal(t, "user1", ar.Username)
	assert.Equal(t, "client1", ar.ClientID)
	assert.Equal(t, "scope3", ar.Scope)
	assert.Equal(t, "https://localhost", ar.RedirectURL)
}

func TestIsValidAuthorization(t *testing.T) {
	authorizedScopes := []string{"user:name", "user:email:main", "user:phone",
		"user:address", "user:memberof:testorg"}
	assert.True(t, IsAuthorizationValid([]string{"user:name", "user:phone"}, authorizedScopes))
	assert.False(t, IsAuthorizationValid([]string{"user:phone:main"}, authorizedScopes))
	assert.False(t, IsAuthorizationValid([]string{"user:address:work"}, authorizedScopes))
	assert.True(t, IsAuthorizationValid([]string{"user:name", "user:email:main", "user:memberof:testorg"}, authorizedScopes))
}

type testClientManager struct {
	clients []*Oauth2Client
}

//AllByClientID just returns the all the clients, is only usefull for testing off course
func (m *testClientManager) AllByClientID(clientID string) (clients []*Oauth2Client, err error) {
	clients = m.clients
	return
}

func TestValidateRedirectURI(t *testing.T) {
	type testcase struct {
		redirectURI string
		valid       bool
	}
	mgr := &testClientManager{
		clients: []*Oauth2Client{&Oauth2Client{CallbackURL: "http://www.url.com/callback"}},
	}
	testcases := []testcase{
		testcase{redirectURI: "", valid: false},
		testcase{redirectURI: "test.com", valid: false},
		testcase{redirectURI: "https://itsyou.online", valid: false},
		testcase{redirectURI: "https://test.itsyou.online", valid: false},
		testcase{redirectURI: "https://test.itsyou.online:443", valid: false},
		testcase{redirectURI: "http://www.url.com/callback/subpath", valid: true},
	}
	for i, test := range testcases {
		valid, err := validateRedirectURI(mgr, test.redirectURI, "clientID")
		assert.NoError(t, err, i)
		assert.Equal(t, test.valid, valid, i)
	}
}
