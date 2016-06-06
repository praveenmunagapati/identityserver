package organization

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationValidation(t *testing.T) {
	type testcase struct {
		org   *Organization
		valid bool
	}
	testcases := []testcase{
		testcase{org: &Organization{Globalid: ""}, valid: false},
		testcase{org: &Organization{Globalid: "ab"}, valid: false},
		//	testcase{org: &Organization{Globalid: "a♥"}, valid: false}, Let's just limit the amount of bytes for now
		testcase{org: &Organization{Globalid: "abc"}, valid: true},
		testcase{org: &Organization{Globalid: strings.Repeat("1", 150)}, valid: true},
		testcase{org: &Organization{Globalid: strings.Repeat("1", 151)}, valid: false},
	}
	for _, test := range testcases {
		assert.Equal(t, test.valid, test.org.IsValid(), test.org.Globalid)
	}
}
