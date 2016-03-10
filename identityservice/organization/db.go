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

//Manager is used to store organizations
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

// All get all organizations.
// TODO: this method can take username(i.e. owner or members?) as filtering parameter.
func (m *Manager) All() ([]Organization, error) {
	organizations := []Organization{}

	if err := m.collection.Find(nil).All(&organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}

// AllByUser get organizations for certain user.
func (m *Manager) AllByUser(username string) ([]Organization, error) {
	var organizations []Organization
	//TODO: handle this a bit smarter, select only the ones where the user is owner first, and take select only the org name
	//do the same for the orgs where the username is member but not owners
	//No need to pull in 1000's of records for this

	condition := []interface{}{
		bson.M{"members": username},
		bson.M{"owners": username},
	}

	err := m.collection.Find(bson.M{"$or": condition}).All(&organizations)

	return organizations, err
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

// SaveDns save or update DNS
func (m *Manager) SaveDns(organization *Organization, dns string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"dns": dns}})
}

// RemoveDns remove DNS
func (m *Manager) RemoveDns(organization *Organization, dns string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"dns": dns}})
}

// SaveMember save or update member
func (m *Manager) SaveMember(organization *Organization, username string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"members": username}})
}

// RemoveMember remove member
func (m *Manager) RemoveMember(organization *Organization, username string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"members": username}})
}

// SaveOwner save or update owners
func (m *Manager) SaveOwner(organization *Organization, owner string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"owners": owner}})
}

// RemoveOwner remove owner
func (m *Manager) RemoveOwner(organization *Organization, owner string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"owners": owner}})
}

func (o *Organization) GetId() string {
	return o.Id.Hex()
}
