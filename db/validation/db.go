package validation

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/itsyouonline/identityserver/db"
	"time"
	"math/big"
	"encoding/base64"
	"github.com/itsyouonline/identityserver/db/user"
	"crypto/rand"
	"gopkg.in/mgo.v2/bson"
)


const (
	mongoOngoingPhonenumberValidationCollectionName = "ongoingphonenumbervalidations"
	mongoValidatedPhonenumbers                      = "validatedphonenumbers"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
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

//Manager is used to store users
type Manager struct {
	session *mgo.Session
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session: session,
	}
}

func (manager *Manager) NewPhonenumberValidationInformation(username string, phonenumber user.Phonenumber) (info *PhonenumberValidationInformation, err error) {
	info = &PhonenumberValidationInformation{CreatedAt: time.Now(), Username: username, Phonenumber: string(phonenumber)}
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

func (manager *Manager) SavePhonenumberValidationInformation(info *PhonenumberValidationInformation) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	err = mgoCollection.Insert(info)
	return
}

func (manager *Manager) RemovePhonenumberValidationInformation(key string) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	_, err = mgoCollection.RemoveAll(bson.M{"key": key})
	return
}

func (manager *Manager) UpdatePhonenumberValidationInformation(key string, confirmed bool) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	err = mgoCollection.Update(bson.M{"key": key}, bson.M{"$set": bson.M{"confirmed": confirmed}})
	return
}


func (manager *Manager) GetByKeyPhonenumberValidationInformation(key string) (info *PhonenumberValidationInformation,  err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	err = mgoCollection.Find(bson.M{"key": key}).One(&info)
	if err == mgo.ErrNotFound {
		info = nil
		err = nil
	}
	return
}

func (manager *Manager) NewValidatedPhonenumber(username string, phonenumber string) (validatedphonenumber *ValidatedPhonenumber) {
	validatedphonenumber = &ValidatedPhonenumber{CreatedAt: time.Now(), Username: username, Phonenumber: string(phonenumber)}
	return
}

func (manager *Manager) SaveValidatedPhonenumber(validated *ValidatedPhonenumber) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	err = mgoCollection.Insert(validated)
	return
}


func (manager *Manager) GetByUsernameValidatedPhonenumbers(username string) (validatednumbers []ValidatedPhonenumber, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	err = mgoCollection.Find(bson.M{"username": username}).All(&validatednumbers)
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
