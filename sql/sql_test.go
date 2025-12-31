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
	sqlContent := "INSERT INTO t1 (id) VALUES (1); INSERT INTO t1 (id) VALUES (2)"
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
	assert.Equal(t, []string{"id"}, table.Headers)
	assert.Equal(t, 2, len(table.Rows))
	assert.Equal(t, [][]string{{"1"}, {"2"}}, table.Rows)
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

// TestUnmarshalNilTable tests error when table is nil
func TestUnmarshalNilTable(t *testing.T) {
	cfg := &common.Config{
		Reader: strings.NewReader("INSERT INTO t VALUES (1);"),
	}

	err := Unmarshal(cfg, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target table pointer cannot be nil")
}

// TestUnmarshalNilReader tests error when reader is nil
func TestUnmarshalNilReader(t *testing.T) {
	table := &common.Table{}
	cfg := &common.Config{
		Reader: nil,
	}

	err := Unmarshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Reader in Config cannot be nil")
}

// TestUnmarshalInvalidSQL tests error handling for malformed SQL
func TestUnmarshalInvalidSQL(t *testing.T) {
	input := "INVALID SQL STATEMENT"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse SQL")
}

// TestUnmarshalUnsupportedStatement tests error for unsupported SQL types
func TestUnmarshalUnsupportedStatement(t *testing.T) {
	// SELECT is not supported
	input := "SELECT * FROM table1"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported SQL statement type")
}

// TestUnmarshalInsertWithNullValues tests NULL value handling
func TestUnmarshalInsertWithNullValues(t *testing.T) {
	input := "INSERT INTO table1 (id, name) VALUES (1, NULL);"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"id", "name"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, "1", table.Rows[0][0])
	assert.Equal(t, "NULL", table.Rows[0][1])
}

// TestUnmarshalInsertWithNumericLiterals tests different numeric literal formats
func TestUnmarshalInsertWithNumericLiterals(t *testing.T) {
	input := "INSERT INTO table1 (int_val, float_val, neg_val) VALUES (123, 3.14, -456);"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"int_val", "float_val", "neg_val"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, "123", table.Rows[0][0])
	assert.Equal(t, "3.14", table.Rows[0][1])
	assert.Equal(t, "-456", table.Rows[0][2])
}

// TestUnmarshalInsertWithMultipleRows tests multiple rows in single INSERT
func TestUnmarshalInsertWithMultipleRows(t *testing.T) {
	input := "INSERT INTO table1 (id, name) VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Charlie');"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"id", "name"}, table.Headers)
	assert.Equal(t, 3, len(table.Rows))
	assert.Equal(t, []string{"1", "Alice"}, table.Rows[0])
	assert.Equal(t, []string{"2", "Bob"}, table.Rows[1])
	assert.Equal(t, []string{"3", "Charlie"}, table.Rows[2])
}

// TestUnmarshalInsertWithBacktickIdentifiers tests backtick-quoted identifiers
func TestUnmarshalInsertWithBacktickIdentifiers(t *testing.T) {
	input := "INSERT INTO `table-name` (`column-name`, `another_col`) VALUES (1, 'test');"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"column-name", "another_col"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
}

// TestUnmarshalInsertColumnCountMismatch tests column count validation
func TestUnmarshalInsertColumnCountMismatch(t *testing.T) {
	input := "INSERT INTO table1 (id, name) VALUES (1, 'a', 'extra');"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	// The parser might handle this differently, but we should get some error
	// or the values should be handled appropriately
	assert.Error(t, err)
}

// TestMarshalNilTable tests error when marshaling nil table
func TestMarshalNilTable(t *testing.T) {
	cfg := &common.Config{
		Writer: &bytes.Buffer{},
	}

	err := Marshal(cfg, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input table pointer cannot be nil")
}

// TestMarshalEmptyHeaders tests error when table has no headers
func TestMarshalEmptyHeaders(t *testing.T) {
	table := &common.Table{
		Headers: []string{},
		Rows:    [][]string{{"1"}},
	}
	cfg := &common.Config{
		Writer: &bytes.Buffer{},
	}

	err := Marshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "table must have at least one header")
}

// TestMarshalInconsistentRowLength tests error when rows have inconsistent column counts
func TestMarshalInconsistentRowLength(t *testing.T) {
	table := &common.Table{
		Headers: []string{"id", "name"},
		Rows: [][]string{
			{"1", "Alice"},
			{"2"}, // Missing column
		},
	}
	cfg := &common.Config{
		Writer: &bytes.Buffer{},
	}

	err := Marshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "row has 1 columns, but table has 2")
}

