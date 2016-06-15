package client

import (
	"gopkg.in/validator.v2"
)

type OrganizationAPIKey struct {
	CallbackURL                string `json:"callbackURL,omitempty" validate:"min=5,max=250"`
	ClientCredentialsGrantType bool   `json:"clientCredentialsGrantType,omitempty"`
	Label                      string `json:"label" validate:"min=2,max=50,nonzero"`
	Secret                     string `json:"secret,omitempty" validate:"max=250"`
}

func (s OrganizationAPIKey) Validate() error {

	return validator.Validate(s)
}
