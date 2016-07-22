package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidMembersPutReqBody struct {
	Role     string `json:"role" validate:"nonzero"`
	Username string `json:"username" validate:"nonzero"`
}

func (s OrganizationsGlobalidMembersPutReqBody) Validate() error {

	return validator.Validate(s)
}
