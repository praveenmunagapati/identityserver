package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernameBanksGetRespBody struct {
	Type []BankAccount `json:"type" validate:"nonzero"`
}

func (s UsersUsernameBanksGetRespBody) Validate() error {

	return validator.Validate(s)
}
