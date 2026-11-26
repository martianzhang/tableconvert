package yaml

import (
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalArray(t *testing.T) {
	input := "- id: 1\n  user:\n    name: Alice\n- id: 2\n  user:\n    name: Bob\n"
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table))
	assert.Equal(t, []string{"id", "user.name"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "Alice"}, {"2", "Bob"}}, table.Rows)
}

func TestUnmarshalColumnAnd2D(t *testing.T) {
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader("a: [1, 2]\nb: [x, y]\n")}, table))
	assert.Equal(t, []string{"a", "b"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "x"}, {"2", "y"}}, table.Rows)

	table = &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader("- [a, b]\n- [1, 2]\n")}, table))
	assert.Equal(t, []string{"a", "b"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "2"}}, table.Rows)
}

func TestMarshalObjectAndColumn(t *testing.T) {
	table := &common.Table{Headers: []string{"a", "b"}, Rows: [][]string{{"1", "2"}}}
	var buf strings.Builder
	require.NoError(t, Marshal(&common.Config{Writer: &buf}, table))
	assert.Contains(t, buf.String(), "a: \"1\"")
	assert.Contains(t, buf.String(), "b: \"2\"")

	buf.Reset()
	require.NoError(t, Marshal(&common.Config{Writer: &buf, Extension: map[string]string{"format": "column", "indent": "4"}}, table))
	assert.Contains(t, buf.String(), "a:")
	assert.Contains(t, buf.String(), "- \"1\"")
}

func TestMarshalFlowAndQuote(t *testing.T) {
	table := &common.Table{Headers: []string{"a", "b"}, Rows: [][]string{{"1", "2"}}}
	var buf strings.Builder
	cfg := &common.Config{
		Writer:    &buf,
		Extension: map[string]string{"flow": "true", "quote": "double"},
	}
	require.NoError(t, Marshal(cfg, table))
	assert.Contains(t, buf.String(), "{")
	assert.Contains(t, buf.String(), "\"1\"")
}

func TestRoundTripPreservesTypes(t *testing.T) {
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader("- a: 1\n  b: true\n")}, table))
	var buf strings.Builder
	cfg := &common.Config{From: "yaml", Writer: &buf}
	require.NoError(t, Marshal(cfg, table))
	assert.Contains(t, buf.String(), "a: 1")
	assert.Contains(t, buf.String(), "b: true")
}
