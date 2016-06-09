package tools

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
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

func RenderTemplate(templatepath string, data interface{}) (message string, err error) {
	log.Print(data)
	templateEngine := template.New("template")
	templateEngine, err = template.ParseFiles(templatepath)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	templateEngine.Execute(buf, data)
	message = buf.String()
	return
}
