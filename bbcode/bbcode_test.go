package bbcode

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
	assert.Equal(t, "[table]\n[tr][th]name[/th][th]age[/th][/tr]\n[tr][td]Alice[/td][td]30[/td][/tr]\n[/table]\n", buf.String())
}
