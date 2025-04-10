package common

import (
	"fmt"
	"io"
	"regexp"
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
	TableFormatMySQL     = "mysql"
	TableFormatHTML      = "html"
	TableFormatMediaWiki = "mediawiki"
	TableFormatLatex     = "latex"
	TableFormatJSON      = "json"
	TableFormatXML       = "xml"
)

// DetectTableFormatByExtension detects the table format by the file extension.
func DetectTableFormatByExtension(filename string) string {
	// Check if the input filename is empty
	if filename == "" {
		return ""
	}
	// Convert the filename to lowercase
	filename = strings.ToLower(filename)

	// Define a mapping from file extensions to table formats
	extensions := map[string]string{
		".xlsx":      TableFormatExcel,
		".xls":       TableFormatExcel,
		".csv":       TableFormatCSV,
		".sql":       TableFormatMySQL,
		".json":      TableFormatJSON,
		".tex":       TableFormatLatex,
		".xml":       TableFormatXML,
		".md":        TableFormatMarkdown,
		".markdown":  TableFormatMarkdown,
		".html":      TableFormatHTML,
		".mediawiki": TableFormatMediaWiki,
		".wiki":      TableFormatMediaWiki,
	}

	// Iterate through the mapping to check the file extension
	for ext, format := range extensions {
		if strings.HasSuffix(filename, ext) {
			return format
		}
	}

	return ""
}

func DetectTableFormatByData(reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	strContent := string(content)

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
		return TableFormatMySQL, nil
	case isJSON(strContent):
		return TableFormatJSON, nil
	case isXML(strContent):
		return TableFormatXML, nil
	default:
		return "", fmt.Errorf("unsupported file format")
	}
}

func isHTML(content string) bool {
	return strings.Contains(content, "<html") ||
		strings.Contains(content, "<!DOCTYPE html") ||
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
	return strings.Contains(content, "\\documentclass") ||
		strings.Contains(content, "\\begin{document}")
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
	return strings.Contains(strings.ToLower(content), "select") ||
		strings.Contains(strings.ToLower(content), "insert") ||
		strings.Contains(strings.ToLower(content), "update") ||
		strings.Contains(strings.ToLower(content), "delete")
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
