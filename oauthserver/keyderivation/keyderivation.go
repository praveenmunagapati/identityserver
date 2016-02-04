package keyderivation

import "github.com/itsyouonline/identityserver/oauthserver/keyderivation/crypt/sha512crypt"

//Hash creates a random 16 character salt and creates a key using this salt
//Key generation function: SHA512 with 5000 iterations
func Hash(password string) (key string, err error) {
	c := sha512crypt.New()
	key, err = c.Generate([]byte(password), []byte(""))

	return
}

//Check takes the password and the encoded key
// and checks if the combination matches
func Check(password, key string) bool {
	err := sha512crypt.New().Verify(key, []byte(password))
	return err == nil
}
