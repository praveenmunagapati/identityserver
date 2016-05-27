package validation

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//SMSService is the interface an sms communicaction channel should have to be used by the IYOPhonenumberValidationService
type SMSService interface {
	Send(phonenumber string, message string) (err error)
}

const (
	mongoOngoingPhonenumberValidationCollectionName = "ongoingphonenumbervalidations"
	mongoValidatedPhonenumbers                      = "validatedphonenumbers"
)

//ValidatedPhonenumber is a record of a phonenumber for a user and when it is validated
type ValidatedPhonenumber struct {
	Username    string
	Phonenumber string
	CreatedAt   time.Time
}

type phonenumberValidationInformation struct {
	Key         string
	SMSCode     string
	Username    string
	Phonenumber string
	CreatedAt   time.Time
}

//IYOPhonenumberValidationService is the itsyou.online implementation of a PhonenumberValidationService
type IYOPhonenumberValidationService struct {
	SMSService SMSService
}

//InitPhonenumberValidationModels initialize models in mongo
func (service *IYOPhonenumberValidationService) InitPhonenumberValidationModels() {
	index := mgo.Index{
		Key:      []string{"key"},
		Unique:   true,
		DropDups: false,
	}

	db.EnsureIndex(mongoOngoingPhonenumberValidationCollectionName, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 10,
		Background:  true,
	}
	db.EnsureIndex(mongoOngoingPhonenumberValidationCollectionName, automaticExpiration)

	index = mgo.Index{
		Key:      []string{"username", "phonenumber"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoValidatedPhonenumbers, index)

}

func newPhonenumberValidationInformation() (info *phonenumberValidationInformation, err error) {
	info = &phonenumberValidationInformation{CreatedAt: time.Now()}
	info.Key, err = generateRandomString()
	if err != nil {
		return
	}
	numbercode, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return
	}
	info.SMSCode = fmt.Sprintf("%06d", numbercode)
	return
}

//RequestValidation validates the phonenumber by sending an SMS
func (service *IYOPhonenumberValidationService) RequestValidation(request *http.Request, username string, phone user.Phonenumber, host string) (key string, err error) {
	info, err := newPhonenumberValidationInformation()
	if err != nil {
		return
	}
	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoOngoingPhonenumberValidationCollectionName)
	err = mgoCollection.Insert(info)
	if err != nil {
		return
	}
	smsmessage := fmt.Sprintf("%s?c=%s&k=%s or enter the code %s in the form", host, info.SMSCode, url.QueryEscape(info.Key), info.SMSCode)

	go service.SMSService.Send(string(phone), smsmessage)
	key = info.Key
	return
}

var (
	//ErrInvalid denotes that the supplied combination for validation is invalid.
	ErrInvalid = errors.New("Invalid")
)

func (service *IYOPhonenumberValidationService) getPhonenumberValidationInformation(request *http.Request, key string) (info *phonenumberValidationInformation, err error) {
	if key == "" {
		return
	}

	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoOngoingPhonenumberValidationCollectionName)
	info = &phonenumberValidationInformation{}
	err = mgoCollection.Find(bson.M{"key": key}).One(info)
	if err == mgo.ErrNotFound {
		info = nil
		err = nil
	}
	return
}

//ConfirmValidation checks if the supplied code matches the username and key
func (service *IYOPhonenumberValidationService) ConfirmValidation(request *http.Request, key, code string) (err error) {
	info, err := service.getPhonenumberValidationInformation(request, key)
	if err != nil {
		return
	}
	if info == nil || info.SMSCode != code {
		err = ErrInvalid
		return
	}
	p := &ValidatedPhonenumber{Username: info.Username, Phonenumber: info.Phonenumber, CreatedAt: time.Now()}
	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoValidatedPhonenumbers)
	err = mgoCollection.Insert(p)

	return
}

func generateRandomString() (randomString string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	randomString = base64.StdEncoding.EncodeToString(b)
	return
}
