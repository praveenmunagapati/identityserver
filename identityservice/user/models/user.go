package models

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

type UserAddress struct {
	City       string `json:"city"`
	Street     string `json:"street"`
	Nr         string `json:"nr"`
	Other      string `json:"other"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type User struct {
	Id       bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	Username string                 `json:"username"`
	Expires  int64                  `json:"expires"`
	Email    map[string]string      `json:"email"`
	Phone    map[string]string      `json:"phone"`
	Address  map[string]UserAddress `json:"address"`

	session    *mgo.Session
	collection *mgo.Collection
}

type UserList []*User

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

// Get all users from DB.
func (um *UserManager) All() (UserList, error) {
	var userList UserList

	err := um.collection.Find(nil).All(&userList)

	return userList, err
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

// Return new blank user.
func NewUser(r *http.Request) *User {
	session := db.GetDBSession(r)
	return &User{
		Id:         "",
		session:    session,
		collection: getUserCollection(session),
	}
}

func (u *User) GetId() string {
	return u.Id.Hex()
}

// Save current user.
func (u *User) Save() error {
	// TODO: Validation!

	if u.Id == "" {
		// New Doc!
		u.Id = bson.NewObjectId()
		err := u.collection.Insert(u)
		return err
	}

	_, err := u.collection.UpsertId(u.Id, u)

	return err
}

// Delete current user.
func (u *User) Delete() error {
	if u.Id == "" {
		return errors.New("User not stored")
	}

	return u.collection.RemoveId(u.Id)
}
