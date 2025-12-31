package html

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

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
			err:      errors.New("Marshal: row 0 has 1 columns, but table has 2"),
		},
		{
			name: "basic table",
			table: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice", "25"},
					{"Bob", "30"},
				},
			},
			expected: `<table>
  <tr><th>Name</th><th>Age</th></tr>
  <tr><td>Alice</td><td>25</td></tr>
  <tr><td>Bob</td><td>30</td></tr>
</table>`,
			err: nil,
		},
		{
			name: "table with no rows",
			table: &common.Table{
				Headers: []string{"Col1", "Col2"},
				Rows:    [][]string{},
			},
			expected: `<table>
  <tr><th>Col1</th><th>Col2</th></tr>
</table>`,
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := &common.Config{Writer: &buf}
			err := Marshal(cfg, tt.table)

			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
				assert.Empty(t, buf.String())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, buf.String())
			}
		})
	}
}

func TestMarshalWithOptions(t *testing.T) {
	baseTable := &common.Table{
		Headers: []string{"Name", "Age"},
		Rows: [][]string{
			{"Alice", "25"},
			{"Bob", "30"},
		},
	}

	tests := []struct {
		name     string
		table    *common.Table
		config   map[string]string
		expected string
	}{
		{
			name:  "with escape",
			table: baseTable,
			config: map[string]string{
				"escape": "true",
			},
			expected: "<table>\n" +
				"  <tr><th>Name</th><th>Age</th></tr>\n" +
				"  <tr><td>Alice</td><td>25</td></tr>\n" +
				"  <tr><td>Bob</td><td>30</td></tr>\n" +
				"</table>",
		},
		{
			name: "with escape and special characters",
			table: &common.Table{
				Headers: []string{"Name", "Description"},
				Rows: [][]string{
					{"Test<1>", "Contains < and >"},
				},
			},
			config: map[string]string{
				"escape": "true",
			},
			expected: "<table>\n" +
				"  <tr><th>Name</th><th>Description</th></tr>\n" +
				"  <tr><td>Test&lt;1&gt;</td><td>Contains &lt; and &gt;</td></tr>\n" +
				"</table>",
		},
		{
			name:  "with minify",
			table: baseTable,
			config: map[string]string{
				"minify": "true",
			},
			expected: "<table><tr><th>Name</th><th>Age</th></tr><tr><td>Alice</td><td>25</td></tr><tr><td>Bob</td><td>30</td></tr></table>",
		},
		{
			name:  "with thead",
			table: baseTable,
			config: map[string]string{
				"thead": "true",
			},
			expected: `<table>
  <thead>
  <tr><th>Name</th><th>Age</th></tr>
  </thead>
  <tbody>
  <tr><td>Alice</td><td>25</td></tr>
  <tr><td>Bob</td><td>30</td></tr>
  </tbody>
</table>`,
		},
		{
			name:  "with div",
			table: baseTable,
			config: map[string]string{
				"div": "true",
			},
			expected: `<div class="table">
  <div class="tr"><div class="th">Name</div><div class="th">Age</div></div>
  <div class="tr"><div class="td">Alice</div><div class="td">25</div></div>
  <div class="tr"><div class="td">Bob</div><div class="td">30</div></div>
</div>`,
		},
		{
			name:  "with minify div",
			table: baseTable,
			config: map[string]string{
				"div":    "true",
				"minify": "true",
			},
			expected: `<div class="table"><div class="tr"><div class="th">Name</div><div class="th">Age</div></div><div class="tr"><div class="td">Alice</div><div class="td">25</div></div><div class="tr"><div class="td">Bob</div><div class="td">30</div></div></div>`,
		},
		{
			name:  "with multiple options",
			table: baseTable,
			config: map[string]string{
				"minify": "true",
				"thead":  "true",
			},
			expected: "<table><thead><tr><th>Name</th><th>Age</th></tr></thead><tbody><tr><td>Alice</td><td>25</td></tr><tr><td>Bob</td><td>30</td></tr></tbody></table>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := &common.Config{
				Writer:    &buf,
				Extension: tt.config,
			}
			err := Marshal(cfg, tt.table)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedTable *common.Table
		expectedError string
	}{
		{
			name:          "nil table",
			input:         "<table><tr><th>Name</th></tr></table>",
			expectedTable: nil,
			expectedError: "Unmarshal: output table cannot be nil",
		},
		{
			name:  "no table in content",
			input: "<div>Some content</div>",
			expectedTable: &common.Table{
				Headers: nil,
				Rows:    nil,
			},
			expectedError: "Unmarshal: no table found in HTML content",
		},
		{
			name:  "no header row",
			input: "<table></table>",
			expectedTable: &common.Table{
				Headers: nil,
				Rows:    nil,
			},
			expectedError: "Unmarshal: no header row found in table",
		},
		{
			name:  "no header cells",
			input: "<table><tr></tr></table>",
			expectedTable: &common.Table{
				Headers: nil,
				Rows:    nil,
			},
			expectedError: "Unmarshal: no header row found in table",
		},
		{
			name: "valid basic table",
			input: "<table>\n" +
				"  <tr><th>Name</th><th>Age</th></tr>\n" +
				"  <tr><td>Alice</td><td>25</td></tr>\n" +
				"  <tr><td>Bob</td><td>30</td></tr>\n" +
				"</table>",
			expectedTable: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice", "25"},
					{"Bob", "30"},
				},
			},
			expectedError: "",
		},
		{
			name: "table with attributes",
			input: `<table class="data-table">
                <tr><th>ID</th><th>Value</th></tr>
                <tr><td>1</td><td>First</td></tr>
            </table>`,
			expectedTable: &common.Table{
				Headers: []string{"ID", "Value"},
				Rows: [][]string{
					{"1", "First"},
				},
			},
			expectedError: "",
		},
		{
			name: "table with empty rows",
			input: "<table>\n" +
				"  <tr><th>Col1</th><th>Col2</th></tr>\n" +
				"  <tr></tr>\n" +
				"  <tr><td>Data1</td><td>Data2</td></tr>\n" +
				"</table>",
			expectedTable: &common.Table{
				Headers: []string{"Col1", "Col2"},
				Rows: [][]string{
					{"Data1", "Data2"},
				},
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "nil table" {
				err := Unmarshal(&common.Config{
					Reader: strings.NewReader(tt.input),
				}, nil)
				assert.EqualError(t, err, tt.expectedError)
				return
			}

			table := &common.Table{}
			cfg := &common.Config{
				Reader: strings.NewReader(tt.input),
			}

			err := Unmarshal(cfg, table)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTable.Headers, table.Headers)
				assert.Equal(t, tt.expectedTable.Rows, table.Rows)
			}
		})
	}
}

