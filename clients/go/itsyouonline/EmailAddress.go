package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type EmailAddress struct {
	Emailaddress string `json:"emailaddress" validate:"nonzero"`
	Label        string `json:"label" validate:"nonzero"`
}

func (s EmailAddress) Validate() error {

	return validator.Validate(s)
}
