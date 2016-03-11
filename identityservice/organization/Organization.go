package organization

import "gopkg.in/mgo.v2/bson"

type Organization struct {
	Id         bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Dns        []string      `json:"dns"`
	Globalid   string        `json:"globalid"`
	Includes   []string      `json:"includes"`
	Members    []string      `json:"members"`
	Owners     []string      `json:"owners"`
	PublicKeys []string      `json:"publicKeys"`
}

// IsValid performs basic validation on the content of an organizations fields
func (c *Organization) IsValid() (valid bool) {
	valid = true
	globalIDLength := len(c.Globalid)
	valid = valid && (globalIDLength >= 3) && (globalIDLength <= 150)
	return
}
