package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernamePhonenumbersLabelActivatePutReqBody struct {
	Smscode       string `json:"smscode" validate:"nonzero"`
	Validationkey string `json:"validationkey" validate:"nonzero"`
}

func (s UsersUsernamePhonenumbersLabelActivatePutReqBody) Validate() error {

	return validator.Validate(s)
}
