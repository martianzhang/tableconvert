package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatsFromCSV(t *testing.T) {
	input := "name,age\nAlice,30\n"
	tests := []struct {
		format string
		marker string
	}{
		{"yaml", "name: Alice"},
		{"toml", "[[row]]"},
		{"bbcode", "[table]"},
		{"jira", "||name|age|"},
		{"asciidoc", "|==="},
		{"rst", "+-------+"},
		{"textile", "|_. name |_. age |"},
		{"php", "$headers = [\"name\", \"age\"];"},
		{"ruby", "headers = [\"name\", \"age\"]"},
		{"pandas", "df = pd.DataFrame(["},
	}
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := &common.Config{
				From:      "csv",
				To:        tt.format,
				Reader:    strings.NewReader(input),
				Writer:    &buf,
				Extension: map[string]string{},
			}
			require.NoError(t, performConversion(cfg))
			assert.Contains(t, buf.String(), tt.marker)
		})
	}
}

func TestNewINIFormat(t *testing.T) {
	input := "section,key,value\nserver,host,localhost\nserver,port,8080\n"
	var buf bytes.Buffer
	cfg := &common.Config{
		From:      "csv",
		To:        "ini",
		Reader:    strings.NewReader(input),
		Writer:    &buf,
		Extension: map[string]string{},
	}
	require.NoError(t, performConversion(cfg))
	assert.Contains(t, buf.String(), "[server]")
	assert.Contains(t, buf.String(), "host = localhost")
}

func TestInputLimits(t *testing.T) {
	var buf bytes.Buffer
	cfg := &common.Config{
		From:      "csv",
		To:        "json",
		Reader:    strings.NewReader("a\n1\n2\n"),
		Writer:    &buf,
		Extension: map[string]string{"max-rows": "1"},
	}
	err := performConversion(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max-rows")

	buf.Reset()
	cfg = &common.Config{
		From:      "csv",
		To:        "json",
		Reader:    strings.NewReader("a\n1\n2\n"),
		Writer:    &buf,
		Extension: map[string]string{"max-bytes": "3"},
	}
	err = performConversion(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max-bytes")
}
