package user

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	COLLECTION_USERS = "users"
)

// Initialize models in DB, if required.
func InitModels() {
	// TODO: Use model tags to ensure indices/constraints.
	index := mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(COLLECTION_USERS, index)
}

type UserManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getUserCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, COLLECTION_USERS)
}

func NewUserManager(r *http.Request) *UserManager {
	session := db.GetDBSession(r)
	return &UserManager{
		session:    session,
		collection: getUserCollection(session),
	}
}

// Get user by ID.
func (um *UserManager) Get(id string) (*User, error) {
	var user User

	objectId := bson.ObjectIdHex(id)

	if err := um.collection.FindId(objectId).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by Name.
func (um *UserManager) GetByName(username string) (*User, error) {
	var user User

	err := um.collection.Find(bson.M{"username": username}).One(&user)

	return &user, err
}

// Check if user exists.
func (um *UserManager) Exists(username string) bool {
	count, _ := um.collection.Find(bson.M{"username": username}).Count()

	return count != 1
}

func (u *User) GetId() string {
	return u.Id.Hex()
}

// Save a user.
func (um *UserManager) Save(u *User) error {
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
func (um *UserManager) Delete(u *User) error {
	if u.Id == "" {
		return errors.New("User not stored")
	}

	return um.collection.RemoveId(u.Id)
}
