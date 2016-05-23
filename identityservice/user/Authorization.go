package user

import "strings"

// Authorization defines what userinformation is authorized to be seen by an organization
// For an explanation about scopes and scopemapping, see https://github.com/itsyouonline/identityserver/blob/master/docs/oauth2/scopes.md
type Authorization struct {
	Address       map[string]string `json:"address,omitempty"`
	Bank          map[string]string `json:"bank,omitempty"`
	Email         map[string]string `json:"email,omitempty"`
	Facebook      bool              `json:"facebook,omitempty"`
	Github        bool              `json:"github,omitempty"`
	GrantedTo     string            `json:"grantedTo"`
	Organizations []string          `json:"organizations"`
	Phone         map[string]string `json:"phone,omitempty"`
	PublicKeys    []string          `json:"publicKeys,omitempty"`
	Username      string            `json:"username"`
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
		if scope == "user:github" {
			authorized = authorized && authorization.Github
		}
		if scope == "user:facebook" {
			authorized = authorized && authorization.Facebook
		}
		if strings.HasPrefix(scope, "user:address") {
			authorized = authorized && labelledPropertyIsAuthorized(scope, "user:address", authorization.Address)
		}
		if strings.HasPrefix(scope, "user:bankaccount") {
			authorized = authorized && labelledPropertyIsAuthorized(scope, "user:bankaccount", authorization.Bank)
		}
		if strings.HasPrefix(scope, "user:email") {
			authorized = authorized && labelledPropertyIsAuthorized(scope, "user:email", authorization.Email)
		}
		if strings.HasPrefix(scope, "user:phone") {
			authorized = authorized && labelledPropertyIsAuthorized(scope, "user:phone", authorization.Phone)
		}
		if !authorized {
			return
		}
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

func labelledPropertyIsAuthorized(scope string, scopePrefix string, authorizedLabels map[string]string) (authorized bool) {
	if authorizedLabels == nil {
		return
	}
	if scope == scopePrefix {
		authorized = len(authorizedLabels) > 0
		return
	}
	if strings.HasPrefix(scope, scopePrefix+":") {
		requestedLabel := strings.TrimPrefix(scope, scopePrefix+":")
		_, authorized = authorizedLabels[requestedLabel]
	}
	return
}
