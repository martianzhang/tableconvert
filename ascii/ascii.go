package ascii

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

// isBorderLine checks if a line is a table border (e.g., "+---+---+").
func isBorderLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "+") && strings.HasSuffix(trimmed, "+") && strings.Contains(trimmed, "-")
}

// isDataLine checks if a line contains table data (e.g., "| val1 | val2 |").
func isDataLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|")
}

// splitAndTrim splits a data/header line by '|' and trims whitespace from each part.
// It skips the empty strings resulting from the leading and trailing '|'.
func splitAndTrim(line string) []string {
	fmt.Println("#+ ", line)
	parts := strings.Split(line, "|")
	if len(parts) < 2 { // Should have at least '|' at start and end
		return []string{}
	}
	// Exclude the first and last empty strings caused by leading/trailing '|'
	relevantParts := parts[1 : len(parts)-1]
	result := make([]string, len(relevantParts))
	for i, part := range relevantParts {
		result[i] = strings.TrimSpace(part)
	}
	fmt.Println("#- ", result)
	return result
}

func Unmarshal(input io.Reader, table *common.Table) error {
	scanner := bufio.NewScanner(input)
	lineNumber := 0
	var headers []string
	var rows [][]string
	parsingState := "start" // states: start, header, header_separator, data, end

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			continue // Skip empty lines
		}

		switch parsingState {
		case "start":
			if isBorderLine(line) {
				parsingState = "header"
			} else {
				continue
			}
		case "header":
			if isDataLine(line) {
				headers = splitAndTrim(line)
				if len(headers) == 0 {
					return &common.ParseError{LineNumber: lineNumber, Message: "failed to parse header line", Line: line}
				}
				parsingState = "header_separator"
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected header data line (| Header |)", Line: line}
			}
		case "header_separator":
			if isBorderLine(line) {
				parsingState = "data"
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected header separator line (+--+)", Line: line}
			}
		case "data":
			if isDataLine(line) {
				rowData := splitAndTrim(line)
				if len(rowData) != len(headers) {
					// Allow parsing even if column count mismatches, but maybe log a warning?
					// For stricter parsing, return an error:
					// return nil, &ParseError{LineNumber: lineNumber, Message: fmt.Sprintf("data column count (%d) does not match header count (%d)", len(rowData), len(headers)), Line: line}
				}
				rows = append(rows, rowData)
				// Stay in 'data' state to parse more rows
			} else if isBorderLine(line) {
				parsingState = "end" // Found the bottom border
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected data line (| Data |) or bottom border line (+--+)", Line: line}
			}
		case "end":
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Final state check
	if parsingState != "end" {
		// This can happen if the input ends abruptly without a bottom border
		if parsingState == "data" && len(rows) > 0 {
			// Tolerate missing bottom border if we have data
			// fmt.Println("Warning: Input ended without a bottom border.")
		} else if parsingState == "header_separator" && len(headers) > 0 {
			// Tolerate missing data rows and bottom border if we only have headers
			// fmt.Println("Warning: Input contained only headers, no data rows or bottom border.")
		} else {
			return &common.ParseError{LineNumber: lineNumber, Message: fmt.Sprintf("input ended unexpectedly in state '%s', missing bottom border?", parsingState), Line: ""}
		}
	}

	table.Headers = headers
	table.Rows = rows
	return nil
}

func Marshal(table *common.Table, writer io.Writer) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}
	// --- Header Row ---
	columnCounts := len(table.Headers)
	if columnCounts == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}
	// Calculate column widths
	columnWidths := make([]int, columnCounts)
	for i, header := range table.Headers {
		columnWidths[i] = len(header)
	}
	// Update widths based on row data
	for j, row := range table.Rows {
		if len(row) != columnCounts {
			return fmt.Errorf("Marshal: %d row has %d columns, but table has %d", j, len(row), columnCounts)
		}
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}
	// --- Separator Row ---
	for _, width := range columnWidths {
		fmt.Fprintf(writer, "+-%s-", strings.Repeat("-", width))
	}
	fmt.Fprintln(writer, "+")
	// Write header row
	for i, header := range table.Headers {
		fmt.Fprintf(writer, "| %-*s ", columnWidths[i], header)
	}
	fmt.Fprintln(writer, "|")
	// --- Separator Row ---
	for _, width := range columnWidths {
		fmt.Fprintf(writer, "+-%s-", strings.Repeat("-", width))
	}
	fmt.Fprintln(writer, "+")
	// --- Data Rows ---
	for _, row := range table.Rows {
		for i, cell := range row {
			fmt.Fprintf(writer, "| %-*s ", columnWidths[i], cell)
		}
		fmt.Fprintln(writer, "|")
	}
	// --- Separator Row ---
	for _, width := range columnWidths {
		fmt.Fprintf(writer, "+-%s-", strings.Repeat("-", width))
	}
	fmt.Fprintln(writer, "+")
	return nil // Success
}
