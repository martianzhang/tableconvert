package markdown

import (
	"bytes"
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
		"```"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.Nil(t, err)

	// Verify headers
	assert.Equal(t, 4, len(table.Headers))
	assert.Equal(t, "DATE", strings.TrimSpace(table.Headers[0]))
	assert.Equal(t, "DESCRIPTION", strings.TrimSpace(table.Headers[1]))

	// Verify row count
	assert.Equal(t, 3, len(table.Rows))

	// Verify first row
	assert.Equal(t, "1/1/2014", strings.TrimSpace(table.Rows[0][0]))
	assert.Equal(t, "Domain name", strings.TrimSpace(table.Rows[0][1]))
}

func TestUmarshalEmptyCells(t *testing.T) {
	input := "| A | B | C |\n|---|---|---|\n| 1 |   | 3 |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, "1", table.Rows[0][0])
	assert.Equal(t, "", table.Rows[0][1])
	assert.Equal(t, "3", table.Rows[0][2])
}

func TestUnmarshalPipesInContent(t *testing.T) {
	input := "| A | B |\n|---|---|\n| x\\|y | z |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, "x|y", table.Rows[0][0])
	assert.Equal(t, "z", table.Rows[0][1])
}

func TestMarshalPipesInContent(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"x|y", "z"}},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.Nil(t, err)

	output := buf.String()
	// Should escape the pipe
	assert.Contains(t, output, "x\\|y")
}

func TestMarshalSpecialChars(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A*B", "C_D"},
		Rows:    [][]string{{"x*y", "z_w"}},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.Nil(t, err)

	output := buf.String()
	// Should escape special chars
	assert.Contains(t, output, "A\\*B")
	assert.Contains(t, output, "C\\_D")
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name     string
		table    *common.Table
		err      string // empty means no error expected
		expected string
	}{
		{
			name:  "nil table",
			table: nil,
			err:   "Marshal: input table pointer cannot be nil",
		},
		{
			name: "empty headers",
			table: &common.Table{
				Headers: []string{},
				Rows:    [][]string{},
			},
			err: "Marshal: table must have at least one header",
		},
		{
			name: "column count mismatch",
			table: &common.Table{
				Headers: []string{"Header1", "Header2"},
				Rows: [][]string{
					{"Cell1"},
				},
			},
			err: "Marshal: 1 row has 1 columns, but table has 2",
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
			err:      "",
			expected: "| Header1 | Header2 |\n|---------|---------|\n| Cell1   | Cell2   |\n| Cell3   | Cell4   |\n",
		},
		{
			name: "no rows",
			table: &common.Table{
				Headers: []string{"ColA", "ColB"},
				Rows:    [][]string{},
			},
			err:      "",
			expected: "| ColA | ColB |\n|------|------|\n",
		},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg.Writer = &buf

			err := Marshal(&cfg, tt.table)

			if tt.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, buf.String())
			}
		})
	}
}

