package models

import (
	"gopkg.in/mgo.v2"
)

// Initialize models in DB, if required.
func InitModels(session *mgo.Session) error {
	// TODO: Use model tags to ensure indices/constraints.
	c := getUserCollection(session)

	err := c.EnsureIndex(mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	})

	return err
}
