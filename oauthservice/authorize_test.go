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
	ar := newAuthorizationRequest("user1", "client1", "state2")
	assert.NotEmpty(t, ar.AuthorizationCode)
	assert.False(t, strings.HasSuffix(ar.AuthorizationCode, "="))
	assert.NotEqual(t, time.Time{}, ar.CreatedAt)
	assert.Equal(t, "user1", ar.Username)
	assert.Equal(t, "state2", ar.State)
	assert.Equal(t, "client1", ar.ClientID)
}
