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
	braceDepth := 0 // Track nested braces to handle & inside braces

	for _, r := range line {
		switch {
		case r == '\\':
			inEscape = true
			currentCell.WriteRune(r)
		case inEscape:
			inEscape = false
			currentCell.WriteRune(r)
		case r == '{':
			braceDepth++
			currentCell.WriteRune(r)
		case r == '}':
			if braceDepth > 0 {
				braceDepth--
			}
			currentCell.WriteRune(r)
		case r == '&' && !inEscape && braceDepth == 0: // Only split if not in escape and not inside braces
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

	var out strings.Builder
	escape := cfg.GetExtensionBool("escape", true)
	boldFirstColumn := cfg.GetExtensionBool("bold-first-column", false)
	boldFirstRow := cfg.GetExtensionBool("bold-first-row", false)
	mwe := cfg.GetExtensionBool("mwe", false)
	ht := cfg.GetExtensionBool("ht", false)
	caption := cfg.GetExtensionString("caption", "")
	label := cfg.GetExtensionString("label", "")
	location := strings.ToLower(cfg.GetExtensionString("location", "above"))
	textAlign := strings.ToLower(cfg.GetExtensionString("text-align", "l"))
	tableAlign := strings.ToLower(cfg.GetExtensionString("table-align", "centering"))
	borders := cfg.GetExtensionString("borders", "1111,1111")

	if textAlign != "l" && textAlign != "c" && textAlign != "r" {
		textAlign = "l"
	}
	if tableAlign != "centering" && tableAlign != "raggedleft" && tableAlign != "raggedright" {
		tableAlign = "centering"
	}

	// Border bits are treated as top, header separator, row separator, bottom.
	// Unknown patterns fall back to the default.
	bitString := strings.ReplaceAll(strings.ReplaceAll(borders, ",", ""), " ", "")
	showTop, showHeader, showRows, showBottom := true, true, true, true
	if len(bitString) >= 4 && strings.Trim(bitString, "01") == "" {
		showTop = bitString[0] == '1'
		showHeader = bitString[1] == '1'
		showRows = bitString[2] == '1'
		showBottom = bitString[3] == '1'
	}

	columns := len(table.Headers)
	if columns == 0 && len(table.Rows) > 0 {
		columns = len(table.Rows[0])
	}
	colSpec := strings.Repeat(textAlign, columns)

	formatCell := func(cell string, bold bool) string {
		if escape {
			cell = common.LaTeXEscape(cell)
		}
		if bold {
			cell = `\textbf{` + cell + `}`
		}
		return cell
	}

	needsFloat := caption != "" || label != "" || ht || mwe
	if _, ok := cfg.Extension["table-align"]; ok {
		needsFloat = true
	}

	if mwe {
		out.WriteString("\\documentclass{article}\n\\begin{document}\n")
	}
	if needsFloat {
		if ht {
			out.WriteString("\\begin{table}[ht]\n")
		} else {
			out.WriteString("\\begin{table}\n")
		}
		out.WriteString("\\" + tableAlign + "\n")
		if caption != "" && location != "below" {
			out.WriteString("\\caption{" + formatCaption(caption, escape) + "}\n")
		}
	}
	out.WriteString("\\begin{tabular}{" + colSpec + "}\n")
	if showTop {
		out.WriteString("\\hline\n")
	}

	if len(table.Headers) > 0 {
		headerLine := ""
		for i, header := range table.Headers {
			if i > 0 {
				headerLine += " & "
			}
			bold := boldFirstRow || (boldFirstColumn && i == 0)
			headerLine += formatCell(header, bold)
		}
		headerLine += " \\\\\n"
		out.WriteString(headerLine)
		if showHeader {
			out.WriteString("\\hline\n")
		}
	}

	for _, row := range table.Rows {
		rowLine := ""
		for i, cell := range row {
			if i > 0 {
				rowLine += " & "
			}
			rowLine += formatCell(cell, boldFirstColumn && i == 0)
		}
		rowLine += " \\\\\n"
		out.WriteString(rowLine)
		if showRows {
			out.WriteString("\\hline\n")
		}
	}

	if showBottom && len(table.Rows) > 0 && !showRows {
		out.WriteString("\\hline\n")
	}
	out.WriteString("\\end{tabular}\n")
	if needsFloat {
		if caption != "" && location == "below" {
			out.WriteString("\\caption{" + formatCaption(caption, escape) + "}\n")
		}
		if label != "" {
			out.WriteString("\\label{" + label + "}\n")
		}
		out.WriteString("\\end{table}\n")
	}
	if mwe {
		out.WriteString("\\end{document}\n")
	}

	_, err := cfg.Writer.Write([]byte(out.String()))
	return err
}

func formatCaption(caption string, escape bool) string {
	if escape {
		return common.LaTeXEscape(caption)
	}
	return caption
}
