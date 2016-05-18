package oauthservice

import (
	"strconv"
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
		testcase{allowed: "user:memberOf:org1", requested: "", valid: true},
		testcase{allowed: "user:memberOf:org2", requested: "user:memberOf:org1", valid: false},
		testcase{allowed: "user:memberOf:org1, user:memberOf:org2", requested: "user:memberOf:org1", valid: true},
		testcase{allowed: "user:memberOf:org1", requested: "user:memberOf:org1, user:memberOf:org2", valid: false},
	}
	for i, test := range testcases {
		valid := jwtScopesAreAllowed(test.allowed, test.requested)
		assert.Equal(t, test.valid, valid, strconv.Itoa(i))
	}
}
