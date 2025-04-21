package tmpl

import (
	"bytes"
	"os"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	// Create a temporary template file for testing
	tmpfile, err := os.CreateTemp("", "template.*.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write simple template content
	templateContent := `{{range .Rows}}{{index . 0}} is {{index . 1}} years old from {{index . 2}}
{{end}}`
	if _, err := tmpfile.WriteString(templateContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

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
