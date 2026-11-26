package common

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyInputLimit(t *testing.T) {
	cfg := &Config{Reader: strings.NewReader("abcdef"), Extension: map[string]string{"max-bytes": "3"}}
	assert.Error(t, cfg.ApplyInputLimit())

	cfg = &Config{Reader: strings.NewReader("abc"), Extension: map[string]string{"max-bytes": "3"}}
	require.NoError(t, cfg.ApplyInputLimit())
	data, err := io.ReadAll(cfg.Reader)
	require.NoError(t, err)
	assert.Equal(t, "abc", string(data))

	cfg = &Config{Reader: strings.NewReader("abcdef"), Extension: map[string]string{"max-bytes": "0"}}
	require.NoError(t, cfg.ApplyInputLimit())
	data, err = io.ReadAll(cfg.Reader)
	require.NoError(t, err)
	assert.Equal(t, "abcdef", string(data))
}

func TestEnforceRowLimit(t *testing.T) {
	cfg := &Config{Extension: map[string]string{"max-rows": "2"}}
	table := &Table{Rows: [][]string{{"1"}, {"2"}, {"3"}}}
	assert.Error(t, cfg.EnforceRowLimit(table))

	cfg.Extension["max-rows"] = "3"
	assert.NoError(t, cfg.EnforceRowLimit(table))

	cfg.Extension["max-rows"] = "0"
	assert.NoError(t, cfg.EnforceRowLimit(table))
}
