package csv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	input := `name,age,city
Alice,30,New York
Bob,25,Los Angeles`

	args := []string{"--from", "csv", "--to", "csv"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)

	assert.Equal(t, []string{"name", "age", "city"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
	}, table.Rows)
}

func TestMarshal(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name", "age", "city"},
		Rows: [][]string{
			{"Alice", "30", "New York"},
			{"Bob", "25", "Los Angeles"},
		},
	}

	args := []string{"--from", "csv", "--to", "csv"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf
	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	expectedCSV := "name,age,city\nAlice,30,New York\nBob,25,Los Angeles\n"
	assert.Equal(t, expectedCSV, buf.String())
}
