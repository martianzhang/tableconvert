package toml

import (
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalArrayOfTables(t *testing.T) {
	input := "[[row]]\nname = \"Alice\"\nage = 30\n[[row]]\nname = \"Bob\"\nage = 25\n"
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table))
	assert.Equal(t, []string{"age", "name"}, table.Headers)
	assert.Equal(t, [][]string{{"30", "Alice"}, {"25", "Bob"}}, table.Rows)
}

func TestUnmarshalColumnShape(t *testing.T) {
	input := "a = [1, 2]\nb = [\"x\", \"y\"]\n"
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table))
	assert.Equal(t, []string{"a", "b"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "x"}, {"2", "y"}}, table.Rows)
}

func TestMarshalArrayOfTables(t *testing.T) {
	table := &common.Table{Headers: []string{"a", "b"}, Rows: [][]string{{"1", "2"}}}
	var buf strings.Builder
	require.NoError(t, Marshal(&common.Config{Writer: &buf}, table))
	assert.Contains(t, buf.String(), "[[row]]")
	assert.Contains(t, buf.String(), "a = '1'")

	buf.Reset()
	require.NoError(t, Marshal(&common.Config{Writer: &buf, Extension: map[string]string{"toml-table": "items"}}, table))
	assert.Contains(t, buf.String(), "[[items]]")
}
