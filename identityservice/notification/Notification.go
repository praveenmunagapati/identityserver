package notification

import (
	"github.com/itsyouonline/identityserver/identityservice/contract"
	"github.com/itsyouonline/identityserver/identityservice/organization"
)

type Notification struct {
	Approvals        []organization.JoinOrganizationRequest `json:"approvals"`
	ContractRequests []contract.ContractSigningRequest      `json:"contractRequests"`
	Invitations      []organization.JoinOrganizationRequest `json:"invitations"`
}
