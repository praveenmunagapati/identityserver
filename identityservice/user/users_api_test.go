package user

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelValidation(t *testing.T) {
	type testcase struct {
		label string
		valid bool
	}
	testcases := []testcase{
		testcase{label: "", valid: false},
		testcase{label: "a", valid: false},
		testcase{label: "ab", valid: true},
		testcase{label: "abc", valid: true},
		testcase{label: strings.Repeat("1", 50), valid: true},
		testcase{label: strings.Repeat("1", 51), valid: false},
	}
	for _, test := range testcases {
		assert.Equal(t, test.valid, isValidLabel(test.label), test.label)
	}
}
