package mysql

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"github.com/mattn/go-runewidth"
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

// findAnchors finds the positions of '+' characters in the header line.
func findAnchors(line string) []int {
	line = strings.TrimSpace(line)
	var anchors []int
	for i, char := range line {
		if char == '+' {
			anchors = append(anchors, i)
		}
	}
	if len(anchors) < 2 { // Should have at least '+' at start and end
		return []int{}
	}
	return anchors
}

// parseFields extracts the field values from a MySQL/ascii style data line.
// It splits by '|' delimiters, which works correctly for UTF-8 content.
// Note: This approach does not support pipes within cell content, as that
// would break the MySQL table format anyway.
func parseFields(line string) []string {
	if !isDataLine(line) {
		return []string{}
	}
	trimmed := strings.TrimSpace(line)
	// Remove leading and trailing pipe
	trimmed = strings.TrimPrefix(trimmed, "|")
	trimmed = strings.TrimSuffix(trimmed, "|")
	raw := strings.Split(trimmed, "|")
	result := make([]string, 0, len(raw))
	for _, cell := range raw {
		result = append(result, strings.TrimSpace(cell))
	}
	return result
}

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Unmarshal: target table pointer cannot be nil")
	}

	scanner := bufio.NewScanner(cfg.Reader)
	lineNumber := 0
	var headers []string
	var rows [][]string
	var anchors []int
	parsingState := "start" // states: start, header, header_separator, data, end
	var expectLineLength int
	var preline string // Handles lines potentially split across buffer reads

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// --- Line Concatenation Logic (Handles lines split by buffer boundaries) ---
		// If we have a previous partial line...
		if len(preline) > 0 {
			line = preline + line // Concatenate
			// If the combined line is still shorter than expected, store it and continue
			if expectLineLength > 0 && len(line) < expectLineLength {
				preline = line
				continue
			} else {
				// Otherwise, we have the full line (or more), clear preline
				preline = ""
			}
		} else { // If no previous partial line...
			// If this line is shorter than expected (and we expect a certain length),
			// store it as a partial line and continue. Only do this after anchors are found.
			if expectLineLength > 0 && len(line) < expectLineLength {
				// Simple check: Does it look like a plausible start of a border/data line?
				// Avoid treating completely unrelated short lines as partial table lines.
				trimmedLookahead := strings.TrimSpace(line)
				if strings.HasPrefix(trimmedLookahead, "+-") || strings.HasPrefix(trimmedLookahead, "|") {
					preline = line
					continue
				}
				// If it doesn't look like part of the table, process it as is (might lead to error later)
			}
			// Otherwise, process the line as is
		}
		// --- End Line Concatenation Logic ---

		trimmedLine := strings.TrimSpace(line)

		// Skip effectively empty lines after potential concatenation
		if trimmedLine == "" {
			continue
		}

		switch parsingState {
		case "start":
			if isBorderLine(line) {
				parsingState = "header"
				anchors = findAnchors(line)
				if len(anchors) < 2 { // Need at least two '+' for one column
					return &common.ParseError{LineNumber: lineNumber, Message: "failed to parse header anchors (need at least two '+')", Line: line}
				}
				expectLineLength = anchors[len(anchors)-1] + 1 // Expect lines to reach the last '+'
			} else {
				// Allow skipping introductory lines before the first border
				// return &common.ParseError{LineNumber: lineNumber, Message: "expected top border line (+--+)", Line: line}
			}
		case "header":
			if isDataLine(line) {
				// Ensure line length matches expectation *before* parsing
				if len(line) < expectLineLength {
					// Should have been caught by concatenation logic, but double-check
					preline = line // Assume it's a partial line
					continue
				}
				headers = parseFields(line)
				if len(headers) == 0 {
					// Check if the line structure *looks* right but parsing failed
					if len(strings.Split(strings.Trim(line, "|"), "|")) >= 1 {
						// It looks like a data line, maybe parsing logic failed?
						// Or perhaps anchors were wrong? Let's assume header parse failure.
						return &common.ParseError{LineNumber: lineNumber, Message: "failed to parse header fields from data line", Line: line}
					}
					// If it doesn't even look like a data line (|...|)
					return &common.ParseError{LineNumber: lineNumber, Message: "expected header data line (| Header |)", Line: line}
				}
				parsingState = "header_separator"
			} else if isBorderLine(line) { // Handle case like +---+ \n +---+ (empty header)
				return &common.ParseError{LineNumber: lineNumber, Message: "expected header data line (| Header |), got another border", Line: line}
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected header data line (| Header |)", Line: line}
			}
		case "header_separator":
			if isBorderLine(line) {
				// Optional: Verify separator anchors match header anchors?
				// sepAnchors := findAnchors(line)
				// if !reflect.DeepEqual(anchors, sepAnchors) {
				//     return &common.ParseError{LineNumber: lineNumber, Message: "header separator anchors mismatch header anchors", Line: line}
				// }
				parsingState = "data"
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected header separator line (+--+)", Line: line}
			}
		case "data":
			// Ensure line length matches expectation *before* parsing
			if len(line) < expectLineLength {
				// Should have been caught by concatenation logic, but double-check
				// Check if it looks like a partial data line before assuming
				if isDataLine(line) || isBorderLine(line) {
					preline = line // Assume it's a partial line
					continue
				}
				// If it's short and doesn't look like table content, treat as error
				return &common.ParseError{LineNumber: lineNumber, Message: "unexpected short line", Line: line}
			}

			if isDataLine(line) {
				rowData := parseFields(line)
				// Check column count consistency (optional, but good practice)
				if len(rowData) != len(headers) {
					// Decide how strict to be. Log a warning or return an error.
					// Using fmt.Sprintf for better error message formatting.
					return &common.ParseError{
						LineNumber: lineNumber,
						Message:    fmt.Sprintf("data column count (%d) does not match header count (%d)", len(rowData), len(headers)),
						Line:       line,
					}
				}
				rows = append(rows, rowData)
				// Stay in 'data' state
			} else if isBorderLine(line) {
				parsingState = "end" // Found the bottom border
			} else {
				return &common.ParseError{LineNumber: lineNumber, Message: "expected data line (| Data |) or bottom border line (+--+)", Line: line}
			}
		case "end":
			// Any non-empty line after the final border is usually an error
			// (e.g., query summary like "1 row in set (0.00 sec)")
			// Allow trimming, but if anything remains, it's unexpected content.
			if trimmedLine != "" {
				// You might want to be more lenient here and just stop parsing successfully
				// return &common.ParseError{LineNumber: lineNumber, Message: "unexpected content after bottom border", Line: line}
				// Option: Break the loop instead of erroring out
				goto endLoop // Use goto to break out cleanly after state is 'end'
			}
			// Otherwise, just ignore trailing empty/whitespace lines
		}
	}
