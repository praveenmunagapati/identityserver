package client

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidDnsDnsnamePostReqBody struct {
	Name string `json:"name" validate:"min=4,max=250,nonzero"`
}

func (s OrganizationsGlobalidDnsDnsnamePostReqBody) Validate() error {

	return validator.Validate(s)
}
