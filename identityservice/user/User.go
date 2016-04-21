package user

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id         bson.ObjectId          `json:"-" bson:"_id,omitempty"`
	Address    map[string]Address     `json:"address"`
	Bank       map[string]BankAccount `json:"bank"`
	Email      map[string]string      `json:"email"`
	Expire     Date                   `json:"expire"`
	Facebook   FacebookAccount        `json:"facebook"`
	Github     GithubAccount          `json:"github"`
	Phone      map[string]Phonenumber `json:"phone"`
	PublicKeys []string               `json:"publicKeys"`
	Username   string                 `json:"username"`
}
