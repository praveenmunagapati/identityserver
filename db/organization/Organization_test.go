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
	owners := []string{"testowner"}
	testCases := []testcase{
		{org: &Organization{Owners: owners, Globalid: ""}, valid: false},
		{org: &Organization{Owners: owners, Globalid: "ab"}, valid: false},
		{org: &Organization{Owners: owners, Globalid: "a♥"}, valid: false},
		{org: &Organization{Owners: owners, Globalid: "abc"}, valid: true},
		{org: &Organization{Owners: owners, Globalid: strings.Repeat("1", 150)}, valid: true},
		{org: &Organization{Owners: owners, Globalid: strings.Repeat("1", 151)}, valid: false},
		{org: &Organization{Owners: []string{}, Globalid: "abc"}, valid: false},
	}
	for _, test := range testCases {
		assert.Equal(t, test.valid, test.org.IsValid(), test.org.Globalid)
	}
}
