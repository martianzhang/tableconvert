package toml

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	tomlv2 "github.com/pelletier/go-toml/v2"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("toml: target table is nil")
	}
	data, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return err
	}
	var root map[string]interface{}
	if err := tomlv2.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("toml: %w", err)
	}

	records, err := recordsFromTOML(root, cfg.GetExtensionString("toml-table", ""))
	if err != nil {
		return err
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
	if table == nil {
		return fmt.Errorf("toml: input table is nil")
	}
	typed := cfg.GetParseTypes() || (common.IsJSONSource(cfg.From) && table.Raw != nil)

	records := make([]map[string]interface{}, 0, len(table.Rows))
	for rowIndex, row := range table.Rows {
		if len(row) != len(table.Headers) {
			return fmt.Errorf("toml: row length %d does not match header length %d", len(row), len(table.Headers))
		}
		record := make(map[string]interface{}, len(table.Headers))
		for i, header := range table.Headers {
			record[header] = cellValue(table, rowIndex, i, typed)
		}
		records = append(records, record)
	}

	var output interface{} = map[string]interface{}{
		tomlTableName(cfg): records,
	}
	if strings.EqualFold(cfg.GetExtensionString("format", ""), "column") {
		columns := make(map[string][]interface{}, len(table.Headers))
		for _, header := range table.Headers {
			columns[header] = []interface{}{}
		}
		for rowIndex := range table.Rows {
			for i, header := range table.Headers {
				columns[header] = append(columns[header], cellValue(table, rowIndex, i, typed))
			}
		}
		output = columns
	}

	data, err := tomlv2.Marshal(output)
	if err != nil {
		return err
	}
	_, err = cfg.Writer.Write(data)
	return err
}

func tomlTableName(cfg *common.Config) string {
	name := strings.TrimSpace(cfg.GetExtensionString("toml-table", ""))
	if name == "" {
		return "row"
	}
	return name
}

func cellValue(table *common.Table, rowIndex, colIndex int, typed bool) interface{} {
	cell := table.Rows[rowIndex][colIndex]
	if typed && table.Raw != nil && rowIndex < len(table.Raw) && colIndex < len(table.Raw[rowIndex]) {
		return table.Raw[rowIndex][colIndex]
	}
	if typed {
		return common.InferType(cell)
	}
	return cell
}

func recordsFromTOML(root map[string]interface{}, tableName string) ([]interface{}, error) {
	if tableName != "" {
		value, ok := root[tableName]
		if !ok {
			return nil, fmt.Errorf("toml: table %q not found", tableName)
		}
		return asRecords(value), nil
	}

	keys := make([]string, 0, len(root))
	for key := range root {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if records := arrayOfRecords(root[key]); records != nil {
			return records, nil
		}
	}
	for _, key := range keys {
		if _, ok := root[key].([]interface{}); ok {
			return columnRecords(root), nil
		}
	}
	return []interface{}{root}, nil
}

func asRecords(value interface{}) []interface{} {
	if records := arrayOfRecords(value); records != nil {
		return records
	}
	if array, ok := value.([]interface{}); ok {
		return array
	}
	if object, ok := value.(map[string]interface{}); ok {
		return []interface{}{object}
	}
	return []interface{}{value}
}

func arrayOfRecords(value interface{}) []interface{} {
	switch array := value.(type) {
	case []map[string]interface{}:
		records := make([]interface{}, len(array))
		for i := range array {
			records[i] = array[i]
		}
		return records
	case []interface{}:
		if len(array) == 0 {
			return nil
		}
		for _, item := range array {
			if _, ok := item.(map[string]interface{}); !ok {
				return nil
			}
		}
		return array
	default:
		return nil
	}
}

func columnRecords(root map[string]interface{}) []interface{} {
	maxRows := 0
	for _, value := range root {
		if array, ok := value.([]interface{}); ok && len(array) > maxRows {
			maxRows = len(array)
		}
	}
	records := make([]interface{}, maxRows)
	for i := 0; i < maxRows; i++ {
		record := make(map[string]interface{}, len(root))
		for key, value := range root {
			if array, ok := value.([]interface{}); ok && i < len(array) {
				record[key] = array[i]
			}
		}
		records[i] = record
	}
	return records
}
