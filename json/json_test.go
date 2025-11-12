package json

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDefaultFormatBasicCase(t *testing.T) {
	// Prepare input JSON
	input := `[{"id":1,"name":"test"}, {"id":2, "name":"John"}]`

	// Create config with default format (no format specified)
	cfg := &common.Config{
		Extension: map[string]string{}, // empty map for default format
		Reader:    strings.NewReader(input),
	}

	// Create empty table to store results
	table := &common.Table{}

	// Call Unmarshal
	err := Unmarshal(cfg, table)

	// Assert no error occurred
	assert.NoError(t, err, "Expected no error during unmarshaling")

	// Assert headers are correct
	expectedHeaders := []string{"id", "name"}
	sort.Strings(expectedHeaders)
	sort.Strings(table.Headers)
	assert.Equal(t, len(expectedHeaders), len(table.Headers), "Headers don't match expected values")

	// Assert row data is correct
	expectedRows := [][]string{{"1", "test"}, {"2", "John"}}
	assert.Equal(t, len(expectedRows), len(table.Rows), "Row data doesn't match expected values")
}

func TestMarshalDefaultObjectFormat(t *testing.T) {
	// Create a buffer to capture the output
	var buf strings.Builder

	// Create config with default format (empty string)
	cfg := &common.Config{
		Extension: map[string]string{},
		Writer:    &buf,
	}

	// Create table with test data
	table := &common.Table{
		Headers: []string{"name", "age"},
		Rows:    [][]string{{"John", "30"}, {"Alice", "25"}},
	}

	// Call Marshal
	err := Marshal(cfg, table)

	// Assert no error occurred
	assert.NoError(t, err, "Expected no error during marshaling")

	var expected, actual []map[string]interface{}
	err = json.Unmarshal([]byte(`[
	{
	  "name": "John",
	  "age": "30"
	},
	{
		"name": "Alice",
		"age": "25"
	}
  ]`), &expected)
	assert.Nil(t, err)

	// Get the actual output
	err = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &actual)
	assert.Nil(t, err)

	// Assert the output matches expected
	assert.Equal(t, expected, actual, "JSON output doesn't match expected format")
}
