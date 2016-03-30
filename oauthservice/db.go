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

//getAuthorizationRequestCollection returns the mongo collection for the authorizationRequests
func (m *Manager) getAuthorizationRequestCollection() *mgo.Collection {
	return db.GetCollection(m.session, requestsCollectionName)
}

//getAccessTokenCollection returns the mongo collection for the accessTokens
func (m *Manager) getAccessTokenCollection() *mgo.Collection {
	return db.GetCollection(m.session, tokensCollectionName)
}

// Get an authorizationRequest by it's authorizationcode.
func (m *Manager) Get(authorizationcode string) (*authorizationRequest, error) {
	var ar authorizationRequest

	err := m.getAuthorizationRequestCollection().Find(bson.M{"authorizationcode": authorizationcode}).One(&ar)

	return &ar, err
}

// saveAuthorizationRequest stores an authorizationRequest, only adding new authorizationRequests is allowed, updating is not
func (m *Manager) saveAuthorizationRequest(ar *authorizationRequest) (err error) {
	// TODO: Validation!

	err = m.getAuthorizationRequestCollection().Insert(ar)

	return
}

// saveAccessToken stores an accessToken, only adding new accessTokens is allowed, updating is not
func (m *Manager) saveAccessToken(at *AccessToken) (err error) {
	// TODO: Validation!

	err = m.getAccessTokenCollection().Insert(at)

	return
}

//GetAccessToken gets an access token by it's actual token string
// If the token is not found or is expired, nil is returned
func (m *Manager) GetAccessToken(token string) (at *AccessToken, err error) {
	at = &AccessToken{}

	err = m.getAccessTokenCollection().Find(bson.M{"accesstoken": token}).One(at)
	if err != nil && err == mgo.ErrNotFound {
		at = nil
		err = nil
		return
	}
	if err != nil {
		at = nil
		return
	}
	if at.IsExpired() {
		at = nil
		err = nil
	}

	return
}

//getClientsCollection returns the mongo collection for the clients
func (m *Manager) getClientsCollection() *mgo.Collection {
	return db.GetCollection(m.session, clientsCollectionName)
}

//GetClientSecretLabels returns a list of labels for which there are secrets registered for a specific client
func (m *Manager) GetClientSecretLabels(clientID string) (labels []string, err error) {
	results := []struct{ Label string }{}
	err = m.getClientsCollection().Find(bson.M{"clientid": clientID}).Select(bson.M{"label": 1}).All(&results)
	labels = make([]string, len(results), len(results))
	for i, value := range results {
		labels[i] = value.Label
	}
	return
}

//CreateClientSecret saves an Oauth2 client secret
func (m *Manager) CreateClientSecret(client *Oauth2Client) (err error) {

	err = m.getClientsCollection().Insert(client)

	if err != nil && mgo.IsDup(err) {
		err = db.ErrDuplicate
	}
	return
}

//RenameClientSecret changes the label for a client secret
func (m *Manager) RenameClientSecret(clientID, oldLabel, newLabel string) (err error) {

	_, err = m.getClientsCollection().UpdateAll(bson.M{"clientid": clientID, "label": oldLabel}, bson.M{"$set": bson.M{"label": newLabel}})

	if err != nil && mgo.IsDup(err) {
		err = db.ErrDuplicate
	}
	return
}

//DeleteClientSecret removes a client secret by it's clientID and label
func (m *Manager) DeleteClientSecret(clientID, label string) (err error) {
	_, err = m.getClientsCollection().RemoveAll(bson.M{"clientid": clientID, "label": label})
	return
}

//GetClientSecret retrieves a clientsecret given a clientid and a label
func (m *Manager) GetClientSecret(clientID, label string) (secret string, err error) {
	c := &Oauth2Client{}
	err = m.getClientsCollection().Find(bson.M{"clientid": clientID, "label": label}).One(c)
	if err == mgo.ErrNotFound {
		err = nil
		return
	}
	secret = c.Secret
	return
}
