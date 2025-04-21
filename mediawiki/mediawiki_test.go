package mediawiki

import (
	"bytes"
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
