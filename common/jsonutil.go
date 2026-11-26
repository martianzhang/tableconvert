package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// OrderedObject is a JSON object that preserves the key order from the input.
// encoding/json decodes objects into maps, which loses that order; table
// headers should follow the source document instead of an arbitrary order.
type OrderedObject struct {
	Keys   []string
	Values map[string]interface{}
}

// Get returns the value for key and whether the key exists.
func (o *OrderedObject) Get(key string) (interface{}, bool) {
	if o == nil {
		return nil, false
	}
	value, ok := o.Values[key]
	return value, ok
}

// DecodeOrderedJSON decodes exactly one JSON value while preserving object key
// order. Numbers are decoded as json.Number so their original text is kept.
func DecodeOrderedJSON(data []byte) (interface{}, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := decodeOrderedValue(decoder)
	if err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("unexpected trailing JSON data")
		}
		return nil, err
	}
	return value, nil
}

func decodeOrderedValue(decoder *json.Decoder) (interface{}, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			object := &OrderedObject{Values: make(map[string]interface{})}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, fmt.Errorf("invalid JSON object key")
				}
				child, err := decodeOrderedValue(decoder)
				if err != nil {
					return nil, err
				}
				if _, exists := object.Values[key]; !exists {
					object.Keys = append(object.Keys, key)
				}
				object.Values[key] = child
			}
			if _, err := decoder.Token(); err != nil { // closing '}'
				return nil, err
			}
			return object, nil
		case '[':
			array := make([]interface{}, 0)
			for decoder.More() {
				child, err := decodeOrderedValue(decoder)
				if err != nil {
					return nil, err
				}
				array = append(array, child)
			}
			if _, err := decoder.Token(); err != nil { // closing ']'
				return nil, err
			}
			return array, nil
		default:
			return nil, fmt.Errorf("unexpected JSON delimiter %q", value)
		}
	case string, json.Number, bool, nil:
		return value, nil
	default:
		return nil, fmt.Errorf("unexpected JSON token %T", token)
	}
}

// JSONPathValue selects a nested value using a small dot/bracket path such as
// "data.items" or "data.items[0]". It is intentionally not a full JSONPath
// implementation; it covers the common "API wrapper object" case.
func JSONPathValue(root interface{}, path string) (interface{}, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return root, nil
	}
	current := root
	for _, segment := range strings.Split(path, ".") {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		name, indexes, err := splitJSONPathSegment(segment)
		if err != nil {
			return nil, err
		}
		if name != "" {
			switch node := current.(type) {
			case *OrderedObject:
				value, ok := node.Get(name)
				if !ok {
					return nil, fmt.Errorf("JSON path %q not found", path)
				}
				current = value
			case map[string]interface{}:
				value, ok := node[name]
				if !ok {
					return nil, fmt.Errorf("JSON path %q not found", path)
				}
				current = value
			default:
				return nil, fmt.Errorf("JSON path %q traverses a non-object", path)
			}
		}
		for _, index := range indexes {
			array, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("JSON path %q indexes a non-array", path)
			}
			if index < 0 || index >= len(array) {
				return nil, fmt.Errorf("JSON path %q index %d out of range", path, index)
			}
			current = array[index]
		}
	}
	return current, nil
}

func splitJSONPathSegment(segment string) (string, []int, error) {
	open := strings.IndexByte(segment, '[')
	if open < 0 {
		return segment, nil, nil
	}
	name := segment[:open]
	rest := segment[open:]
	var indexes []int
	for len(rest) > 0 {
		if rest[0] != '[' {
			return "", nil, fmt.Errorf("invalid JSON path segment %q", segment)
		}
		end := strings.IndexByte(rest, ']')
		if end < 0 {
			return "", nil, fmt.Errorf("invalid JSON path segment %q", segment)
		}
		index, err := strconv.Atoi(rest[1:end])
		if err != nil {
			return "", nil, fmt.Errorf("invalid JSON path index in %q", segment)
		}
		indexes = append(indexes, index)
		rest = rest[end+1:]
	}
	return name, indexes, nil
}

// JSONFlattenOptions controls how nested values are written into string cells.
type JSONFlattenOptions struct {
	// Arrays is "json" (serialize the whole array into one cell) or "index"
	// (expand into tags[0], tags[1], ...). Empty means "json".
	Arrays string
}

// FlatJSONRecord is the flattened string view of one JSON row plus the key
// order in which the keys first appeared.
type FlatJSONRecord struct {
	Keys      []string
	Values    map[string]string
	RawValues map[string]interface{}
}

// FlattenJSONValue flattens a JSON object into dot/bracket paths. It returns
// nil when the value is not an object (callers handle top-level scalars/arrays).
func FlattenJSONValue(value interface{}, opts JSONFlattenOptions) *FlatJSONRecord {
	if opts.Arrays == "" {
		opts.Arrays = "json"
	}
	record := &FlatJSONRecord{Values: make(map[string]string), RawValues: make(map[string]interface{})}
	if object, ok := value.(*OrderedObject); ok {
		for _, key := range object.Keys {
			flattenJSONEntry(joinJSONPath("", key), object.Values[key], record, opts)
		}
		return record
	}
	if object, ok := value.(map[string]interface{}); ok {
		keys := make([]string, 0, len(object))
		for key := range object {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			flattenJSONEntry(joinJSONPath("", key), object[key], record, opts)
		}
		return record
	}
	return nil
}

