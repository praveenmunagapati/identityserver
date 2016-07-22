package oauthservice

import "strings"

func splitScopeString(scopestring string) (scopeList []string) {
	scopeList = []string{}
	for _, value := range strings.Split(scopestring, ",") {
		scope := strings.TrimSpace(value)
		if scope != "" {
			scopeList = append(scopeList, scope)
		}
	}
	return
}
