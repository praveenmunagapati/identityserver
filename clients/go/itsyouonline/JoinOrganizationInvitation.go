package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type JoinOrganizationInvitation struct {
	Created      DateTime `json:"created,omitempty"`
	Organization string   `json:"organization" validate:"nonzero"`
	Role         string   `json:"role" validate:"nonzero"`
	User         string   `json:"user" validate:"nonzero"`
}

func (s JoinOrganizationInvitation) Validate() error {

	return validator.Validate(s)
}
