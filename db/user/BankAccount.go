package user

type BankAccount struct {
	Bic     string `json:"bic"`
	Country string `json:"country"`
	Iban    string `json:"iban"`
	Label   string `json:"label"`
}
