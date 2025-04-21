package common

import (
	"html"
	"strconv"
	"strings"
)

// CSVForceQuote returns a string with quotes
func CSVForceQuote(field string) string {
	// Escape quotes by doubling them
	field = strings.ReplaceAll(field, "\"", "\"\"")
	// Wrap the field in quotes
	return "\"" + field + "\""
}

// CSVQuoteEscape function is used to escape quotes in CSV fields
func CSVQuoteEscape(field string) string {
	// Check if the field contains commas, quotes, or newlines
	if strings.ContainsAny(field, ",\"\n") {
		// Escape quotes by doubling them
		field = strings.ReplaceAll(field, "\"", "\"\"")
		// Wrap the field in quotes
		return "\"" + field + "\""
	}
	return field
}

func MarkdownEscape(s string) string {
	// Escape special Markdown characters
	specialChars := []string{"\\", "`", "*", "_", "{", "}", "[", "]", "(", ")", "#", "+", "-", ".", "!", "|", "~"}
	for _, char := range specialChars {
		s = strings.ReplaceAll(s, char, "\\"+char)
	}
	return s
}

func HtmlEscape(s string) string {
	return html.EscapeString(s)
}

// LaTeXEscape escapes special LaTeX characters in the content
func LaTeXEscape(s string) string {
	// List of LaTeX special characters that need escaping
	specialChars := []string{"\\", "&", "%", "$", "#", "_", "{", "}", "~", "^"}
	for _, char := range specialChars {
		s = strings.ReplaceAll(s, char, "\\"+char)
	}
	return s
}

func LaTeXUnescape(s string) string {
	// First handle special characters with {}
	s = strings.ReplaceAll(s, `\^{}`, "^")
	s = strings.ReplaceAll(s, `\~{}`, "~")
	s = strings.ReplaceAll(s, `\textasciitilde{}`, "~")

	// Handle basic escape characters
	replacements := map[string]string{
		`\&`:              "&",
		`\%`:              "%",
		`\$`:              "$",
		`\#`:              "#",
		`\_`:              "_",
		`\~`:              "~",
		`\{`:              "{",
		`\}`:              "}",
		`\textasciitilde`: "~",
		`\textbackslash`:  `\`,
		`\ `:              " ", // Handle LaTeX space escape
	}

	for from, to := range replacements {
		s = strings.ReplaceAll(s, from, to)
	}

	// Trim extra whitespace
	s = strings.TrimSpace(s)

	return s
}

// SQLValueEscape values for SQL insertion (handle quotes, NULLs, etc.)
func SQLValueEscape(s string) string {
	if strings.EqualFold(s, "NULL") {
		return "NULL"
	}

	// Double Quote string
	s = strconv.Quote(s)
	// Convert to Single Quote string
	s = s[1 : len(s)-1]
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "'", "\\'")
	return "'" + s + "'"
}

// SQLIdentifierEscape mysql identifier escape
func SQLIdentifierEscape(s string) string {
	return "`" + strings.ReplaceAll(s, "`", "\\`") + "`"
}

func OracleIdentifierEscape(s string) string {
	return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
}

func PostgreSQLIdentifierEscape(s string) string {
	return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""
}

func MssqlIdentifierEscape(s string) string {
	return "[" + strings.ReplaceAll(s, "[", "[]") + "]"
}
