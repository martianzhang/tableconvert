package common

import "fmt"

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
)

// DetectTableFormatByExtension detects the table format by the file extension.
func DetectTableFormatByExtension(filename string) string {
	filename = strings.ToLower(filename)

	// 定义文件后缀到格式的映射
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

	// 遍历映射检查后缀
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
		return TableFormatLaTeX, nil
	case isMediaWiki(strContent):
		return TableFormatMediaWiki, nil
		// TODO: csv, sql, json, xml, mysql
	default:
		return "", fmt.Errorf("unsupported file format")
	}
}

func isHTML(content string) bool {
	return strings.Contains(content, "<html") ||
		strings.Contains(content, "<!DOCTYPE html")
}

func isMarkdown(content string) bool {
	return regexp.MustCompile(`^#+\s|^-\s|^\*\s`).MatchString(content)
}

func isLaTeX(content string) bool {
	return strings.Contains(content, "\\documentclass") ||
		strings.Contains(content, "\\begin{document}")
}

func isMediaWiki(content string) bool {
	return strings.Contains(content, "{|") ||
		strings.Contains(content, "|}")
}
