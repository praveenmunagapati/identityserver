package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type DigitalAssetAddress struct {
	Address        string   `json:"address" validate:"nonzero"`
	Currencysymbol string   `json:"currencysymbol" validate:"nonzero"`
	Expire         DateTime `json:"expire" validate:"nonzero"`
	Label          string   `json:"label" validate:"nonzero"`
}

func (s DigitalAssetAddress) Validate() error {

	return validator.Validate(s)
}
