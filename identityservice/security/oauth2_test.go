package security

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGetScopestringFromJWT(t *testing.T) {
	originaltoken := jwt.New(jwt.SigningMethodHS256)
	originaltoken.Claims["username"] = "rob"
	originaltoken.Claims["scope"] = []string{"1", "2"}
	encodedjwt, _ := originaltoken.SignedString([]byte("abcde"))

	token, err := jwt.Parse(encodedjwt, func(token *jwt.Token) (interface{}, error) {
		return []byte("abcde"), nil
	})

	scopestring := GetScopestringFromJWT(token)

	assert.NoError(t, err, "")
	assert.Equal(t, "1,2", scopestring, "")
}
