package json

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

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
			if header == nil {
				table.Headers[i] = "NULL"
			} else {
				table.Headers[i] = fmt.Sprint(header)
			}
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
				} else {
					// Handle missing values
					stringRow[i] = "NULL"
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

		// Set headers (sorted for determinism)
		for k := range columnData {
			table.Headers = append(table.Headers, k)
		}
		sort.Strings(table.Headers)

		// Fill rows
		for i := 0; i < numRows; i++ {
			row := make([]string, len(table.Headers))
			for j, header := range table.Headers {
				col := columnData[header]
				if i < len(col) {
					if col[i] == nil {
						row[j] = "NULL"
					} else {
						row[j] = fmt.Sprint(col[i])
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

		// Collect all headers and sort them for deterministic order
		headerSet := make(map[string]struct{})
		for _, obj := range input {
			for key := range obj {
				headerSet[key] = struct{}{}
			}
		}
		// Sort all headers
		for key := range headerSet {
			table.Headers = append(table.Headers, key)
		}
		sort.Strings(table.Headers)

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
						record[i] = common.InferType(row[i])
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
				header := table.Headers[i]
				if parsing {
					columns[header] = append(columns[header], common.InferType(cell))
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
						record[header] = common.InferType(row[i])
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
