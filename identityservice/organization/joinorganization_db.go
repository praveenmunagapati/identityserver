package organization

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoOrganizationRequestCollectionName = "organization-requests"
)

//OrganizationRequestManager is used to store organizations
type OrganizationRequestManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getOrganizationRequestCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoOrganizationRequestCollectionName)
}

//NewOrganizationRequestManager creates and initializes a new OrganizationRequestManager
func NewOrganizationRequestManager(r *http.Request) *OrganizationRequestManager {
	session := db.GetDBSession(r)
	return &OrganizationRequestManager{
		session:    session,
		collection: getOrganizationRequestCollection(session),
	}
}

// GetByUser get all requests for a user.
func (o *OrganizationRequestManager) GetByUser(username string) ([]JoinOrganizationRequest, error) {
	orgRequests := []JoinOrganizationRequest{}

	err := o.collection.Find(bson.M{"user": username}).All(&orgRequests)

	return orgRequests, err
}

func (o *OrganizationRequestManager) Get(username string, organization string, role string) (*JoinOrganizationRequest, error) {
	var orgRequest JoinOrganizationRequest

	query := bson.M{
		"user":         username,
		"role":         role,
		"organization": organization,
	}

	err := o.collection.Find(query).One(&orgRequest)

	return &orgRequest, err
}

// Save save/update join request
func (o *OrganizationRequestManager) Save(joinRequest *JoinOrganizationRequest) error {
	if joinRequest.Id == "" {
		// New Doc!
		joinRequest.Id = bson.NewObjectId()
		err := o.collection.Insert(joinRequest)
		return err
	}

	_, err := o.collection.UpsertId(joinRequest.Id, joinRequest)

	return err
}
