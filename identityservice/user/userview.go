package user

import "github.com/itsyouonline/identityserver/db/user"

type Userview struct {
	Address       map[string]user.Address     `json:"address"`
	Bank          map[string]user.BankAccount `json:"bank"`
	Email         map[string]string      `json:"email"`
	Facebook      string                 `json:"facebook"`
	Github        string                 `json:"github"`
	Organizations []string               `json:"organizations"`
	Phone         map[string]user.Phonenumber `json:"phone"`
	PublicKeys    []string               `json:"publicKeys"`
	Username      string                 `json:"username"`
}
