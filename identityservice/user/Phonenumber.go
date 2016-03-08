package user

import "regexp"

type Phonenumber string

var (
	phoneRegex = regexp.MustCompile(`^\+[0-9]+$`)
)

func IsValidPhonenumber(phonenumber Phonenumber) bool {
	return phoneRegex.Match([]byte(phonenumber))
}
