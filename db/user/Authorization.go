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
	Name          bool              `json:"name"`
}

//FilterAuthorizedScopes filters the requested scopes to the ones this Authorization covers
func (authorization Authorization) FilterAuthorizedScopes(requestedscopes []string) (authorizedScopes []string) {
	authorizedScopes = make([]string, 0, len(requestedscopes))
	for _, rawscope := range requestedscopes {
		scope := strings.TrimSpace(rawscope)
		if scope == "user:name" && authorization.Name {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:memberof:") {
			requestedorgid := strings.TrimPrefix(scope, "user:memberof:")
			if authorization.containsOrganization(requestedorgid) {
				authorizedScopes = append(authorizedScopes, scope)
			}
		}
		if scope == "user:github" && authorization.Github {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if scope == "user:facebook" && authorization.Facebook {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:address") && labelledPropertyIsAuthorized(scope, "user:address", authorization.Address) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:bankaccount") && labelledPropertyIsAuthorized(scope, "user:bankaccount", authorization.Bank) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:email") && labelledPropertyIsAuthorized(scope, "user:email", authorization.Email) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:phone") && labelledPropertyIsAuthorized(scope, "user:phone", authorization.Phone) {
			authorizedScopes = append(authorizedScopes, scope)
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
