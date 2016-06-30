package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type UsersUsernamePhonenumbersLabelActivatePostRespBody struct {
	Validationkey string `json:"validationkey" validate:"nonzero"`
}

func (s UsersUsernamePhonenumbersLabelActivatePostRespBody) Validate() error {

	return validator.Validate(s)
}
