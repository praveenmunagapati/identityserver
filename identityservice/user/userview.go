package user

type userview struct {
	Address       map[string]Address     `json:"address"`
	Bank          map[string]BankAccount `json:"bank"`
	Email         map[string]string      `json:"email"`
	Facebook      string                 `json:"facebook"`
	Github        string                 `json:"github"`
	Organizations []string               `json:"organizations"`
	Phone         map[string]Phonenumber `json:"phone"`
	PublicKeys    []string               `json:"publicKeys"`
	Username      string                 `json:"username"`
}