// Test centerPad function
func TestCenterPad(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		width    int
		expected string
	}{
		{
			name:     "string already wider than width",
			s:        "abcdefgh",
			width:    4,
			expected: "abcdefgh",
		},
		{
			name:     "string same width",
			s:        "abcd",
			width:    4,
			expected: "abcd",
		},
		{
			name:     "odd total padding",
			s:        "ab",
			width:    5,
			expected: " ab  ", // total=3, left=1, right=2
		},
		{
			name:     "even total padding",
			s:        "a",
			width:    4,
			expected: " a  ", // total=3, left=1, right=2
		},
		{
			name:     "UTF-8 characters",
			s:        "你好",
			width:    6,
			expected: " 你好 ", // runewidth=4, total=2, left=1, right=1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := centerPad(tt.s, tt.width)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test isSeparatorLine function
func TestIsSeparatorLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "valid separator",
			line:     "|---|---|",
			expected: true,
		},
		{
			name:     "valid separator with spaces",
			line:     "| --- | --- |",
			expected: true,
		},
		{
			name:     "valid separator with alignment",
			line:     "|:---|---:|",
			expected: true,
		},
		{
			name:     "valid separator centered",
			line:     "|:---:|:---:|",
			expected: true,
		},
		{
			name:     "no pipes at start",
			line:     "---|---|",
			expected: false,
		},
		{
			name:     "no pipes at end",
			line:     "|---|---",
			expected: false,
		},
		{
			name:     "contains invalid characters",
			line:     "|---|abc|",
			expected: false,
		},
		{
			name:     "no dashes",
			line:     "|   |   |",
			expected: false,
		},
		{
			name:     "empty line",
			line:     "",
			expected: false,
		},
		{
			name:     "only pipes",
			line:     "|||",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSeparatorLine(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Unmarshal error paths
func TestUnmarshalErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name: "no separator after header",
			input: `|a|b|
|1|2|`,
			expectError: true,
			errorMsg:    "expected separator",
		},
		{
			name: "separator column mismatch",
			input: `|a|b|
|---|---|---|
|1|2|`,
			expectError: true,
			errorMsg:    "separator line has 3 columns, but header has 2",
		},
		{
			name: "data row column mismatch",
			input: `|a|b|
|---|---|
|1|`,
			expectError: true,
			errorMsg:    "data row has 1 columns, but header has 2",
		},
		{
			name: "no header found",
			input: `some text
more text`,
			expectError: true,
			errorMsg:    "no valid header row found",
		},
		{
			name: "no separator found",
			input: `|a|b|
|1|2|`,
			expectError: true,
			errorMsg:    "expected separator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"--from", "markdown", "--to", "markdown"}
			cfg, err := common.ParseConfig(args)
			assert.Nil(t, err)
			cfg.Reader = strings.NewReader(tt.input)

			var table common.Table
			err = Unmarshal(&cfg, &table)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test Unmarshal with nil table
func TestUnmarshalNilTable(t *testing.T) {
	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader("|a|b|\n|---|---|\n|1|2|")

	err = Unmarshal(&cfg, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "output table cannot be nil")
}

// Test Unmarshal with empty/whitespace input
func TestUnmarshalEmptyInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"empty", "", false},
		{"only whitespace", "\n\n", true},      // lineNumber > 0 but no table content
		{"only code fences", "```\n```", true}, // lineNumber > 0 but no table content
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"--from", "markdown", "--to", "markdown"}
			cfg, _ := common.ParseConfig(args)
			cfg.Reader = strings.NewReader(tt.input)

			var table common.Table
			err := Unmarshal(&cfg, &table)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				// Empty input should succeed with empty table
				assert.NoError(t, err)
			}
		})
	}
}

// Test Marshal with extension parameters
func TestMarshalExtensionParams(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1", "2"}},
	}

	tests := []struct {
		name     string
		options  map[string]string
		contains []string // strings that should be in output
	}{
		{
			name:     "pretty=false",
			options:  map[string]string{"pretty": "false"},
			contains: []string{"|A|B|", "|---|---|", "|1|2|"},
		},
		{
			name:     "bold-header=true",
			options:  map[string]string{"bold-header": "true"},
			contains: []string{"**A**", "**B**"},
		},
		{
			name:     "bold-first-column=true",
			options:  map[string]string{"bold-first-column": "true"},
			contains: []string{"**A**", "**1**"},
		},
		{
			name:     "escape=false",
			options:  map[string]string{"escape": "false"},
			contains: []string{"| A | B |"},
		},
		{
			name:     "align=right",
			options:  map[string]string{"align": "r"},
			contains: []string{"|--:|"},
		},
		{
			name:     "align=center",
			options:  map[string]string{"align": "c"},
			contains: []string{"|:--:|"},
		},
		{
			name:     "align=multiple",
			options:  map[string]string{"align": "l,r"},
			contains: []string{"|---|", "|--:|"},
		},
		{
			name:     "align=invalid",
			options:  map[string]string{"align": "x"},
			contains: []string{"|---|"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"--from", "markdown", "--to", "markdown"}
			cfg, err := common.ParseConfig(args)
			assert.Nil(t, err)

			// Set extension parameters
			for k, v := range tt.options {
				cfg.Extension[k] = v
			}

			var buf bytes.Buffer
			cfg.Writer = &buf

			err = Marshal(&cfg, table)
			assert.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected, "Output should contain '%s'", expected)
			}
		})
	}
}

