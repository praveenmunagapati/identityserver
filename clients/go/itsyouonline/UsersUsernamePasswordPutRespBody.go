package client

import (
	"gopkg.in/validator.v2"
)

type UsersUsernamePasswordPutRespBody struct {
	Error string `json:"error" validate:"nonzero"`
}

func (s UsersUsernamePasswordPutRespBody) Validate() error {

	return validator.Validate(s)
}
