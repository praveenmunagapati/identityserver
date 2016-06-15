package client

import (
	"gopkg.in/validator.v2"
)

type Phonenumber struct {
	Label       string `json:"label" validate:"nonzero"`
	Phonenumber string `json:"phonenumber" validate:"regexp=^\+?[0-9]+$,nonzero"`
}

func (s Phonenumber) Validate() error {

	return validator.Validate(s)
}
