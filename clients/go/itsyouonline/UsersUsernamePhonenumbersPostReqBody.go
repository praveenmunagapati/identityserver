package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernamePhonenumbersPostReqBody struct {
	Label       Label       `json:"label" validate:"nonzero"`
	Phonenumber Phonenumber `json:"phonenumber" validate:"nonzero"`
}

func (s UsersUsernamePhonenumbersPostReqBody) Validate() error {

	return validator.Validate(s)
}
