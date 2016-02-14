package models

import (
	"gopkg.in/mgo.v2"

	"github.com/itsyouonline/identityserver/db"
)

// Initialize models in DB, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"globalid"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(COLLECTION_COMPANIES, index)
}
