package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

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
	format := strings.ToLower(cfg.GetExtensionString("format", ""))

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
		root, err := common.DecodeOrderedJSON(data)
		if err != nil {
			return err
		}

		columnData := make(map[string][]interface{})
		headerOrder := make([]string, 0)
		addColumns := func(object *common.OrderedObject) error {
			for _, key := range object.Keys {
				arr, ok := object.Values[key].([]interface{})
				if !ok {
					return fmt.Errorf("invalid column format for key %s", key)
				}
				if _, exists := columnData[key]; !exists {
					headerOrder = append(headerOrder, key)
				}
				columnData[key] = arr
			}
			return nil
		}

		switch value := root.(type) {
		case *common.OrderedObject:
			if err := addColumns(value); err != nil {
				return err
			}
		case []interface{}:
			for _, item := range value {
				switch object := item.(type) {
				case *common.OrderedObject:
					if err := addColumns(object); err != nil {
						return err
					}
				case map[string]interface{}:
					keys := make([]string, 0, len(object))
					for key := range object {
						keys = append(keys, key)
					}
					sort.Strings(keys)
					ordered := &common.OrderedObject{Keys: keys, Values: object}
					if err := addColumns(ordered); err != nil {
						return err
					}
				default:
					return fmt.Errorf("invalid column format: expected an object")
				}
			}
		default:
			return fmt.Errorf("invalid column format: expected an object or an array of objects")
		}

		table.Headers = headerOrder
		numRows := 0
		for _, column := range columnData {
			if len(column) > numRows {
				numRows = len(column)
			}
		}
		for i := 0; i < numRows; i++ {
			row := make([]string, len(table.Headers))
			for j, header := range table.Headers {
				column := columnData[header]
				if i < len(column) {
					if column[i] == nil {
						row[j] = "NULL"
					} else {
						row[j] = common.JSONScalarString(column[i])
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
		headers, rows, raw := common.FlattenJSONRecordsWithRaw(
			root,
			common.JSONFlattenOptions{Arrays: cfg.GetExtensionString("flatten-arrays", "json")},
			cfg.GetExtensionBool("sort-headers", false),
		)
		table.Headers = headers
		table.Rows = rows
		table.RawRows = root
		table.Raw = raw
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
	format := strings.ToLower(cfg.GetExtensionString("format", ""))
	parsing := cfg.GetParseTypes()
	typed := parsing || (common.IsJSONSource(cfg.From) && table.Raw != nil)

	var output interface{}
	switch format {
	// 2D Array
	case "2d":
		rowsOutput := make([][]interface{}, 0, len(table.Rows)+1)
		headers := make([]interface{}, len(table.Headers))
		for i, header := range table.Headers {
			headers[i] = header
		}
		rowsOutput = append(rowsOutput, headers)

		for rowIndex, row := range table.Rows {
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			record := make([]interface{}, len(table.Headers))
			for i := range table.Headers {
				record[i] = jsonCellValue(table, rowIndex, i, typed)
			}
			rowsOutput = append(rowsOutput, record)
		}
		output = rowsOutput

	// Column-oriented object: {"a":[...],"b":[...]}
	case "column":
		columns := make(map[string][]interface{}, len(table.Headers))
		for _, header := range table.Headers {
			columns[header] = []interface{}{}
		}
		for rowIndex, row := range table.Rows {
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			for i := range row {
				header := table.Headers[i]
				columns[header] = append(columns[header], jsonCellValue(table, rowIndex, i, typed))
			}
		}
		output = columns

	// Legacy column array: [{"a":[...]},{"b":[...]}]
	case "column-array":
		columns := make(map[string][]interface{}, len(table.Headers))
		for _, header := range table.Headers {
			columns[header] = []interface{}{}
		}
		for rowIndex, row := range table.Rows {
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			for i := range row {
				header := table.Headers[i]
				columns[header] = append(columns[header], jsonCellValue(table, rowIndex, i, typed))
			}
		}
		legacy := make([]map[string]interface{}, 0, len(table.Headers))
		for _, header := range table.Headers {
			legacy = append(legacy, map[string]interface{}{header: columns[header]})
		}
		output = legacy

	// Keyed by first column
	case "keyed":
		if len(table.Headers) == 0 {
			output = map[string]interface{}{}
			break
		}
		keyed := make(map[string]map[string]interface{}, len(table.Rows))
		for rowIndex, row := range table.Rows {
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			key := row[0]
			if _, exists := keyed[key]; exists {
				return fmt.Errorf("duplicate key %q in keyed JSON output", key)
			}
			record := make(map[string]interface{}, len(table.Headers)-1)
			for i := 1; i < len(table.Headers); i++ {
				record[table.Headers[i]] = jsonCellValue(table, rowIndex, i, typed)
			}
			keyed[key] = record
		}
		output = keyed

	// Array of Object
	default:
		if len(table.RawRows) == len(table.Rows) && common.IsJSONSource(cfg.From) {
			output = common.ToPlainJSON(table.RawRows)
			break
		}
		objects := make([]map[string]interface{}, 0, len(table.Rows))
		for rowIndex := range table.Rows {
			row := table.Rows[rowIndex]
			if len(row) != len(table.Headers) {
				return fmt.Errorf("row length %d does not match header length %d", len(row), len(table.Headers))
			}
			record := make(map[string]interface{})
			for i, header := range table.Headers {
				record[header] = jsonCellValue(table, rowIndex, i, typed)
			}
			objects = append(objects, record)
		}
		output = objects
	}

	data, err := marshalJSONOutput(cfg, output)
	if err != nil {
		return err
	}
	_, err = cfg.Writer.Write(data)
	return err
}

func marshalJSONOutput(cfg *common.Config, output interface{}) ([]byte, error) {
	if root := strings.TrimSpace(cfg.GetExtensionString("root", "")); root != "" {
		output = map[string]interface{}{root: output}
	}
	if cfg.GetExtensionBool("minify", false) {
		return json.Marshal(output)
	}
	indent := cfg.GetExtensionInt("indent", 2)
	if indent <= 0 {
		return json.Marshal(output)
	}
	return json.MarshalIndent(output, "", strings.Repeat(" ", indent))
}

func jsonCellValue(table *common.Table, rowIndex, colIndex int, typed bool) interface{} {
	cell := table.Rows[rowIndex][colIndex]
	if typed && table.Raw != nil && rowIndex < len(table.Raw) && colIndex < len(table.Raw[rowIndex]) {
		return table.Raw[rowIndex][colIndex]
	}
	if typed {
		return common.InferType(cell)
	}
	return cell
}
