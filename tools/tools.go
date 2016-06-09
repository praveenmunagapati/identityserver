package tools

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString() (randomString string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	randomString = base64.StdEncoding.EncodeToString(b)
	return
}
