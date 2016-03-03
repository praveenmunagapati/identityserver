package organization

import "gopkg.in/mgo.v2/bson"

type Organization struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Dns        []string      `json:"dns"`
	Globalid   string        `json:"globalid"`
	Includes   []string      `json:"includes"`
	Members    []string      `json:"members"`
	Owners     []string      `json:"owners"`
	PublicKeys []string      `json:"publicKeys"`
}
