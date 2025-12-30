package ascii

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"github.com/mattn/go-runewidth"
)

const ASCIIDefaultStyle = "box"

var tableStyleMap = map[string]string{
	"dot":    "·",
	"bubble": "◌",
	"plus":   "+",
}

func Unmarshal(cfg *common.Config, table *common.Table) error {
	style := cfg.GetExtensionString("style", ASCIIDefaultStyle)

	switch style {
	case ASCIIDefaultStyle:
		return boxUnmarshal(cfg, table)
	default:
		return omniUnmarshal(cfg, table)
	}
}

func omniUnmarshal(cfg *common.Config, table *common.Table) error {
	style := cfg.GetExtensionString("style", "dot")
	if len(style) != 1 {
		if v, ok := tableStyleMap[style]; ok {
			style = v
		} else {
			return fmt.Errorf("unknown style: %s", style)
		}
	}

	scanner := bufio.NewScanner(cfg.Reader)
	lineNumber := 0
	var headers []string
	var rows [][]string
	parsingState := "start" // states: start, header, header_separator, data, end

	done := false
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			continue // Skip empty lines
		}

		// Inline separator line check
		isSeparator := func() bool {
			if !strings.HasPrefix(trimmedLine, style) || !strings.HasSuffix(trimmedLine, style) {
				return false
			}
			for _, c := range trimmedLine {
				if c != rune(style[0]) && c != '+' {
					return false
				}
			}
			return true
		}

		// Inline data line check and parsing
		parseDataLine := func() []string {
			trimmed := strings.TrimPrefix(strings.TrimSuffix(trimmedLine, style), style)
			parts := strings.Split(trimmed, style)
			var cells []string
			for _, part := range parts {
				cells = append(cells, strings.TrimSpace(part))
			}
			return cells
		}

		switch parsingState {
		case "start":
			if isSeparator() {
				parsingState = "header"
			}
		case "header":
			if strings.HasPrefix(trimmedLine, style) && strings.HasSuffix(trimmedLine, style) {
				headers = parseDataLine()
				if len(headers) == 0 {
					return &common.ParseError{
						LineNumber: lineNumber,
						Message:    "failed to parse header line",
						Line:       line,
					}
				}
				parsingState = "header_separator"
			}
		case "header_separator":
			if isSeparator() {
				parsingState = "data"
			}
		case "data":
			if isSeparator() {
				parsingState = "end"
				done = true
			} else if strings.HasPrefix(trimmedLine, style) && strings.HasSuffix(trimmedLine, style) {
				row := parseDataLine()
				if len(row) != len(headers) {
					return &common.ParseError{
						LineNumber: lineNumber,
						Message:    fmt.Sprintf("row has %d columns, expected %d", len(row), len(headers)),
						Line:       line,
					}
				}
				rows = append(rows, row)
			}
		case "end":
			done = true
		}

		if done {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if parsingState != "end" {
		return fmt.Errorf("incomplete table data")
	}

	table.Headers = headers
	table.Rows = rows
	return nil
}

