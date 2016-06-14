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
	mongoUsersCollectionName          = "users"
	mongoAuthorizationsCollectionName = "authorizations"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoUsersCollectionName, index)

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

func (m *Manager) getUserCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoUsersCollectionName)
}

func (m *Manager) getAuthorizationCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoAuthorizationsCollectionName)
}

// Get user by ID.
func (m *Manager) Get(id string) (*User, error) {
	var user User

	objectID := bson.ObjectIdHex(id)

	if err := m.getUserCollection().FindId(objectID).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

//GetByName gets a user by it's username.
func (m *Manager) GetByName(username string) (*User, error) {
	var user User

	err := m.getUserCollection().Find(bson.M{"username": username}).One(&user)

	return &user, err
}

//Exists checks if a user with this username already exists.
func (m *Manager) Exists(username string) (bool, error) {
	count, err := m.getUserCollection().Find(bson.M{"username": username}).Count()

	return count >= 1, err
}

// Save a user.
func (m *Manager) Save(u *User) error {
	// TODO: Validation!

	if u.ID == "" {
		// New Doc!
		u.ID = bson.NewObjectId()
		err := m.getUserCollection().Insert(u)
		return err
	}

	_, err := m.getUserCollection().UpsertId(u.ID, u)

	return err
}

// Delete a user.
func (m *Manager) Delete(u *User) error {
	if u.ID == "" {
		return errors.New("User not stored")
	}

	return m.getUserCollection().RemoveId(u.ID)
}

// SaveEmail save or update email along with its label
func (m *Manager) SaveEmail(username string, email EmailAddress) error {

	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"emailaddresses": email}})
}

// RemoveEmail remove email associated with label
func (m *Manager) RemoveEmail(username string, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"emailaddreses": bson.M{"label": label}}})
}

// SavePhone save or update phone along with its label
func (m *Manager) SavePhone(username string, phonenumber Phonenumber) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"phonenumbers": phonenumber}})
}

// RemovePhone remove phone associated with label
func (m *Manager) RemovePhone(username string, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"phonenumbers": bson.M{"label": label}}})
}

// SaveAddress save or update address along with its label
func (m *Manager) SaveAddress(username, label string, address Address) error {
	//TODO: is this safe to do?
	addressLabel := fmt.Sprintf("address.%s", label)

	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$set": bson.M{addressLabel: address}})
}

// RemoveAddress remove address associated with label
func (m *Manager) RemoveAddress(username, label string) error {
	//TODO: is this safe to do?
	addressLabel := fmt.Sprintf("address.%s", label)

	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$unset": bson.M{addressLabel: ""}})
}

// SaveBank save or update bank account along with its label
func (m *Manager) SaveBank(u *User, label string, bank BankAccount) error {
	bankLabel := fmt.Sprintf("bank.%s", label)

	return m.getUserCollection().Update(
		bson.M{"username": u.Username},
		bson.M{"$set": bson.M{bankLabel: bank}})
}

// RemoveBank remove bank associated with label
func (m *Manager) RemoveBank(u *User, label string) error {
	bankLabel := fmt.Sprintf("bank.%s", label)

	return m.getUserCollection().Update(
		bson.M{"username": u.Username},
		bson.M{"$unset": bson.M{bankLabel: ""}})
}

func (m *Manager) UpdateGithubAccount(username string, githubaccount GithubAccount) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"github": githubaccount}})
	return
}

func (m *Manager) DeleteGithubAccount(username string) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"github": bson.M{}}})
	return
}

func (m *Manager) UpdateFacebookAccount(username string, facebookaccount FacebookAccount) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"facebook": facebookaccount}})
	return
}

func (m *Manager) DeleteFacebookAccount(username string) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"facebook": bson.M{}}})
	return
}

// GetAuthorizationsByUser returns all authorizations for a specific user
func (m *Manager) GetAuthorizationsByUser(username string) (authorizations []Authorization, err error) {
	err = m.getAuthorizationCollection().Find(bson.M{"username": username}).All(&authorizations)
	return
}

//GetAuthorization returns the authorization for a specific organization, nil if no such auhorization exists
func (m *Manager) GetAuthorization(username, organization string) (authorization *Authorization, err error) {
	authorization = &Authorization{}
	err = m.getAuthorizationCollection().Find(bson.M{"username": username, "grantedto": organization}).One(authorization)
	if err == mgo.ErrNotFound {
		err = nil
	} else if err != nil {
		authorization = nil
	}
	return
}

//UpdateAuthorization inserts or updates an authorization
func (m *Manager) UpdateAuthorization(authorization *Authorization) (err error) {
	_, err = m.getAuthorizationCollection().Upsert(bson.M{"username": authorization.Username, "grantedto": authorization.GrantedTo}, authorization)
	return
}

//DeleteAuthorization removes an authorization
func (m *Manager) DeleteAuthorization(username, organization string) (err error) {
	_, err = m.getAuthorizationCollection().RemoveAll(bson.M{"username": username, "grantedto": organization})
	return
}

func (u *User) getID() string {
	return u.ID.Hex()
}

func (m *Manager) UpdateName(username string, firstname string, lastname string) (err error) {
	values := bson.M{
		"firstname": firstname,
		"lastname":  lastname,
	}
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": values})
	return
}
