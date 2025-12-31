package markdown

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	input := "```txt\n" + // Added code fence for realism
		"|   DATE   |         DESCRIPTION         | CV2  | AMOUNT |\n" +
		"|----------|--------------------------|------|--------|\n" +
		"| 1/1/2014 | Domain name              | 2233 | $10.98 |\n" +
		"| 1/1/2014 | January Hosting          | 2233 | $54.95 |\n" +
		"| 1/4/2014 | February Hosting         | 2233 | $51.00 |\n" +
		"| 1/4/2014 | February Extra Bandwidth | 2233 | $30.00 |\n" +
		"```" // Added closing fence

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	cfg.Reader = strings.NewReader(input)
	assert.Nil(t, err)

	var table common.Table
	err = Unmarshal(&cfg, &table)

	assert.Nil(t, err)
	assert.Equal(t, []string{"DATE", "DESCRIPTION", "CV2", "AMOUNT"}, table.Headers)
	assert.Equal(t, 4, len(table.Rows))
	assert.Equal(t, []string{"1/4/2014", "February Hosting", "2233", "$51.00"}, table.Rows[2])
}

func TestUmarshalEmptyCells(t *testing.T) {
	input := `
|FIELD|TYPE|NULL|KEY|DEFAULT|EXTRA|
|---|---|---|---|---|---|
|user_id|smallint(5)|NO|PRI|NULL|auto_increment|
|username|varchar(10)|NO||NULL||
|password|varchar(100)|NO||NULL||
`

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table

	err = Unmarshal(&cfg, &table)
	assert.Nil(t, err)
}

func TestUnmarshalPipesInContent(t *testing.T) {
	// Test escaped pipes in content
	input := "| A | B |\n|---|---|\n| x | y\\|z |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)

	assert.Nil(t, err)
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"x", "y|z"}, table.Rows[0])
}

func TestMarshalPipesInContent(t *testing.T) {
	// Test that pipes in data are properly escaped
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"x", "y|z"}},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.Nil(t, err)

	output := buf.String()
	// Should contain escaped pipe
	assert.Contains(t, output, "y\\|z")

	// Test round-trip
	cfg.Reader = strings.NewReader(output)
	var table2 common.Table
	err = Unmarshal(&cfg, &table2)
	assert.Nil(t, err)
	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)
}

func TestMarshalSpecialChars(t *testing.T) {
	// Test that various special characters are properly escaped
	table := &common.Table{
		Headers: []string{"A|B", "C*D", "E_F"},
		Rows:    [][]string{{"x|y", "z*w", "u_v"}},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.Nil(t, err)

	output := buf.String()
	// Should contain escaped special chars
	assert.Contains(t, output, "A\\|B")
	assert.Contains(t, output, "C\\*D")
	assert.Contains(t, output, "E\\_F")
	assert.Contains(t, output, "x\\|y")
	assert.Contains(t, output, "z\\*w")
	assert.Contains(t, output, "u\\_v")

	// Test round-trip
	cfg.Reader = strings.NewReader(output)
	var table2 common.Table
	err = Unmarshal(&cfg, &table2)
	assert.Nil(t, err)
	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name     string
		table    *common.Table
		err      error
		expected string
	}{
		{
			name:     "nil table",
			table:    nil,
			err:      errors.New("Marshal: input table pointer cannot be nil"), // Use errors.New
			expected: "",                                                       // Output buffer should be empty on error
		},
		{
			name: "empty headers",
			table: &common.Table{
				Headers: []string{},
				Rows:    [][]string{},
			},
			err:      errors.New("Marshal: table must have at least one header"), // Use errors.New
			expected: "",                                                         // Output buffer should be empty on error
		},
		{
			name: "column count mismatch",
			table: &common.Table{
				Headers: []string{"Header1", "Header2"},
				Rows: [][]string{
					{"Cell1"}, // Row 1 has 1 column, headers have 2
				},
			},
			// Adjusted error message to reflect 0-based row index if that's how Marshal reports it
			err:      errors.New("Marshal: 1 row has 1 columns, but table has 2"), // Use errors.New
			expected: "",                                                          // Output buffer should be empty on error
		},
		{
			name: "successful marshal",
			table: &common.Table{
				Headers: []string{"Header1", "Header2"},
				Rows: [][]string{
					{"Cell1", "Cell2"},
					{"Cell3", "Cell4"},
				},
			},
			err:      nil,                                                                                            // No error expected
			expected: "| Header1 | Header2 |\n|---------|---------|\n| Cell1   | Cell2   |\n| Cell3   | Cell4   |\n", // Pretty-printed expected output
		},
		{
			name: "no rows", // Added test case for table with headers but no rows
			table: &common.Table{
				Headers: []string{"ColA", "ColB"},
				Rows:    [][]string{},
			},
			err:      nil,
			expected: "| ColA | ColB |\n|------|------|\n", // Should output headers and separator line
		},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a bytes.Buffer for each test run to capture output
			var buf bytes.Buffer
			cfg.Writer = &buf

			// Call the Marshal function, passing the table and the buffer as the writer
			err := Marshal(&cfg, tt.table)

			// Check the error status
			if tt.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err.Error())
				// Optionally, assert buffer is empty on error if that's the contract
				// assert.Empty(t, buf.String(), "Output buffer should be empty when an error occurs")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, buf.String())
			}
		})
	}
}
