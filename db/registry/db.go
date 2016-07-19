package registry

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoRegistryCollectionName = "registry"
)

//ErrUsernameOrGlobalIDRequired is used to indicate that no username of globalid were specified
var ErrUsernameOrGlobalIDRequired = errors.New("Username or globalid is required")

//ErrUsernameAndGlobalIDAreMutuallyExclusive is the error given when both a username and a globalid were given
var ErrUsernameAndGlobalIDAreMutuallyExclusive = errors.New("Username and globalid can not both be specified")

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}

	db.EnsureIndex(mongoRegistryCollectionName, index)

	index = mgo.Index{
		Key:    []string{"globalid"},
		Unique: true,
	}

	db.EnsureIndex(mongoRegistryCollectionName, index)

	index = mgo.Index{
		Key:    []string{"entries.key"},
		Unique: true,
	}

	db.EnsureIndex(mongoRegistryCollectionName, index)

}

//Manager is used to store KeyValuePairs in a user or organization registry
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

func (m *Manager) getRegistryCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoRegistryCollectionName)
}

func validateUsernameAndGlobalID(username, globalid string) (err error) {
	if username == "" && globalid == "" {
		err = ErrUsernameOrGlobalIDRequired
	}
	if username != "" && globalid != "" {
		err = ErrUsernameAndGlobalIDAreMutuallyExclusive
	}
	return
}

func createSelector(username, globalid, key string) (selector bson.M, err error) {
	err = validateUsernameAndGlobalID(username, globalid)
	if err != nil {
		return
	}
	if username != "" {
		selector = bson.M{"username": username, "entries.key": key}
	} else {
		selector = bson.M{"globalid": globalid, "entries.key": key}
	}
	return
}

//DeleteRegistryEntry deletes a registry entry
// Either a username or a globalid needs to be given
// If the key does not exist, no error is returned
func (m *Manager) DeleteRegistryEntry(username string, globalid string, key string) (err error) {
	selector, err := createSelector(username, globalid, key)
	if err != nil {
		return
	}
	m.getRegistryCollection().RemoveAll(selector)
	return
}

//UpsertRegistryEntry updates or inserts a registry entry
// Either a username or a globalid needs to be given
func (m *Manager) UpsertRegistryEntry(username string, globalid string, registryEntry RegistryEntry) (err error) {
	selector, err := createSelector(username, globalid, registryEntry.Key)
	if err != nil {
		return
	}
	m.getRegistryCollection().Upsert(selector, &registryEntry)
	return
}
