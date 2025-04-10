package csv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	csvData := `name,age,city
Alice,30,New York
Bob,25,Los Angeles`

	reader := strings.NewReader(csvData)
	var table common.Table

	err := Unmarshal(reader, &table)
	assert.NoError(t, err)

	expectedHeaders := []string{"name", "age", "city"}
	expectedRows := [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
	}

	assert.Equal(t, expectedHeaders, table.Headers)
	assert.Equal(t, expectedRows, table.Rows)
}

func TestMarshal(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name", "age", "city"},
		Rows: [][]string{
			{"Alice", "30", "New York"},
			{"Bob", "25", "Los Angeles"},
		},
	}

	var buf bytes.Buffer
	err := Marshal(table, &buf)
	assert.NoError(t, err)

	expectedCSV := "name,age,city\nAlice,30,New York\nBob,25,Los Angeles\n"
	assert.Equal(t, expectedCSV, buf.String())
}
