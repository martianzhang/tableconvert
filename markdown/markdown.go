package markdown

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"github.com/mattn/go-runewidth"
)

// centerPad centers s into a field of given display width using spaces.
func centerPad(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return s
	}
	total := width - w
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// parseLine parses one Markdown table row. Leading/trailing pipes are
// optional (GFM), and escaped pipes are not treated as separators.
func parseLine(line string, lineNumber int) ([]string, *common.ParseError) {
	trimmedLine := strings.TrimSpace(line)
	hadOuterPipe := strings.HasPrefix(trimmedLine, "|") || strings.HasSuffix(trimmedLine, "|")
	if strings.HasPrefix(trimmedLine, "|") {
		trimmedLine = trimmedLine[1:]
	}
	if hasUnescapedTrailingPipe(trimmedLine) {
		trimmedLine = trimmedLine[:len(trimmedLine)-1]
	}

	if !hadOuterPipe && !containsUnescapedPipe(trimmedLine) {
		return nil, &common.ParseError{
			LineNumber: lineNumber,
			Message:    "line does not contain an unescaped '|'",
			Line:       line,
		}
	}

	parts := splitMarkdownCells(trimmedLine)
	cells := make([]string, len(parts))
	for i, cell := range parts {
		cells[i] = common.MarkdownUnescape(strings.TrimSpace(cell))
	}
	return cells, nil
}

func containsUnescapedPipe(s string) bool {
	escaped := false
	for _, r := range s {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '|' {
			return true
		}
	}
	return false
}

func hasUnescapedTrailingPipe(s string) bool {
	if !strings.HasSuffix(s, "|") {
		return false
	}
	backslashes := 0
	for i := len(s) - 2; i >= 0 && s[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 0
}

// splitMarkdownCells splits on unescaped pipes but keeps escapes in the cell
// text for MarkdownUnescape.
func splitMarkdownCells(s string) []string {
	var parts []string
	var b strings.Builder
	escaped := false
	for _, r := range s {
		if escaped {
			b.WriteRune('\\')
			b.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '|' {
			parts = append(parts, b.String())
			b.Reset()
			continue
		}
		b.WriteRune(r)
	}
	if escaped {
		b.WriteRune('\\')
	}
	return append(parts, b.String())
}

// isSeparatorLine checks if a line looks like a Markdown table separator.
func isSeparatorLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "|") {
		trimmedLine = trimmedLine[1:]
	}
	if strings.HasSuffix(trimmedLine, "|") {
		trimmedLine = trimmedLine[:len(trimmedLine)-1]
	}
	if trimmedLine == "" {
		return false
	}
	// Check if the content between pipes consists only of '-', ':', '|' and spaces
	for _, r := range trimmedLine {
		if r != '-' && r != ':' && r != '|' && r != ' ' {
			return false
		}
	}
	// Ensure there's at least one '-' for it to be a valid separator
	return strings.Contains(trimmedLine, "-")
}

