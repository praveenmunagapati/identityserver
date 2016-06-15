package user

import (
	"errors"
	"regexp"

	"github.com/itsyouonline/identityserver/db"
	"gopkg.in/mgo.v2/bson"
)

type EmailAddress struct {
	EmailAddress string `json:"emailaddress"`
	Label        string `json:"label"`
}

type PublicKey struct {
	PublicKey string `json:"publickey"`
	Label     string `json:"label"`
}

type User struct {
	ID             bson.ObjectId   `json:"-" bson:"_id,omitempty"`
	Addresses      []Address       `json:"addresses"`
	BankAccounts   []BankAccount   `json:"bankaccounts"`
	EmailAddresses []EmailAddress  `json:"emailaddresses"`
	Expire         db.Date         `json:"expire"`
	Facebook       FacebookAccount `json:"facebook"`
	Github         GithubAccount   `json:"github"`
	Phonenumbers   []Phonenumber   `json:"phonenumbers"`
	PublicKeys     []PublicKey     `json:"publicKeys"`
	Username       string          `json:"username"`
	Firstname      string          `json:"firstname"`
	Lastname       string          `json:"lastname"`
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

func (u *User) GetPhonenumberByLabel(label string) (phonenumber Phonenumber, err error) {
	for _, phonenumber = range u.Phonenumbers {
		if phonenumber.Label == label {
			return
		}
	}
	err = errors.New("Could not find Phonenumber with Label")
	return
}

func (u *User) GetBankAccountByLabel(label string) (bankaccount BankAccount, err error) {
	for _, bankaccount = range u.BankAccounts {
		if bankaccount.Label == label {
			return
		}
	}
	err = errors.New("Could not find Phonenumber with Label")
	return
}

func (u *User) GetAddressByLabel(label string) (address Address, err error) {
	for _, address = range u.Addresses {
		if address.Label == label {
			return
		}
	}
	err = errors.New("Could not find Phonenumber with Label")
	return
}

func ValidateUsername(username string) (valid bool) {
	regex, _ := regexp.Compile(`^[a-z0-9\s-_]+$`)
	matches := regex.FindAllString(username, 2)
	return len(matches) == 1
}
