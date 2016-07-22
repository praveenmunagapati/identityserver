package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type Organization struct {
	Dns        []string `json:"dns" validate:"max=100,nonzero"`
	Globalid   string   `json:"globalid" validate:"min=3,max=150,regexp=^[a-z0-9]{3,150}$,nonzero"`
	Includes   []string `json:"includes" validate:"max=100,nonzero"`
	Members    []string `json:"members" validate:"max=2000,nonzero"`
	Owners     []string `json:"owners" validate:"max=20,nonzero"`
	PublicKeys []string `json:"publicKeys" validate:"max=20,nonzero"`
}

func (s Organization) Validate() error {

	return validator.Validate(s)
}
