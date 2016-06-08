package password

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/credentials/password/keyderivation"
	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName = "password"
)

type userPass struct {
	Username string
	Password string
}

//InitModels initializes models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoCollectionName, index)
}

//Manager stores and validates passwords
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getPasswordCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

//NewManager creates a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: getPasswordCollection(session),
	}
}

//Validate checks the password for a specific username
func (pwm *Manager) Validate(username, password string) (bool, error) {
	var storedPassword userPass
	if err := pwm.collection.Find(bson.M{"username": username}).One(&storedPassword); err != nil {
		if err == mgo.ErrNotFound {
			log.Debug("No password found for this user")
			return false, nil
		}
		log.Debug(err)
		return false, err
	}
	match := keyderivation.Check(password, storedPassword.Password)
	return match, nil
}

// Save stores a password for a specific username.
func (pwm *Manager) Save(username, password string) error {
	//TODO: username and password validation
	passwordHash, err := keyderivation.Hash(password)
	if err != nil {
		log.Error("ERROR hashing password")
		log.Debug("ERROR hashing password: ", err)
		return errors.New("internal_error")
	}
	if len(password) < 6 {
		return errors.New("invalid_password")
	}
	storedPassword := userPass{Username: username, Password: passwordHash}

	_, err = pwm.collection.Upsert(bson.M{"username": username}, storedPassword)

	return err
}
