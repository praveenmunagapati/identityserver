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
	clientsCollectionName  = "oauth_clients"
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
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 10,
		Background:  true,
	}
	db.EnsureIndex(requestsCollectionName, automaticExpiration)

	index = mgo.Index{
		Key:    []string{"accesstoken"},
		Unique: true,
	} //Do not drop duplicates since it would hijack another authorizationrequest, better to error out

	db.EnsureIndex(tokensCollectionName, index)

	//TODO: unique username/clientid combination

	automaticExpiration = mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: AccessTokenExpiration,
		Background:  true,
	}
	db.EnsureIndex(tokensCollectionName, automaticExpiration)

	index = mgo.Index{
		Key:    []string{"clientid", "label"},
		Unique: true,
	}
	db.EnsureIndex(clientsCollectionName, index)

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

//GetClientsCollection returns the mongo collection for the clients
func (m *Manager) GetClientsCollection() *mgo.Collection {
	return db.GetCollection(m.session, clientsCollectionName)
}

//GetClientSecretLabels returns a list of labels for which there are secrets registered for a specific client
func (m *Manager) GetClientSecretLabels(clientID string) (labels []string, err error) {
	results := []struct{ Label string }{}
	err = m.GetClientsCollection().Find(bson.M{"clientid": clientID}).Select(bson.M{"label": 1}).All(&results)
	labels = make([]string, len(results), len(results))
	for i, value := range results {
		labels[i] = value.Label
	}
	return
}

//CreateClientSecret saves an Oauth2 client secret
func (m *Manager) CreateClientSecret(client *Oauth2Client) (err error) {

	err = m.GetClientsCollection().Insert(client)

	if err != nil && mgo.IsDup(err) {
		err = db.ErrDuplicate
	}
	return
}

//RenameClientSecret changes the label for a client secret
func (m *Manager) RenameClientSecret(clientID, oldLabel, newLabel string) (err error) {

	_, err = m.GetClientsCollection().UpdateAll(bson.M{"clientid": clientID, "label": oldLabel}, bson.M{"$set": bson.M{"label": newLabel}})

	if err != nil && mgo.IsDup(err) {
		err = db.ErrDuplicate
	}
	return
}

//DeleteClientSecret removes a client secret by it's clientID and label
func (m *Manager) DeleteClientSecret(clientID, label string) (err error) {
	_, err = m.GetClientsCollection().RemoveAll(bson.M{"clientid": clientID, "label": label})
	return
}

//GetClientSecret retrieves a clientsecret given a clientid and a label
func (m *Manager) GetClientSecret(clientID, label string) (secret string, err error) {
	c := &Oauth2Client{}
	err = m.GetClientsCollection().Find(bson.M{"clientid": clientID, "label": label}).One(c)
	if err == mgo.ErrNotFound {
		err = nil
		return
	}
	secret = c.Secret
	return
}
