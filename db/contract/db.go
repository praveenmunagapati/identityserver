package contract

import (
	"errors"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName = "contracts"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:    []string{"ContractId"},
		Unique: true,
	}

	db.EnsureIndex(mongoCollectionName, index)
}

//Manager is used to store organizations
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: db.GetCollection(session, mongoCollectionName),
	}
}

//Save contract
func (m *Manager) Save(contract *Contract) (err error) {
	if contract.ContractId == "" {
		err = errors.New("Contractid can not be empty")
		return
	}
	err = m.collection.Insert(contract)
	return
}

//Get contract
func (m *Manager) Get(contractid string) (contract *Contract, err error) {
	contract = &Contract{}
	err = m.collection.Find(bson.M{"contractid": contractid}).One(contract)
	return
}

//Delete  contract
func (m *Manager) Delete(contractid string) (err error) {
	_, err = m.collection.RemoveAll(bson.M{"contractid": contractid})
	return
}

//GetByIncludedParty Get contracts that include the included party
func (m *Manager) GetByIncludedParty(party *Party, start int, max int, includeExpired bool) (contracts []Contract, err error) {
	contracts = make([]Contract, 0)
	query := bson.M{"parties.type": party.Type, "parties.name": party.Name}
	if !includeExpired {
		query["$or"] = []bson.M{bson.M{"expired": bson.M{"$lt": time.Now()}}, bson.M{"expired": nil}}
	}
	err = m.collection.Find(query).Skip(start).Limit(max).All(&contracts)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return
}
