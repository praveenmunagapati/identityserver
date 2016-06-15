package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameAddressesPostRespBody struct {
	Address Address `json:"address" validate:"nonzero"`
	Label   Label   `json:"label" validate:"nonzero"`
}

func (s UsersUsernameAddressesPostRespBody) Validate() error {

	return validator.Validate(s)
}
