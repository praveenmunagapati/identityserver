package user

import "strings"

// Authorization defines what userinformation is authorized to be seen by an organization
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

//ScopesAreAuthorized checks if this Authorization covers all the requested scopes
func (authorization Authorization) ScopesAreAuthorized(scopes string) (authorized bool) {
	authorized = true
	for _, rawscope := range strings.Split(scopes, ",") {
		scope := strings.TrimSpace(rawscope)
		if strings.HasPrefix(scope, "user:memberof:") {
			requestedorgid := strings.TrimPrefix(scope, "user:memberof:")
			authorized = authorized && authorization.containsOrganization(requestedorgid)
		}
		//TODO: authorization of other properties besides organization
	}

	return
}

func (authorization Authorization) containsOrganization(globalid string) bool {
	for _, orgid := range authorization.Organizations {
		if orgid == globalid {
			return true
		}
	}
	return false
}
