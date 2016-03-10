package organization

import "gopkg.in/mgo.v2/bson"

type RequestStatus int

const (
	RequestPending RequestStatus = iota
	RequestAccepted
	RequestRejected
)

const (
	RoleMember = "member"
	RoleOwner  = "owner"
)

type JoinOrganizationRequest struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Organization string        `json:"organization"`
	Role         []string      `json:"role"`
	User         string        `json:"user"`
	Status       RequestStatus `json:"status"`
}
