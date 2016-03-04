package company

import "gopkg.in/mgo.v2/bson"

type Company struct {
	Id            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Expire        Date          `json:"expire"`
	Globalid      string        `json:"globalid"`
	Info          []string      `json:"info"`
	Organizations []string      `json:"organizations"`
	PublicKeys    []string      `json:"publicKeys"`
	Taxnr         string        `json:"taxnr"`
}
