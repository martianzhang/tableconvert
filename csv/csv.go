package csv

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	csvReader := csv.NewReader(cfg.Reader)

	// Set custom delimiter (default: comma)
	if vd, ok := cfg.Extension["delimiter"]; ok {
		switch strings.ToUpper(vd) {
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

	table.Headers = records[0]
	table.Rows = records[1:]

	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	// Write UTF-8 BOM
	if bom, ok := cfg.Extension["bom"]; ok && strings.ToLower(bom) != "false" {
		if _, err := cfg.Writer.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
			return fmt.Errorf("failed to write BOM: %w", err)
		}
	}
	// Config CSV writer
	csvWriter := csv.NewWriter(cfg.Writer)

	// Set custom delimiter (default: comma)
	if vd, ok := cfg.Extension["delimiter"]; ok {
		switch strings.ToUpper(vd) {
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
	}

	// CSV header
	if err := csvWriter.Write(table.Headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// CSV records
	for _, row := range table.Rows {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
