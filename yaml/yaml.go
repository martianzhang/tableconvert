package yaml

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	yamlv3 "gopkg.in/yaml.v3"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("yaml: target table is nil")
	}
	data, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return err
	}
	var root interface{}
	if err := yamlv3.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("yaml: %w", err)
	}
	root = normalize(root)
	if is2D(root) {
		headers, rows := flatten2D(root)
		table.Headers = headers
		table.Rows = rows
		return nil
	}

	records := recordsFromYAML(root)
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
		return fmt.Errorf("yaml: input table is nil")
	}
	typed := cfg.GetParseTypes() || (common.IsJSONSource(cfg.From) && table.Raw != nil)

	output := make([]map[string]interface{}, 0, len(table.Rows))
	for rowIndex, row := range table.Rows {
		if len(row) != len(table.Headers) {
			return fmt.Errorf("yaml: row length %d does not match header length %d", len(row), len(table.Headers))
		}
		record := make(map[string]interface{}, len(table.Headers))
		for i, header := range table.Headers {
			record[header] = cellValue(table, rowIndex, i, typed)
		}
		output = append(output, record)
	}

	var value interface{} = output
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
		value = columns
	}

	encoder := yamlv3.NewEncoder(cfg.Writer)
	indent := cfg.GetExtensionInt("indent", 2)
	if indent > 0 {
		encoder.SetIndent(indent)
	}
	flow := cfg.GetExtensionBool("flow", false)
	quote := strings.ToLower(strings.TrimSpace(cfg.GetExtensionString("quote", "")))
	if flow || quote != "" {
		if err := encoder.Encode(toYAMLNode(value, flow, quote)); err != nil {
			return err
		}
	} else if err := encoder.Encode(value); err != nil {
		return err
	}
	return encoder.Close()
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

func normalize(value interface{}) interface{} {
	switch node := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(node))
		for key, child := range node {
			result[key] = normalize(child)
		}
		return result
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(node))
		for key, child := range node {
			result[fmt.Sprint(key)] = normalize(child)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(node))
		for i, child := range node {
			result[i] = normalize(child)
		}
		return result
	default:
		return value
	}
}

func is2D(root interface{}) bool {
	array, ok := root.([]interface{})
	if !ok || len(array) == 0 {
		return false
	}
	_, ok = array[0].([]interface{})
	return ok
}

func flatten2D(root interface{}) ([]string, [][]string) {
	rows := root.([]interface{})
	headers := make([]string, 0)
	if len(rows) > 0 {
		if first, ok := rows[0].([]interface{}); ok {
			for _, value := range first {
				headers = append(headers, common.JSONScalarString(value))
			}
		}
	}
	data := make([][]string, 0, len(rows)-1)
	for _, rawRow := range rows[1:] {
		rowValues, ok := rawRow.([]interface{})
		if !ok {
			continue
		}
		row := make([]string, len(headers))
		for i := range headers {
			if i < len(rowValues) {
				row[i] = common.JSONScalarString(rowValues[i])
			}
		}
		data = append(data, row)
	}
	return headers, data
}

func recordsFromYAML(root interface{}) []interface{} {
	switch node := root.(type) {
	case []interface{}:
		return node
	case map[string]interface{}:
		allArrays := len(node) > 0
		maxRows := 0
		for _, value := range node {
			array, ok := value.([]interface{})
			if !ok {
				allArrays = false
				break
			}
			if len(array) > maxRows {
				maxRows = len(array)
			}
		}
		if allArrays {
			records := make([]interface{}, maxRows)
			for i := 0; i < maxRows; i++ {
				record := make(map[string]interface{}, len(node))
				for key, value := range node {
					array := value.([]interface{})
					if i < len(array) {
						record[key] = array[i]
					}
				}
				records[i] = record
			}
			return records
		}
		return []interface{}{node}
	default:
		return []interface{}{node}
	}
}

func toYAMLNode(value interface{}, flow bool, quote string) *yamlv3.Node {
	switch node := value.(type) {
	case map[string]interface{}:
		mapping := &yamlv3.Node{Kind: yamlv3.MappingNode}
		if flow {
			mapping.Style = yamlv3.FlowStyle
		}
		keys := make([]string, 0, len(node))
		for key := range node {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			mapping.Content = append(mapping.Content, scalarNode(key, quote), toYAMLNode(node[key], flow, quote))
		}
		return mapping
	case map[string][]interface{}:
		converted := make(map[string]interface{}, len(node))
		for key, values := range node {
			converted[key] = values
		}
		return toYAMLNode(converted, flow, quote)
	case []map[string]interface{}:
		sequence := &yamlv3.Node{Kind: yamlv3.SequenceNode}
		if flow {
			sequence.Style = yamlv3.FlowStyle
		}
		for _, item := range node {
			sequence.Content = append(sequence.Content, toYAMLNode(item, flow, quote))
		}
		return sequence
	case []interface{}:
		sequence := &yamlv3.Node{Kind: yamlv3.SequenceNode}
		if flow {
			sequence.Style = yamlv3.FlowStyle
		}
		for _, item := range node {
			sequence.Content = append(sequence.Content, toYAMLNode(item, flow, quote))
		}
		return sequence
	default:
		return scalarNode(fmt.Sprint(node), quote)
	}
}

func scalarNode(value, quote string) *yamlv3.Node {
	node := &yamlv3.Node{Kind: yamlv3.ScalarNode, Value: value}
	switch strings.ToLower(quote) {
	case "double":
		node.Style = yamlv3.DoubleQuotedStyle
	case "single":
		node.Style = yamlv3.SingleQuotedStyle
	}
	return node
}
