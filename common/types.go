package common

import (
	"fmt"
	"io"
	"path/filepath"
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

// DetectTableFormatByExtension detects the table format by the file extension.
func DetectTableFormatByExtension(filename string) string {
	// Check if the input filename is empty
	if filename == "" {
		return ""
	}

	// Get the file extension in lowercase
	ext := strings.ToLower(filepath.Ext(filename))

	// Define a mapping from file extensions to table formats
	extensions := map[string]string{
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
