package ascii

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	input := `
+----------+--------------+------+-----+---------+----------------+
|  FIELD   |     TYPE     | NULL | KEY | DEFAULT |     EXTRA      |
+----------+--------------+------+-----+---------+----------------+
| user_id  | smallint(5)  | NO  | PRI | NULL    | auto_increment |
| username | varchar(10)  | NO  |     | NULL    |                |
| password | varchar(100) | NO  |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
` + "\n" // Add trailing newline like real output often has

	args := []string{"--from", "ascii", "--to", "ascii"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.Nil(t, err)
	assert.Equal(t, []string{"FIELD", "TYPE", "NULL", "KEY", "DEFAULT", "EXTRA"}, table.Headers)
	assert.Equal(t, 3, len(table.Rows))
}

func TestUnmarshalPlus(t *testing.T) {
	input := `
	+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	+ FIELD    + TYPE         + NULL + KEY + DEFAULT + EXTRA          +
	+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	+ user_id  + smallint(5)  + NO   + PRI + NULL    + auto_increment +
	+ username + varchar(10)  + NO   +     + NULL    +                +
	+ password + varchar(100) + NO   +     +         +                +
	+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
` + "\n" // Add trailing newline like real output often has

	args := []string{"--from", "ascii", "--to", "ascii", "--style", "plus"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.Nil(t, err)
	assert.Equal(t, []string{"FIELD", "TYPE", "NULL", "KEY", "DEFAULT", "EXTRA"}, table.Headers)
	assert.Equal(t, 3, len(table.Rows))
}

func TestParseASCIIArtTable(t *testing.T) {
	input := `

+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
|  A   |       The Good        |    500 |
|  B   | The Very very Bad Man |    288 |
|  C   |     The Ugly        |    120 |
|  D   |    The Gopher       |    800 |
+------+-----------------------+--------+

anything else after the table is ignored

`
	args := []string{"--from", "ascii", "--to", "ascii"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)

	assert.Nil(t, err)
	assert.Equal(t, []string{"NAME", "SIGN", "RATING"}, table.Headers)
	assert.Equal(t, 4, len(table.Rows))
	assert.Equal(t, []string{"C", "The Ugly", "120"}, table.Rows[2])
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name     string
		table    *common.Table
		expected string
		err      error
	}{
		{
			name:     "nil table",
			table:    nil,
			expected: "",
			err:      errors.New("Marshal: input table pointer cannot be nil"),
		},
		{
			name: "empty headers",
			table: &common.Table{
				Headers: []string{},
				Rows:    [][]string{},
			},
			expected: "",
			err:      errors.New("Marshal: table must have at least one header"),
		},
		{
			name: "column count mismatch",
			table: &common.Table{
				Headers: []string{"A", "B"},
				Rows: [][]string{
					{"1"},
				},
			},
			expected: "",
			err:      errors.New("Marshal: 0 row has 1 columns, but table has 2"),
		},
		{
			name: "single row",
			table: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice", "25"},
				},
			},
			expected: "+-------+-----+\n| Name  | Age |\n+-------+-----+\n| Alice | 25  |\n+-------+-----+\n",
			err:      nil,
		},
		{
			name: "multiple rows with different widths",
			table: &common.Table{
				Headers: []string{"ID", "Name", "Description"},
				Rows: [][]string{
					{"1", "Alice", "Developer"},
					{"2", "Bob", "Senior Developer with long title"},
				},
			},
			expected: "+----+-------+----------------------------------+\n" +
				"| ID | Name  | Description                      |\n" +
				"+----+-------+----------------------------------+\n" +
				"| 1  | Alice | Developer                        |\n" +
				"| 2  | Bob   | Senior Developer with long title |\n" +
				"+----+-------+----------------------------------+\n",
			err: nil,
		},
	}

	args := []string{"--from", "mysql", "--to", "mysql"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg.Writer = &buf
			err := Marshal(&cfg, tt.table)

			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, buf.String())
			}
		})
	}
}
