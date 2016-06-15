package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameAddressesPostReqBody struct {
	Address Address `json:"address" validate:"nonzero"`
	Label   Label   `json:"label" validate:"nonzero"`
}

func (s UsersUsernameAddressesPostReqBody) Validate() error {

	return validator.Validate(s)
}
