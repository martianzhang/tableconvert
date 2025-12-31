package mediawiki

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

	// Read all content from reader
	content, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return fmt.Errorf("Unmarshal: failed to read input: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	table.Headers = []string{}
	table.Rows = [][]string{}

	var currentRow []string
	inTable := false
	headerCount := 0
	lineNumber := 0

	for _, line := range lines {
		lineNumber++
		line = strings.TrimSpace(line)

		// Check table start/end
		if strings.HasPrefix(line, "{|") {
			inTable = true
			continue
		} else if line == "|}" {
			inTable = false
			// Add the last row if exists
			if len(currentRow) > 0 {
				// Validate column count
				if len(currentRow) != headerCount {
					return fmt.Errorf("parse error on line %d: row has %d columns, but header has %d (line: %q)",
						lineNumber, len(currentRow), headerCount, line)
				}
				table.Rows = append(table.Rows, currentRow)
				currentRow = nil
			}
			continue
		}

		if !inTable {
			continue
		}

		// Process table content
		if strings.HasPrefix(line, "|-") {
			// Row separator - commit current row
			if len(currentRow) > 0 {
				// Validate column count
				if len(currentRow) != headerCount {
					return fmt.Errorf("parse error on line %d: row has %d columns, but header has %d (line: %q)",
						lineNumber, len(currentRow), headerCount, line)
				}
				table.Rows = append(table.Rows, currentRow)
				currentRow = nil
			}
		} else if strings.HasPrefix(line, "!") {
			// Header row - handle both !! and | separators
			headerLine := strings.TrimPrefix(line, "!")
			headerCells := strings.Split(headerLine, "!!")
			if len(headerCells) == 1 {
				headerCells = strings.Split(headerLine, "|")
			}
			for i, cell := range headerCells {
				cell = strings.TrimSpace(cell)
				// Handle {{!}} template
				cell = strings.ReplaceAll(cell, "{{!}}", "|")
				headerCells[i] = cell
			}
			table.Headers = headerCells
			headerCount = len(headerCells)
		} else if strings.HasPrefix(line, "|") {
			// Data row - handle both || and | separators
			dataLine := strings.TrimPrefix(line, "|")
			cells := strings.Split(dataLine, "||")
			if len(cells) == 1 {
				cells = strings.Split(dataLine, "|")
			}
			for i, cell := range cells {
				cell = strings.TrimSpace(cell)
				// Handle {{!}} template
				cell = strings.ReplaceAll(cell, "{{!}}", "|")
				cells[i] = cell
			}
			currentRow = cells
		}
	}

	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}
	if cfg.Writer == nil {
		return fmt.Errorf("Marshal: config writer cannot be nil")
	}

	// Validate table structure
	columnCounts := len(table.Headers)
	if columnCounts == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}
	for i, row := range table.Rows {
		if len(row) != columnCounts {
			return fmt.Errorf("Marshal: row %d has %d columns, but table has %d", i, len(row), columnCounts)
		}
	}

	writer := cfg.Writer

	// Write table start with default wikitable class
	writer.Write([]byte("{| class=\"wikitable\"\n"))

	// Write headers (MediaWiki uses ! for headers)
	if len(table.Headers) > 0 {
		headerLine := "! "
		for i, header := range table.Headers {
			if i > 0 {
				headerLine += " !! "
			}
			// Escape pipes in header content using {{!}} template
			escapedHeader := strings.ReplaceAll(header, "|", "{{!}}")
			headerLine += escapedHeader
		}
		writer.Write([]byte(headerLine + "\n"))
	}

	// Write data rows
	for i, row := range table.Rows {
		// Add row separator before each row
		if i == 0 {
			// First row needs separator if there are headers
			if len(table.Headers) > 0 {
				writer.Write([]byte("|-\n"))
			}
		} else {
			// Subsequent rows always need separator
			writer.Write([]byte("|-\n"))
		}

		// Write row data (MediaWiki uses | for data cells)
		rowLine := "| "
		for j, cell := range row {
			if j > 0 {
				rowLine += " || "
			}
			// Escape pipes in cell content using {{!}} template
			escapedCell := strings.ReplaceAll(cell, "|", "{{!}}")
			rowLine += escapedCell
		}
		writer.Write([]byte(rowLine + "\n"))
	}

	// Write table end
	writer.Write([]byte("|}\n"))

	return nil
}
