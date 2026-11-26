package ascii

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestUnmarshalDotAndBubbleRoundTrip(t *testing.T) {
	for _, style := range []string{"dot", "bubble", "plus"} {
		t.Run(style, func(t *testing.T) {
			table := &common.Table{
				Headers: []string{"h", "two"},
				Rows: [][]string{
					{"a|b", "c" + map[string]string{"dot": "·", "bubble": "◌", "plus": "+"}[style] + "d"},
				},
			}

			var buf bytes.Buffer
			cfg := &common.Config{
				Writer:    &buf,
				Extension: map[string]string{"style": style},
			}
			assert.NoError(t, Marshal(cfg, table))

			got := &common.Table{}
			readCfg := &common.Config{
				Reader:    bytes.NewReader(buf.Bytes()),
				Extension: map[string]string{"style": style},
			}
			assert.NoError(t, Unmarshal(readCfg, got))
			assert.Equal(t, table.Headers, got.Headers)
			assert.Equal(t, table.Rows, got.Rows)
		})
	}
}

func TestMarshalAndUnmarshalBoxPipe(t *testing.T) {
	table := &common.Table{
		Headers: []string{"h|1", `path\name`},
		Rows: [][]string{
			{"a|b", `C:\temp`},
		},
	}
	var buf bytes.Buffer
	assert.NoError(t, Marshal(&common.Config{Writer: &buf}, table))

	got := &common.Table{}
	assert.NoError(t, Unmarshal(&common.Config{Reader: bytes.NewReader(buf.Bytes())}, got))
	assert.Equal(t, table.Headers, got.Headers)
	assert.Equal(t, table.Rows, got.Rows)
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
			err:      errors.New("Marshal: row 0 has 1 columns, but table header has 2 columns"),
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

func TestMarshal_CustomStyleTable(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create config with custom style
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"style": "*",
		},
	}

	// Create test table
	table := &common.Table{
		Headers: []string{"ID", "Name", "Value"},
		Rows: [][]string{
			{"1", "Item A", "100"},
			{"2", "Item B", "200"},
		},
	}

	// Expected output
	expected := `***********************
* ID * Name   * Value *
***********************
* 1  * Item A * 100   *
* 2  * Item B * 200   *
***********************
`

	// Call the function
	err := Marshal(cfg, table)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}

func TestMarshalStyleCaseInsensitive(t *testing.T) {
	table := &common.Table{
		Headers: []string{"h"},
		Rows:    [][]string{{"v"}},
	}
	var upper, lower bytes.Buffer
	require.NoError(t, Marshal(&common.Config{Writer: &upper, Extension: map[string]string{"style": "PLUS"}}, table))
	require.NoError(t, Marshal(&common.Config{Writer: &lower, Extension: map[string]string{"style": "plus"}}, table))
	assert.Equal(t, lower.String(), upper.String())
}
