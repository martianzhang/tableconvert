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
	assert.Equal(t, []string{"id", "big"}, table.Headers)
	assert.Equal(t, [][]string{{"1234567", "9007199254740993"}}, table.Rows)
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

// TestUnmarshalFirstSeenHeaders tests that headers follow first-seen key order.
func TestUnmarshalFirstSeenHeaders(t *testing.T) {
	input := `{"z": 1, "a": 2}
{"m": 3, "z": 4}`
	reader := bytes.NewReader([]byte(input))
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"z", "a", "m"}, table.Headers)
	assert.Equal(t, 2, len(table.Rows), "Row count doesn't match expected value")
	assert.Equal(t, "1", table.Rows[0][0]) // z=1
	assert.Equal(t, "2", table.Rows[0][1]) // a=2
	assert.Equal(t, "", table.Rows[0][2])  // m missing
	assert.Equal(t, "4", table.Rows[1][0]) // z=4
	assert.Equal(t, "", table.Rows[1][1])  // a missing
	assert.Equal(t, "3", table.Rows[1][2]) // m=3
}

func TestUnmarshalSortHeaders(t *testing.T) {
	input := `{"z": 1, "a": 2}
{"m": 3, "z": 4}`
	cfg := &common.Config{
		Extension: map[string]string{"sort-headers": "true"},
		Reader:    bytes.NewReader([]byte(input)),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "m", "z"}, table.Headers)
}

func TestUnmarshalNestedObjectFlattens(t *testing.T) {
	input := `{"id":1,"user":{"name":"Alice"}}
{"id":2,"user":{"name":"Bob"}}
`
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table))

	assert.Equal(t, []string{"id", "user.name"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "Alice"}, {"2", "Bob"}}, table.Rows)
}

func TestUnmarshalArrayLine(t *testing.T) {
	input := `[{"a":1},{"a":2}]
{"a":3}
`
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table))

	assert.Equal(t, []string{"a"}, table.Headers)
	assert.Equal(t, [][]string{{"1"}, {"2"}, {"3"}}, table.Rows)
}

func TestMarshalRawRowsRoundTrip(t *testing.T) {
	input := `{"id":1,"ok":true,"x":null,"nested":{"a":1},"big":9007199254740993,"f":1.50}
`
	ingest := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(input)}, ingest))

	var buf strings.Builder
	require.NoError(t, Marshal(&common.Config{From: "jsonl", Writer: &buf}, ingest))

	output := buf.String()
	assert.Contains(t, output, `"id":1`)
	assert.Contains(t, output, `"ok":true`)
	assert.Contains(t, output, `"x":null`)
	assert.Contains(t, output, `"nested":{"a":1}`)
	assert.Contains(t, output, `"big":9007199254740993`)
	assert.Contains(t, output, `"f":1.50`)
}

func TestUnmarshalScalarLines(t *testing.T) {
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader("1\n2\n")}, table))
	assert.Equal(t, []string{"value"}, table.Headers)
	assert.Equal(t, [][]string{{"1"}, {"2"}}, table.Rows)
}
