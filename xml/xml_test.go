package xml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			name: "marshal with XML declaration",
			cfg: &common.Config{
				Writer: &bytes.Buffer{},
				Extension: map[string]string{
					"minify":      "true",
					"declaration": "true",
				},
			},
			table: &common.Table{
				Headers: []string{"Name"},
				Rows: [][]string{
					{"Alice"},
				},
			},
			expectedXML: `<?xml version="1.0" encoding="UTF-8" ?>
<dataset><record><Name>Alice</Name></record></dataset>`,
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

func TestMarshalErrorCases(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		table := &common.Table{
			Headers: []string{"Name"},
			Rows:    [][]string{{"Alice"}},
		}
		err := Marshal(nil, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config cannot be nil")
	})

	t.Run("nil table", func(t *testing.T) {
		cfg := &common.Config{
			Writer: &bytes.Buffer{},
		}
		err := Marshal(cfg, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input table pointer cannot be nil")
	})

	t.Run("nil writer", func(t *testing.T) {
		cfg := &common.Config{
			Writer: nil,
		}
		table := &common.Table{
			Headers: []string{"Name"},
			Rows:    [][]string{{"Alice"}},
		}
		err := Marshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "writer is nil")
	})

	t.Run("mismatched row length", func(t *testing.T) {
		cfg := &common.Config{
			Writer: &bytes.Buffer{},
		}
		table := &common.Table{
			Headers: []string{"Name", "Age"},
			Rows: [][]string{
				{"Alice"}, // Only 1 column, but headers have 2
			},
		}
		err := Marshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "has 1 columns but headers have 2")
	})

	t.Run("invalid root element name", func(t *testing.T) {
		cfg := &common.Config{
			Writer: &bytes.Buffer{},
			Extension: map[string]string{
				"root-element": "123invalid", // starts with number
			},
		}
		table := &common.Table{
			Headers: []string{"Name"},
			Rows:    [][]string{{"Alice"}},
		}
		err := Marshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid root element name")
	})

	t.Run("invalid row element name", func(t *testing.T) {
		cfg := &common.Config{
			Writer: &bytes.Buffer{},
			Extension: map[string]string{
				"row-element": "invalid-name!", // contains invalid char
			},
		}
		table := &common.Table{
			Headers: []string{"Name"},
			Rows:    [][]string{{"Alice"}},
		}
		err := Marshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid row element name")
	})

	t.Run("invalid header name", func(t *testing.T) {
		cfg := &common.Config{
			Writer: &bytes.Buffer{},
		}
		table := &common.Table{
			Headers: []string{"Name", "Invalid@Header"}, // @ is invalid
			Rows:    [][]string{{"Alice", "Test"}},
		}
		err := Marshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid XML element name")
	})
}

func TestUnmarshalErrorCases(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		table := &common.Table{}
		err := Unmarshal(nil, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config cannot be nil")
	})

	t.Run("nil reader", func(t *testing.T) {
		cfg := &common.Config{
			Reader: nil,
		}
		table := &common.Table{}
		err := Unmarshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reader is nil")
	})

	t.Run("nil table", func(t *testing.T) {
		cfg := &common.Config{
			Reader: strings.NewReader("<dataset></dataset>"),
		}
		err := Unmarshal(cfg, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "output table cannot be nil")
	})

	t.Run("invalid xml", func(t *testing.T) {
		cfg := &common.Config{
			Reader: strings.NewReader("<invalid>"),
		}
		table := &common.Table{}
		err := Unmarshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode XML")
	})

	t.Run("mismatched column count is padded", func(t *testing.T) {
		xmlData := `
<dataset>
	<record>
		<name>John</name>
		<age>30</age>
	</record>
	<record>
		<name>Jane</name>
	</record>
</dataset>`
		cfg := &common.Config{
			Reader: strings.NewReader(xmlData),
		}
		table := &common.Table{}
		require.NoError(t, Unmarshal(cfg, table))
		assert.Equal(t, []string{"name", "age"}, table.Headers)
		assert.Equal(t, [][]string{{"John", "30"}, {"Jane", ""}}, table.Rows)
	})
}

func TestRoundTrip(t *testing.T) {
	// Test that we can marshal and then unmarshal back to the same data
	original := &common.Table{
		Headers: []string{"Name", "Age", "City"},
		Rows: [][]string{
			{"Alice", "30", "New York"},
			{"Bob", "25", "Los Angeles"},
		},
	}

	// Marshal to XML
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"minify": "true",
		},
	}
	err := Marshal(cfg, original)
	require.NoError(t, err)

	// Unmarshal back to table
	cfg.Reader = &buf
	result := &common.Table{}
	err = Unmarshal(cfg, result)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, original.Headers, result.Headers)
	assert.Equal(t, original.Rows, result.Rows)
}

func TestUTF8Handling(t *testing.T) {
	// Test UTF-8 characters in data
	table := &common.Table{
		Headers: []string{"Name", "Message"},
		Rows: [][]string{
			{"测试", "你好世界"},
			{"🎉", "emoji支持"},
		},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"minify": "true",
		},
	}

	err := Marshal(cfg, table)
	require.NoError(t, err)

	// Unmarshal and verify
	cfg.Reader = &buf
	result := &common.Table{}
	err = Unmarshal(cfg, result)
	require.NoError(t, err)

	assert.Equal(t, table.Headers, result.Headers)
	assert.Equal(t, table.Rows, result.Rows)
}

