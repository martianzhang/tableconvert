package php

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
	assert.Contains(t, buf.String(), "$headers = [\"name\", \"age\"];")
	assert.Contains(t, buf.String(), "    [\"Alice\", \"30\"],")
}
