package latex

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	input := `\begin{tabular}{|l|c|r|}
\hline
Name & Age & City \\
\hline
Alice & 30 & New York \\
Bob & 25 & Los Angeles \\
\hline
\end{tabular}`

	args := []string{"--from", "latex", "--to", "latex"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)

	assert.Equal(t, []string{"Name", "Age", "City"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
	}, table.Rows)
}

func TestMarshal_MultiColumnTable(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a valid config with the buffer as writer
	cfg := &common.Config{
		Writer: &buf,
	}

	// Create a table with multiple columns
	table := &common.Table{
		Headers: []string{"Name", "Age", "City", "Country"},
		Rows: [][]string{
			{"Alice", "30", "New York", "USA"},
			{"Bob", "25", "Los Angeles", "USA"},
		},
	}

	// Call the Marshal function
	err := Marshal(cfg, table)

	// Verify no error occurred
	assert.NoError(t, err)

	// Verify the output contains the expected LaTeX tabular format
	expectedOutput := `\begin{tabular}{llll}
\hline
Name & Age & City & Country \\
\hline
Alice & 30 & New York & USA \\
\hline
Bob & 25 & Los Angeles & USA \\
\hline
\end{tabular}
`
	assert.Equal(t, expectedOutput, buf.String())
}

func TestUnmarshal_AmpersandInsideBraces(t *testing.T) {
	input := `\begin{tabular}{|l|c|}
\hline
Name & Info \\
\hline
Alice & {30 & New York} \\
Bob & {\textbf{25} & Los Angeles} \\
\hline
\end{tabular}`

	args := []string{"--from", "latex", "--to", "latex"}
	cfg, err := common.ParseConfig(args)
	assert.Nil(t, err)
	cfg.Reader = strings.NewReader(input)

	var table common.Table
	err = Unmarshal(&cfg, &table)
	assert.NoError(t, err)

	assert.Equal(t, []string{"Name", "Info"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Alice", "{30 & New York}"},            // Braces are preserved
		{"Bob", "{\\textbf{25} & Los Angeles}"}, // \textbf{} is preserved
	}, table.Rows)
}
