package contract

type Signature struct {
	Date      Date   `json:"date"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	SignedBy  string `json:"signedBy"`
}