// Test Marshal with UTF-8 characters
func TestMarshalUTF8(t *testing.T) {
	table := &common.Table{
		Headers: []string{"姓名", "年龄", "城市"},
		Rows: [][]string{
			{"张三", "25", "北京"},
			{"李四", "30", "上海"},
		},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "姓名")
	assert.Contains(t, output, "张三")
	assert.Contains(t, output, "北京")

	// Test round-trip
	cfg.Reader = strings.NewReader(output)
	var table2 common.Table
	err = Unmarshal(&cfg, &table2)
	assert.NoError(t, err)
	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)
}

// Test Unmarshal with code fences
func TestUnmarshalCodeFences(t *testing.T) {
	input := "```markdown\n| A | B |\n|---|---|\n| 1 | 2 |\n```"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"1", "2"}, table.Rows[0])
}

// Test Unmarshal with text before table
func TestUnmarshalTextBeforeTable(t *testing.T) {
	input := `Some introduction text
| A | B |
|---|---|
| 1 | 2 |
More text after`

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, table.Headers)
}

// Test Unmarshal with multiple tables (should stop at first table end)
func TestUnmarshalMultipleTables(t *testing.T) {
	input := `| A | B |
|---|---|
| 1 | 2 |
Some text
| C | D |
|---|---|
| 3 | 4 |`

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	// Should only parse the first table
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
}

// Test Marshal with special characters that need escaping
func TestMarshalEscaping(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A*B", "C_D", "E{F}", "G[H]", "I(J)", "K#L", "M+N", "O-P", "P.Q", "Q!R", "S~T"},
		Rows: [][]string{
			{"x*y", "u_v", "w{x}", "y[z]", "z(w)", "a#b", "c+d", "e-f", "g.h", "i!j", "k~l"},
		},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Check that special characters are escaped
	assert.Contains(t, output, "A\\*B")
	assert.Contains(t, output, "C\\_D")
	assert.Contains(t, output, "E\\{F\\}")
	assert.Contains(t, output, "G\\[H\\]")
	assert.Contains(t, output, "I\\(J\\)")
	assert.Contains(t, output, "K\\#L")
	assert.Contains(t, output, "M\\+N")
	assert.Contains(t, output, "O\\-P")
	assert.Contains(t, output, "P\\.Q")
	assert.Contains(t, output, "Q\\!R")
	assert.Contains(t, output, "S\\~T")
}

// Test Marshal with pretty=false and various options
func TestMarshalNoPretty(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1", "2"}},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Extension["pretty"] = "false"
	cfg.Extension["align"] = "c"

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// With pretty=false, alignment markers should be simple
	assert.Contains(t, output, "|:---:|")
}

// Test Marshal with empty rows and pretty mode
func TestMarshalEmptyRowsPretty(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	assert.Equal(t, "| A | B |\n|---|---|\n", output)
}

