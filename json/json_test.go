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

// TestColumnFormatUnmarshal tests the column format unmarshal (Bug 1 fix)
func TestColumnFormatUnmarshal(t *testing.T) {
	// This test verifies that col[i] is used instead of col[j]
	input := `[{"col1": [1, 2, 3]}, {"col2": ["a", "b", "c"]}]`
	cfg := &common.Config{
		Extension: map[string]string{"format": "column"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	// Headers are sorted for determinism
	assert.Equal(t, []string{"col1", "col2"}, table.Headers)
	assert.Equal(t, 3, len(table.Rows))
	// Row 0: col1=1, col2=a
	assert.Equal(t, "1", table.Rows[0][0])
	assert.Equal(t, "a", table.Rows[0][1])
	// Row 1: col1=2, col2=b
	assert.Equal(t, "2", table.Rows[1][0])
	assert.Equal(t, "b", table.Rows[1][1])
	// Row 2: col1=3, col2=c
	assert.Equal(t, "3", table.Rows[2][0])
	assert.Equal(t, "c", table.Rows[2][1])
}

// TestDefaultFormatDeterministicHeaders tests header order determinism (Bug 2 fix)
func TestDefaultFormatDeterministicHeaders(t *testing.T) {
	// Run multiple times to ensure deterministic order
	for i := 0; i < 10; i++ {
		input := `[{"z": 1, "a": 2, "m": 3}]`
		cfg := &common.Config{
			Extension: map[string]string{},
			Reader:    strings.NewReader(input),
		}
		table := &common.Table{}
		err := Unmarshal(cfg, table)

		assert.NoError(t, err)
		// Headers should always be sorted
		assert.Equal(t, []string{"a", "m", "z"}, table.Headers)
	}
}

// Test2dFormatMissingValues tests handling of missing values (Bug 4 fix)
func Test2dFormatMissingValues(t *testing.T) {
	// Row 2 has only 1 value but headers has 2
	input := `[[1, 2], [3]]`
	cfg := &common.Config{
		Extension: map[string]string{"format": "2d"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	// Missing value should be "NULL"
	assert.Equal(t, []string{"3", "NULL"}, table.Rows[0])
}

// Test2dFormatNilHeaders tests handling of nil headers (Bug 5 fix)
func Test2dFormatNilHeaders(t *testing.T) {
	// Header with nil value
	input := `[[null, "col2"], [1, 2]]`
	cfg := &common.Config{
		Extension: map[string]string{"format": "2d"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"NULL", "col2"}, table.Headers)
}

// Test2dFormatNilValues tests handling of nil values in data
func Test2dFormatNilValues(t *testing.T) {
	input := `[["a", "b"], [null, 2]]`
	cfg := &common.Config{
		Extension: map[string]string{"format": "2d"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"NULL", "2"}, table.Rows[0])
}

// TestColumnFormatNilValues tests handling of nil values in column format
func TestColumnFormatNilValues(t *testing.T) {
	input := `[{"col1": [1, null, 3]}]`
	cfg := &common.Config{
		Extension: map[string]string{"format": "column"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"col1"}, table.Headers)
	assert.Equal(t, 3, len(table.Rows))
	assert.Equal(t, "1", table.Rows[0][0])
	assert.Equal(t, "NULL", table.Rows[1][0])
	assert.Equal(t, "3", table.Rows[2][0])
}

// TestDefaultFormatNilValues tests handling of nil values in default format
func TestDefaultFormatNilValues(t *testing.T) {
	input := `[{"a": 1, "b": null}, {"a": null, "b": 2}]`
	cfg := &common.Config{
		Extension: map[string]string{},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, table.Headers)
	assert.Equal(t, 2, len(table.Rows))
	assert.Equal(t, []string{"1", "NULL"}, table.Rows[0])
	assert.Equal(t, []string{"NULL", "2"}, table.Rows[1])
}

// TestColumnFormatRoundTrip tests column format round-trip
func TestColumnFormatRoundTrip(t *testing.T) {
	// Unmarshal
	input := `[{"col1": [1, 2, 3]}, {"col2": ["a", "b", "c"]}]`
	cfgIn := &common.Config{
		Extension: map[string]string{"format": "column"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfgIn, table)
	assert.NoError(t, err)

	// Marshal
	var buf strings.Builder
	cfgOut := &common.Config{
		Extension: map[string]string{"format": "column"},
		Writer:    &buf,
	}
	err = Marshal(cfgOut, table)
	assert.NoError(t, err)

	// Parse output
	var output []map[string]interface{}
	err = json.Unmarshal([]byte(buf.String()), &output)
	assert.NoError(t, err)

	// Verify structure
	assert.Equal(t, 2, len(output))
	// Output order may vary due to map iteration, so check both possibilities
	hasCol1 := false
	hasCol2 := false
	for _, obj := range output {
		if col, ok := obj["col1"]; ok {
			hasCol1 = true
			arr, _ := col.([]interface{})
			assert.Equal(t, 3, len(arr))
		}
		if col, ok := obj["col2"]; ok {
			hasCol2 = true
			arr, _ := col.([]interface{})
			assert.Equal(t, 3, len(arr))
		}
	}
	assert.True(t, hasCol1 && hasCol2, "Should have both col1 and col2")
}

// Test2dFormatRoundTrip tests 2d format round-trip
func Test2dFormatRoundTrip(t *testing.T) {
	// Unmarshal
	input := `[[1, 2], [3, 4], [5, 6]]`
	cfgIn := &common.Config{
		Extension: map[string]string{"format": "2d"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfgIn, table)
	assert.NoError(t, err)

	// Marshal
	var buf strings.Builder
	cfgOut := &common.Config{
		Extension: map[string]string{"format": "2d"},
		Writer:    &buf,
	}
	err = Marshal(cfgOut, table)
	assert.NoError(t, err)

	// Parse output
	var output [][]interface{}
	err = json.Unmarshal([]byte(buf.String()), &output)
	assert.NoError(t, err)

	// Verify
	assert.Equal(t, 3, len(output))
	assert.Equal(t, []interface{}{"1", "2"}, output[0])
	assert.Equal(t, []interface{}{"3", "4"}, output[1])
	assert.Equal(t, []interface{}{"5", "6"}, output[2])
}
