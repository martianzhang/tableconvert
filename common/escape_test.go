package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLValueEscape(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		// Case when input is "NULL"
		{"NULL", "NULL"},
		// Case for a normal string
		{"test", "'test'"},
		// Case for a string containing a single quote
		{"it's a test", "'it\\'s a test'"},
		// Case for a string containing a double quote
		{"a \"quote\" here", "'a \"quote\" here'"},
		// Case for an empty string
		{"", "''"},
	}

	for _, tc := range testCases {
		result := SQLValueEscape(tc.input)
		assert.Equal(t, tc.expected, result, "Input: %s", tc.input)
	}
}
