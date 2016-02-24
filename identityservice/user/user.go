package user

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id       bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	Address  map[string]Address     `json:"address"`
	Email    map[string]string      `json:"email"`
	Expire   Date                   `json:"expire"`
	Phone    map[string]Phonenumber `json:"phone"`
	Username string                 `json:"username"`
}