func TestSpecialCharactersInData(t *testing.T) {
	// Test special XML characters in data values
	table := &common.Table{
		Headers: []string{"Name", "Value"},
		Rows: [][]string{
			{"Test", "Value with <special> & characters"},
		},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"minify": "true",
		},
	}

	err := Marshal(cfg, table)
	require.NoError(t, err)

	// The XML encoder should handle escaping automatically
	output := buf.String()
	// Check that the output contains properly escaped XML
	assert.Contains(t, output, "<Value>")
	assert.Contains(t, output, "</Value>")
}

func TestUnmarshalNestedTextAndSkipsMetadata(t *testing.T) {
	xmlData := `<dataset><meta><name>ignored</name></meta><record><note>a <b>bold</b> c</note><age>30</age></record></dataset>`
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: strings.NewReader(xmlData)}, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"note", "age"}, table.Headers)
	assert.Equal(t, [][]string{{"a bold c", "30"}}, table.Rows)
}

func TestUnicodeElementNamesRoundTrip(t *testing.T) {
	table := &common.Table{
		Headers: []string{"姓名", "年龄"},
		Rows:    [][]string{{"张三", "30"}},
	}
	var buf bytes.Buffer
	require.NoError(t, Marshal(&common.Config{Writer: &buf}, table))

	got := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: bytes.NewReader(buf.Bytes())}, got))
	assert.Equal(t, table.Headers, got.Headers)
	assert.Equal(t, table.Rows, got.Rows)
}

func TestUnmarshalWithNestedElements(t *testing.T) {
	// Test that nested elements are handled correctly
	xmlData := `
<dataset>
	<record>
		<name>John</name>
		<details>
			<age>30</age>
		</details>
	</record>
</dataset>`
	cfg := &common.Config{
		Reader:    strings.NewReader(xmlData),
		Extension: map[string]string{},
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	// This should work but might flatten nested elements
	// The current implementation uses xml:",any" which captures all child elements
	assert.NoError(t, err)
	// Verify we get the expected structure
	assert.Equal(t, 2, len(table.Headers))
	assert.Equal(t, 1, len(table.Rows))
}

func TestMarshalWithEmptyHeaders(t *testing.T) {
	// Test that empty headers are handled as NULL
	table := &common.Table{
		Headers: []string{"", "Name", ""},
		Rows: [][]string{
			{"Alice", "Bob", "Charlie"},
		},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"minify": "true",
		},
	}

	err := Marshal(cfg, table)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "<NULL>Alice</NULL>")
	assert.Contains(t, output, "<Name>Bob</Name>")
	assert.Contains(t, output, "<NULL>Charlie</NULL>")
}

func TestUnmarshalWithEmptyData(t *testing.T) {
	// Test that empty data values are preserved
	xmlData := `
<dataset>
	<record>
		<name></name>
		<age>30</age>
	</record>
</dataset>`
	cfg := &common.Config{
		Reader:    strings.NewReader(xmlData),
		Extension: map[string]string{},
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	require.NoError(t, err)
	assert.Equal(t, []string{"name", "age"}, table.Headers)
	assert.Equal(t, [][]string{{"", "30"}}, table.Rows)
}

func TestUnmarshalAutoDetectRecords(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		headers []string
		rows    [][]string
	}{
		{
			name:    "generic item elements",
			input:   `<root><item><a>1</a><b>2</b></item><item><a>3</a><b>4</b></item></root>`,
			headers: []string{"a", "b"},
			rows:    [][]string{{"1", "2"}, {"3", "4"}},
		},
		{
			name:    "single item element",
			input:   `<root><item><a>1</a><b>2</b></item></root>`,
			headers: []string{"a", "b"},
			rows:    [][]string{{"1", "2"}},
		},
		{
			name:    "root fields",
			input:   `<root><a>1</a><b>2</b></root>`,
			headers: []string{"a", "b"},
			rows:    [][]string{{"1", "2"}},
		},
		{
			name:    "repeated scalar elements",
			input:   `<root><value>1</value><value>2</value></root>`,
			headers: []string{"value"},
			rows:    [][]string{{"1"}, {"2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &common.Table{}
			require.NoError(t, Unmarshal(&common.Config{Reader: strings.NewReader(tt.input)}, table))
			assert.Equal(t, tt.headers, table.Headers)
			assert.Equal(t, tt.rows, table.Rows)
		})
	}
}

func TestUnmarshalExplicitRowElementMissing(t *testing.T) {
	table := &common.Table{}
	cfg := &common.Config{
		Reader:    strings.NewReader(`<root><item><a>1</a></item></root>`),
		Extension: map[string]string{"row-element": "record"},
	}
	err := Unmarshal(cfg, table)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no <record> elements")
}

func TestUnmarshalAmbiguousRepeatedRows(t *testing.T) {
	table := &common.Table{}
	input := `<root><a><x>1</x></a><a><x>2</x></a><b><y>3</y></b><b><y>4</y></b></root>`
	err := Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "multiple repeated row elements")
}
