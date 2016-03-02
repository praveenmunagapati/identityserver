package user

type Address struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	Nr         string `json:"nr"`
	Other      string `json:"other"`
	Postalcode string `json:"postalcode"`
	Street     string `json:"street"`
}
