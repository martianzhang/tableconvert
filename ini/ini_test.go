package ini

import (
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	input := "[server]\nhost = localhost\nport = 8080\n"
	table := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table))
	assert.Equal(t, []string{"section", "key", "value"}, table.Headers)
	assert.Contains(t, table.Rows, []string{"server", "host", "localhost"})
	assert.Contains(t, table.Rows, []string{"server", "port", "8080"})
}

func TestMarshal(t *testing.T) {
	table := &common.Table{
		Headers: []string{"section", "key", "value"},
		Rows:    [][]string{{"server", "host", "localhost"}, {"server", "port", "8080"}},
	}
	var buf strings.Builder
	require.NoError(t, Marshal(&common.Config{Writer: &buf}, table))
	assert.Contains(t, buf.String(), "[server]")
	assert.Contains(t, buf.String(), "host = localhost")
	assert.Contains(t, buf.String(), "port = 8080")
}
