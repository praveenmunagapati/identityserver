package user

import (
	"errors"
	"regexp"

	"github.com/itsyouonline/identityserver/db"
	"gopkg.in/mgo.v2/bson"
)

type EmailAddress struct {
	EmailAddress string
	Label        string
}

type User struct {
	ID             bson.ObjectId          `json:"-" bson:"_id,omitempty"`
	Address        map[string]Address     `json:"address"`
	Bank           map[string]BankAccount `json:"bank"`
	EmailAddresses []EmailAddress         `json:"emailaddresses"`
	Expire         db.Date                `json:"expire"`
	Facebook       FacebookAccount        `json:"facebook"`
	Github         GithubAccount          `json:"github"`
	Phone          map[string]Phonenumber `json:"phone"`
	PublicKeys     []string               `json:"publicKeys"`
	Username       string                 `json:"username"`
	TwoFAMethod    string                 `json:"twofamethod"`
	Firstname      string                 `json:"firstname"`
	Lastname       string                 `json:"lastname"`
}

func (u *User) GetEmailAddressByLabel(label string) (email EmailAddress, err error) {
	for _, email = range u.EmailAddresses {
		if email.Label == label {
			return
		}
	}
	err = errors.New("Could not find EmailAddress with Label")
	return
}

func ValidateUsername(username string) (valid bool) {
	regex, _ := regexp.Compile(`^[a-zA-Z0-9\s-_]+$`)
	matches := regex.FindAllString(username, 2)
	return len(matches) == 1
}
