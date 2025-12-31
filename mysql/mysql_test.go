package mysql

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	mysqlOutput := `
+----------+--------------+------+-----+---------+----------------+
|  FIELD   |     TYPE     | NULL | KEY | DEFAULT |     EXTRA      |
+----------+--------------+------+-----+---------+----------------+
| user_id  | smallint(5)  | NO   | PRI | NULL    | auto_increment |
| username | varchar(10)  | NO   |     | NULL    |                |
| password | varchar(100) | NO   |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
` + "\n" // Add trailing newline like real output often has

	args := []string{"--from", "mysql", "--to", "mysql"}
	cfg, err := common.ParseConfig(args)
	cfg.Reader = strings.NewReader(mysqlOutput)
	assert.Nil(t, err)

	var table common.Table
	err = Unmarshal(&cfg, &table)

	assert.Nil(t, err)
	assert.Equal(t, []string{"FIELD", "TYPE", "NULL", "KEY", "DEFAULT", "EXTRA"}, table.Headers)
	assert.Equal(t, 3, len(table.Rows))
}

// TestUnmarshalUTF8Student verifies UTF-8 handling with a generic student info table (no business data).
func TestUnmarshalUTF8Student(t *testing.T) {
	mysqlOutput := `
+----------+--------------------+--------+----------+--------+----------------+
|  学生ID  |        姓名        | 年级   |   班级   | 状态   |     备注       |
+----------+--------------------+--------+----------+--------+----------------+
| stu001   | 张三               | 高一   | 一班     | 在读   | 喜欢数学       |
| stu002   | 李四               | 高一   | 一班     | 在读   |                |
| stu003   | 王五               | 高二   | 二班     | 休学   | 转专业申请中   |
+----------+--------------------+--------+----------+--------+----------------+
` + "\n"
	args := []string{"--from", "mysql", "--to", "mysql"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(mysqlOutput)
	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.Nil(t, err)
	assert.Equal(t, []string{"学生ID", "姓名", "年级", "班级", "状态", "备注"}, table.Headers)
	assert.Equal(t, 3, len(table.Rows))
	assert.Equal(t, []string{"stu001", "张三", "高一", "一班", "在读", "喜欢数学"}, table.Rows[0])
	assert.Equal(t, []string{"stu003", "王五", "高二", "二班", "休学", "转专业申请中"}, table.Rows[2])
}

func TestParseASCIIArtTable(t *testing.T) {
	asciiArtTable := `

+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
|  A   |       The Good        |    500 |
|  B   | The Very very Bad Man |    288 |
|  C   |       The Ugly        |    120 |
|  D   |      The Gopher       |    800 |
+------+-----------------------+--------+

anything else after the table is ignored

`
	args := []string{"--from", "mysql", "--to", "mysql"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(asciiArtTable)

	var table common.Table
	err = Unmarshal(&cfg, &table)

	assert.Nil(t, err)
	assert.Equal(t, []string{"NAME", "SIGN", "RATING"}, table.Headers)
	assert.Equal(t, 4, len(table.Rows))
	assert.Equal(t, []string{"C", "The Ugly", "120"}, table.Rows[2])
}

func TestUnmarshalEmptyCells(t *testing.T) {
	// Test empty cells are preserved
	input := `+----+----+----+
| A  | B  | C  |
+----+----+----+
| 1  |    | 3  |
+----+----+----+`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B", "C"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"1", "", "3"}, table.Rows[0])
}

func TestUnmarshalHeadersOnly(t *testing.T) {
	// Test table with headers but no data rows
	input := `+----+----+
| A  | B  |
+----+----+`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, 0, len(table.Rows))
}

func TestUnmarshalMissingBottomBorder(t *testing.T) {
	// Test table without bottom border (should be allowed)
	input := `+----+----+
| A  | B  |
+----+----+
| 1  | 2  |`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"1", "2"}, table.Rows[0])
}

func TestUnmarshalColumnMismatch(t *testing.T) {
	// Test that column count mismatches are detected
	input := `+----+----+
| A  | B  |
+----+----+
| 1  | 2  | 3  |
+----+----+`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "column count")
}

func TestUnmarshalPipesInContent(t *testing.T) {
	// Test that pipes in content are detected as errors
	// MySQL format doesn't support pipes in content
	input := `+----+----+
| A  | B  |
+----+----+
| 1  | 2|3|
+----+----+`

	cfg := &common.Config{
		Reader: bytes.NewBufferString(input),
	}
	var table common.Table
	err := Unmarshal(cfg, &table)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "column count")
}

func TestMarshalEmptyRows(t *testing.T) {
	// Test that marshal with no rows doesn't output extra separator
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	expected := `+---+---+
| A | B |
+---+---+
`

	assert.Equal(t, expected, output)
}

func TestMarshalNilTable(t *testing.T) {
	// Test nil table
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestMarshalEmptyHeaders(t *testing.T) {
	// Test empty headers
	table := &common.Table{
		Headers: []string{},
		Rows:    [][]string{},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one header")
}

func TestMarshalColumnMismatch(t *testing.T) {
	// Test column count mismatch
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1"}},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "columns")
}

func TestMarshalPipesInContent(t *testing.T) {
	// Test that pipes in content are output as-is
	// Note: This will create output that cannot be parsed back
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1", "2|3"}},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Should contain the pipe in the output
	assert.Contains(t, output, "2|3")
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
