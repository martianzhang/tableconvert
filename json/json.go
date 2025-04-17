package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	format := cfg.GetExtensionString("format", "")

	data, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return err
	}

	switch format {
	case "2d":
		var input [][]interface{}
		if err := json.Unmarshal(data, &input); err != nil {
			return err
		}
		if len(input) == 0 {
			return nil
		}
		// Extract headers
		table.Headers = make([]string, len(input[0]))
		for i, header := range input[0] {
			table.Headers[i] = fmt.Sprint(header)
		}
		// Extract rows
		for _, row := range input[1:] {
			stringRow := make([]string, len(table.Headers))
			for i := range table.Headers {
				if i < len(row) {
					if row[i] == nil {
						stringRow[i] = "NULL"
					} else {
						stringRow[i] = fmt.Sprint(row[i])
					}
				}
			}
			table.Rows = append(table.Rows, stringRow)
		}

	case "column":
		var input []map[string]interface{}
		if err := json.Unmarshal(data, &input); err != nil {
			return err
		}

		columnData := make(map[string][]interface{})
		var numRows int
		for _, obj := range input {
			for k, v := range obj {
				switch arr := v.(type) {
				case []interface{}:
					columnData[k] = arr
					if len(arr) > numRows {
						numRows = len(arr)
					}
				default:
					return fmt.Errorf("invalid column format for key %s", k)
				}
			}
		}

		// Set headers
		for k := range columnData {
			table.Headers = append(table.Headers, k)
		}

		// Fill rows
		for i := 0; i < numRows; i++ {
			row := make([]string, len(table.Headers))
			for j, header := range table.Headers {
				col := columnData[header]
				if i < len(col) {
					if col[j] == nil {
						row[j] = "NULL"
					} else {
						row[j] = fmt.Sprint(col[j])
					}
				}
			}
			table.Rows = append(table.Rows, row)
		}

	default: // Array of Object
		var input []map[string]interface{}
		if err := json.Unmarshal(data, &input); err != nil {
			return err
		}

		// Collect headers in consistent order
		headerSet := make(map[string]struct{})
		for _, obj := range input {
			for key := range obj {
				headerSet[key] = struct{}{}
			}
		}
		// Use the first object to get header order
		if len(input) > 0 {
			for key := range input[0] {
				table.Headers = append(table.Headers, key)
			}
			for key := range headerSet {
				if !contains(table.Headers, key) {
					table.Headers = append(table.Headers, key)
				}
			}
		}

		for _, obj := range input {
			row := make([]string, len(table.Headers))
			for i, header := range table.Headers {
				if val, ok := obj[header]; ok {
					if val == nil {
						row[i] = "NULL"
					} else {
						row[i] = fmt.Sprint(val)
					}
				}
			}
			table.Rows = append(table.Rows, row)
		}
	}

	return nil
}

// contains checks if a string exists in a slice
func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func Marshal(cfg *common.Config, table *common.Table) error {
	format := cfg.GetExtensionString("format", "")
	parsing := cfg.GetExtensionBool("parsing-json", false)
	minify := cfg.GetExtensionBool("minify", false)

	var data []byte
	var err error
	switch format {
	// 2D Array
	case "2d":
		var output [][]interface{}
		// Header
		headers := make([]interface{}, len(table.Headers))
		for i, header := range table.Headers {
			headers[i] = header
		}
		output = append(output, headers)

		// Rows
		for _, row := range table.Rows {
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			record := make([]interface{}, len(table.Headers))
			for i := range table.Headers {
				if i < len(row) {
					if parsing {
						record[i] = inferType(row[i])
					} else {
						record[i] = row[i]
					}
				}
			}
			output = append(output, record)
		}

		if minify {
			data, err = json.Marshal(output)
		} else {
			data, err = json.MarshalIndent(output, "", "  ")
		}
	// Column Array
	case "column":
		columns := make(map[string][]interface{}, len(table.Headers))
		for _, header := range table.Headers {
			columns[header] = []interface{}{}
		}
		for _, row := range table.Rows {
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			for i, cell := range row {
				if i >= len(table.Headers) {
					continue
				}
				header := table.Headers[i]
				if parsing {
					columns[header] = append(columns[header], inferType(cell))
				} else {
					columns[header] = append(columns[header], cell)
				}
			}
		}
		// convert map[string][]interface{} to []map[string]interface{}
		var output []map[string]interface{}
		for _, header := range table.Headers {
			output = append(output, map[string]interface{}{header: columns[header]})
		}
		if minify {
			data, err = json.Marshal(output)
		} else {
			data, err = json.MarshalIndent(output, "", "  ")
		}
	// Array of Object
	default:
		var output []map[string]interface{}
		for _, row := range table.Rows {
			record := make(map[string]interface{})
			for i, header := range table.Headers {
				if i < len(row) {
					if parsing {
						record[header] = inferType(row[i])
					} else {
						record[header] = row[i]
					}
				}
			}
			output = append(output, record)
		}
		// minify json output
		if minify {
			data, err = json.Marshal(output)
		} else {
			data, err = json.MarshalIndent(output, "", "  ")
		}
	}
	// deal with json marshal error
	if err != nil {
		return err
	}

	_, err = cfg.Writer.Write(data)
	return err
}

// inferType attempts to convert a string value to a more specific type (bool, int64, float64, nil)
// If no conversion is successful, it returns the original string.
func inferType(value string) interface{} {
	trimmedValue := strings.TrimSpace(value)

	// 1. Check for explicit null (case-insensitive)
	if strings.ToLower(trimmedValue) == "null" {
		return nil
	}

	// 2. Check for boolean (case-insensitive)
	lowerValue := strings.ToLower(trimmedValue)
	if lowerValue == "true" {
		return true
	}
	if lowerValue == "false" {
		return false
	}

	// 3. Check for integer
	// Use ParseInt for potentially larger numbers and better base control
	if intVal, err := strconv.ParseInt(trimmedValue, 10, 64); err == nil {
		// Optional: Double-check if the string representation truly matches
		// to avoid interpreting things like "0xf" as integers if that's not desired.
		// Here, we assume a valid ParseInt result means it's an integer.
		return intVal
	}

	// 4. Check for float
	// It must contain ".", "e", or "E" to be considered float by some stricter definitions,
	// but ParseFloat is more general. We'll parse anything that looks like a float.
	if floatVal, err := strconv.ParseFloat(trimmedValue, 64); err == nil {
		// Check if it's actually an integer represented as float (e.g., "123.0")
		// If you want "123.0" to become integer 123, you might need extra logic here.
		// For simplicity, we'll let ParseFloat decide.
		return floatVal
	}

	// 5. Default: return the original string (or trimmed, depending on preference)
	// Returning the original 'value' preserves leading/trailing whitespace if needed.
	// If you always want trimmed strings, return 'trimmedValue'.
	return value
}
