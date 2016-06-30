package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidDnsDnsnamePutReqBody struct {
	Newname string `json:"newname" validate:"min=4,max=250,nonzero"`
	Oldname string `json:"oldname" validate:"min=4,max=250,nonzero"`
}

func (s OrganizationsGlobalidDnsDnsnamePutReqBody) Validate() error {

	return validator.Validate(s)
}
