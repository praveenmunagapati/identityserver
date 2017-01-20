package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidOrgmembersInvitePostReqBody struct {
	Searchstring string `json:"searchstring" validate:"nonzero"`
}

func (s OrganizationsGlobalidOrgmembersInvitePostReqBody) Validate() error {

	return validator.Validate(s)
}
