package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameEmailaddressesPostRespBody struct {
	Emailaddress string `json:"emailaddress" validate:"nonzero"`
	Label        Label  `json:"label" validate:"nonzero"`
}

func (s UsersUsernameEmailaddressesPostRespBody) Validate() error {

	return validator.Validate(s)
}
