package common

import (
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
)

// Define a mapping from file extensions to table formats
var extensions = map[string]string{
	".xlsx":      TableFormatExcel,
	".xls":       TableFormatExcel,
	".csv":       TableFormatCSV,
	".sql":       TableFormatSQL,
	".json":      TableFormatJSON,
	".tex":       TableFormatLatex,
	".xml":       TableFormatXML,
	".md":        TableFormatMarkdown,
	".markdown":  TableFormatMarkdown,
	".html":      TableFormatHTML,
	".mediawiki": TableFormatMediaWiki,
	".wiki":      TableFormatMediaWiki,
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
	case isCSV(strContent):
		return TableFormatCSV, nil
	case isSQL(strContent):
		return TableFormatSQL, nil
	case isJSON(strContent):
		return TableFormatJSON, nil
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
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, ",") {
			return true
		}
	}
	return false
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

// InferType attempts to convert a string value to a more specific type (bool, int64, float64, nil)
// If no conversion is successful, it returns the original string.
func InferType(value string) interface{} {
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
