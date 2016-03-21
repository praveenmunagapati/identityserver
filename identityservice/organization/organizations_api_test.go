package organization

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPISecretLabelValidation(t *testing.T) {
	type testcase struct {
		label string
		valid bool
	}
	testcases := []testcase{
		testcase{label: "", valid: false},
		testcase{label: "ab", valid: false},
		testcase{label: "abc", valid: true},
		testcase{label: strings.Repeat("1", 50), valid: true},
		testcase{label: strings.Repeat("1", 51), valid: false},
	}
	for _, test := range testcases {
		assert.Equal(t, test.valid, isValidAPISecretLabel(test.label), test.label)
	}
}
