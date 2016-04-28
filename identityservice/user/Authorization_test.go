package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopesAreAuthorized(t *testing.T) {
	type testcase struct {
		authorization Authorization
		scopes        string
		authorized    bool
	}
	testcases := []testcase{
		testcase{authorization: Authorization{}, scopes: "user:memberof:orgid1", authorized: false},
		testcase{authorization: Authorization{Organizations: []string{"orgid"}}, scopes: "user:memberof:orgid", authorized: true},
		testcase{authorization: Authorization{Organizations: []string{"orgid.suborg"}}, scopes: "user:memberof:orgid.suborg", authorized: true},
		testcase{authorization: Authorization{Organizations: []string{"orgid1", "orgid2"}}, scopes: "user:memberof:orgid1, user:memberof:orgid2", authorized: true},
		testcase{authorization: Authorization{Organizations: []string{"orgid1", "orgid3"}}, scopes: "user:memberof:orgid1, user:memberof:orgid2", authorized: false},
	}
	for _, test := range testcases {
		assert.Equal(t, test.authorized, test.authorization.ScopesAreAuthorized(test.scopes), test.scopes)
	}
}
