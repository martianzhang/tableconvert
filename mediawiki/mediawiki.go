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

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check table start/end
		if strings.HasPrefix(line, "{|") {
			inTable = true
			continue
		} else if line == "|}" {
			inTable = false
			// Add the last row if exists
			if len(currentRow) > 0 {
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
				headerCells[i] = strings.TrimSpace(cell)
			}
			table.Headers = headerCells
		} else if strings.HasPrefix(line, "|") {
			// Data row - handle both || and | separators
			dataLine := strings.TrimPrefix(line, "|")
			cells := strings.Split(dataLine, "||")
			if len(cells) == 1 {
				cells = strings.Split(dataLine, "|")
			}
			for i, cell := range cells {
				cells[i] = strings.TrimSpace(cell)
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
			headerLine += header
		}
		writer.Write([]byte(headerLine + "\n"))
	}

	// Write data rows
	for _, row := range table.Rows {
		// Add row separator before each row except the first one
		if len(table.Headers) > 0 || len(table.Rows) > 1 {
			writer.Write([]byte("|-\n"))
		}

		// Write row data (MediaWiki uses | for data cells)
		rowLine := "| "
		for i, cell := range row {
			if i > 0 {
				rowLine += " || "
			}
			rowLine += cell
		}
		writer.Write([]byte(rowLine + "\n"))
	}

	// Write table end
	writer.Write([]byte("|}\n"))

	return nil
}
