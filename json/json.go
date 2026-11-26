package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/martianzhang/tableconvert/common"
)

func decodeJSON(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return err
	}
	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON data")
		}
		return err
	}
	return nil
}

func Unmarshal(cfg *common.Config, table *common.Table) error {
	format := cfg.GetExtensionString("format", "")

	data, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return err
	}
	data = common.TrimUTF8BOM(data)

	switch format {
	case "2d":
		var input [][]interface{}
		if err := decodeJSON(data, &input); err != nil {
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
				table.Headers[i] = common.JSONScalarString(header)
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
						stringRow[i] = common.JSONScalarString(row[i])
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
		if err := decodeJSON(data, &input); err != nil {
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
						row[j] = common.JSONScalarString(col[i])
					}
				}
			}
			table.Rows = append(table.Rows, row)
		}

	case "keyed":
		var input map[string]map[string]interface{}
		if err := decodeJSON(data, &input); err != nil {
			return err
		}
		if len(input) == 0 {
			return nil
		}

		// Object values provide the column names; the object key is the first
		// column. We use "key" as the synthetic first-column header because the
		// keyed shape does not carry the original header name.
		headerSet := make(map[string]struct{})
		for _, obj := range input {
			for k := range obj {
				headerSet[k] = struct{}{}
			}
		}
		headers := []string{"key"}
		for k := range headerSet {
			headers = append(headers, k)
		}
		sort.Strings(headers[1:])
		table.Headers = headers

		keys := make([]string, 0, len(input))
		for k := range input {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			row := make([]string, len(headers))
			row[0] = key
			for i, header := range headers[1:] {
				if val, ok := input[key][header]; ok && val != nil {
					row[i+1] = common.JSONScalarString(val)
				}
			}
			table.Rows = append(table.Rows, row)
		}

	default: // Array of Object
		root, err := decodeJSONRoot(data, cfg)
		if err != nil {
			return err
		}
		headers, rows := common.FlattenJSONRecords(
			root,
			common.JSONFlattenOptions{Arrays: cfg.GetExtensionString("flatten-arrays", "json")},
			cfg.GetExtensionBool("sort-headers", false),
		)
		table.Headers = headers
		table.Rows = rows
	}

	return nil
}

// decodeJSONRoot decodes a JSON document and returns the rows selected by
// --json-path. A top-level object is treated as one row, so typical API
// payloads like {"id":1} work without wrapping them in an array.
func decodeJSONRoot(data []byte, cfg *common.Config) ([]interface{}, error) {
	root, err := common.DecodeOrderedJSON(data)
	if err != nil {
		return nil, err
	}
	if path := cfg.GetExtensionString("json-path", ""); path != "" {
		root, err = common.JSONPathValue(root, path)
		if err != nil {
			return nil, err
		}
	}
	switch value := root.(type) {
	case []interface{}:
		return value, nil
	case *common.OrderedObject, map[string]interface{}:
		return []interface{}{value}, nil
	default:
		// Scalar roots become a one-column table; FlattenJSONRecords handles
		// the "value" column.
		return []interface{}{value}, nil
	}
}

func Marshal(cfg *common.Config, table *common.Table) error {
	format := cfg.GetExtensionString("format", "")
	parsing := cfg.GetParseTypes()
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
	// Keyed by first column
	case "keyed":
		if len(table.Headers) == 0 {
			data, err = json.Marshal(map[string]interface{}{})
			break
		}
		output := make(map[string]map[string]interface{}, len(table.Rows))
		for _, row := range table.Rows {
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			key := row[0]
			if _, exists := output[key]; exists {
				return fmt.Errorf("duplicate key %q in keyed JSON output", key)
			}
			record := make(map[string]interface{}, len(table.Headers)-1)
			for i := 1; i < len(table.Headers); i++ {
				if parsing {
					record[table.Headers[i]] = common.InferType(row[i])
				} else {
					record[table.Headers[i]] = row[i]
				}
			}
			output[key] = record
		}
		if minify {
			data, err = json.Marshal(output)
		} else {
			data, err = json.MarshalIndent(output, "", "  ")
		}
	// Array of Object
	default:
		output := make([]map[string]interface{}, 0, len(table.Rows))
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
