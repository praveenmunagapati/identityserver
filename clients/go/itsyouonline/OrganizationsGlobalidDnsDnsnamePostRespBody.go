package client

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidDnsDnsnamePostRespBody struct {
	Name string `json:"name" validate:"nonzero"`
}

func (s OrganizationsGlobalidDnsDnsnamePostRespBody) Validate() error {

	return validator.Validate(s)
}
