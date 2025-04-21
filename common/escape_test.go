package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLValueEscape(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		// Case when input is "NULL"
		{"NULL", "NULL"},
		// Case for a normal string
		{"test", "'test'"},
		// Case for a string containing a single quote
		{"it's a test", "'it\\'s a test'"},
		// Case for a string containing a double quote
		{"a \"quote\" here", "'a \"quote\" here'"},
		// Case for an empty string
		{"", "''"},
	}

	for _, tc := range testCases {
		result := SQLValueEscape(tc.input)
		assert.Equal(t, tc.expected, result, "Input: %s", tc.input)
	}
}

func TestCSVForceQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
		{
			name:     "no quotes",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "contains single quote",
			input:    `he"llo`,
			expected: `"he""llo"`,
		},
		{
			name:     "contains multiple quotes",
			input:    `he"ll"o`,
			expected: `"he""ll""o"`,
		},
		{
			name:     "already quoted",
			input:    `"hello"`,
			expected: `"""hello"""`,
		},
		{
			name:     "contains special characters",
			input:    `he,llo`,
			expected: `"he,llo"`,
		},
		{
			name:     "contains newline",
			input:    "he\nllo",
			expected: "\"he\nllo\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CSVForceQuote(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCSVQuoteEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "normal text",
			expected: "normal text",
		},
		{
			name:     "contains comma",
			input:    "text,with,comma",
			expected: "\"text,with,comma\"",
		},
		{
			name:     "contains quote",
			input:    "text\"with\"quote",
			expected: "\"text\"\"with\"\"quote\"",
		},
		{
			name:     "contains newline",
			input:    "text\nwith\nnewline",
			expected: "\"text\nwith\nnewline\"",
		},
		{
			name:     "contains all special characters",
			input:    "text,\"\nwith all",
			expected: "\"text,\"\"\nwith all\"",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only quote",
			input:    "\"",
			expected: "\"\"\"\"",
		},
		{
			name:     "multiple quotes",
			input:    "\"\"",
			expected: "\"\"\"\"\"\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CSVQuoteEscape(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestMarkdownEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no special characters",
			input:    "normal text",
			expected: "normal text",
		},
		{
			name:     "single special character",
			input:    "text with *star",
			expected: "text with \\*star",
		},
		{
			name:     "multiple special characters",
			input:    "a*b_c[d]e{f}g(h)i#j+k-l.m!n|o~p",
			expected: "a\\*b\\_c\\[d\\]e\\{f\\}g\\(h\\)i\\#j\\+k\\-l\\.m\\!n\\|o\\~p",
		},
		{
			name:     "consecutive special characters",
			input:    "***bold***",
			expected: "\\*\\*\\*bold\\*\\*\\*",
		},
		{
			name:     "backslash at start",
			input:    "\\start",
			expected: "\\\\start",
		},
		{
			name:     "backslash at end",
			input:    "end\\",
			expected: "end\\\\",
		},
		{
			name:     "mixed special and normal",
			input:    "Hello [World]! This_is_a_test.",
			expected: "Hello \\[World\\]\\! This\\_is\\_a\\_test\\.",
		},
		{
			name:     "unicode characters",
			input:    "中文*测试",
			expected: "中文\\*测试",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkdownEscape(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHtmlEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no special characters",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "with HTML special characters",
			input:    "<div>Hello & World</div>",
			expected: "&lt;div&gt;Hello &amp; World&lt;/div&gt;",
		},
		{
			name:     "with quotes",
			input:    `"Hello" 'World'`,
			expected: "&#34;Hello&#34; &#39;World&#39;",
		},
		{
			name:     "mixed content",
			input:    "Normal text <script>alert('xss')</script>",
			expected: "Normal text &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := HtmlEscape(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLaTeXEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no special characters",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "backslash",
			input:    "\\",
			expected: "\\\\",
		},
		{
			name:     "ampersand",
			input:    "&",
			expected: "\\&",
		},
		{
			name:     "percent",
			input:    "%",
			expected: "\\%",
		},
		{
			name:     "dollar",
			input:    "$",
			expected: "\\$",
		},
		{
			name:     "hash",
			input:    "#",
			expected: "\\#",
		},
		{
			name:     "underscore",
			input:    "_",
			expected: "\\_",
		},
		{
			name:     "open brace",
			input:    "{",
			expected: "\\{",
		},
		{
			name:     "close brace",
			input:    "}",
			expected: "\\}",
		},
		{
			name:     "tilde",
			input:    "~",
			expected: "\\~",
		},
		{
			name:     "caret",
			input:    "^",
			expected: "\\^",
		},
		{
			name:     "multiple special characters",
			input:    "a\\b&c%d$e#f_g{h}i~j^k",
			expected: "a\\\\b\\&c\\%d\\$e\\#f\\_g\\{h\\}i\\~j\\^k",
		},
		{
			name:     "mixed characters",
			input:    "Hello\\World & 100%",
			expected: "Hello\\\\World \\& 100\\%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, LaTeXEscape(tt.input))
		})
	}
}

func TestSQLIdentifierEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "``",
		},
		{
			name:     "no backticks",
			input:    "table_name",
			expected: "`table_name`",
		},
		{
			name:     "single backtick",
			input:    "table`name",
			expected: "`table\\`name`",
		},
		{
			name:     "multiple backticks",
			input:    "table`name`with`ticks",
			expected: "`table\\`name\\`with\\`ticks`",
		},
		{
			name:     "backtick at start",
			input:    "`tablename",
			expected: "`\\`tablename`",
		},
		{
			name:     "backtick at end",
			input:    "tablename`",
			expected: "`tablename\\``",
		},
		{
			name:     "only backticks",
			input:    "```",
			expected: "`\\`\\`\\``",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SQLIdentifierEscape(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestOracleIdentifierEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "\"\"",
		},
		{
			name:     "no quotes",
			input:    "table_name",
			expected: "\"table_name\"",
		},
		{
			name:     "single quote",
			input:    `"column"`,
			expected: `"""column"""`,
		},
		{
			name:     "multiple quotes",
			input:    `"col""umn"`,
			expected: `"""col""""umn"""`,
		},
		{
			name:     "mixed characters",
			input:    `abc"123"xyz`,
			expected: `"abc""123""xyz"`,
		},
		{
			name:     "only quote",
			input:    `"`,
			expected: `""""`,
		},
		{
			name:     "multiple consecutive quotes",
			input:    `""""`,
			expected: `""""""""""`, // 修正为10个引号
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := OracleIdentifierEscape(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPostgreSQLIdentifierEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "\"\"",
		},
		{
			name:     "no quotes",
			input:    "table_name",
			expected: "\"table_name\"",
		},
		{
			name:     "with single quote",
			input:    "table\"name",
			expected: "\"table\\\"name\"",
		},
		{
			name:     "multiple quotes",
			input:    "\"table\"\"name\"",
			expected: "\"\\\"table\\\"\\\"name\\\"\"",
		},
		{
			name:     "special characters",
			input:    "table$name#123",
			expected: "\"table$name#123\"",
		},
		{
			name:     "spaces in name",
			input:    "table name",
			expected: "\"table name\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := PostgreSQLIdentifierEscape(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestMssqlIdentifierEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "[]",
		},
		{
			name:     "no special characters",
			input:    "table_name",
			expected: "[table_name]",
		},
		{
			name:     "contains single [",
			input:    "table[name",
			expected: "[table[]name]",
		},
		{
			name:     "contains multiple [",
			input:    "[table][name]",
			expected: "[[]table][]name]]", // 修正为正确的转义
		},
		{
			name:     "contains ] without [",
			input:    "table]name",
			expected: "[table]name]",
		},
		{
			name:     "contains both [ and ]",
			input:    "table[name]value",
			expected: "[table[]name]value]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MssqlIdentifierEscape(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLaTeXUnescape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no replacements",
			input:    "normal text",
			expected: "normal text",
		},
		{
			name:     "single replacement",
			input:    "\\&",
			expected: "&",
		},
		{
			name:     "multiple replacements",
			input:    "\\% \\$ \\#",
			expected: "% $ #",
		},
		{
			name:     "mixed content",
			input:    "a\\_b\\^{}c\\~d",
			expected: "a_b^c~d",
		},
		{
			name:     "text commands",
			input:    "\\textasciitilde \\textbackslash",
			expected: "~ \\",
		},
		{
			name:     "with spaces",
			input:    "  hello\\ world  ",
			expected: "hello world",
		},
		{
			name:     "nested replacements",
			input:    "\\{\\textasciitilde\\}",
			expected: "{~}",
		},
		{
			name:     "complex case",
			input:    "\\$100 \\%50 \\& \\_underscore\\^{} \\textasciitilde{}",
			expected: "$100 %50 & _underscore^ ~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := LaTeXUnescape(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
