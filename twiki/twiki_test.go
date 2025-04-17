package twiki

import (
	"bytes"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	table := &common.Table{
		Headers: []string{"Header"},
		Rows: [][]string{
			{"R1"},
			{"R2"},
		},
	}

	// Act
	err := Marshal(cfg, table)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "|=Header=|\n|R1|\n|R2|\n", buf.String())
}

func TestUnmarshal(t *testing.T) {
	// Arrange
	input := "|=Header1=|=Header2=|\n|Data1|Data2|\n|Data3|Data4|"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	// Act
	err := Unmarshal(cfg, table)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Data1", "Data2"},
		{"Data3", "Data4"},
	}, table.Rows)
}
