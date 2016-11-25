package tools

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/templates/packaged"
	"html/template"
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

func RenderTemplate(templateName string, data interface{}) (message string, err error) {
	log.Debug("Email template: ", templateName)
	log.Debug("Email Data: ", data)
	htmlData, err := templates.Asset(templateName)
	if err != nil {
		log.Error("Could not get email asset: ", err)
		return
	}
	templateEngine := template.New("template")
	templateEngine, err = templateEngine.Parse(string(htmlData))
	if err != nil {
		log.Error("Could parse template: ", err)
		return
	}
	buf := new(bytes.Buffer)
	templateEngine.Execute(buf, data)
	message = buf.String()
	return
}
