package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectTableFormatByExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic extension",
			input:    ".csv",
			expected: TableFormatCSV,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := DetectTableFormatByExtension(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseError(t *testing.T) {
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

func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "uppercase null",
			input:    "NULL",
			expected: nil,
		},
		{
			name:     "lowercase null",
			input:    "null",
			expected: nil,
		},
		{
			name:     "mixed case null",
			input:    "Null",
			expected: nil,
		},
		{
			name:     "random case null",
			input:    "nUlL",
			expected: nil,
		},
		{
			name:     "true boolean",
			input:    "true",
			expected: true,
		},
		{
			name:     "false boolean",
			input:    "false",
			expected: false,
		},
		{
			name:     "mixed case true",
			input:    "TrUe",
			expected: true,
		},
		{
			name:     "mixed case false",
			input:    "FaLsE",
			expected: false,
		},
		{
			name:     "positive integer",
			input:    "123",
			expected: int64(123),
		},
		{
			name:     "negative integer",
			input:    "-456",
			expected: int64(-456),
		},
		{
			name:     "zero",
			input:    "0",
			expected: int64(0),
		},
		{
			name:     "positive float",
			input:    "123.456",
			expected: 123.456,
		},
		{
			name:     "negative float",
			input:    "-789.123",
			expected: -789.123,
		},
		{
			name:     "scientific notation",
			input:    "1.23e4",
			expected: 12300.0,
		},
		{
			name:     "float with leading zero",
			input:    "0.123",
			expected: 0.123,
		},
		{
			name:     "integer as float",
			input:    "123.0",
			expected: 123.0,
		},
		{
			name:     "invalid number",
			input:    "123abc",
			expected: "123abc",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace string",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "regular string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with numbers",
			input:    "123 main st",
			expected: "123 main st",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := InferType(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInferPrintType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase null",
			input:    "NULL",
			expected: "null",
		},
		{
			name:     "lowercase null",
			input:    "null",
			expected: "null",
		},
		{
			name:     "mixed case null",
			input:    "Null",
			expected: "null",
		},
		{
			name:     "random case null",
			input:    "nUlL",
			expected: "null",
		},
		{
			name:     "true boolean",
			input:    "true",
			expected: "true",
		},
		{
			name:     "false boolean",
			input:    "false",
			expected: "false",
		},
		{
			name:     "mixed case true",
			input:    "TrUe",
			expected: "true",
		},
		{
			name:     "mixed case false",
			input:    "FaLsE",
			expected: "false",
		},
		{
			name:     "positive integer",
			input:    "123",
			expected: "123",
		},
		{
			name:     "negative integer",
			input:    "-456",
			expected: "-456",
		},
		{
			name:     "zero",
			input:    "0",
			expected: "0",
		},
		{
			name:     "positive float",
			input:    "123.456",
			expected: "123.456",
		},
		{
			name:     "negative float",
			input:    "-789.123",
			expected: "-789.123",
		},
		{
			name:     "scientific notation",
			input:    "1.23e4",
			expected: "1.23e4",
		},
		{
			name:     "float with leading zero",
			input:    "0.123",
			expected: "0.123",
		},
		{
			name:     "integer as float",
			input:    "123.0",
			expected: "123.0",
		},
		{
			name:     "invalid number",
			input:    "123abc",
			expected: "123abc",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace string",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "regular string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with numbers",
			input:    "123 main st",
			expected: "123 main st",
		},
		{
			name:     "string with leading/trailing spaces",
			input:    "  test  ",
			expected: "  test  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := InferPrintType(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
