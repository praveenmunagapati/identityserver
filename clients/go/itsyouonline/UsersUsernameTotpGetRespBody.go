package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameTotpGetRespBody struct {
	Totpsecret string `json:"totpsecret" validate:"nonzero"`
}

func (s UsersUsernameTotpGetRespBody) Validate() error {

	return validator.Validate(s)
}
