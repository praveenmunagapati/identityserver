package invitations

import (
	"github.com/itsyouonline/identityserver/db"
	"gopkg.in/mgo.v2/bson"
)

type InvitationStatus string

const (
	RequestPending  InvitationStatus = "pending"
	RequestAccepted InvitationStatus = "accepted"
	RequestRejected InvitationStatus = "rejected"
)

const (
	RoleMember = "member"
	RoleOwner  = "owner"
)

type InviteMethod string

const (
	MethodWebsite InviteMethod = "website"
	MethodEmail   InviteMethod = "email"
	MethodPhone   InviteMethod = "phone"
)
const (
	InviteUrl = "https://%s/login#/organizationinvite/%s"
)

//JoinOrganizationInvitation defines an invitation to join an organization
type JoinOrganizationInvitation struct {
	ID           bson.ObjectId    `json:"-" bson:"_id,omitempty"`
	Organization string           `json:"organization"`
	Role         string           `json:"role"`
	User         string           `json:"user"`
	Status       InvitationStatus `json:"status"`
	Created      db.DateTime      `json:"created"`
	Method       InviteMethod     `json:"method"`
	EmailAddress string           `json:"emailaddress"`
	PhoneNumber  string           `json:"phonenumber"`
	Code         string           `json:"-"`
}
