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

func TestUnmarshalWithDifferentDelimiters(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		input     string
		expected  *common.Table
	}{
		{
			name:      "tab delimiter",
			delimiter: "TAB",
			input:     "name\tage\tcity\nAlice\t30\tNew York\nBob\t25\tLos Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "semicolon delimiter",
			delimiter: "SEMICOLON",
			input:     "name;age;city\nAlice;30;New York\nBob;25;Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "pipe delimiter",
			delimiter: "PIPE",
			input:     "name|age|city\nAlice|30|New York\nBob|25|Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "slash delimiter",
			delimiter: "SLASH",
			input:     "name/age/city\nAlice/30/New York\nBob/25/Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "hash delimiter",
			delimiter: "HASH",
			input:     "name#age#city\nAlice#30#New York\nBob#25#Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "literal tab character",
			delimiter: "\t",
			input:     "name\tage\tcity\nAlice\t30\tNew York\nBob\t25\tLos Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "literal semicolon character",
			delimiter: ";",
			input:     "name;age;city\nAlice;30;New York\nBob;25;Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "literal pipe character",
			delimiter: "|",
			input:     "name|age|city\nAlice|30|New York\nBob|25|Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "literal slash character",
			delimiter: "/",
			input:     "name/age/city\nAlice/30/New York\nBob/25/Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:      "literal hash character",
			delimiter: "#",
			input:     "name#age#city\nAlice#30#New York\nBob#25#Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"--from", "csv", "--to", "csv", "--delimiter", tt.delimiter}
			cfg, err := common.ParseConfig(args)
			assert.NoError(t, err)

			cfg.Reader = strings.NewReader(tt.input)

			var table common.Table
			err = Unmarshal(&cfg, &table)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Headers, table.Headers)
			assert.Equal(t, tt.expected.Rows, table.Rows)
		})
	}
}

func TestMarshalWithDifferentDelimiters(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		expected  string
	}{
		{
			name:      "tab delimiter",
			delimiter: "TAB",
			expected:  "name\tage\tcity\nAlice\t30\tNew York\nBob\t25\tLos Angeles\n",
		},
		{
			name:      "semicolon delimiter",
			delimiter: "SEMICOLON",
			expected:  "name;age;city\nAlice;30;New York\nBob;25;Los Angeles\n",
		},
		{
			name:      "pipe delimiter",
			delimiter: "PIPE",
			expected:  "name|age|city\nAlice|30|New York\nBob|25|Los Angeles\n",
		},
		{
			name:      "slash delimiter",
			delimiter: "SLASH",
			expected:  "name/age/city\nAlice/30/New York\nBob/25/Los Angeles\n",
		},
		{
			name:      "hash delimiter",
			delimiter: "HASH",
			expected:  "name#age#city\nAlice#30#New York\nBob#25#Los Angeles\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			}

			args := []string{"--from", "csv", "--to", "csv", "--delimiter", tt.delimiter}
			cfg, err := common.ParseConfig(args)
			assert.NoError(t, err)

			var buf bytes.Buffer
			cfg.Writer = &buf
			err = Marshal(&cfg, table)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestUnmarshalWithFirstColumnHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *common.Table
	}{
		{
			name:  "standard case",
			input: "name,Alice,Bob\nage,30,25\ncity,New York,Los Angeles\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", "Los Angeles"},
				},
			},
		},
		{
			name:  "empty first column values",
			input: ",header1,header2\nrow1,val1,val2\nrow2,val3,val4\n",
			expected: &common.Table{
				Headers: []string{"", "row1", "row2"},
				Rows: [][]string{
					{"header1", "val1", "val3"},
					{"header2", "val2", "val4"},
				},
			},
		},
		{
			name:  "variable row lengths",
			input: "name,Alice,Bob\nage,30,25,extra\ncity,New York\n",
			expected: &common.Table{
				Headers: []string{"name", "age", "city"},
				Rows: [][]string{
					{"Alice", "30", "New York"},
					{"Bob", "25", ""},
					{"", "extra", ""},
				},
			},
		},
		{
			name:  "single column",
			input: "name\nAlice\nBob\n",
			expected: &common.Table{
				Headers: []string{"name"},
				Rows: [][]string{
					{"Alice"},
					{"Bob"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"--from", "csv", "--to", "csv", "--first-column-header"}
			cfg, err := common.ParseConfig(args)
			assert.NoError(t, err)

			cfg.Reader = strings.NewReader(tt.input)

			var table common.Table
			err = Unmarshal(&cfg, &table)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Headers, table.Headers)
			assert.Equal(t, tt.expected.Rows, table.Rows)
		})
	}
}

func TestMarshalWithBOM(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name", "age"},
		Rows: [][]string{
			{"Alice", "30"},
		},
	}

	args := []string{"--from", "csv", "--to", "csv", "--bom"}
	cfg, err := common.ParseConfig(args)
	assert.NoError(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf
	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	output := buf.String()
	// Check for UTF-8 BOM (0xEF 0xBB 0xBF)
	assert.True(t, len(output) >= 3)
	assert.Equal(t, byte(0xEF), output[0])
	assert.Equal(t, byte(0xBB), output[1])
	assert.Equal(t, byte(0xBF), output[2])

	// The rest should be the CSV content
	csvContent := output[3:]
	expectedCSV := "name,age\nAlice,30\n"
	assert.Equal(t, expectedCSV, csvContent)
}

func TestUnmarshalEmptyCSV(t *testing.T) {
	args := []string{"--from", "csv", "--to", "csv"}
	cfg, err := common.ParseConfig(args)
	assert.NoError(t, err)

	cfg.Reader = strings.NewReader("")

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty CSV file")
}

func TestUnmarshalWithQuotedValues(t *testing.T) {
	input := `"name","age","city"
"Alice","30","New York"
"Bob","25","Los Angeles, CA"`

	args := []string{"--from", "csv", "--to", "csv"}
	cfg, err := common.ParseConfig(args)
	assert.NoError(t, err)

	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)

	assert.Equal(t, []string{"name", "age", "city"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles, CA"},
	}, table.Rows)
}

func TestMarshalWithEmptyRows(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name", "age"},
		Rows: [][]string{
			{"Alice", "30"},
			{"", ""},
			{"Bob", "25"},
		},
	}

	args := []string{"--from", "csv", "--to", "csv"}
	cfg, err := common.ParseConfig(args)
	assert.NoError(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf
	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	expectedCSV := "name,age\nAlice,30\n,\nBob,25\n"
	assert.Equal(t, expectedCSV, buf.String())
}

func TestMarshalWithSpecialCharacters(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name", "description"},
		Rows: [][]string{
			{"Alice", "Likes, coffee"},
			{"Bob", "Hates; semicolons"},
		},
	}

	args := []string{"--from", "csv", "--to", "csv"}
	cfg, err := common.ParseConfig(args)
	assert.NoError(t, err)

	var buf bytes.Buffer
	cfg.Writer = &buf
	err = Marshal(&cfg, table)
	assert.NoError(t, err)

	// CSV should quote values containing commas
	output := buf.String()
	assert.Contains(t, output, `"Likes, coffee"`)
	// Note: semicolons don't require quoting in standard CSV, only commas do
	// So "Hates; semicolons" may not be quoted
	assert.Contains(t, output, "Hates; semicolons")
}
