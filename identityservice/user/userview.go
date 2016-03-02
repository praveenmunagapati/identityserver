package user

type userview struct {
	Address  []Address     `json:"address"`
	Bank     []BankAccount `json:"bank"`
	Email    []string      `json:"email"`
	Facebook string        `json:"facebook"`
	Github   string        `json:"github"`
	Phone    []Phonenumber `json:"phone"`
	Username string        `json:"username"`
}
