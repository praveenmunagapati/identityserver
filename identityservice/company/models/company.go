package models

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

type Company struct {
	Id            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	GlobalId      string        `json:"globalId"`
	Expires       int           `json:"expires"`
	TaxNr         string        `json:"taxNr"`
	Organizations []string      `json:"organizations"`
	Info          []string      `json:"info"`

	session    *mgo.Session
	collection *mgo.Collection
}

type CompanyList []*Company

type CompanyManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getCompanyCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, COLLECTION_COMPANIES)
}

func NewCompanyManager(r *http.Request) *CompanyManager {
	session := db.GetDBSession(r)
	return &CompanyManager{
		session:    session,
		collection: getCompanyCollection(session),
	}
}

// Get all companies from DB.
func (cm *CompanyManager) All() (CompanyList, error) {
	var companyList CompanyList

	err := cm.collection.Find(nil).All(&companyList)

	return companyList, err
}

// Get company by ID.
func (cm *CompanyManager) Get(id string) (*Company, error) {
	var company Company

	objectId := bson.ObjectIdHex(id)

	if err := cm.collection.FindId(objectId).One(&company); err != nil {
		return nil, err
	}

	return &company, nil
}

// Get company by Name.
func (cm *CompanyManager) GetByName(globalId string) (*Company, error) {
	var company Company

	err := cm.collection.Find(bson.M{"globalid": globalId}).One(&company)

	return &company, err
}

// Check if company exists.
func (cm *CompanyManager) Exists(globalId string) bool {
	count, _ := cm.collection.Find(bson.M{"globalid": globalId}).Count()

	return count != 1
}

func NewCompany(r *http.Request) *Company {
	session := db.GetDBSession(r)
	return &Company{
		Id:            "",
		Organizations: []string{},
		Info:          []string{},
		session:       session,
		collection:    getCompanyCollection(session),
	}
}

func (c *Company) GetId() string {
	return c.Id.Hex()
}

// Save current company.
func (c *Company) Save() error {
	// TODO: Validation!

	if c.Id == "" {
		// New Doc!
		c.Id = bson.NewObjectId()
		err := c.collection.Insert(c)
		return err
	}

	_, err := c.collection.UpsertId(c.Id, c)

	return err
}

// Delete current company.
func (c *Company) Delete() error {
	if c.Id == "" {
		return errors.New("Company not stored")
	}

	return c.collection.RemoveId(c.Id)
}
