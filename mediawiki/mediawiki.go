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
	var currentHeader []string
	inTable := false
	headerCount := 0
	headerDone := false
	firstRowHeader := cfg.GetExtensionBool("first-row-header", false)
	lineNumber := 0

	flushRow := func(lineNo int) error {
		if len(currentRow) == 0 {
			return nil
		}
		if headerCount == 0 && firstRowHeader {
			table.Headers = append([]string(nil), currentRow...)
			headerCount = len(currentRow)
			headerDone = true
			currentRow = nil
			return nil
		}
		if len(currentRow) != headerCount {
			return fmt.Errorf("parse error on line %d: row has %d columns, but header has %d (line: %q)",
				lineNo, len(currentRow), headerCount, "")
		}
		table.Rows = append(table.Rows, currentRow)
		currentRow = nil
		return nil
	}

	finishHeader := func() {
		if headerDone {
			return
		}
		headerDone = true
		if len(currentHeader) > 0 {
			table.Headers = currentHeader
			headerCount = len(currentHeader)
		}
	}

	for _, line := range lines {
		lineNumber++
		line = strings.TrimSpace(line)

		// Check table start/end
		if strings.HasPrefix(line, "{|") {
			inTable = true
			continue
		} else if line == "|}" {
			if inTable {
				finishHeader()
				if err := flushRow(lineNumber); err != nil {
					return err
				}
			}
			inTable = false
			continue
		}

		if !inTable {
			continue
		}

		// Captions are not part of the table data.
		if strings.HasPrefix(line, "|+") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "|-"):
			// A row separator commits the previous row and ends the header
			// section if we were still reading it.
			finishHeader()
			if err := flushRow(lineNumber); err != nil {
				return err
			}
		case strings.HasPrefix(line, "!"):
			// MediaWiki allows each header cell on its own line. Before the
			// first data row, append to the header row; afterwards treat it as
			// a regular cell in the current row.
			cells := parseMediaWikiCells(line, true)
			if !headerDone {
				currentHeader = append(currentHeader, cells...)
			} else {
				currentRow = append(currentRow, cells...)
			}
		case strings.HasPrefix(line, "|"):
			cells := parseMediaWikiCells(line, false)
			currentRow = append(currentRow, cells...)
		}
	}

	if inTable {
		finishHeader()
		if err := flushRow(lineNumber); err != nil {
			return err
		}
	}

	return nil
}

// parseMediaWikiCells parses one MediaWiki header/data line. Cells on the
// same line use "!!" (headers) or "||" (data); otherwise each "|" starts a
// cell, which is how MediaWiki writes one cell per line.
func parseMediaWikiCells(line string, header bool) []string {
	content := line
	if header {
		content = strings.TrimPrefix(content, "!")
	} else {
		content = strings.TrimPrefix(content, "|")
	}
	separator := "||"
	if header {
		separator = "!!"
	}
	parts := strings.Split(content, separator)
	if len(parts) == 1 {
		parts = strings.Split(content, "|")
		// The leading marker was already removed for headers, but data lines
		// start with "|" and the first split element is empty.
		if !header && len(parts) > 0 && parts[0] == "" {
			parts = parts[1:]
		}
	}
	cells := make([]string, 0, len(parts))
	for _, cell := range parts {
		cell = strings.TrimSpace(cell)
		cell = stripMediaWikiAttributes(cell)
		cell = strings.ReplaceAll(cell, "{{!}}", "|")
		cell = strings.ReplaceAll(cell, "&#33;", "!")
		cells = append(cells, cell)
	}
	return cells
}

func stripMediaWikiAttributes(cell string) string {
	if idx := strings.Index(cell, "|"); idx >= 0 {
		attrs := strings.TrimSpace(cell[:idx])
		if strings.Contains(attrs, "=") {
			return strings.TrimSpace(cell[idx+1:])
		}
	}
	return cell
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
	minify := cfg.GetExtensionBool("minify", false)
	sortable := cfg.GetExtensionBool("sort", false)

	class := "wikitable"
	if sortable {
		class += " sortable"
	}
	// Write table start with default wikitable class
	writer.Write([]byte("{| class=\"" + class + "\"\n"))

	// Write headers (MediaWiki uses ! for headers)
	if len(table.Headers) > 0 {
		headerLine := "! "
		cellSeparator := " !! "
		if minify {
			headerLine = "!"
			cellSeparator = "!!"
		}
		for i, header := range table.Headers {
			if i > 0 {
				headerLine += cellSeparator
			}
			// Escape header separators and pipes in header content.
			escapedHeader := strings.ReplaceAll(header, "!", "&#33;")
			escapedHeader = strings.ReplaceAll(escapedHeader, "|", "{{!}}")
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
		cellSeparator := " || "
		if minify {
			rowLine = "|"
			cellSeparator = "||"
		}
		for j, cell := range row {
			if j > 0 {
				rowLine += cellSeparator
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
