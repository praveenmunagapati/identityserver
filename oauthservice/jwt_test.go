package oauthservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTScopesAreAllowed(t *testing.T) {
	type testcase struct {
		allowed   string
		requested string
		valid     bool
	}
	testcases := []testcase{
		testcase{allowed: "", requested: "", valid: true},
		testcase{allowed: "", requested: "user:memberof:org1", valid: false},
		testcase{allowed: "user:memberof:org1", requested: "", valid: true},
		testcase{allowed: "user:memberof:org2", requested: "user:memberof:org1", valid: false},
		testcase{allowed: "user:memberof:org1, user:memberof:org2", requested: "user:memberof:org1", valid: true},
		testcase{allowed: "user:memberof:org1", requested: "user:memberof:org1, user:memberof:org2", valid: false},
		testcase{allowed: "user:admin", requested: "user:memberof:org1", valid: true},
	}
	for _, test := range testcases {
		valid := jwtScopesAreAllowed(splitScopeString(test.allowed), splitScopeString(test.requested))
		assert.Equal(t, test.valid, valid, "Allowed: \"%s\" - Requested: \"%s\"", test.allowed, test.requested)
	}
}
