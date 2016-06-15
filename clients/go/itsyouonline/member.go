package client

import (
	"gopkg.in/validator.v2"
)

type member struct {
	Username string `json:"username" validate:"nonzero"`
}

func (s member) Validate() error {

	return validator.Validate(s)
}
