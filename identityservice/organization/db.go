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

// AllRoot get all root organizations.
// TODO: this method can take username(i.e. owner or members?) as filtering parameter.
func (m *Manager) AllRoot() ([]Organization, error) {
	organizations := []Organization{}
	//query in mongodb shell: db.organizations.find({"globalid": {"$not":/\./}})
	var qry = bson.M{"globalid":  bson.M{"$not":bson.RegEx{`/\./`, ""}}}
	if err := m.collection.Find(qry).All(&organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}

// AllByRoot returns all organizations which have {globalID} as parent
func (m *Manager) AllByRoot(globalID string) ([]Organization, error) {
	organizations := []Organization{}
	var qry = bson.M{"globalid":  bson.M{"$regex":bson.RegEx{"/^(" + globalID + ").*/", ""}}}
	if err := m.collection.Find(qry).All(&organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}

//IsOwner checks if a specific user is in the owners list of an organization
func (m *Manager) IsOwner(globalID, username string) (isowner bool, err error) {
	matches, err := m.collection.Find(bson.M{"globalid": globalID, "owners": username}).Count()
	isowner = (matches > 0)
	return
}

//IsMember checks if a specific user is in the members list of an organization
func (m *Manager) IsMember(globalID, username string) (isowner bool, err error) {
	matches, err := m.collection.Find(bson.M{"globalid": globalID, "members": username}).Count()
	isowner = (matches > 0)
	return
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

	objectID := bson.ObjectIdHex(id)

	if err := m.collection.FindId(objectID).One(&organization); err != nil {
		return nil, err
	}

	return &organization, nil
}

// GetByName gets an organization by Name.
func (m *Manager) GetByName(globalID string) (*Organization, error) {
	var organization Organization

	err := m.collection.Find(bson.M{"globalid": globalID}).One(&organization)

	return &organization, err
}

// Exists checks if an organization exists.
func (m *Manager) Exists(globalID string) bool {
	count, _ := m.collection.Find(bson.M{"globalid": globalID}).Count()

	return count != 1
}

// Create a new organization.
func (m *Manager) Create(organization *Organization) error {
	// TODO: Validation!

	organization.Id = bson.NewObjectId()
	err := m.collection.Insert(organization)
	if mgo.IsDup(err) {
		return db.ErrDuplicate
	}
	return err
}

// Save an organization.
func (m *Manager) Save(organization *Organization) error {
	// TODO: Validation!

	// TODO: Save
	return errors.New("Save is not implemented yet")
}

// Delete a organization.
func (m *Manager) Delete(organization *Organization) error {
	if organization.Id == "" {
		return errors.New("organization not stored")
	}

	return m.collection.RemoveId(organization.Id)
}

// SaveDNS save or update DNS
func (m *Manager) SaveDNS(organization *Organization, dns string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"dns": dns}})
}

// RemoveDNS remove DNS
func (m *Manager) RemoveDNS(organization *Organization, dns string) error {
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