func flattenJSONEntry(path string, value interface{}, record *FlatJSONRecord, opts JSONFlattenOptions) {
	if path == "" {
		path = "value"
	}
	switch node := value.(type) {
	case *OrderedObject:
		for _, key := range node.Keys {
			flattenJSONEntry(joinJSONPath(path, key), node.Values[key], record, opts)
		}
	case map[string]interface{}:
		keys := make([]string, 0, len(node))
		for key := range node {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			flattenJSONEntry(joinJSONPath(path, key), node[key], record, opts)
		}
	case []interface{}:
		if opts.Arrays == "index" && len(node) > 0 {
			for i, child := range node {
				flattenJSONEntry(path+"["+strconv.Itoa(i)+"]", child, record, opts)
			}
			return
		}
		record.set(path, JSONCompactString(node), node)
	default:
		record.set(path, JSONScalarString(node), node)
	}
}

func (r *FlatJSONRecord) set(key, value string, raw interface{}) {
	if _, exists := r.Values[key]; !exists {
		r.Keys = append(r.Keys, key)
	}
	r.Values[key] = value
	r.RawValues[key] = raw
}

func joinJSONPath(prefix, key string) string {
	if !isSimpleJSONKey(key) {
		key = "[" + strconv.Quote(key) + "]"
		if prefix == "" {
			return key
		}
		return prefix + key
	}
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func isSimpleJSONKey(key string) bool {
	if key == "" {
		return false
	}
	for i, r := range key {
		if r == '_' || unicode.IsLetter(r) {
			continue
		}
		if i > 0 && unicode.IsDigit(r) {
			continue
		}
		return false
	}
	return true
}

// JSONScalarString renders a scalar without producing Go's fmt.Sprint syntax
// for composite values.
func JSONScalarString(value interface{}) string {
	switch node := value.(type) {
	case nil:
		return "NULL"
	case string:
		return node
	case json.Number:
		return node.String()
	case bool:
		if node {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(node, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(node), 'g', -1, 32)
	case int:
		return strconv.Itoa(node)
	case int64:
		return strconv.FormatInt(node, 10)
	case int32:
		return strconv.FormatInt(int64(node), 10)
	default:
		return JSONCompactString(value)
	}
}

// JSONCompactString serializes composite values as compact JSON instead of
// Go's "map[k:v]" / "[a b]" representation.
func JSONCompactString(value interface{}) string {
	plain := ToPlainJSON(value)
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(plain); err != nil {
		return JSONScalarString(value)
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

// ToPlainJSON converts OrderedObject values into regular maps/slices so they
// can be passed to encoding/json.
func ToPlainJSON(value interface{}) interface{} {
	switch node := value.(type) {
	case *OrderedObject:
		object := make(map[string]interface{}, len(node.Keys))
		for _, key := range node.Keys {
			object[key] = ToPlainJSON(node.Values[key])
		}
		return object
	case map[string]interface{}:
		object := make(map[string]interface{}, len(node))
		for key, child := range node {
			object[key] = ToPlainJSON(child)
		}
		return object
	case []interface{}:
		array := make([]interface{}, len(node))
		for i, child := range node {
			array[i] = ToPlainJSON(child)
		}
		return array
	default:
		return value
	}
}

// IsJSONSource reports whether format is one of the JSON-family input formats.
func IsJSONSource(format string) bool {
	switch strings.ToLower(format) {
	case "json", "jsonl", "jsonlines":
		return true
	default:
		return false
	}
}

// FlattenJSONRecords builds a table from decoded JSON rows. It preserves the
// first-seen key order unless sortHeaders is true.
func FlattenJSONRecords(records []interface{}, opts JSONFlattenOptions, sortHeaders bool) ([]string, [][]string) {
	headers, rows, _ := FlattenJSONRecordsWithRaw(records, opts, sortHeaders)
	return headers, rows
}

// FlattenJSONRecordsWithRaw is FlattenJSONRecords plus a typed raw view for
// JSON/SQL output. Raw cells use json.Number for numbers, bool for booleans,
// nil for nulls, and the original composite value for arrays/objects.
func FlattenJSONRecordsWithRaw(records []interface{}, opts JSONFlattenOptions, sortHeaders bool) ([]string, [][]string, [][]interface{}) {
	flatRecords := make([]*FlatJSONRecord, 0, len(records))
	seen := make(map[string]bool)
	headers := make([]string, 0)
	for _, record := range records {
		flat := FlattenJSONValue(record, opts)
		if flat == nil {
			// A scalar/array row is represented as a single "value" column so
			// that arrays like [1,2,3] still produce a usable table.
			flat = &FlatJSONRecord{
				Keys:      []string{"value"},
				Values:    map[string]string{"value": JSONScalarString(record)},
				RawValues: map[string]interface{}{"value": record},
			}
		}
		flatRecords = append(flatRecords, flat)
		for _, key := range flat.Keys {
			if !seen[key] {
				seen[key] = true
				headers = append(headers, key)
			}
		}
	}
	if sortHeaders {
		sort.Strings(headers)
	}
	rows := make([][]string, len(flatRecords))
	raw := make([][]interface{}, len(flatRecords))
	for i, flat := range flatRecords {
		row := make([]string, len(headers))
		rawRow := make([]interface{}, len(headers))
		for j, header := range headers {
			row[j] = flat.Values[header]
			if value, ok := flat.RawValues[header]; ok {
				rawRow[j] = value
			} else {
				// Missing keys are empty strings in the string view; keep the
				// raw view consistent instead of emitting JSON null.
				rawRow[j] = ""
			}
		}
		rows[i] = row
		raw[i] = rawRow
	}
	return headers, rows, raw
}
