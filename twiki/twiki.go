package twiki

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

// Unmarshal parses TWiki table content from an io.Reader and populates the given Table struct.
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
	headerCount := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines
		if trimmedLine == "" {
			continue
		}

		// Skip lines that don't start with '|' (potential non-table content)
		if !strings.HasPrefix(trimmedLine, "|") {
			// If we've already found the header, treat this as end of table
			if foundHeader {
				break
			}
			continue
		}

		// --- State Machine: Header -> Data Rows ---
		if !foundHeader {
			// --- Expecting Header ---
			cells, err := parseTWikiHeaderLine(trimmedLine, lineNumber)
			if err != nil {
				return err
			}
			if len(cells) == 0 {
				return &common.ParseError{
					LineNumber: lineNumber,
					Message:    "header line contains no columns",
					Line:       line,
				}
			}
			table.Headers = cells
			headerCount = len(cells)
			foundHeader = true
		} else {
			// --- Expecting Data Row ---
			cells, parseErr := parseLine(trimmedLine, lineNumber)
			if parseErr != nil {
				// If parseLine fails, treat it as end of table data
				break
			}
			if len(cells) != headerCount {
				return &common.ParseError{
					LineNumber: lineNumber,
					Message:    fmt.Sprintf("data row has %d columns, but header has %d", len(cells), headerCount),
					Line:       line,
				}
			}
			table.Rows = append(table.Rows, cells)
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Final validation
	if !foundHeader {
		if lineNumber > 0 {
			return fmt.Errorf("parsing failed: no valid header row found in input")
		}
		return nil
	}

	return nil
}

// parseTWikiHeaderLine parses a TWiki table header line (e.g. "|=Header1|=Header2=|")
func parseTWikiHeaderLine(line string, lineNumber int) ([]string, error) {
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return nil, &common.ParseError{
			LineNumber: lineNumber,
			Message:    "header line does not start and end with '|'",
			Line:       line,
		}
	}

	// Remove leading/trailing pipes and split
	content := strings.TrimSuffix(strings.TrimPrefix(line, "|"), "|")
	parts := strings.Split(content, "|")

	cells := make([]string, 0, len(parts))
	for _, part := range parts {
		// Skip empty parts (can happen with consecutive pipes)
		if part == "" {
			continue
		}
		// TWiki headers are wrapped in = signs
		// Check for = signs in the original part (before trimming spaces)
		if !strings.HasPrefix(part, "=") || !strings.HasSuffix(part, "=") {
			return nil, &common.ParseError{
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("header cell '%s' is not wrapped in '=' signs", part),
				Line:       line,
			}
		}
		// Remove = signs and trim spaces to get the actual content
		cellContent := strings.Trim(part, "=")
		cellContent = strings.TrimSpace(cellContent)
		cells = append(cells, cellContent)
	}

	return cells, nil
}

// parseLine parses a regular TWiki table line (either header or data row)
func parseLine(line string, lineNumber int) ([]string, error) {
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return nil, &common.ParseError{
			LineNumber: lineNumber,
			Message:    "line does not start and end with '|'",
			Line:       line,
		}
	}

	// Remove leading/trailing pipes and split
	content := strings.TrimSuffix(strings.TrimPrefix(line, "|"), "|")
	parts := strings.Split(content, "|")

	cells := make([]string, 0, len(parts))
	for _, part := range parts {
		cells = append(cells, strings.TrimSpace(part))
	}

	return cells, nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}

	writer := cfg.Writer

	// Write TWiki table header
	headerLine := "|"
	for _, header := range table.Headers {
		// Escape pipe characters in header content
		escapedHeader := strings.ReplaceAll(header, "|", "\\|")
		headerLine += "=" + escapedHeader + "=|"
	}
	headerLine += "\n"
	if _, err := writer.Write([]byte(headerLine)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data rows
	for _, row := range table.Rows {
		rowLine := "|"
		for _, cell := range row {
			// Escape pipe characters in cell content
			escapedCell := strings.ReplaceAll(cell, "|", "\\|")
			rowLine += escapedCell + "|"
		}
		rowLine += "\n"
		if _, err := writer.Write([]byte(rowLine)); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}
