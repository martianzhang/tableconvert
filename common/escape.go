package common

import (
	"html"
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
	specialChars := "\\`*_{}[]()#+-.!|~"
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\n', '\r':
			b.WriteString("<br>")
		case '\\':
			b.WriteString(`\\`)
		default:
			if strings.ContainsRune(specialChars, r) {
				b.WriteByte('\\')
			}
			b.WriteRune(r)
		}
	}
	return b.String()
}

func MarkdownUnescape(s string) string {
	const backslashMarker = "\x00"
	s = strings.ReplaceAll(s, `\\`, backslashMarker)

	replacements := map[rune]rune{
		' ': ' ', '!': '!', '#': '#', '$': '$', '%': '%', '&': '&',
		'(': '(', ')': ')', '*': '*', '+': '+', ',': ',', '-': '-',
		'.': '.', '/': '/', ':': ':', ';': ';', '<': '<', '=': '=',
		'>': '>', '?': '?', '@': '@', '[': '[', ']': ']', '^': '^',
		'_': '_', '`': '`', '{': '{', '|': '|', '}': '}', '~': '~',
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			next := rune(s[i+1])
			if _, ok := replacements[next]; ok {
				b.WriteRune(next)
				i++
				continue
			}
		}
		b.WriteByte(s[i])
	}

	out := strings.ReplaceAll(b.String(), backslashMarker, `\`)
	return strings.ReplaceAll(out, "<br>", "\n")
}

func HtmlEscape(s string) string {
	return html.EscapeString(s)
}

// LaTeXEscape escapes special LaTeX characters in the content
func LaTeXEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\textbackslash{}`)
		case '&':
			b.WriteString(`\&`)
		case '%':
			b.WriteString(`\%`)
		case '$':
			b.WriteString(`\$`)
		case '#':
			b.WriteString(`\#`)
		case '_':
			b.WriteString(`\_`)
		case '{':
			b.WriteString(`\{`)
		case '}':
			b.WriteString(`\}`)
		case '~':
			b.WriteString(`\textasciitilde{}`)
		case '^':
			b.WriteString(`\^{}`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func LaTeXUnescape(s string) string {
	// Use placeholders for characters produced by \text... commands so that
	// later unescaping passes cannot reinterpret the inserted backslashes.
	const (
		backslashMarker = "\x00"
		tildeMarker     = "\x01"
		caretMarker     = "\x02"
	)

	s = strings.ReplaceAll(s, `\textbackslash{}`, backslashMarker)
	s = strings.ReplaceAll(s, `\textasciitilde{}`, tildeMarker)
	s = strings.ReplaceAll(s, `\^{}`, caretMarker)
	s = strings.ReplaceAll(s, `\~{}`, tildeMarker)
	s = strings.ReplaceAll(s, `\textbackslash`, backslashMarker)
	s = strings.ReplaceAll(s, `\textasciitilde`, tildeMarker)

	// Handle basic escape characters. Replace longer forms before shorter ones.
	replacements := []struct{ from, to string }{
		{`\^`, "^"},
		{`\&`, "&"},
		{`\%`, "%"},
		{`\$`, "$"},
		{`\#`, "#"},
		{`\_`, "_"},
		{`\~`, "~"},
		{`\{`, "{"},
		{`\}`, "}"},
		{`\ `, " "},
	}
	for _, replacement := range replacements {
		s = strings.ReplaceAll(s, replacement.from, replacement.to)
	}

	s = strings.ReplaceAll(s, backslashMarker, `\`)
	s = strings.ReplaceAll(s, tildeMarker, "~")
	s = strings.ReplaceAll(s, caretMarker, "^")

	return s
}

// SQLValueEscape escapes values for SQL insertion using MySQL syntax.
// Kept for template helpers and backwards compatibility.
func SQLValueEscape(s string) string {
	return SQLValueEscapeDialect(s, "mysql")
}

// SQLValueEscapeDialect escapes a SQL string value using the requested
// dialect. It treats the text "NULL" as SQL NULL for backwards compatibility.
func SQLValueEscapeDialect(s, dialect string) string {
	return SQLValueEscapeDialectWithNull(s, dialect, true)
}

// SQLValueEscapeDialectWithNull escapes a SQL string value using the requested
// dialect. When nullSentinel is true, the text "NULL" is emitted as SQL NULL;
// otherwise it is quoted as an ordinary string.
func SQLValueEscapeDialectWithNull(s, dialect string, nullSentinel bool) string {
	if nullSentinel && strings.EqualFold(s, "NULL") {
		return "NULL"
	}

	switch dialect {
	case "postgres", "postgresql", "oracle", "mssql", "none":
		return "'" + strings.ReplaceAll(s, "'", "''") + "'"
	default: // mysql (also fallback)
		var b strings.Builder
		b.WriteByte('\'')
		for _, r := range s {
			switch r {
			case '\'':
				b.WriteString(`\'`)
			case '\\':
				b.WriteString(`\\`)
			case 0:
				b.WriteString(`\0`)
			case '\n':
				b.WriteString(`\n`)
			case '\r':
				b.WriteString(`\r`)
			case '\t':
				b.WriteString(`\t`)
			case 0x1a:
				b.WriteString(`\Z`)
			default:
				b.WriteRune(r)
			}
		}
		b.WriteByte('\'')
		return b.String()
	}
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
