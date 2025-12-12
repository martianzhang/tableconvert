package jsonl

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
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
