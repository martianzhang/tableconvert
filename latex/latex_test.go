package latex

import (
	"bytes"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalOptions(t *testing.T) {
	table := &common.Table{
		Headers: []string{"Name", "Age"},
		Rows:    [][]string{{"Alice", "30"}},
	}

	t.Run("text align and bold", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"text-align":        "c",
				"bold-first-row":    "true",
				"bold-first-column": "true",
			},
		}
		require.NoError(t, Marshal(cfg, table))
		out := buf.String()
		assert.Contains(t, out, `\begin{tabular}{cc}`)
		assert.Contains(t, out, `\textbf{Name}`)
		assert.Contains(t, out, `\textbf{Alice}`)
	})

	t.Run("caption label and float", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"caption":     "Results",
				"label":       "tab:results",
				"table-align": "raggedleft",
				"location":    "below",
			},
		}
		require.NoError(t, Marshal(cfg, table))
		out := buf.String()
		assert.Contains(t, out, `\begin{table}`)
		assert.Contains(t, out, `\raggedleft`)
		assert.Contains(t, out, `\caption{Results}`)
		assert.Contains(t, out, `\label{tab:results}`)
		assert.True(t, strings.Index(out, `\caption{Results}`) > strings.Index(out, `\end{tabular}`))
	})

	t.Run("mwe and no borders", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &common.Config{
			Writer: &buf,
			Extension: map[string]string{
				"mwe":     "true",
				"borders": "0000,0000",
			},
		}
		require.NoError(t, Marshal(cfg, table))
		out := buf.String()
		assert.Contains(t, out, `\documentclass{article}`)
		assert.Contains(t, out, `\begin{document}`)
		assert.NotContains(t, out, `\hline`)
	})
}

func TestUnmarshalColumnSpecsAndBooktabs(t *testing.T) {
	input := `\begin{tabular}{p{3cm}l}
\toprule
A & B \\
\midrule
1 & 2 \\
\bottomrule
\end{tabular}`
	table := &common.Table{}
	err := Unmarshal(&common.Config{Reader: strings.NewReader(input)}, table)
	require.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "2"}}, table.Rows)
}

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
