package jsonl

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	scanner := bufio.NewScanner(cfg.Reader)
	var records []map[string]interface{}

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if line == "" {
			continue // skip empty lines
		}

		var record map[string]interface{}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return &common.ParseError{
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("invalid JSON: %v", err),
				Line:       line,
			}
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read JSONL: %w", err)
	}

	if len(records) == 0 {
		return &common.ParseError{
			LineNumber: 0,
			Message:    "empty JSONL file",
			Line:       "",
		}
	}

	// Collect all unique keys from all records to form headers
	headerMap := make(map[string]bool)
	for _, record := range records {
		for key := range record {
			headerMap[key] = true
		}
	}

	// Convert map to sorted slice of headers
	headers := make([]string, 0, len(headerMap))
	for key := range headerMap {
		headers = append(headers, key)
	}

	// Convert records to rows
	rows := make([][]string, len(records))
	for i, record := range records {
		row := make([]string, len(headers))
		for j, header := range headers {
			if val, ok := record[header]; ok {
				row[j] = fmt.Sprintf("%v", val)
			} else {
				row[j] = "" // empty string for missing fields
			}
		}
		rows[i] = row
	}

	table.Headers = headers
	table.Rows = rows
	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	parsing := cfg.GetExtensionBool("parsing-json", false)

	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()

	// Each row becomes a JSON object with headers as keys
	for _, row := range table.Rows {
		if len(row) != len(table.Headers) {
			return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
		}
		record := make(map[string]interface{})
		for i, header := range table.Headers {
			if parsing {
				record[header] = common.InferType(row[i])
			} else {
				record[header] = row[i]
			}
		}

		jsonData, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		if _, err := writer.Write(jsonData); err != nil {
			return fmt.Errorf("failed to write JSON: %w", err)
		}
		if _, err := writer.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	return nil
}
