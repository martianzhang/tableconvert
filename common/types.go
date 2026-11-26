package common

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Table represents the parsed data from the MySQL output.
type Table struct {
	Headers []string
	Rows    [][]string
	// RawRows optionally holds the original JSON rows. It is used by JSON/JSONL
	// output to preserve nested objects, nulls, and number text when no
	// destructive transformation has been applied.
	RawRows []interface{}
	// Raw optionally holds the flattened scalar values for each cell, in the
	// same shape as Rows. It lets JSON/SQL output keep number/bool/null types
	// after a transformation that only reorders rows.
	Raw [][]interface{}
}

// ParseError represents an error during parsing.
type ParseError struct {
	LineNumber int
	Message    string
	Line       string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error on line %d: %s (line: %q)", e.LineNumber, e.Message, e.Line)
}

const (
	TableFormatExcel     = "excel"
	TableFormatCSV       = "csv"
	TableFormatMarkdown  = "markdown"
	TableFormatHTML      = "html"
	TableFormatMediaWiki = "mediawiki"
	TableFormatLatex     = "latex"
	TableFormatJSON      = "json"
	TableFormatXML       = "xml"
	TableFormatASCII     = "ascii"
	TableFormatSQL       = "sql"
	TableFormatTracWiki  = "tracwiki"
	TableFormatTWiki     = "twiki"
	TableFormatTSV       = "tsv"
	TableFormatJSONL     = "jsonl"
	TableFormatMySQL     = "mysql"
)

var (
	integerPattern = regexp.MustCompile(`^[+-]?(0|[1-9][0-9]*)$`)
	floatPattern   = regexp.MustCompile(`^[+-]?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$`)
)

// Define a mapping from file extensions to table formats
var extensions = map[string]string{
	".xlsx":      TableFormatExcel,
	".xls":       TableFormatExcel,
	".csv":       TableFormatCSV,
	".tsv":       TableFormatTSV,
	".tab":       TableFormatTSV,
	".sql":       TableFormatSQL,
	".json":      TableFormatJSON,
	".jsonl":     TableFormatJSONL,
	".jsonlines": TableFormatJSONL,
	".tex":       TableFormatLatex,
	".latex":     TableFormatLatex,
	".xml":       TableFormatXML,
	".md":        TableFormatMarkdown,
	".markdown":  TableFormatMarkdown,
	".html":      TableFormatHTML,
	".htm":       TableFormatHTML,
	".mediawiki": TableFormatMediaWiki,
	".wiki":      TableFormatMediaWiki,
	".twiki":     TableFormatTWiki,
	".tracwiki":  TableFormatTWiki,
	".tmpl":      "tmpl",
	".template":  "tmpl",
	".ascii":     TableFormatASCII,
}

// DetectTableFormatByExtension detects the table format by the file extension.
func DetectTableFormatByExtension(filename string) string {
	// Check if the input filename is empty
	if filename == "" {
		return ""
	}

	// Get the file extension in lowercase
	ext := strings.ToLower(filepath.Ext(filename))

	// Look up the extension in the map
	if format, ok := extensions[ext]; ok {
		return format
	}

	return ""
}

func DetectTableFormatByData(reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	strContent := string(content)

	// Attention: file type detection order is important
	switch {
	case isHTML(strContent):
		return TableFormatHTML, nil
	case isMarkdown(strContent):
		return TableFormatMarkdown, nil
	case isLaTeX(strContent):
		return TableFormatLatex, nil
	case isMediaWiki(strContent):
		return TableFormatMediaWiki, nil
	case isJSONLines(strContent):
		return TableFormatJSONL, nil
	case isJSON(strContent):
		return TableFormatJSON, nil
	case isMySQL(strContent):
		return TableFormatMySQL, nil
	case isSQL(strContent):
		return TableFormatSQL, nil
	case isTSV(strContent):
		return TableFormatTSV, nil
	case isCSV(strContent):
		return TableFormatCSV, nil
	case isXML(strContent):
		return TableFormatXML, nil
	default:
		return "", fmt.Errorf("unsupported file format")
	}
}

func isHTML(content string) bool {
	return strings.Contains(content, "<thead") ||
		strings.Contains(content, "<table")
}