// TestMarshalWithDifferentDialects tests different SQL dialects
func TestMarshalWithDifferentDialects(t *testing.T) {
	tests := []struct {
		name     string
		dialect  string
		expected string
	}{
		{
			name:     "mysql",
			dialect:  "mysql",
			expected: "`id`",
		},
		{
			name:     "oracle",
			dialect:  "oracle",
			expected: `"id"`,
		},
		{
			name:     "postgres",
			dialect:  "postgres",
			expected: `"id"`,
		},
		{
			name:     "mssql",
			dialect:  "mssql",
			expected: `[id]`,
		},
		{
			name:     "none",
			dialect:  "none",
			expected: `id`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &common.Table{
				Headers: []string{"id"},
				Rows:    [][]string{{"1"}},
			}
			var buf bytes.Buffer
			cfg := &common.Config{
				Writer: &buf,
				Extension: map[string]string{
					"table":   "test_table",
					"dialect": tt.dialect,
				},
			}

			err := Marshal(cfg, table)
			assert.NoError(t, err)

			output := buf.String()
			assert.Contains(t, output, tt.expected)
		})
	}
}

// TestMarshalWithOneInsert tests all-in-one INSERT mode
func TestMarshalWithOneInsert(t *testing.T) {
	table := &common.Table{
		Headers: []string{"id", "name"},
		Rows: [][]string{
			{"1", "Alice"},
			{"2", "Bob"},
		},
	}
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"table":      "test_table",
			"one-insert": "true",
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Should have one INSERT with multiple VALUES
	assert.Contains(t, output, "INSERT INTO")
	assert.Contains(t, output, "VALUES")
	// Note: SQLValueEscape quotes numbers and uses backslash escaping
	assert.Contains(t, output, "('1', 'Alice')")
	assert.Contains(t, output, "('2', 'Bob')")
	// Should end with semicolon
	assert.Contains(t, output, ";\n")
}

// TestMarshalWithSpecialCharacters tests escaping of special characters
func TestMarshalWithSpecialCharacters(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name", "description"},
		Rows: [][]string{
			{"Alice", "She's happy"},
			{"Bob", "Has 'quotes' in text"},
		},
	}
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"table": "test_table",
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Note: SQLValueEscape uses backslash escaping
	assert.Contains(t, output, "'She\\'s happy'")
	assert.Contains(t, output, "'Has \\'quotes\\' in text'")
}

// TestMarshalWithEmptyString tests empty string handling
func TestMarshalWithEmptyString(t *testing.T) {
	table := &common.Table{
		Headers: []string{"id", "name"},
		Rows: [][]string{
			{"1", ""},
		},
	}
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"table": "test_table",
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "''")
}

// TestEscapeIdentifierDialects tests escapeIdentifier with all dialects
func TestEscapeIdentifierDialects(t *testing.T) {
	tests := []struct {
		name     string
		dialect  string
		input    string
		expected string
	}{
		{"mysql", "mysql", "col-name", "`col-name`"},
		{"oracle", "oracle", "col-name", `"col-name"`},
		{"postgres", "postgres", "col-name", `"col-name"`},
		{"mssql", "mssql", "col-name", `[col-name]`},
		{"none", "none", "col-name", "col-name"},
		{"default", "", "col-name", "`col-name`"}, // default to mysql
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeIdentifier(tt.input, tt.dialect)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestUnmarshalInsertWithQuotedIdentifiers tests quoted identifiers in INSERT
func TestUnmarshalInsertWithQuotedIdentifiers(t *testing.T) {
	input := "INSERT INTO `my-table` (`my-col`, `another-col`) VALUES (1, 'test');"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"my-col", "another-col"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
}

// TestUnmarshalInsertWithComplexLiterals tests various literal types
func TestUnmarshalInsertWithComplexLiterals(t *testing.T) {
	input := "INSERT INTO table1 (str, num, bool, null_val) VALUES ('hello', 42, true, NULL);"
	cfg := &common.Config{
		Reader: strings.NewReader(input),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"str", "num", "bool", "null_val"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, "hello", table.Rows[0][0])
	assert.Equal(t, "42", table.Rows[0][1])
	assert.Equal(t, "true", table.Rows[0][2])
	assert.Equal(t, "NULL", table.Rows[0][3])
}
