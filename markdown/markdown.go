package markdown

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

// Helper function to parse a single line into cells
func parseLine(line string, lineNumber int) ([]string, *common.ParseError) {
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return nil, &common.ParseError{
			LineNumber: lineNumber,
			Message:    "line does not start and end with '|'",
			Line:       line,
		}
	}
	// Trim leading/trailing '|' and then split by '|'
	trimmedLine := strings.TrimPrefix(strings.TrimSuffix(line, "|"), "|")
	rawCells := strings.Split(trimmedLine, "|")
	cells := make([]string, 0, len(rawCells))
	for _, cell := range rawCells {
		// Trim whitespace from each cell
		cells = append(cells, strings.TrimSpace(cell))
	}
	return cells, nil
}

// isSeparatorLine checks if a line looks like a Markdown table separator.
func isSeparatorLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmedLine, "|") || !strings.HasSuffix(trimmedLine, "|") {
		return false
	}
	// Check if the content between pipes consists only of '-', ':', '|' and spaces
	inner := strings.Trim(trimmedLine, "|")
	for _, r := range inner {
		if r != '-' && r != ':' && r != '|' && r != ' ' {
			return false
		}
	}
	// Ensure there's at least one '-' for it to be a valid separator
	return strings.Contains(inner, "-")
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
			if !strings.HasPrefix(trimmedLine, "|") {
				// Ignore lines before the first '|' line (could be text before table)
				// A stricter parser might error here immediately.
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
			if !strings.HasPrefix(trimmedLine, "|") {
				// Stop parsing data rows if a line doesn't start with '|' after the separator.
				// This treats subsequent non-table lines as the end of the table block.
				break
			}

			cells, parseErr := parseLine(trimmedLine, lineNumber)
			if parseErr != nil {
				// If parseLine indicates it's just not a valid table row (e.g., missing pipes)
				// treat it as the end of the table data. Check the specific error if possible.
				// The placeholder parseLine returns a specific message for this.
				if strings.Contains(parseErr.Message, "does not start and end with '|'") {
					break // Assume end of table data
				}
				return parseErr // Propagate other actual parsing errors within a potential row
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
	// --- Header Row ---
	columnCounts := len(table.Headers)
	if columnCounts == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}

	writer := cfg.Writer
	// Write header row with pipes and alignment markers
	headerRow := "|"
	for _, header := range table.Headers {
		headerRow += header + "|" // Use the header text directly
	}
	writer.Write([]byte(headerRow + "\n")) // Write the header row with newline

	// --- Separator Row ---
	separatorRow := "|"
	for i := 0; i < columnCounts; i++ {
		separatorRow += "---|" // Use '---' for alignment markers
	}
	writer.Write([]byte(separatorRow + "\n")) // Write the separator row with newline

	// --- Data Rows ---
	for _, row := range table.Rows {
		if len(row) != columnCounts {
			return fmt.Errorf("Marshal: %d row has %d columns, but table has %d", len(table.Rows), len(row), columnCounts)
		}
	}
	for _, row := range table.Rows {
		dataRow := "|"
		for _, cell := range row {
			dataRow += cell + "|" // Use the cell text directly
		}
		writer.Write([]byte(dataRow + "\n")) // Write the data row with newline
	}
	return nil
}
