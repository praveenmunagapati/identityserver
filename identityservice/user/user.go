package user

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Address  string        `json:"address"`
	Email    string        `json:"email"`
	Expire   Date          `json:"expire"`
	Phone    string        `json:"phone"`
	Username string        `json:"username"`
}
