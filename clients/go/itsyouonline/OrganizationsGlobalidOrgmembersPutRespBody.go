package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidOrgmembersPutRespBody struct {
	Org Organization `json:"org" validate:"nonzero"`
}

func (s OrganizationsGlobalidOrgmembersPutRespBody) Validate() error {

	return validator.Validate(s)
}
