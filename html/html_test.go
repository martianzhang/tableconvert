package html

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"
	"github.com/stretchr/testify/assert"
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
