package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernamePhonenumbersLabelActivatePutRespBody struct {
	Error string `json:"error" validate:"nonzero"`
}

func (s UsersUsernamePhonenumbersLabelActivatePutRespBody) Validate() error {

	return validator.Validate(s)
}
