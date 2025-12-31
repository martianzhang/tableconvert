package mediawiki

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	// Create a simple table markup
	tableMarkup := `{| class="wikitable"
! Header1 !! Header2
|-
| Row1Cell1 || Row1Cell2
|-
| Row2Cell1 || Row2Cell2
|}
`

	// Create a config with the table markup as reader
	cfg := &common.Config{
		Reader: bytes.NewBufferString(tableMarkup),
	}

	// Create an empty table to be populated
	table := &common.Table{}

	// Call the Unmarshal function
	err := Unmarshal(cfg, table)
	assert.NoError(t, err, "Unmarshal failed")

	// Verify the headers
	expectedHeaders := []string{"Header1", "Header2"}
	assert.Equal(t, expectedHeaders, table.Headers, "Headers are incorrect")

	// Verify the rows
	expectedRows := [][]string{
		{"Row1Cell1", "Row1Cell2"},
		{"Row2Cell1", "Row2Cell2"},
	}
	assert.Equal(t, expectedRows, table.Rows, "Rows are incorrect")
}

func TestUnmarshalPipesInContent(t *testing.T) {
	// Test that pipes in content are handled correctly
	input := `{| class="wikitable"
! Header1 !! Header2
|-
| Cell1|Cell2 || Cell3|Cell4
|}`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"Cell1|Cell2", "Cell3|Cell4"}, table.Rows[0])
}

func TestUnmarshalTemplate(t *testing.T) {
	// Test that {{!}} template is converted to pipes
	input := `{| class="wikitable"
! A{{!}}B !! C
|-
| x{{!}}y || z
|}`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"A|B", "C"}, table.Headers)
	assert.Equal(t, []string{"x|y", "z"}, table.Rows[0])
}

func TestUnmarshalColumnValidation(t *testing.T) {
	// Test that inconsistent column counts are rejected
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "row too many columns",
			input: `{| class="wikitable"
! A !! B
|-
| 1 || 2 || 3
|}`,
		},
		{
			name: "row too few columns",
			input: `{| class="wikitable"
! A !! B !! C
|-
| 1 || 2
|}`,
		},
		{
			name: "mixed column counts",
			input: `{| class="wikitable"
! A !! B
|-
| 1 || 2
|-
| 3 || 4 || 5
|}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &common.Config{
				Reader: bytes.NewBufferString(tt.input),
			}
			var table common.Table
			err := Unmarshal(cfg, &table)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "columns")
		})
	}
}

func TestUnmarshalEmptyCells(t *testing.T) {
	// Test empty cells
	input := `{| class="wikitable"
! A !! B !! C
|-
| 1 ||  || 3
|}`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B", "C"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"1", "", "3"}, table.Rows[0])
}

func TestUnmarshalUTF8(t *testing.T) {
	// Test UTF-8 characters
	input := `{| class="wikitable"
! 名字 !! 年龄
|-
| 小明 || 25
|}`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"名字", "年龄"}, table.Headers)
	assert.Equal(t, []string{"小明", "25"}, table.Rows[0])
}

func TestMarshalPipesInContent(t *testing.T) {
	// Test that pipes in content are escaped with {{!}}
	table := &common.Table{
		Headers: []string{"A|B", "C"},
		Rows:    [][]string{{"x|y", "z"}},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "A{{!}}B")
	assert.Contains(t, output, "x{{!}}y")

	// Test round-trip
	cfg.Reader = strings.NewReader(output)
	var table2 common.Table
	err = Unmarshal(cfg, &table2)
	assert.NoError(t, err)
	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)
}

func TestMarshalColumnValidation(t *testing.T) {
	// Test that column count mismatches are rejected
	tests := []struct {
		name  string
		table *common.Table
	}{
		{
			name: "row too few columns",
			table: &common.Table{
				Headers: []string{"A", "B"},
				Rows:    [][]string{{"1"}},
			},
		},
		{
			name: "row too many columns",
			table: &common.Table{
				Headers: []string{"A", "B"},
				Rows:    [][]string{{"1", "2", "3"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := &common.Config{
				Writer: &buf,
			}
			err := Marshal(cfg, tt.table)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "columns")
		})
	}
}

func TestMarshalEmptyHeaders(t *testing.T) {
	// Test that empty headers are rejected
	table := &common.Table{
		Headers: []string{},
		Rows:    [][]string{{"1"}},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one header")
}

func TestMarshalNilTable(t *testing.T) {
	// Test nil table
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestMarshalEmptyRows(t *testing.T) {
	// Test table with headers but no rows
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	expected := `{| class="wikitable"
! A !! B
|}
`

	assert.Equal(t, expected, output)
}

func TestMarshal(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a valid config with the buffer as writer
	cfg := &common.Config{
		Writer: &buf,
	}

	// Create a table with 1 header and multiple rows
	table := &common.Table{
		Headers: []string{"Header1"},
		Rows: [][]string{
			{"Row1Cell1"},
			{"Row2Cell1"},
			{"Row3Cell1"},
		},
	}

	// Call the Marshal function
	err := Marshal(cfg, table)
	assert.NoError(t, err, "Marshal failed")

	// Get the output
	output := buf.String()

	// Check for proper row separators between each row
	// Should have 3 separators: after header, and between each row
	expectedRowSeparators := 3
	actualRowSeparators := bytes.Count([]byte(output), []byte("|-\n"))
	assert.Equal(t, expectedRowSeparators, actualRowSeparators, "Incorrect number of row separators")

	// Verify the full output format
	expectedPattern := `{| class="wikitable"
! Header1
|-
| Row1Cell1
|-
| Row2Cell1
|-
| Row3Cell1
|}
`
	assert.Equal(t, expectedPattern, output, "Output format is incorrect")
}