// Test Marshal with mixed length columns in pretty mode
func TestMarshalMixedLengthColumns(t *testing.T) {
	table := &common.Table{
		Headers: []string{"Short", "VeryLongHeader"},
		Rows: [][]string{
			{"a", "b"},
			{"c", "d"},
		},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Should have proper padding
	assert.Contains(t, output, "Short")
	assert.Contains(t, output, "VeryLongHeader")
}

// Test Marshal with right alignment and pretty mode
func TestMarshalRightAlignPretty(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1", "10"}},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Extension["align"] = "r"

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Right aligned should have colons on right
	assert.Contains(t, output, "|--:|")
}

// Test Marshal with center alignment and pretty mode
func TestMarshalCenterAlignPretty(t *testing.T) {
	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1", "10"}},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Extension["align"] = "c"

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Center aligned should have colons on both sides
	assert.Contains(t, output, "|:--:|")
}

// Test Unmarshal with escaped backslashes
func TestUnmarshalEscapedBackslashes(t *testing.T) {
	input := "| A |\n|---|\n| \\\\ |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"\\"}, table.Rows[0])
}

// Test Unmarshal with mixed escaped characters
func TestUnmarshalMixedEscaped(t *testing.T) {
	input := "| A | B | C |\n|---|---|---|\n| x\\|y | z\\*w | u\\_v |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(table.Rows))
	assert.Equal(t, []string{"x|y", "z*w", "u_v"}, table.Rows[0])
}

// Test Unmarshal with trailing content after table
func TestUnmarshalTrailingContent(t *testing.T) {
	input := `| A | B |
|---|---|
| 1 | 2 |
This is not part of the table
| C | D |
|---|---|
| 3 | 4 |`

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	// Should stop at first non-table line
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
}

// Test Unmarshal with empty cells in header
func TestUnmarshalEmptyHeaderCells(t *testing.T) {
	input := "| A |  | C |\n|---|---|---|\n| 1 | 2 | 3 |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "", "C"}, table.Headers)
}

// Test Unmarshal with whitespace in cells
func TestUnmarshalWhitespaceCells(t *testing.T) {
	input := "|  A  |  B  |\n|-----|-----|\n|  1  |  2  |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	// Whitespace should be trimmed
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, []string{"1", "2"}, table.Rows[0])
}

// Test Unmarshal with Windows line endings
func TestUnmarshalWindowsLineEndings(t *testing.T) {
	input := "| A | B |\r\n|---|---|\r\n| 1 | 2 |"

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, 1, len(table.Rows))
}

// Test Marshal with all extension parameters combined
func TestMarshalAllParams(t *testing.T) {
	table := &common.Table{
		Headers: []string{"Name", "Value"},
		Rows: [][]string{
			{"Test", "123"},
		},
	}

	args := []string{"--from", "markdown", "--to", "markdown"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Extension["pretty"] = "true"
	cfg.Extension["escape"] = "true"
	cfg.Extension["bold-header"] = "true"
	cfg.Extension["bold-first-column"] = "true"
	cfg.Extension["align"] = "l,r"

	var buf bytes.Buffer
	cfg.Writer = &buf

	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "**Name**") // bold header
	assert.Contains(t, output, "**Test**") // bold first column
	// Right align second column - with pretty mode and column width 5 for "Value" and 3 for "123"
	// The separator should be right-aligned for second column
	assert.Contains(t, output, "|----------:|") // separator with right align on second column
}

// Test round-trip with various inputs
func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		table *common.Table
	}{
		{
			name: "simple",
			table: &common.Table{
				Headers: []string{"A", "B"},
				Rows:    [][]string{{"1", "2"}},
			},
		},
		{
			name: "with special chars",
			table: &common.Table{
				Headers: []string{"A|B", "C*D"},
				Rows:    [][]string{{"x|y", "z*w"}},
			},
		},
		{
			name: "empty rows",
			table: &common.Table{
				Headers: []string{"A", "B"},
				Rows:    [][]string{},
			},
		},
		{
			name: "multiple rows",
			table: &common.Table{
				Headers: []string{"A", "B", "C"},
				Rows: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
					{"7", "8", "9"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			args := []string{"--from", "markdown", "--to", "markdown"}
			cfg, _ := common.ParseConfig(args)
			var buf bytes.Buffer
			cfg.Writer = &buf
			err := Marshal(&cfg, tt.table)
			assert.NoError(t, err)

			// Unmarshal
			cfg.Reader = strings.NewReader(buf.String())
			var table2 common.Table
			err = Unmarshal(&cfg, &table2)
			assert.NoError(t, err)

			// Compare
			assert.Equal(t, tt.table.Headers, table2.Headers)
			// Handle nil vs empty slice for rows
			if len(tt.table.Rows) == 0 {
				assert.True(t, len(table2.Rows) == 0 || table2.Rows == nil)
			} else {
				assert.Equal(t, tt.table.Rows, table2.Rows)
			}
		})
	}
}