func TestUnmarshalFirstColumnHeader(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		extension     map[string]string
		expectedTable *common.Table
		expectedError string
	}{
		{
			name: "first-column-header with explicit extension",
			input: "<table>\n" +
				"  <tr><td>Name</td><td>Alice</td></tr>\n" +
				"  <tr><td>Age</td><td>25</td></tr>\n" +
				"</table>",
			extension: map[string]string{
				"first-column-header": "true",
			},
			expectedTable: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice"},
					{"25"},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header auto-detect with all caps",
			input: "<table>\n" +
				"  <tr><td>NAME</td><td>Alice</td></tr>\n" +
				"  <tr><td>AGE</td><td>25</td></tr>\n" +
				"</table>",
			extension: map[string]string{},
			expectedTable: &common.Table{
				Headers: []string{"NAME", "AGE"},
				Rows: [][]string{
					{"Alice"},
					{"25"},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header auto-detect with colon",
			input: "<table>\n" +
				"  <tr><td>Name:</td><td>Alice</td></tr>\n" +
				"  <tr><td>Age:</td><td>25</td></tr>\n" +
				"</table>",
			extension: map[string]string{},
			expectedTable: &common.Table{
				Headers: []string{"Name:", "Age:"},
				Rows: [][]string{
					{"Alice"},
					{"25"},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header with empty cells",
			input: "<table>\n" +
				"  <tr><td></td><td>Alice</td></tr>\n" +
				"  <tr><td>Age</td><td>25</td></tr>\n" +
				"</table>",
			extension: map[string]string{
				"first-column-header": "true",
			},
			expectedTable: &common.Table{
				Headers: []string{"", "Age"},
				Rows: [][]string{
					{"Alice"},
					{"25"},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header with only header column",
			input: "<table>\n" +
				"  <tr><td>Name</td></tr>\n" +
				"  <tr><td>Age</td></tr>\n" +
				"</table>",
			extension: map[string]string{
				"first-column-header": "true",
			},
			expectedTable: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{},
					{},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header with mixed td/th",
			input: "<table>\n" +
				"  <tr><th>Name</th><td>Alice</td></tr>\n" +
				"  <tr><th>Age</th><td>25</td></tr>\n" +
				"</table>",
			extension: map[string]string{
				"first-column-header": "true",
			},
			expectedTable: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice"},
					{"25"},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header with thead should not trigger",
			input: "<table>\n" +
				"  <thead>\n" +
				"    <tr><th>Name</th><th>Age</th></tr>\n" +
				"  </thead>\n" +
				"  <tbody>\n" +
				"    <tr><td>Alice</td><td>25</td></tr>\n" +
				"  </tbody>\n" +
				"</table>",
			extension: map[string]string{},
			expectedTable: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice", "25"},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header with th in first row should not trigger",
			input: "<table>\n" +
				"  <tr><th>Name</th><th>Age</th></tr>\n" +
				"  <tr><td>Alice</td><td>25</td></tr>\n" +
				"</table>",
			extension: map[string]string{},
			expectedTable: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice", "25"},
				},
			},
			expectedError: "",
		},
		{
			name: "first-column-header with no header-like content should not auto-detect",
			input: "<table>\n" +
				"  <tr><td>alice</td><td>25</td></tr>\n" +
				"  <tr><td>bob</td><td>30</td></tr>\n" +
				"</table>",
			extension: map[string]string{},
			expectedTable: &common.Table{
				Headers: nil,
				Rows:    nil,
			},
			expectedError: "Unmarshal: no header row found in table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &common.Table{}
			cfg := &common.Config{
				Reader:    strings.NewReader(tt.input),
				Extension: tt.extension,
			}

			err := Unmarshal(cfg, table)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTable.Headers, table.Headers)
				assert.Equal(t, tt.expectedTable.Rows, table.Rows)
			}
		})
	}
}

