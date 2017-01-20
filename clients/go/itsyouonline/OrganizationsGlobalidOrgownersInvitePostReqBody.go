package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidOrgownersInvitePostReqBody struct {
	Searchstring string `json:"searchstring" validate:"nonzero"`
}

func (s OrganizationsGlobalidOrgownersInvitePostReqBody) Validate() error {

	return validator.Validate(s)
}
