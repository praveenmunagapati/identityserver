package contract

import "github.com/itsyouonline/identityserver/db"

type Contract struct {
	Content      string      `json:"content"`
	ContractId   string      `json:"contractId"`
	ContractType string      `json:"contractType"`
	Expires      db.Date     `json:"expires"`
	Extends      []string    `json:"extends"`
	Invalidates  []string    `json:"invalidates"`
	Parties      []string    `json:"parties"`
	Signatures   []Signature `json:"signatures"`
}
