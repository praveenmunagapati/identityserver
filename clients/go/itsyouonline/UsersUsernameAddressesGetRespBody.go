package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameAddressesGetRespBody struct {
	Type []Address `json:"type" validate:"nonzero"`
}

func (s UsersUsernameAddressesGetRespBody) Validate() error {

	return validator.Validate(s)
}