func boxUnmarshal(cfg *common.Config, table *common.Table) error {
	// isBorderLine checks if a line represents a table border (e.g., "+---+---+")
	isBorderLine := func(line string) bool {
		trimmed := strings.TrimSpace(line)
		return strings.HasPrefix(trimmed, "+") && strings.HasSuffix(trimmed, "+") && strings.Contains(trimmed, "-")
	}

	// isDataLine checks if a line contains table data (e.g., "| val1 | val2 |")
	isDataLine := func(line string) bool {
		trimmed := strings.TrimSpace(line)
		return strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|")
	}

	// splitAndTrim splits a data/header line by '|' and trims whitespace from each part.
	// It skips empty strings resulting from leading/trailing '|' characters.
	splitAndTrim := func(line string) []string {
		parts := strings.Split(line, "|")
		if len(parts) < 2 { // Must have at least leading and trailing '|'
			return []string{}
		}
		// Exclude first and last empty strings from leading/trailing '|'
		relevantParts := parts[1 : len(parts)-1]
		result := make([]string, len(relevantParts))
		for i, part := range relevantParts {
			result[i] = strings.TrimSpace(part)
		}
		return result
	}

	scanner := bufio.NewScanner(cfg.Reader)
	lineNumber := 0
	var headers []string
	var rows [][]string
	parsingState := "start" // Possible states: start, header, header_separator, data, end

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			continue // Skip empty lines
		}

		done := false

		switch parsingState {
		case "start":
			if isBorderLine(line) {
				parsingState = "header" // Found top border, expect header next
			} else {
				continue
			}
		case "header":
			if isDataLine(line) {
				headers = splitAndTrim(line)
				if len(headers) == 0 {
					return &common.ParseError{LineNumber: lineNumber, Message: "failed to parse header line", Line: line}
				}
				parsingState = "header_separator" // Expect separator line after header
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected header data line (| Header |)", Line: line}
			}
		case "header_separator":
			if isBorderLine(line) {
				parsingState = "data" // Found separator, expect data rows next
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected header separator line (+--+)", Line: line}
			}
		case "data":
			if isDataLine(line) {
				rowData := splitAndTrim(line)
				if len(rowData) != len(headers) {
					// Note: Currently lenient about column count mismatch
					// Could add warning or make this an error
				}
				rows = append(rows, rowData)
				// Remain in 'data' state to process more rows
			} else if isBorderLine(line) {
				parsingState = "end" // Found bottom border
				done = true
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected data line (| Data |) or bottom border line (+--+)", Line: line}
			}
		case "end":
			done = true // Parsing complete, break out of the loop
		}

		if done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Validate final parsing state
	if parsingState != "end" {
		// Handle cases where input ends without proper termination
		if parsingState == "data" && len(rows) > 0 {
			// Lenient: Accept missing bottom border if we have data
		} else if parsingState == "header_separator" && len(headers) > 0 {
			// Lenient: Accept missing data if we have headers
		} else {
			return &common.ParseError{LineNumber: lineNumber, Message: fmt.Sprintf("input ended unexpectedly in state '%s', missing bottom border?", parsingState), Line: ""}
		}
	}

	table.Headers = headers
	table.Rows = rows
	return nil

}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}
	// --- Header Row ---
	columnCounts := len(table.Headers)
	if columnCounts == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}
	// Calculate column widths using runewidth for proper CJK support
	columnWidths := make([]int, columnCounts)
	for i, header := range table.Headers {
		columnWidths[i] = runewidth.StringWidth(header)
	}
	// Update widths based on row data
	for j, row := range table.Rows {
		if len(row) != columnCounts {
			return fmt.Errorf("Marshal: row %d has %d columns, but table header has %d columns", j, len(row), columnCounts)
		}
		for i, cell := range row {
			cellWidth := runewidth.StringWidth(cell)
			if cellWidth > columnWidths[i] {
				columnWidths[i] = cellWidth
			}
		}
	}
	writer := cfg.Writer

	// Table Style
	style := cfg.GetExtensionString("style", ASCIIDefaultStyle)
	if style != ASCIIDefaultStyle {
		if v, ok := tableStyleMap[style]; ok {
			style = v
		} else if len(style) != 1 {
			style = ASCIIDefaultStyle
		}
	}

	// Draw table
	switch style {
	case ASCIIDefaultStyle:
		separator := func() {
			for _, width := range columnWidths {
				fmt.Fprintf(writer, "+-%s-", strings.Repeat("-", width))
			}
			fmt.Fprintln(writer, "+")
		}

		separator() // Top separator

		// Write header row
		for i, header := range table.Headers {
			paddedHeader := runewidth.FillRight(header, columnWidths[i])
			fmt.Fprintf(writer, "| %s ", paddedHeader)
		}
		fmt.Fprintln(writer, "|")

		separator() // Middle separator

		// --- Data Rows ---
		for _, row := range table.Rows {
			for i, cell := range row {
				paddedCell := runewidth.FillRight(cell, columnWidths[i])
				fmt.Fprintf(writer, "| %s ", paddedCell)
			}
			fmt.Fprintln(writer, "|")
		}

		separator() // Bottom separator

	default:
		// --- Separator Row ---
		separator := func() {
			for _, width := range columnWidths {
				fmt.Fprintf(writer, "%s", strings.Repeat(style, width+3))
			}
			fmt.Fprintln(writer, style)
		}

		separator() // Top separator

		// Write header row
		for i, header := range table.Headers {
			paddedHeader := runewidth.FillRight(header, columnWidths[i])
			fmt.Fprintf(writer, "%s %s ", style, paddedHeader)
		}
		fmt.Fprintln(writer, style)

		separator() // Middle separator

		// --- Data Rows ---
		for _, row := range table.Rows {
			for i, cell := range row {
				paddedCell := runewidth.FillRight(cell, columnWidths[i])
				fmt.Fprintf(writer, "%s %s ", style, paddedCell)
			}
			fmt.Fprintln(writer, style)
		}

		separator() // Bottom separator

	}
	return nil // Success
}
