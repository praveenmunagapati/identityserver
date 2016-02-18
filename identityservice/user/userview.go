package user

type Userview struct {
	Address  []Address     `json:"address"`
	Email    []string      `json:"email"`
	Phone    []Phonenumber `json:"phone"`
	Username string        `json:"username"`
}
