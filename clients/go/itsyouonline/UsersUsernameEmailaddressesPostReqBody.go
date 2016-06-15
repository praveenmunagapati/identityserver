package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameEmailaddressesPostReqBody struct {
	Emailaddress string `json:"emailaddress" validate:"nonzero"`
	Label        Label  `json:"label" validate:"nonzero"`
}

func (s UsersUsernameEmailaddressesPostReqBody) Validate() error {

	return validator.Validate(s)
}
