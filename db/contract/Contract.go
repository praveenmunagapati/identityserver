package contract

import "github.com/itsyouonline/identityserver/db"

type Party struct {
	Type string
	Name string
}

type Contract struct {
	Content      string      `json:"content"`
	ContractType string      `json:"contractType"`
	Expires      db.Date     `json:"expires"`
	Extends      []string    `json:"extends"`
	Invalidates  []string    `json:"invalidates"`
	Parties      []Party     `json:"parties"`
	ContractId   string      `json:"contractId"`
	Signatures   []Signature `json:"signatures"`
}
