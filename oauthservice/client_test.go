package oauthservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOauth2Client(t *testing.T) {
	c := NewOauth2Client("client1", "main")
	assert.Equal(t, "client1", c.ClientID)
	assert.Equal(t, "main", c.Label)
	assert.NotEmpty(t, c.Secret)

	c2 := NewOauth2Client("clientid", "")
	assert.NotEqual(t, c.Secret, c2.Secret)
}
