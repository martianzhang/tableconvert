package latex

import (
	"fmt"
	"io"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if cfg == nil || cfg.Reader == nil {
		return fmt.Errorf("Unmarshal: config or reader cannot be nil")
	}

	content, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return fmt.Errorf("Unmarshal: failed to read input: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	table.Headers = []string{}
	table.Rows = [][]string{}

	inTabular := false
	var headersProcessed bool
	var expectedColumns int

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "\\begin{tabular}") {
			inTabular = true
			// Parse column specification to determine expected column count
			start := strings.LastIndex(line, "{") + 1
			end := strings.LastIndex(line, "}")
			if start > 0 && end > start {
				colSpec := line[start:end]
				// Count the number of column specifiers (l, c, r etc.)
				// Remove | characters first
				cleanSpec := strings.ReplaceAll(colSpec, "|", "")
				expectedColumns = len(cleanSpec)
			}
			continue
		}

		if strings.HasPrefix(line, "\\end{tabular}") {
			inTabular = false
			continue
		}

		if !inTabular {
			continue
		}

		if strings.HasPrefix(line, "\\hline") || line == "" {
			continue
		}

		// Process table content
		cells := []string{}
		// Split the line properly, handling escaped & characters
		parts := splitLaTeXLine(line)
		for _, part := range parts {
			cell := strings.TrimSpace(part)
			// Handle LaTeX special cases
			if cell == "~" || cell == "\\textasciitilde{}" {
				cell = "" // Empty cell
			}
			cell = common.LaTeXUnescape(cell)
			cells = append(cells, cell)
		}

		// Pad or truncate cells to match expected column count
		if expectedColumns > 0 {
			if len(cells) > expectedColumns {
				cells = cells[:expectedColumns]
			} else if len(cells) < expectedColumns {
				for i := len(cells); i < expectedColumns; i++ {
					cells = append(cells, "")
				}
			}
		}

		if !headersProcessed && strings.Contains(line, "\\\\") {
			table.Headers = cells
			headersProcessed = true
			if expectedColumns == 0 {
				expectedColumns = len(cells)
			}
		} else {
			if len(cells) > 0 {
				table.Rows = append(table.Rows, cells)
			}
		}
	}

	return nil
}

// splitLaTeXLine properly splits a LaTeX table line into cells
func splitLaTeXLine(line string) []string {
	// Remove the line terminator if present
	line = strings.Split(line, "\\\\")[0]

	var cells []string
	var currentCell strings.Builder
	inEscape := false

	for _, r := range line {
		switch {
		case r == '\\':
			inEscape = true
			currentCell.WriteRune(r)
		case inEscape:
			inEscape = false
			currentCell.WriteRune(r)
		case r == '&' && !inEscape:
			cells = append(cells, currentCell.String())
			currentCell.Reset()
		default:
			currentCell.WriteRune(r)
		}
	}

	// Add the last cell
	if currentCell.Len() > 0 {
		cells = append(cells, currentCell.String())
	}

	return cells
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}
	if cfg == nil || cfg.Writer == nil {
		return fmt.Errorf("Marshal: config or writer cannot be nil")
	}

	writer := cfg.Writer

	// Begin LaTeX tabular environment
	// Default to left-aligned columns (l) for each column
	colSpec := strings.Repeat("l", len(table.Headers))
	if len(table.Headers) == 0 && len(table.Rows) > 0 {
		colSpec = strings.Repeat("l", len(table.Rows[0]))
	}
	writer.Write([]byte("\\begin{tabular}{" + colSpec + "}\n"))
	writer.Write([]byte("\\hline\n"))

	// Write headers if they exist
	if len(table.Headers) > 0 {
		headerLine := ""
		for i, header := range table.Headers {
			if i > 0 {
				headerLine += " & "
			}
			headerLine += common.LaTeXEscape(header)
		}
		headerLine += " \\\\\n\\hline\n"
		writer.Write([]byte(headerLine))
	}

	// Write data rows
	for _, row := range table.Rows {
		rowLine := ""
		for i, cell := range row {
			if i > 0 {
				rowLine += " & "
			}
			rowLine += common.LaTeXEscape(cell)
		}
		rowLine += " \\\\\n\\hline\n"
		writer.Write([]byte(rowLine))
	}

	// End LaTeX tabular environment
	writer.Write([]byte("\\end{tabular}\n"))

	return nil
}