endLoop: // Label for the goto statement

	if err := scanner.Err(); err != nil {
		// If we have a partial line buffered, include it in the context
		if len(preline) > 0 {
			return fmt.Errorf("error reading input (last partial line: %q): %w", preline, err)
		}
		return fmt.Errorf("error reading input: %w", err)
	}

	// Handle potential partial line at EOF
	if len(preline) > 0 {
		// Treat leftover preline as an error - indicates incomplete input table
		return &common.ParseError{
			LineNumber: lineNumber,
			Message:    fmt.Sprintf("input ended with incomplete line in state '%s'", parsingState),
			Line:       preline,
		}
	}

	// Final state check - Did we reach a valid end state?
	if parsingState != "end" {
		// Allow ending in 'data' state (tolerates missing bottom border)
		if parsingState == "data" {
			// Successfully parsed headers and possibly data rows, missing final border is acceptable.
			// fmt.Println("Warning: Input ended without a bottom border.") // Optional warning
			// Allow ending in 'header_separator' if headers parsed but no rows (empty table)
		} else if parsingState == "header_separator" && len(headers) > 0 && len(rows) == 0 {
			// Successfully parsed headers, separator, but no data rows found before EOF.
			// This represents a valid empty table.
			// fmt.Println("Warning: Input contained only headers, no data rows or bottom border.") // Optional warning
		} else {
			// Any other state means the table is malformed or incomplete.
			return &common.ParseError{
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("input ended unexpectedly in state '%s', table possibly incomplete or malformed", parsingState),
				Line:       "", // No specific line to point to at EOF
			}
		}
	}

	// Populate the provided table struct
	table.Headers = headers
	table.Rows = rows

	return nil // Success
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}
	// --- Header Row ---
	columnCount := len(table.Headers)
	if columnCount == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}
	// Calculate column widths using runewidth for proper CJK support
	columnWidths := make([]int, columnCount)
	for i, header := range table.Headers {
		columnWidths[i] = runewidth.StringWidth(header)
	}
	// Update widths based on row data
	for j, row := range table.Rows {
		if len(row) != columnCount {
			return fmt.Errorf("Marshal: %d row has %d columns, but table has %d", j, len(row), columnCount)
		}
		for i, cell := range row {
			cellWidth := runewidth.StringWidth(cell)
			if cellWidth > columnWidths[i] {
				columnWidths[i] = cellWidth
			}
		}
	}

	writer := cfg.Writer
	// --- Separator Row ---
	for _, width := range columnWidths {
		fmt.Fprintf(writer, "+-%s-", strings.Repeat("-", width))
	}
	fmt.Fprintln(writer, "+")
	// Write header row
	for i, header := range table.Headers {
		paddedHeader := runewidth.FillRight(header, columnWidths[i])
		fmt.Fprintf(writer, "| %s ", paddedHeader)
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
			paddedCell := runewidth.FillRight(cell, columnWidths[i])
			fmt.Fprintf(writer, "| %s ", paddedCell)
		}
		fmt.Fprintln(writer, "|")
	}
	// --- Bottom Separator Row ---
	// Only output if there are rows, otherwise the header separator is the bottom
	if len(table.Rows) > 0 {
		for _, width := range columnWidths {
			fmt.Fprintf(writer, "+-%s-", strings.Repeat("-", width))
		}
		fmt.Fprintln(writer, "+")
	}
	return nil // Success
}