func isMarkdown(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, "|") || !strings.Contains(trimmed, "-") {
			continue
		}
		valid := true
		for _, r := range trimmed {
			if r != '|' && r != '-' && r != ':' && r != ' ' {
				valid = false
				break
			}
		}
		if valid {
			return true
		}
	}
	re, err := regexp.Compile(`^#+\s|^-\s|^\*\s`)
	if err != nil {
		return false
	}
	return re.MatchString(content)
}

func isLaTeX(content string) bool {
	return strings.Contains(content, "\\hline") ||
		strings.Contains(content, "\\begin{tabular}")
}

func isMediaWiki(content string) bool {
	return strings.Contains(content, "{|") ||
		strings.Contains(content, "|}")
}

func isCSV(content string) bool {
	lines := nonEmptyLines(content)
	if len(lines) < 2 {
		return false
	}
	commaCount := strings.Count(lines[0], ",")
	if commaCount == 0 {
		return false
	}
	for _, line := range lines[1:] {
		if strings.Count(line, ",") != commaCount {
			// Quoted fields can contain commas; fall back to a loose check.
			if !strings.Contains(line, ",") {
				return false
			}
		}
	}
	return true
}

func isTSV(content string) bool {
	lines := nonEmptyLines(content)
	if len(lines) < 2 {
		return false
	}
	tabCount := strings.Count(lines[0], "\t")
	if tabCount == 0 {
		return false
	}
	for _, line := range lines[1:] {
		if strings.Count(line, "\t") != tabCount {
			return false
		}
	}
	return true
}

func isJSONLines(content string) bool {
	lines := nonEmptyLines(content)
	if len(lines) < 2 {
		return false
	}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
			return false
		}
		if !json.Valid([]byte(trimmed)) {
			return false
		}
	}
	return true
}

func isMySQL(content string) bool {
	hasBorder := false
	hasPipe := false
	for _, line := range nonEmptyLines(content) {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "+-") || strings.HasPrefix(trimmed, "+=") {
			hasBorder = true
		}
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			hasPipe = true
		}
	}
	return hasBorder && hasPipe
}

func nonEmptyLines(content string) []string {
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimRight(line, "\r"))
		}
	}
	return lines
}

func isSQL(content string) bool {
	return strings.Contains(strings.ToLower(content), "replace") ||
		strings.Contains(strings.ToLower(content), "insert")
}

func isJSON(content string) bool {
	content = strings.TrimSpace(content)
	return (strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}")) ||
		(strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]"))
}

func isXML(content string) bool {
	return strings.Contains(content, "<") && strings.Contains(content, ">") &&
		strings.Contains(content, "</")
}

// InferType attempts to convert a string value to a more specific type (bool, int64, float64, nil).
// It is intentionally conservative: values with leading zeros, phone numbers,
// long IDs, or more than 15 digits stay strings so data such as ZIP codes and
// account numbers is not silently corrupted.
func InferType(value string) interface{} {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return value
	}

	// 1. Explicit null (case-insensitive).
	if strings.EqualFold(trimmedValue, "null") {
		return nil
	}

	// 2. Boolean (case-insensitive).
	lowerValue := strings.ToLower(trimmedValue)
	if lowerValue == "true" {
		return true
	}
	if lowerValue == "false" {
		return false
	}

	// 3. Integer. Reject leading zeros and numbers that lose precision when
	// represented as a float64 (more than 15 significant digits).
	if integerPattern.MatchString(trimmedValue) && digitCount(trimmedValue) <= 15 {
		if intVal, err := strconv.ParseInt(trimmedValue, 10, 64); err == nil {
			return intVal
		}
	}

	// 4. Float. Require a decimal point or exponent so "123" is handled above.
	if floatPattern.MatchString(trimmedValue) &&
		strings.ContainsAny(trimmedValue, ".eE") &&
		digitCount(trimmedValue) <= 15 {
		if floatVal, err := strconv.ParseFloat(trimmedValue, 64); err == nil {
			return floatVal
		}
	}

	// 5. Default: keep the original string.
	return value
}

func digitCount(value string) int {
	count := 0
	for _, r := range value {
		if r >= '0' && r <= '9' {
			count++
		}
	}
	return count
}

func InferPrintType(s string) string {
	switch InferType(s).(type) {
	case nil:
		return "null"
	case bool:
		return strings.ToLower(s)
	default:
		return fmt.Sprint(s)
	}
}
