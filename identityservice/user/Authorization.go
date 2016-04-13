package user

// For an explanation about scopes and scopemapping, see https://github.com/itsyouonline/identityserver/blob/master/docs/oauth2/scopes.md
type Authorization struct {
	Address       string   `json:"address,omitempty"`
	Bank          string   `json:"bank,omitempty"`
	Email         string   `json:"email,omitempty"`
	Facebook      bool     `json:"facebook,omitempty"`
	Github        bool     `json:"github,omitempty"`
	GrantedTo     string   `json:"grantedTo"`
	Organizations []string `json:"organizations"`
	Phone         string   `json:"phone,omitempty"`
	PublicKeys    []string `json:"publicKeys,omitempty"`
	Username      string   `json:"username"`
}