// Unmarshal parses Markdown table content from an io.Reader and populates the given Table struct.
// It expects the standard GitHub Flavored Markdown table format.
func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("output table cannot be nil")
	}

	// Reset the table fields to ensure clean population
	table.Headers = nil
	table.Rows = nil

	scanner := bufio.NewScanner(cfg.Reader)
	lineNumber := 0
	foundHeader := false
	foundSeparator := false
	headerCount := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines
		if trimmedLine == "" {
			continue
		}

		// Allow leading/trailing ``` fences for code blocks containing tables
		if strings.HasPrefix(trimmedLine, "```") {
			// If we've already found the separator, treat trailing fence as end of table
			if foundSeparator {
				break
			}
			// Otherwise, just skip the leading fence line
			continue
		}

		// --- State Machine: Header -> Separator -> Data Rows ---

		if !foundHeader {
			// --- Expecting Header ---
			if !containsUnescapedPipe(trimmedLine) {
				// Ignore lines before the first table row (could be text before table)
				continue
			}

			cells, err := parseLine(trimmedLine, lineNumber)
			if err != nil {
				// If parseLine fails on the potential header line, return the error
				return err
			}
			if len(cells) == 0 {
				return &common.ParseError{
					LineNumber: lineNumber,
					Message:    "header line contains no columns",
					Line:       line,
				}
			}
			table.Headers = cells // Assign the parsed cells
			headerCount = len(cells)
			foundHeader = true
			// continue // Don't need continue, next loop iteration handles separator check

		} else if !foundSeparator {
			// --- Expecting Separator ---
			if !isSeparatorLine(trimmedLine) {
				return &common.ParseError{
					LineNumber: lineNumber,
					Message:    "expected separator line (e.g., |---|---|) after header",
					Line:       line,
				}
			}
			// Validate separator column count (optional but good practice)
			// Use parseLine to count columns robustly, ignore content error here
			// as isSeparatorLine already validated the basic format.
			sepCells, _ := parseLine(trimmedLine, lineNumber)
			if len(sepCells) != headerCount {
				return &common.ParseError{
					LineNumber: lineNumber,
					Message:    fmt.Sprintf("separator line has %d columns, but header has %d", len(sepCells), headerCount),
					Line:       line,
				}
			}
			foundSeparator = true
			// continue // Don't need continue, next loop iteration handles data check

		} else {
			// --- Expecting Data Row or End ---
			if !containsUnescapedPipe(trimmedLine) {
				// Stop parsing data rows at the first non-table line.
				break
			}

			cells, parseErr := parseLine(trimmedLine, lineNumber)
			if parseErr != nil {
				break
			}
			if len(cells) != headerCount {
				return &common.ParseError{
					LineNumber: lineNumber,
					Message:    fmt.Sprintf("data row has %d columns, but header has %d", len(cells), headerCount),
					Line:       line,
				}
			}
			// Append the valid row. `append` handles the case where table.Rows is initially nil.
			table.Rows = append(table.Rows, cells)
		}
	} // end scanner loop

	// Check for scanning errors (e.g., IO errors)
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Final validation after processing all lines
	if !foundHeader {
		// If we never even found a line starting with '|', it wasn't a table.
		// Check if *any* content was processed. If lineNumber is 0, input was empty.
		// If lineNumber > 0 but no header, the content didn't contain a valid table start.
		if lineNumber > 0 {
			// We could return nil error here if non-table content is acceptable,
			// or an error if a table was strictly expected. Let's return an error.
			return fmt.Errorf("parsing failed: no valid header row found in input")
		}
		// If input was completely empty or only whitespace/fences, return success (empty table).
		return nil // Or return an error if an empty table isn't valid
	}
	if !foundSeparator {
		return fmt.Errorf("parsing failed: no separator row found after header (line %d)", lineNumber) // lineNumber might be slightly off if EOF hit
	}
	// It's okay if no data rows are found (table.Rows might be nil or empty)

	return nil // Success
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}

	// Validate table structure
	columnCounts := len(table.Headers)
	if columnCounts == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}
	for _, row := range table.Rows {
		if len(row) != columnCounts {
			return fmt.Errorf("Marshal: %d row has %d columns, but table has %d", len(table.Rows), len(row), columnCounts)
		}
	}

	// Get configuration with defaults
	align := cfg.GetExtensionString("align", "l")
	boldHeader := cfg.GetExtensionBool("bold-header", false)
	boldFirstColumn := cfg.GetExtensionBool("bold-first-column", false)
	escape := cfg.GetExtensionBool("escape", true)
	pretty := cfg.GetExtensionBool("pretty", true)

	writer := cfg.Writer
	headers := table.Headers

	// deal with align config: trim, default to "l", pad/truncate to headers length
	aligns := strings.Split(align, ",")
	for i := range aligns {
		aligns[i] = strings.ToLower(strings.TrimSpace(aligns[i]))
		switch aligns[i] {
		case "l", "c", "r":
			// valid
		default:
			aligns[i] = "l" // default fallback
		}
	}
	if len(aligns) < len(headers) {
		for len(aligns) < len(headers) {
			aligns = append(aligns, "l")
		}
	} else if len(aligns) > len(headers) {
		aligns = aligns[:len(headers)]
	}

	// Prepare display copies of headers and rows where escape/bold are applied
	displayHeaders := make([]string, len(headers))
	for i := range headers {
		h := headers[i]
		if escape {
			h = common.MarkdownEscape(h)
		}
		if boldHeader {
			h = fmt.Sprintf("**%s**", h)
		} else if boldFirstColumn && i == 0 {
			h = fmt.Sprintf("**%s**", h)
		}
		displayHeaders[i] = h
	}

	displayRows := make([][]string, len(table.Rows))
	for ri, row := range table.Rows {
		displayRows[ri] = make([]string, len(row))
		for i, cell := range row {
			c := cell
			if escape {
				c = common.MarkdownEscape(c)
			}
			if boldFirstColumn && i == 0 {
				c = fmt.Sprintf("**%s**", c)
			}
			displayRows[ri][i] = c
		}
	}

	// Calculate column widths for pretty mode using display widths
	columnWidths := make([]int, columnCounts)
	if pretty {
		for i, header := range displayHeaders {
			w := runewidth.StringWidth(header)
			if w > columnWidths[i] {
				columnWidths[i] = w
			}
		}
		for _, row := range displayRows {
			for i, cell := range row {
				w := runewidth.StringWidth(cell)
				if w > columnWidths[i] {
					columnWidths[i] = w
				}
			}
		}
	}

	// Note: headers/rows content for output are taken from displayHeaders/displayRows

	// --- Header Row ---
	headerRow := "|"
	for i := range displayHeaders {
		header := displayHeaders[i]
		alignChar := aligns[i]
		if pretty {
			var cellStr string
			switch alignChar {
			case "c":
				cellStr = centerPad(header, columnWidths[i])
			case "r":
				cellStr = runewidth.FillLeft(header, columnWidths[i])
			default:
				cellStr = runewidth.FillRight(header, columnWidths[i])
			}
			headerRow += fmt.Sprintf(" %s |", cellStr)
		} else {
			headerRow += header + "|"
		}
	}
	writer.Write([]byte(headerRow + "\n"))

	// --- Separator Row ---
	separator := "|"
	for i := 0; i < columnCounts; i++ {
		alignChar := aligns[i]
		var alignMarker string
		switch alignChar {
		case "c":
			if pretty {
				alignMarker = fmt.Sprintf(":%s-:", strings.Repeat("-", columnWidths[i]))
			} else {
				alignMarker = ":---:"
			}
		case "r":
			if pretty {
				alignMarker = fmt.Sprintf("-%s:", strings.Repeat("-", columnWidths[i]))
			} else {
				alignMarker = "---:"
			}
		default: // "l"
			if pretty {
				alignMarker = fmt.Sprintf("-%s-", strings.Repeat("-", columnWidths[i]))
			} else {
				alignMarker = "---"
			}
		}

		separator += alignMarker + "|"
	}
	writer.Write([]byte(separator + "\n"))

	// --- Data Rows ---
	for ri := range displayRows {
		row := displayRows[ri]
		// Write row
		rowLine := "|"
		for i, cell := range row {
			alignChar := aligns[i]
			if pretty {
				var cellStr string
				switch alignChar {
				case "c":
					cellStr = centerPad(cell, columnWidths[i])
				case "r":
					cellStr = runewidth.FillLeft(cell, columnWidths[i])
				default:
					cellStr = runewidth.FillRight(cell, columnWidths[i])
				}
				rowLine += fmt.Sprintf(" %s |", cellStr)
			} else {
				rowLine += cell + "|"
			}
		}
		writer.Write([]byte(rowLine + "\n"))
	}

	return nil
}
