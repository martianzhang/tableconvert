package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"vitess.io/vitess/go/vt/sqlparser"
)

func TestParser(t *testing.T) {
	sqls := []string{
		"INSERT INTO `table_name` (`column1`, `column2`) VALUES (1, 'test');",
	}

	parser, err := sqlparser.New(sqlparser.Options{})
	assert.Nil(t, err)

	for _, sql := range sqls {
		_, err = parser.Parse(sql)
		assert.Nil(t, err)
	}
}

func TestEscapeIdentifier(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"normal_identifier", "`normal_identifier`"},
		{"`identifier_with_backtick`", "`\\`identifier_with_backtick\\``"},
		{"", "``"},
	}

	for _, tc := range testCases {
		result := escapeIdentifier(tc.input, "mysql")
		assert.Equal(t, tc.expected, result)
	}
}

func TestEscapeValue(t *testing.T) {
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
		result := escapeValue(tc.input)
		assert.Equal(t, tc.expected, result, "Input: %s", tc.input)
	}
}
