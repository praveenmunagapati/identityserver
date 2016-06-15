package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameEmailaddressesLabelPutReqBody struct {
	Emailaddress string `json:"emailaddress" validate:"nonzero"`
	Label        Label  `json:"label" validate:"nonzero"`
}

func (s UsersUsernameEmailaddressesLabelPutReqBody) Validate() error {

	return validator.Validate(s)
}
