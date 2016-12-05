package user

import (
	"strings"
	"testing"

	"github.com/itsyouonline/identityserver/db/user"
	"github.com/stretchr/testify/assert"
)

func TestLabelValidation(t *testing.T) {
	type testcase struct {
		label string
		valid bool
	}
	testcases := []testcase{
		{label: "", valid: false},
		{label: "a", valid: false},
		{label: "ab", valid: true},
		{label: "abc", valid: true},
		{label: "abc- _", valid: true},
		{label: "abc%", valid: false},
		{label: strings.Repeat("1", 50), valid: true},
		{label: strings.Repeat("1", 51), valid: false},
	}
	for _, test := range testcases {
		assert.Equal(t, test.valid, isValidLabel(test.label), test.label)
	}
}

func TestUsernameValidation(t *testing.T) {
	type testcase struct {
		label string
		valid bool
	}
	testcases := []testcase{
		{label: "", valid: false},
		{label: "a", valid: false},
		{label: "ab", valid: true},
		{label: "abc", valid: true},
		{label: "ABC", valid: false},
		{label: "abc- _", valid: true},
		{label: "abb%", valid: false},
		{label: strings.Repeat("1", 30), valid: true},
		{label: strings.Repeat("1", 31), valid: false},
	}
	for _, test := range testcases {
		assert.Equal(t, test.valid, user.ValidateUsername(test.label), test.label)
	}
}
