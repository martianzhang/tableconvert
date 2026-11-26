package jsonl

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	scanner := bufio.NewScanner(common.StripUTF8BOM(cfg.Reader))
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	var records []interface{}

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if line == "" {
			continue // skip empty lines
		}

		value, err := common.DecodeOrderedJSON([]byte(line))
		if err != nil {
			return &common.ParseError{
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("invalid JSON: %v", err),
				Line:       line,
			}
		}
		if array, ok := value.([]interface{}); ok {
			records = append(records, array...)
			continue
		}
		records = append(records, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read JSONL: %w", err)
	}

	if len(records) == 0 {
		return nil // An empty JSONL file is an empty table.
	}

	headers, rows, raw := common.FlattenJSONRecordsWithRaw(
		records,
		common.JSONFlattenOptions{Arrays: cfg.GetExtensionString("flatten-arrays", "json")},
		cfg.GetExtensionBool("sort-headers", false),
	)
	table.Headers = headers
	table.Rows = rows
	table.RawRows = records
	table.Raw = raw
	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	parsing := cfg.GetParseTypes()
	typed := parsing || (common.IsJSONSource(cfg.From) && table.Raw != nil)

	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()

	if len(table.RawRows) == len(table.Rows) && common.IsJSONSource(cfg.From) {
		for _, rawRow := range table.RawRows {
			jsonData, err := json.Marshal(common.ToPlainJSON(rawRow))
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

	// Each row becomes a JSON object with headers as keys
	for rowIndex, row := range table.Rows {
		if len(row) != len(table.Headers) {
			return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
		}
		record := make(map[string]interface{})
		for i, header := range table.Headers {
			record[header] = jsonlCellValue(table, rowIndex, i, typed)
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

func jsonlCellValue(table *common.Table, rowIndex, colIndex int, typed bool) interface{} {
	cell := table.Rows[rowIndex][colIndex]
	if typed && table.Raw != nil && rowIndex < len(table.Raw) && colIndex < len(table.Raw[rowIndex]) {
		return table.Raw[rowIndex][colIndex]
	}
	if typed {
		return common.InferType(cell)
	}
	return cell
}
