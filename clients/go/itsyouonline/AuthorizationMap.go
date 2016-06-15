package client

import (
	"gopkg.in/validator.v2"
)

// Mapping between requested labels and real labels
type AuthorizationMap struct {
	Reallabel      string `json:"reallabel" validate:"nonzero"`
	Requestedlabel string `json:"requestedlabel" validate:"nonzero"`
}

func (s AuthorizationMap) Validate() error {

	return validator.Validate(s)
}
