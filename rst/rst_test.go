package rst

import (
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	table := &common.Table{Headers: []string{"name", "age"}, Rows: [][]string{{"Alice", "30"}}}
	var buf strings.Builder
	require.NoError(t, Marshal(&common.Config{Writer: &buf}, table))
	output := buf.String()
	assert.Contains(t, output, "+-------+-----+")
	assert.Contains(t, output, "| name  | age |")
	assert.Contains(t, output, "| Alice | 30  |")
}
