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
	logoCollectionName = "organizationLogos"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	// TODO: Use model tags to ensure indices/constraints.
	index := mgo.Index{
		Key:    []string{"globalid"},
		Unique: true,
	}

	db.EnsureIndex(mongoCollectionName, index)

	// Index the logo collection
	index = mgo.Index {
		Key:		[]string{"globalid"},
		Unique:	true,
	}

	db.EnsureIndex(logoCollectionName, index)
}

//Manager is used to store organizations
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//LogoManager is used to save the logo for an organization
type LogoManager struct {
	session		 *mgo.Session
	collection *mgo.Collection
}

func getCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

//get the logo collection
func getLogoCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, logoCollectionName)
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: getCollection(session),
	}
}

//NewLogoManager creates and initializes a new LogoManager
func NewLogoManager(r *http.Request) *LogoManager {
	session := db.GetDBSession(r)
	return &LogoManager{
		session:		session,
		collection:	getLogoCollection(session),
	}
}

// GetOrganizations gets a list of organizations.
func (m *Manager) GetOrganizations(organizationIDs []string) ([]Organization, error) {
	var organizations []Organization

	err := m.collection.Find(bson.M{"globalid": bson.M{"$in": organizationIDs}}).All(&organizations)

	return organizations, err
}

// GetSubOrganizations returns all organizations which have {globalID} as parent (including the organization with {globalID} as globalid)
//TODO: validate globalID since it is appended in the query
//TODO: put an index on the globalid field
func (m *Manager) GetSubOrganizations(globalID string) ([]Organization, error) {
	var organizations = make([]Organization, 0, 0)
	var qry = bson.M{"globalid": bson.M{"$regex": bson.RegEx{"^" + globalID + `\.`, ""}}}
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
func (m *Manager) IsMember(globalID, username string) (ismember bool, err error) {
	matches, err := m.collection.Find(bson.M{"globalid": globalID, "members": username}).Count()
	ismember = (matches > 0)
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
func (m *Manager) GetByName(globalID string) (organization *Organization, err error) {
	err = m.collection.Find(bson.M{"globalid": globalID}).One(&organization)
	return
}

func (m *LogoManager) GetByName(globalID string) (organization *Organization, err error) {
	err = m.collection.Find(bson.M{"globalid": globalID}).One(&organization)
	return
}

// Exists checks if an organization exists.
func (m *Manager) Exists(globalID string) bool {
	count, _ := m.collection.Find(bson.M{"globalid": globalID}).Count()

	return count == 1
}

// Exists checks if an organization and logo entry exists.
func (m *LogoManager) Exists(globalID string) bool {
	count, _ := m.collection.Find(bson.M{"globalid": globalID}).Count()

	return count == 1
}

// Create a new organization.
func (m *Manager) Create(organization *Organization) error {
	// TODO: Validation!

	err := m.collection.Insert(organization)
	if mgo.IsDup(err) {
		return db.ErrDuplicate
	}
	return err
}

// Create a new organization entry in the organization logo collection
func (m *LogoManager) Create(organization *Organization) error {
	var orgLogo OrganizationLogo

	orgLogo.Globalid = organization.Globalid

	err := m.collection.Insert(orgLogo)
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

func (m *Manager) AddDNS(organization *Organization, dnsName string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"dns": dnsName}})
}

func (m *Manager) UpdateDNS(organization *Organization, oldDNSName string, newDNSName string) error {
	err := m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"dns": oldDNSName}})
	if err != nil {
		return err
	}
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"dns": newDNSName}})
}

// RemoveDNS remove DNS
func (m *Manager) RemoveDNS(organization *Organization, dns string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"dns": dns}})
}

// Remove removes the organization
func (m *Manager) Remove(globalid string) error {
	return m.collection.Remove(bson.M{"globalid": globalid})
}

// Remove the organization logo
func (m *LogoManager) Remove(globalid string) error {
	return m.collection.Remove(bson.M{"globalid": globalid})
}

// UpdateMembership Updates a user his role in an organization
func (m *Manager) UpdateMembership(globalid string, username string, oldrole string, newrole string) error {
	qry := bson.M{"globalid": globalid}
	pull := bson.M{
		"$pull": bson.M{oldrole: username},
	}
	push := bson.M{
		"$addToSet": bson.M{newrole: username},
	}
	err := m.collection.Update(qry, pull)
	if err != nil {
		return err
	}
	return m.collection.Update(qry, push)
}

// CountByUser counts the amount of organizations by user
func (m *Manager) CountByUser(username string) (int, error) {
	qry := bson.M{"owners": username}
	return m.collection.Find(qry).Count()
}

// RemoveUser Removes a user from an organization
func (m *Manager) RemoveUser(globalId string, username string) error {
	qry := bson.M{"globalid": globalId}
	update := bson.M{"$pull": bson.M{"owners": username, "members": username}}
	return m.collection.Update(qry, update)
}

// SaveLogo save or update logo
func (m *LogoManager) SaveLogo(globalId string, logo string) (*mgo.ChangeInfo, error) {
	return m.collection.Upsert(
		bson.M{"globalid": globalId},
		bson.M{"$set": bson.M{"logo": logo}})
}

// GetLogo Gets the logo from an organization
func(m *LogoManager) GetLogo(globalId string) (string, error) {
	var org *OrganizationLogo
	err := m.collection.Find(bson.M{"globalid": globalId}).One(&org)
	if err != nil {
		return "", err
	}
	return org.Logo, err
}

// RemoveLogo Removes the logo from an organization
func (m *LogoManager) RemoveLogo(globalId string) error {
	qry := bson.M{"globalid": globalId}
	update := bson.M{"$unset": bson.M{"logo": 1}}
	return m.collection.Update(qry, update)
}
