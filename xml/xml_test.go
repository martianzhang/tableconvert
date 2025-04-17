package xml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *common.Config
		table       *common.Table
		expectedXML string
		expectError bool
	}{
		{
			name: "successful marshal with default elements",
			cfg: &common.Config{
				Writer: &bytes.Buffer{},
				Extension: map[string]string{
					"minify": "true",
				},
			},
			table: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice", "30"},
					{"Bob", "25"},
				},
			},
			expectedXML: `<dataset><record><Name>Alice</Name><Age>30</Age></record><record><Name>Bob</Name><Age>25</Age></record></dataset>`,
		},
		{
			name: "successful marshal with custom elements",
			cfg: &common.Config{
				Writer: &bytes.Buffer{},
				Extension: map[string]string{
					"root-element": "people",
					"row-element":  "person",
					"minify":       "false",
				},
			},
			table: &common.Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]string{
					{"Alice", "30"},
				},
			},
			expectedXML: `<people>
  <person>
    <Name>Alice</Name>
    <Age>30</Age>
  </person>
</people>`,
		},
		{
			name: "successful marshal with empty headers",
			cfg: &common.Config{
				Writer: &bytes.Buffer{},
				Extension: map[string]string{
					"minify": "true",
				},
			},
			table: &common.Table{
				Headers: []string{"", "Age"},
				Rows: [][]string{
					{"Alice", "30"},
				},
			},
			expectedXML: `<dataset><record><NULL>Alice</NULL><Age>30</Age></record></dataset>`,
		},
		{
			name: "successful marshal with empty table",
			cfg: &common.Config{
				Writer: &bytes.Buffer{},
				Extension: map[string]string{
					"minify": "true",
				},
			},
			table: &common.Table{
				Headers: []string{"Name"},
				Rows:    [][]string{},
			},
			expectedXML: `<dataset></dataset>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Marshal(tt.cfg, tt.table)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if buf, ok := tt.cfg.Writer.(*bytes.Buffer); ok {
				actual := strings.TrimSpace(buf.String())
				expected := strings.TrimSpace(tt.expectedXML)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestUnmarshal_EmptyTable(t *testing.T) {
	xmlData := `<dataset></dataset>`
	cfg := &common.Config{
		Reader:    strings.NewReader(xmlData),
		Extension: make(map[string]string),
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.NoError(t, err)
	assert.Empty(t, table.Headers)
	assert.Empty(t, table.Rows)
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		xmlInput  string
		cfg       *common.Config
		wantTable *common.Table
		wantErr   bool
	}{
		{
			name: "basic xml with headers and rows",
			xmlInput: `
			<dataset>
				<record>
					<name>John</name>
					<age>30</age>
				</record>
				<record>
					<name>Jane</name>
					<age>25</age>
				</record>
			</dataset>`,
			cfg: &common.Config{
				Reader:    strings.NewReader(""),
				Extension: map[string]string{},
			},
			wantTable: &common.Table{
				Headers: []string{"name", "age"},
				Rows: [][]string{
					{"John", "30"},
					{"Jane", "25"},
				},
			},
			wantErr: false,
		},
		{
			name: "custom root and row elements",
			xmlInput: `
			<customroot>
				<customrow>
					<field1>Value1</field1>
					<field2>Value2</field2>
				</customrow>
			</customroot>`,
			cfg: &common.Config{
				Reader: strings.NewReader(""),
				Extension: map[string]string{
					"root-element": "customroot",
					"row-element":  "customrow",
				},
			},
			wantTable: &common.Table{
				Headers: []string{"field1", "field2"},
				Rows: [][]string{
					{"Value1", "Value2"},
				},
			},
			wantErr: false,
		},
		{
			name:     "empty xml",
			xmlInput: `<dataset></dataset>`,
			cfg: &common.Config{
				Reader:    strings.NewReader(""),
				Extension: map[string]string{},
			},
			wantTable: &common.Table{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup reader with actual XML content
			tt.cfg.Reader = strings.NewReader(tt.xmlInput)

			table := &common.Table{}
			err := Unmarshal(tt.cfg, table)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTable.Headers, table.Headers)
				assert.Equal(t, tt.wantTable.Rows, table.Rows)
			}
		})
	}
}
