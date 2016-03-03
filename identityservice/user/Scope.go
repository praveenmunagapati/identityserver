package user

// For an explanation about scopes and scopemapping, see https://github.com/itsyouonline/identityserver/blob/master/docs/oauth2/scopes.md
type Scope struct {
	Address       map[string]string `json:"address"`
	Bank          map[string]string `json:"bank"`
	Email         map[string]string `json:"email"`
	Facebook      bool              `json:"facebook"`
	Github        bool              `json:"github"`
	GrantedTo     string            `json:"grantedTo"`
	Organizations []string          `json:"organizations"`
	Phone         map[string]string `json:"phone"`
	PublicKeys    []string          `json:"publicKeys"`
	Username      string            `json:"username"`
}
