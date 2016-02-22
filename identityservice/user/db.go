package user

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName = "users"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	// TODO: Use model tags to ensure indices/constraints.
	index := mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoCollectionName, index)
}

//Manager is used to store users
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getUserCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: getUserCollection(session),
	}
}

// Get user by ID.
func (um *Manager) Get(id string) (*User, error) {
	var user User

	objectID := bson.ObjectIdHex(id)

	if err := um.collection.FindId(objectID).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

//GetByName gets a user by it's username.
func (um *Manager) GetByName(username string) (*User, error) {
	var user User

	err := um.collection.Find(bson.M{"username": username}).One(&user)

	return &user, err
}

//Exists checks if a user with this username already exists.
func (um *Manager) Exists(username string) bool {
	count, _ := um.collection.Find(bson.M{"username": username}).Count()

	return count != 1
}

func (u *User) getID() string {
	return u.Id.Hex()
}

// Save a user.
func (um *Manager) Save(u *User) error {
	// TODO: Validation!

	if u.Id == "" {
		// New Doc!
		u.Id = bson.NewObjectId()
		err := um.collection.Insert(u)
		return err
	}

	_, err := um.collection.UpsertId(u.Id, u)

	return err
}

// Delete a user.
func (um *Manager) Delete(u *User) error {
	if u.Id == "" {
		return errors.New("User not stored")
	}

	return um.collection.RemoveId(u.Id)
}
