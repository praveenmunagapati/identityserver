package user

import (
	"errors"
	"fmt"
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
func (m *Manager) Get(id string) (*User, error) {
	var user User

	objectID := bson.ObjectIdHex(id)

	if err := m.collection.FindId(objectID).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

//GetByName gets a user by it's username.
func (m *Manager) GetByName(username string) (*User, error) {
	var user User

	err := m.collection.Find(bson.M{"username": username}).One(&user)

	return &user, err
}

//Exists checks if a user with this username already exists.
func (m *Manager) Exists(username string) (bool, error) {
	count, err := m.collection.Find(bson.M{"username": username}).Count()

	return count >= 1, err
}

// Save a user.
func (m *Manager) Save(u *User) error {
	// TODO: Validation!

	if u.Id == "" {
		// New Doc!
		u.Id = bson.NewObjectId()
		err := m.collection.Insert(u)
		return err
	}

	_, err := m.collection.UpsertId(u.Id, u)

	return err
}

// Delete a user.
func (m *Manager) Delete(u *User) error {
	if u.Id == "" {
		return errors.New("User not stored")
	}

	return m.collection.RemoveId(u.Id)
}

// SaveEmail save or update email along with its label
func (m *Manager) SaveEmail(u *User, label string, email string) error {
	emailLabel := fmt.Sprintf("email.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$set": bson.M{emailLabel: email}})
}

// RemoveEmail remove email associated with label
func (m *Manager) RemoveEmail(u *User, label string) error {
	emailLabel := fmt.Sprintf("email.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$unset": bson.M{emailLabel: ""}})
}

// SavePhone save or update phone along with its label
func (m *Manager) SavePhone(u *User, label string, phonenumber Phonenumber) error {
	phoneLabel := fmt.Sprintf("phone.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$set": bson.M{phoneLabel: phonenumber}})
}

// RemovePhone remove phone associated with label
func (m *Manager) RemovePhone(u *User, label string) error {
	phoneLabel := fmt.Sprintf("phone.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$unset": bson.M{phoneLabel: ""}})
}

// SaveAddress save or update address along with its label
func (m *Manager) SaveAddress(u *User, label string, address Address) error {
	addressLabel := fmt.Sprintf("address.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$set": bson.M{addressLabel: address}})
}

// RemoveAddress remove address associated with label
func (m *Manager) RemoveAddress(u *User, label string) error {
	addressLabel := fmt.Sprintf("address.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$unset": bson.M{addressLabel: ""}})
}

// SaveBank save or update bank account along with its label
func (m *Manager) SaveBank(u *User, label string, bank BankAccount) error {
	bankLabel := fmt.Sprintf("bank.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$set": bson.M{bankLabel: bank}})
}

// RemoveBank remove bank associated with label
func (m *Manager) RemoveBank(u *User, label string) error {
	bankLabel := fmt.Sprintf("bank.%s", label)

	return m.collection.Update(
		bson.M{"username": u.Username},
		bson.M{"$unset": bson.M{bankLabel: ""}})
}

func (u *User) getID() string {
	return u.Id.Hex()
}
