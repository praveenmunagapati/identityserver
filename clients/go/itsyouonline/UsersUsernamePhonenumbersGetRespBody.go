package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernamePhonenumbersGetRespBody struct {
	Type []Phonenumber `json:"type" validate:"nonzero"`
}

func (s UsersUsernamePhonenumbersGetRespBody) Validate() error {

	return validator.Validate(s)
}
