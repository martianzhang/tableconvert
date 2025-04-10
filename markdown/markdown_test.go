package markdown

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	input := "```txt\n" + // Added code fence for realism
		"|   DATE   |         DESCRIPTION         | CV2  | AMOUNT |\n" +
		"|----------|--------------------------|------|--------|\n" +
		"| 1/1/2014 | Domain name              | 2233 | $10.98 |\n" +
		"| 1/1/2014 | January Hosting          | 2233 | $54.95 |\n" +
		"| 1/4/2014 | February Hosting         | 2233 | $51.00 |\n" +
		"| 1/4/2014 | February Extra Bandwidth | 2233 | $30.00 |\n" +
		"```" // Added closing fence

	reader := strings.NewReader(input)
	var table common.Table

	err := Unmarshal(reader, &table)

	assert.Nil(t, err)
	assert.Equal(t, []string{"DATE", "DESCRIPTION", "CV2", "AMOUNT"}, table.Headers)
	assert.Equal(t, 4, len(table.Rows))
	assert.Equal(t, []string{"1/4/2014", "February Hosting", "2233", "$51.00"}, table.Rows[2])
}

func TestUmarshalEmptyCells(t *testing.T) {
	input := `
|FIELD|TYPE|NULL|KEY|DEFAULT|EXTRA|
|---|---|---|---|---|---|
|user_id|smallint(5)|NO|PRI|NULL|auto_increment|
|username|varchar(10)|NO||NULL||
|password|varchar(100)|NO||NULL||
`

	reader := strings.NewReader(input)
	var table common.Table

	err := Unmarshal(reader, &table)
	assert.Nil(t, err)
}

type mockWriter struct {
	buf        bytes.Buffer
	failAt     int
	writeCount int
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.writeCount++
	if m.failAt > 0 && m.writeCount >= m.failAt {
		return 0, errors.New("mock write error")
	}
	return m.buf.Write(p)
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name        string
		table       *common.Table
		writer      io.Writer
		expectedErr string
		expectedOut string
	}{
		{
			name:        "nil table",
			table:       nil,
			writer:      &mockWriter{},
			expectedErr: "Marshal: input table pointer cannot be nil",
		},
		{
			name: "empty headers",
			table: &common.Table{
				Headers: []string{},
				Rows:    [][]string{},
			},
			writer:      &mockWriter{},
			expectedErr: "Marshal: table must have at least one header",
		},
		{
			name: "column count mismatch",
			table: &common.Table{
				Headers: []string{"Header1", "Header2"},
				Rows: [][]string{
					{"Cell1"},
				},
			},
			writer:      &mockWriter{},
			expectedErr: "Marshal: 1 row has 1 columns, but table has 2",
		},
		{
			name: "successful marshal",
			table: &common.Table{
				Headers: []string{"Header1", "Header2"},
				Rows: [][]string{
					{"Cell1", "Cell2"},
					{"Cell3", "Cell4"},
				},
			},
			writer:      &mockWriter{},
			expectedOut: "|Header1|Header2|\n|---|---|\n|Cell1|Cell2|\n|Cell3|Cell4|\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, ok := tt.writer.(*mockWriter)
			if ok {
				mock.buf.Reset()
				mock.writeCount = 0
			}

			err := Marshal(tt.table, tt.writer)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				if mock != nil {
					assert.Equal(t, tt.expectedOut, mock.buf.String())
				}
			}
		})
	}
}
