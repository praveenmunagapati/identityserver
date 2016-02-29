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

	index = mgo.Index{
		Key:    []string{"AccessToken"},
		Unique: true,
	} //Do not drop duplicates since it would hijack another authorizationrequest, better to error out

	db.EnsureIndex(tokensCollectionName, index)

	//TODO: unique username/clientid combination

	automaticExpiration = mgo.Index{
		Key:         []string{"CreatedAt"},
		ExpireAfter: AccessTokenExpiration,
		Background:  true,
	}
	db.EnsureIndex(tokensCollectionName, automaticExpiration)

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

//GetAuthorizationRequestCollection returns the mongo collection for the authorizationRequests
func (m *Manager) GetAuthorizationRequestCollection() *mgo.Collection {
	return db.GetCollection(m.session, requestsCollectionName)
}

//GetAccessTokenCollection returns the mongo collection for the accessTokens
func (m *Manager) GetAccessTokenCollection() *mgo.Collection {
	return db.GetCollection(m.session, tokensCollectionName)
}

// Get an authorizationRequest by it's authorizationcode.
func (m *Manager) Get(authorizationcode string) (*authorizationRequest, error) {
	var ar authorizationRequest

	err := m.GetAuthorizationRequestCollection().Find(bson.M{"authorizationcode": authorizationcode}).One(&ar)

	return &ar, err
}

// SaveAuthorizationRequest stores an authorizationRequest, only adding new authorizationRequests is allowed, updating is not
func (m *Manager) SaveAuthorizationRequest(ar *authorizationRequest) (err error) {
	// TODO: Validation!

	err = m.GetAuthorizationRequestCollection().Insert(ar)

	return
}

// SaveAccessToken stores an accessToken, only adding new accessTokens is allowed, updating is not
func (m *Manager) SaveAccessToken(at *accessToken) (err error) {
	// TODO: Validation!

	err = m.GetAccessTokenCollection().Insert(at)

	return
}
