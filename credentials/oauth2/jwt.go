package oauth2

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//GetValidJWT returns a validated ES384 signed jwt from the authorization header that needs to start with "bearer "
// If no jwt is found in the authorization header, nil is returned
// Validation against the supplied publickey is performed
func GetValidJWT(r *http.Request, publicKey ecdsa.PublicKey) (token *jwt.Token, err error) {
	authorizationHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorizationHeader, "bearer ") {
		return
	}
	jwtstring := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "bearer"))

	token, err = jwt.Parse(jwtstring, func(token *jwt.Token) (interface{}, error) {

		m, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if token.Header["alg"] != m.Alg() {
			return nil, fmt.Errorf("Unexpected signing algorithm: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err == nil && !token.Valid {
		err = errors.New("Invalid jwt supplied:" + jwtstring)
	}
	return
}
