package user

import (
	"github.com/itsyouonline/identityserver/db"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

type User struct {
	ID          bson.ObjectId          `json:"-" bson:"_id,omitempty"`
	Address     map[string]Address     `json:"address"`
	Bank        map[string]BankAccount `json:"bank"`
	Email       map[string]string      `json:"email"`
	Expire      db.Date                `json:"expire"`
	Facebook    FacebookAccount        `json:"facebook"`
	Github      GithubAccount          `json:"github"`
	Phone       map[string]Phonenumber `json:"phone"`
	PublicKeys  []string               `json:"publicKeys"`
	Username    string                 `json:"username"`
	TwoFAMethod string                 `json:"twofamethod"`
	Firstname   string                 `json:"firstname"`
	Lastname    string                 `json:"lastname"`
}

func ValidateUsername(username string) (valid bool) {
	regex, _ := regexp.Compile(`^[a-zA-Z0-9\s-_]+$`)
	matches := regex.FindAllString(username, 2)
	return len(matches) == 1
}
