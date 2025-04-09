package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		input    *ParseError
		expected string
	}{
		{
			name: "basic error",
			input: &ParseError{
				LineNumber: 1,
				Message:    "invalid format",
				Line:       "header1,header2",
			},
			expected: `parse error on line 1: invalid format (line: "header1,header2")`,
		},
		{
			name: "empty line",
			input: &ParseError{
				LineNumber: 5,
				Message:    "empty line",
				Line:       "",
			},
			expected: `parse error on line 5: empty line (line: "")`,
		},
		{
			name: "special characters in line",
			input: &ParseError{
				LineNumber: 10,
				Message:    "invalid character",
				Line:       "data\twith\ttabs",
			},
			expected: `parse error on line 10: invalid character (line: "data\twith\ttabs")`,
		},
		{
			name: "large line number",
			input: &ParseError{
				LineNumber: 99999,
				Message:    "file too large",
				Line:       "very long line...",
			},
			expected: `parse error on line 99999: file too large (line: "very long line...")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.input.Error()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
