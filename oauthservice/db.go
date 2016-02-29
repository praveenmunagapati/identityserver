package oauthservice

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	requestsCollectionName = "oauth_authorizationrequests"
	tokensCollectionName   = "oauth_accesstokens"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:    []string{"authorizationcode"},
		Unique: true,
	} //Do not drop duplicates since it would hijack another authorizationrequest, better to error out

	db.EnsureIndex(requestsCollectionName, index)

	//TODO: unique username/clientid combination

	automaticExpiration := mgo.Index{
		Key:         []string{"CreatedAt"},
		ExpireAfter: time.Second * 10,
		Background:  true,
	}
	db.EnsureIndex(requestsCollectionName, automaticExpiration)

}

//Manager is used to store users
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: db.GetCollection(session, requestsCollectionName),
	}
}

// Get an authorizationRequest by it's authorizationcode.
func (m *Manager) Get(authorizationcode string) (*authorizationRequest, error) {
	var ar authorizationRequest

	err := m.collection.Find(bson.M{"authorizationcode": authorizationcode}).One(&ar)

	return &ar, err
}

// Save stores an authorizationRequest, only adding new authorizationRequests is allowed, updating is not
func (m *Manager) Save(ar *authorizationRequest) (err error) {
	// TODO: Validation!

	err = m.collection.Insert(ar)

	return
}
