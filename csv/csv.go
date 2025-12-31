package csv

import (
	"encoding/csv"
	"fmt"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	csvReader := csv.NewReader(cfg.Reader)
	// Allow variable number of fields per record
	csvReader.FieldsPerRecord = -1

	// Set custom delimiter (default: comma)
	delimiter := cfg.GetExtensionString("delimiter", ",")
	switch delimiter {
	case "TAB", "\t":
		csvReader.Comma = '\t'
	case "SEMICOLON", ";":
		csvReader.Comma = ';'
	case "PIPE", "|":
		csvReader.Comma = '|'
	case "SLASH", "/":
		csvReader.Comma = '/'
	case "HASH", "#":
		csvReader.Comma = '#'
		// default remains comma
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return &common.ParseError{
			LineNumber: 0,
			Message:    "empty CSV file",
			Line:       "",
		}
	}

	// Handle first-column-header option
	firstColHeader := cfg.GetExtensionBool("first-column-header", false)
	if firstColHeader {
		// Find max column count across all records
		maxCols := 0
		for _, record := range records {
			if len(record) > maxCols {
				maxCols = len(record)
			}
		}

		if maxCols == 1 {
			// Single column: first row is header, remaining rows are data
			if len(records) > 0 {
				table.Headers = records[0]
			}
			if len(records) > 1 {
				table.Rows = records[1:]
			}
		} else {
			// Multi-column: transpose
			// Use first column as headers
			headers := make([]string, len(records))
			rows := make([][]string, maxCols-1)

			// Extract headers from first column
			for i, record := range records {
				if len(record) == 0 {
					headers[i] = ""
				} else {
					headers[i] = record[0]
				}
			}

			// Extract data from remaining columns
			for i := 1; i < maxCols; i++ {
				row := make([]string, len(records))
				for j := 0; j < len(records); j++ {
					if i < len(records[j]) {
						row[j] = records[j][i]
					} else {
						row[j] = "" // pad with empty string if column is missing
					}
				}
				rows[i-1] = row
			}

			table.Headers = headers
			table.Rows = rows
		}
	} else {
		// Default behavior: first row is headers
		table.Headers = records[0]
		table.Rows = records[1:]
	}

	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	// Write UTF-8 BOM
	bom := cfg.GetExtensionBool("bom", false)
	if bom {
		if _, err := cfg.Writer.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
			return fmt.Errorf("failed to write BOM: %w", err)
		}
	}

	// Config CSV writer
	csvWriter := csv.NewWriter(cfg.Writer)

	// Set custom delimiter (default: comma)
	delimiter := cfg.GetExtensionString("delimiter", ",")
	switch delimiter {
	case "TAB", "\t":
		csvWriter.Comma = '\t'
	case "SEMICOLON", ";":
		csvWriter.Comma = ';'
	case "PIPE", "|":
		csvWriter.Comma = '|'
	case "SLASH", "/":
		csvWriter.Comma = '/'
	case "HASH", "#":
		csvWriter.Comma = '#'
		// default remains comma
	}

	// Default behavior: first row is headers
	if err := csvWriter.Write(table.Headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data rows
	for _, row := range table.Rows {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
