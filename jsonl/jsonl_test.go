package jsonl

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	// Prepare test input
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	table := &common.Table{
		Headers: []string{"name", "age"},
		Rows:    [][]string{{"John", "30"}},
	}

	// Execute the function
	err := Marshal(cfg, table)

	// Check for errors
	if err != nil {
		t.Errorf("Marshal returned unexpected error: %v", err)
	}

	// Check the output
	var expected, actual map[string]string
	json.Unmarshal([]byte("{\"name\":\"John\",\"age\":\"30\"}\n"), &expected)
	json.Unmarshal(buf.Bytes(), &actual)

	assert.Equal(t, expected, actual, "Output doesn't match expected values")

}

func TestUnmarshalValidSingleJSONLine(t *testing.T) {
	// Prepare test input
	input := `{"name": "John", "age": 30}`
	reader := bytes.NewReader([]byte(input))
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	// Execute the function
	err := Unmarshal(cfg, table)

	// Assert no error occurred
	assert.NoError(t, err)

	assert.Equal(t, 2, len(table.Headers), "Header count doesn't match expected value")
	assert.Equal(t, 1, len(table.Rows), "Row count doesn't match expected value")
}

func TestUnmarshalWithBOM(t *testing.T) {
	input := "\xEF\xBB\xBF{\"a\":\"1\"}\n"
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: bytes.NewReader([]byte(input))}, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, table.Headers)
	assert.Equal(t, [][]string{{"1"}}, table.Rows)
}

func TestUnmarshalPreservesNumberPrecision(t *testing.T) {
	input := `{"id":1234567,"big":9007199254740993}` + "\n"
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: bytes.NewReader([]byte(input))}, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"big", "id"}, table.Headers)
	assert.Equal(t, [][]string{{"9007199254740993", "1234567"}}, table.Rows)
}

func TestEmptyJSONLIsEmptyTable(t *testing.T) {
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: bytes.NewReader(nil)}, table)
	require.NoError(t, err)
	assert.Empty(t, table.Headers)
	assert.Empty(t, table.Rows)
}

func TestUnmarshalLongLine(t *testing.T) {
	long := strings.Repeat("x", 70000)
	input := `{"a":"` + long + `"}` + "\n"
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table)
	require.NoError(t, err)
	require.Len(t, table.Rows, 1)
	assert.Equal(t, long, table.Rows[0][0])
}

// TestUnmarshalDeterministicHeaders tests that headers are in sorted order
func TestUnmarshalDeterministicHeaders(t *testing.T) {
	// Prepare test input with fields in different orders
	input := `{"z": 1, "a": 2}
{"m": 3, "z": 4}`
	reader := bytes.NewReader([]byte(input))
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	// Headers should be sorted alphabetically
	assert.Equal(t, []string{"a", "m", "z"}, table.Headers, "Headers should be sorted for deterministic order")
	assert.Equal(t, 2, len(table.Rows), "Row count doesn't match expected value")
	// Check row values
	assert.Equal(t, "2", table.Rows[0][0]) // a=2
	assert.Equal(t, "4", table.Rows[1][2]) // z=4
}
