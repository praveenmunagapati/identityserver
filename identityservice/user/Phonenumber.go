package user

import "regexp"

type Phonenumber string

var (
	phoneRegex = regexp.MustCompile(`^\+?[0-9]+$`)
)

//IsValidPhonenumber checks if a phonenumber is in a valid format
func IsValidPhonenumber(phonenumber Phonenumber) (valid bool) {
	valid = true
	valid = valid && len(phonenumber) < 51
	return phoneRegex.Match([]byte(phonenumber))
}
