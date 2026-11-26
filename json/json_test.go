package json

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestUnmarshalWithBOM(t *testing.T) {
	input := "\xEF\xBB\xBF[{\"a\":\"1\"}]"
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, table.Headers)
	assert.Equal(t, [][]string{{"1"}}, table.Rows)
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

// TestDefaultFormatFirstSeenHeaders tests that headers follow the source key order.
func TestDefaultFormatFirstSeenHeaders(t *testing.T) {
	for i := 0; i < 10; i++ {
		input := `[{"z": 1, "a": 2, "m": 3}]`
		cfg := &common.Config{
			Extension: map[string]string{},
			Reader:    strings.NewReader(input),
		}
		table := &common.Table{}
		err := Unmarshal(cfg, table)

		assert.NoError(t, err)
		assert.Equal(t, []string{"z", "a", "m"}, table.Headers)
	}
}

func TestDefaultFormatSortHeaders(t *testing.T) {
	input := `[{"z": 1, "a": 2, "m": 3}]`
	cfg := &common.Config{
		Extension: map[string]string{"sort-headers": "true"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "m", "z"}, table.Headers)
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

func TestObjectFormatEmptyRows(t *testing.T) {
	var buf bytes.Buffer
	err := Marshal(&common.Config{Writer: &buf}, &common.Table{Headers: []string{"a", "b"}})
	assert.NoError(t, err)
	assert.JSONEq(t, `[]`, buf.String())
}

func TestObjectFormatRejectsRaggedRows(t *testing.T) {
	var buf bytes.Buffer
	err := Marshal(&common.Config{Writer: &buf}, &common.Table{
		Headers: []string{"a", "b"},
		Rows:    [][]string{{"1"}},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "row length 1 does not match header length 2")
}

func TestUnmarshalPreservesNumberPrecision(t *testing.T) {
	input := `[{"id":1234567,"big":9007199254740993}]`
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"id", "big"}, table.Headers)
	assert.Equal(t, [][]string{{"1234567", "9007199254740993"}}, table.Rows)
}

func TestKeyedFormat(t *testing.T) {
	table := &common.Table{
		Headers: []string{"id", "name", "age"},
		Rows: [][]string{
			{"1", "Alice", "30"},
			{"2", "Bob", "25"},
		},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"format": "keyed",
			"minify": "true",
		},
	}
	err := Marshal(cfg, table)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"1":{"name":"Alice","age":"30"},"2":{"name":"Bob","age":"25"}}`, buf.String())

	got := &common.Table{}
	err = Unmarshal(&common.Config{
		Reader:    bytes.NewReader(buf.Bytes()),
		Extension: map[string]string{"format": "keyed"},
	}, got)
	assert.NoError(t, err)
	assert.Equal(t, []string{"key", "age", "name"}, got.Headers)
	assert.Equal(t, [][]string{
		{"1", "30", "Alice"},
		{"2", "25", "Bob"},
	}, got.Rows)
}

func TestKeyedFormatDuplicateKey(t *testing.T) {
	table := &common.Table{
		Headers: []string{"id", "name"},
		Rows: [][]string{
			{"1", "Alice"},
			{"1", "Bob"},
		},
	}
	var buf bytes.Buffer
	err := Marshal(&common.Config{
		Writer:    &buf,
		Extension: map[string]string{"format": "keyed"},
	}, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
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

func TestUnmarshalNestedObjectFlattens(t *testing.T) {
	input := `[{"id":1,"user":{"name":"Alice"},"tags":["a","b"]}]`
	cfg := &common.Config{Extension: map[string]string{}, Reader: strings.NewReader(input)}
	table := &common.Table{}
	require.NoError(t, Unmarshal(cfg, table))

	assert.Equal(t, []string{"id", "user.name", "tags"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "Alice", `["a","b"]`}}, table.Rows)
}

func TestUnmarshalNestedArrayIndexFlattens(t *testing.T) {
	input := `[{"id":1,"tags":["a","b"]}]`
	cfg := &common.Config{
		Extension: map[string]string{"flatten-arrays": "index"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	require.NoError(t, Unmarshal(cfg, table))

	assert.Equal(t, []string{"id", "tags[0]", "tags[1]"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "a", "b"}}, table.Rows)
}

func TestUnmarshalSingleObject(t *testing.T) {
	input := `{"a":1,"b":{"c":2}}`
	cfg := &common.Config{Extension: map[string]string{}, Reader: strings.NewReader(input)}
	table := &common.Table{}
	require.NoError(t, Unmarshal(cfg, table))

	assert.Equal(t, []string{"a", "b.c"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "2"}}, table.Rows)
}

func TestUnmarshalJSONPath(t *testing.T) {
	input := `{"data":{"items":[{"a":1},{"a":2}]}}`
	cfg := &common.Config{
		Extension: map[string]string{"json-path": "data.items"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	require.NoError(t, Unmarshal(cfg, table))

	assert.Equal(t, []string{"a"}, table.Headers)
	assert.Equal(t, [][]string{{"1"}, {"2"}}, table.Rows)
}

func TestUnmarshalNestedObjectAsCompactJSON(t *testing.T) {
	input := `[{"meta":{"a":1,"b":[2,3]}}]`
	cfg := &common.Config{
		Extension: map[string]string{"flatten-arrays": "json"},
		Reader:    strings.NewReader(input),
	}
	table := &common.Table{}
	require.NoError(t, Unmarshal(cfg, table))

	assert.Equal(t, []string{"meta.a", "meta.b"}, table.Headers)
	assert.Equal(t, [][]string{{"1", `[2,3]`}}, table.Rows)
}

func TestMarshalParseTypesStrict(t *testing.T) {
	table := &common.Table{
		Headers: []string{"zip", "id", "ok", "missing", "amount"},
		Rows:    [][]string{{"007", "9007199254740993", "true", "null", "1.50"}},
	}
	var buf strings.Builder
	cfg := &common.Config{
		Writer:    &buf,
		Extension: map[string]string{"parse-types": "true"},
	}
	require.NoError(t, Marshal(cfg, table))

	var output []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &output))
	assert.Equal(t, "007", output[0]["zip"])
	assert.Equal(t, "9007199254740993", output[0]["id"])
	assert.Equal(t, true, output[0]["ok"])
	assert.Nil(t, output[0]["missing"])
	assert.Equal(t, 1.5, output[0]["amount"])
}