func TestCollectRows(t *testing.T) {
	// This tests the internal collectRows function indirectly through Unmarshal
	// collectRows is used by isFirstColumnHeader for auto-detection
	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name: "basic rows",
			input: "<table>\n" +
				"  <tr><td>A1</td><td>A2</td></tr>\n" +
				"  <tr><td>B1</td><td>B2</td></tr>\n" +
				"</table>",
			expected: [][]string{
				{"A1", "A2"},
				{"B1", "B2"},
			},
		},
		{
			name: "rows with empty cells",
			input: "<table>\n" +
				"  <tr><td></td><td>A2</td></tr>\n" +
				"  <tr><td>B1</td><td></td></tr>\n" +
				"</table>",
			expected: [][]string{
				{"", "A2"},
				{"B1", ""},
			},
		},
		{
			name: "rows with mixed td and th",
			input: "<table>\n" +
				"  <tr><th>A1</th><td>A2</td></tr>\n" +
				"  <tr><td>B1</td><th>B2</th></tr>\n" +
				"</table>",
			expected: [][]string{
				{"A1", "A2"},
				{"B1", "B2"},
			},
		},
		{
			name: "rows with whitespace",
			input: "<table>\n" +
				"  <tr><td>  A1  </td><td>  A2  </td></tr>\n" +
				"  <tr><td>  B1  </td><td>  B2  </td></tr>\n" +
				"</table>",
			expected: [][]string{
				{"A1", "A2"},
				{"B1", "B2"},
			},
		},
		{
			name: "rows with nested elements",
			input: "<table>\n" +
				"  <tr><td><span>A1</span></td><td><b>A2</b></td></tr>\n" +
				"  <tr><td>B1</td><td>B2</td></tr>\n" +
				"</table>",
			expected: [][]string{
				{"span", "b"},
				{"B1", "B2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to parse the HTML and get the table node
			doc, err := html.Parse(strings.NewReader(tt.input))
			assert.NoError(t, err)

			// Find the table node
			var tableNode *html.Node
			var findTable func(*html.Node)
			findTable = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "table" {
					tableNode = n
					return
				}
				for c := n.FirstChild; c != nil && tableNode == nil; c = c.NextSibling {
					findTable(c)
				}
			}
			findTable(doc)
			assert.NotNil(t, tableNode)

			// Call collectRows indirectly through isFirstColumnHeader
			// which uses collectRows internally
			result := collectRows(tableNode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsFirstColumnHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "table with thead - should be false",
			input:    "<table><thead><tr><th>Name</th></tr></thead><tr><td>Alice</td></tr></table>",
			expected: false,
		},
		{
			name:     "table with th in first row - should be false",
			input:    "<table><tr><th>Name</th><th>Age</th></tr><tr><td>Alice</td><td>25</td></tr></table>",
			expected: false,
		},
		{
			name:     "table with all caps first column - should be true",
			input:    "<table><tr><td>NAME</td><td>Alice</td></tr><tr><td>AGE</td><td>25</td></tr></table>",
			expected: true,
		},
		{
			name:     "table with colon in first column - should be true",
			input:    "<table><tr><td>Name:</td><td>Alice</td></tr><tr><td>Age:</td><td>25</td></tr></table>",
			expected: true,
		},
		{
			name:     "table with mixed header-like content - should be true",
			input:    "<table><tr><td>NAME:</td><td>Alice</td></tr><tr><td>AGE</td><td>25</td></tr><tr><td>City</td><td>NYC</td></tr></table>",
			expected: true,
		},
		{
			name:     "table with no header-like content - should be false",
			input:    "<table><tr><td>alice</td><td>25</td></tr><tr><td>bob</td><td>30</td></tr></table>",
			expected: false,
		},
		{
			name:     "empty table - should be false",
			input:    "<table></table>",
			expected: false,
		},
		{
			name:     "table with empty rows - should be false",
			input:    "<table><tr></tr><tr></tr></table>",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tt.input))
			assert.NoError(t, err)

			var tableNode *html.Node
			var findTable func(*html.Node)
			findTable = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "table" {
					tableNode = n
					return
				}
				for c := n.FirstChild; c != nil && tableNode == nil; c = c.NextSibling {
					findTable(c)
				}
			}
			findTable(doc)
			assert.NotNil(t, tableNode)

			result := isFirstColumnHeader(tableNode)
			assert.Equal(t, tt.expected, result)
		})
	}
}
