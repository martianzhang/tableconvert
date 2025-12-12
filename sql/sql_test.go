package sql

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

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

func TestMarshal(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create config with replace=true
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"replace": "",
			"table":   "test_table",
		},
	}

	// Create a valid table with one header and one row
	table := &common.Table{
		Headers: []string{"id"},
		Rows: [][]string{
			{"1"},
		},
	}

	// Call the function
	err := Marshal(cfg, table)
	assert.Nil(t, err)

	// Verify the output contains REPLACE instead of INSERT
	output := buf.String()
	assert.Contains(t, output, "REPLACE INTO")
	assert.NotContains(t, output, "INSERT INTO")
}

func TestUnmarshalMultipleValidInsertStatements(t *testing.T) {
	// Create a config with Reader containing multiple INSERT statements
	sqlContent := "INSERT INTO t1 VALUES (1); INSERT INTO t2 VALUES (2)"
	cfg := &common.Config{
		Reader: bytes.NewBufferString(sqlContent),
	}

	// Create an empty table
	table := &common.Table{
		Headers: []string{},
		Rows:    [][]string{},
	}

	// Call the Unmarshal function
	err := Unmarshal(cfg, table)

	// Verify no error occurred
	assert.Nil(t, err)

	// Verify the table was updated with data from both INSERT statements
	// Note: The exact handling of multiple INSERTs depends on the handleInsert implementation
	// which isn't shown in the provided code. This test assumes it adds rows to the table.
	assert.True(t, len(table.Rows) >= 2, "Table should contain at least 2 rows from the INSERT statements")
}

// TestUnmarshalInsertWithEscapedQuotes tests that escaped quotes are handled correctly
func TestUnmarshalInsertWithEscapedQuotes(t *testing.T) {
	input := "INSERT INTO table1 (name) VALUES ('Don\\'t worry');\n"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"name"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, "Don't worry", table.Rows[0][0]) // Should preserve escaped quotes
}

// TestUnmarshalInsertColumnMismatch tests error handling for column mismatches
func TestUnmarshalInsertColumnMismatch(t *testing.T) {
	input := "INSERT INTO table1 (id, name) VALUES (1, 'a');\nINSERT INTO table1 (name, id) VALUES ('b', 2);\n"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.Error(t, err) // Should return error for column order mismatch
	assert.Contains(t, err.Error(), "column order mismatch")
}
