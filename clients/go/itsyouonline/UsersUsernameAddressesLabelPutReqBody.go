package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameAddressesLabelPutReqBody struct {
	Label       Label   `json:"label" validate:"nonzero"`
	Phonenumber Address `json:"phonenumber" validate:"nonzero"`
}

func (s UsersUsernameAddressesLabelPutReqBody) Validate() error {

	return validator.Validate(s)
}
