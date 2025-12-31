package tmpl

import (
	"bytes"
	"os"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	// Create a temporary template file for testing
	tmpfile, err := os.CreateTemp("", "template.*.tmpl")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write simple template content
	templateContent := `{{range .Rows}}{{index . 0}} is {{index . 1}} years old from {{index . 2}}
{{end}}`
	_, err = tmpfile.WriteString(templateContent)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	// Setup test table
	table := common.Table{
		Headers: []string{"name", "age", "city"},
		Rows: [][]string{
			{"Alice", "30", "New York"},
			{"Bob", "25", "Los Angeles"},
		},
	}

	var buf bytes.Buffer
	cfg := common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"template": tmpfile.Name(),
		},
	}

	// Execute Marshal
	err = Marshal(&cfg, &table)
	assert.NoError(t, err)

	// Verify output matches template processing
	expectedOutput := "Alice is 30 years old from New York\nBob is 25 years old from Los Angeles\n"
	assert.Equal(t, expectedOutput, buf.String())
}

func TestMarshalWithHeaders(t *testing.T) {
	// Create a temporary template file for testing headers access
	tmpfile, err := os.CreateTemp("", "template.*.tmpl")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write template that uses headers
	templateContent := `{{range .Headers}}{{.}}|{{end}}
{{range .Rows}}{{range .}}{{.}}|{{end}}
{{end}}`
	_, err = tmpfile.WriteString(templateContent)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	// Setup test table
	table := common.Table{
		Headers: []string{"Name", "Age"},
		Rows: [][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	}

	var buf bytes.Buffer
	cfg := common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"template": tmpfile.Name(),
		},
	}

	// Execute Marshal
	err = Marshal(&cfg, &table)
	assert.NoError(t, err)

	// Verify output
	expectedOutput := "Name|Age|\nAlice|30|\nBob|25|\n"
	assert.Equal(t, expectedOutput, buf.String())
}

func TestMarshalWithHelperFunctions(t *testing.T) {
	// Create a temporary template file for testing helper functions
	tmpfile, err := os.CreateTemp("", "template.*.tmpl")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write template that uses helper functions
	templateContent := `{{range .Rows}}{{Upper (index . 0)}} | {{Lower (index . 1)}} | {{Capitalize (index . 2)}}
{{end}}`
	_, err = tmpfile.WriteString(templateContent)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	// Setup test table
	table := common.Table{
		Headers: []string{"name", "status", "city"},
		Rows: [][]string{
			{"alice", "ACTIVE", "new york"},
		},
	}

	var buf bytes.Buffer
	cfg := common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"template": tmpfile.Name(),
		},
	}

	// Execute Marshal
	err = Marshal(&cfg, &table)
	assert.NoError(t, err)

	// Verify output
	expectedOutput := "ALICE | active | New York\n"
	assert.Equal(t, expectedOutput, buf.String())
}

func TestMarshalWithEscapeFunctions(t *testing.T) {
	// Create a temporary template file for testing escape functions
	tmpfile, err := os.CreateTemp("", "template.*.tmpl")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write template that uses escape functions
	templateContent := `{{MarkdownEscape (index .Rows 0 0)}}`
	_, err = tmpfile.WriteString(templateContent)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	// Setup test table with special characters
	table := common.Table{
		Headers: []string{"text"},
		Rows: [][]string{
			{"hello *world* _test_"},
		},
	}

	var buf bytes.Buffer
	cfg := common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"template": tmpfile.Name(),
		},
	}

	// Execute Marshal
	err = Marshal(&cfg, &table)
	assert.NoError(t, err)

	// Verify output - should escape markdown special chars
	expectedOutput := `hello \*world\* \_test\_`
	assert.Equal(t, expectedOutput, buf.String())
}

func TestMarshalErrorCases(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name"},
		Rows:    [][]string{{"test"}},
	}

	t.Run("nil table", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"template": "somefile.tmpl",
			},
		}
		err := Marshal(&cfg, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target table pointer cannot be nil")
	})

	t.Run("empty template path", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"template": "",
			},
		}
		err := Marshal(&cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template file path is required")
	})

	t.Run("non-existent template file", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"template": "/nonexistent/file.tmpl",
			},
		}
		err := Marshal(&cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read template file")
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "template.*.tmpl")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		// Write invalid template syntax
		_, err = tmpfile.WriteString("{{invalid syntax")
		require.NoError(t, err)
		require.NoError(t, tmpfile.Close())

		var buf bytes.Buffer
		cfg := common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"template": tmpfile.Name(),
			},
		}
		err = Marshal(&cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse template")
	})

	t.Run("template execution error", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "template.*.tmpl")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		// Write template that will cause execution error (accessing undefined field)
		_, err = tmpfile.WriteString("{{.UndefinedField}}")
		require.NoError(t, err)
		require.NoError(t, tmpfile.Close())

		var buf bytes.Buffer
		cfg := common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"template": tmpfile.Name(),
			},
		}
		err = Marshal(&cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute template")
	})
}

func TestMarshalEmptyTable(t *testing.T) {
	// Create a temporary template file
	tmpfile, err := os.CreateTemp("", "template.*.tmpl")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write template that handles empty rows
	templateContent := `{{range .Rows}}{{.}}{{end}}`
	_, err = tmpfile.WriteString(templateContent)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	// Setup empty table
	table := common.Table{
		Headers: []string{"name"},
		Rows:    [][]string{},
	}

	var buf bytes.Buffer
	cfg := common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"template": tmpfile.Name(),
		},
	}

	// Execute Marshal
	err = Marshal(&cfg, &table)
	assert.NoError(t, err)

	// Should produce empty output
	assert.Equal(t, "", buf.String())
}

func TestMarshalUTF8Handling(t *testing.T) {
	// Create a temporary template file
	tmpfile, err := os.CreateTemp("", "template.*.tmpl")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write template
	templateContent := `{{range .Rows}}{{index . 0}}: {{index . 1}}
{{end}}`
	_, err = tmpfile.WriteString(templateContent)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	// Setup table with UTF-8 characters
	table := common.Table{
		Headers: []string{"name", "message"},
		Rows: [][]string{
			{"测试", "你好世界"},
			{"🎉", "emoji支持"},
		},
	}

	var buf bytes.Buffer
	cfg := common.Config{
		Writer: &buf,
		Extension: map[string]string{
			"template": tmpfile.Name(),
		},
	}

	// Execute Marshal
	err = Marshal(&cfg, &table)
	assert.NoError(t, err)

	// Verify UTF-8 is preserved
	expectedOutput := "测试: 你好世界\n🎉: emoji支持\n"
	assert.Equal(t, expectedOutput, buf.String())
}
