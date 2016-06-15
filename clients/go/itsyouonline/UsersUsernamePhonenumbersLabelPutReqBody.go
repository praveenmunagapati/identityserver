package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernamePhonenumbersLabelPutReqBody struct {
	Label       Label       `json:"label" validate:"nonzero"`
	Phonenumber Phonenumber `json:"phonenumber" validate:"nonzero"`
}

func (s UsersUsernamePhonenumbersLabelPutReqBody) Validate() error {

	return validator.Validate(s)
}
