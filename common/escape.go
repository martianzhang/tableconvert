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

func MarkdownUnescape(s string) string {
	// Unescape special Markdown characters
	// First, replace escaped backslashes with a temporary marker
	s = strings.ReplaceAll(s, `\\`, "\x00")

	// Then replace other escapes
	replacements := map[string]string{
		`\ `: ` `,
		`\!`: `!`,
		`\#`: `#`,
		`\$`: `$`,
		`\%`: `%`,
		`\&`: `&`,
		`\(`: `(`,
		`\)`: `)`,
		`\*`: `*`,
		`\+`: `+`,
		`\,`: `,`,
		`\-`: `-`,
		`\.`: `.`,
		`\/`: `/`,
		`\:`: `:`,
		`\;`: `;`,
		`\<`: `<`,
		`\=`: `=`,
		`\>`: `>`,
		`\?`: `?`,
		`\@`: `@`,
		`\[`: `[`,
		`\]`: `]`,
		`\^`: `^`,
		`\_`: `_`,
		`\{`: `{`,
		`\|`: `|`,
		`\}`: `}`,
		`\~`: `~`,
	}
	for from, to := range replacements {
		s = strings.ReplaceAll(s, from, to)
	}

	// Handle backtick, double quote, and single quote
	s = strings.ReplaceAll(s, "`", "`")
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `'`, `'`)

	// Finally restore the single backslashes
	s = strings.ReplaceAll(s, "\x00", `\`)

	return s
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
