package organization

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName = "organizations"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	// TODO: Use model tags to ensure indices/constraints.
	index := mgo.Index{
		Key:    []string{"globalid"},
		Unique: true,
	}

	db.EnsureIndex(mongoCollectionName, index)
}

//Manager is used to store users
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: getCollection(session),
	}
}

// Get organization by ID.
func (m *Manager) Get(id string) (*Organization, error) {
	var organization Organization

	objectId := bson.ObjectIdHex(id)

	if err := m.collection.FindId(objectId).One(&organization); err != nil {
		return nil, err
	}

	return &organization, nil
}

// Get organization by Name.
func (m *Manager) GetByName(globalId string) (*Organization, error) {
	var organization Organization

	err := m.collection.Find(bson.M{"globalid": globalId}).One(&organization)

	return &organization, err
}

// Check if organization exists.
func (m *Manager) Exists(globalId string) bool {
	count, _ := m.collection.Find(bson.M{"globalid": globalId}).Count()

	return count != 1
}

func (o *Organization) GetId() string {
	return o.Id.Hex()
}

// Save a organization.
func (m *Manager) Save(organization *Organization) error {
	// TODO: Validation!

	if organization.Id == "" {
		// New Doc!
		organization.Id = bson.NewObjectId()
		err := m.collection.Insert(organization)
		return err
	}

	_, err := m.collection.UpsertId(organization.Id, organization)

	return err
}

// Delete a organization.
func (m *Manager) Delete(organization *Organization) error {
	if organization.Id == "" {
		return errors.New("organization not stored")
	}

	return m.collection.RemoveId(organization.Id)
}
