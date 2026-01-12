package common

import (
	"io"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectTableFormatByExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic extension",
			input:    ".csv",
			expected: TableFormatCSV,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := DetectTableFormatByExtension(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name     string
		input    *ParseError
		expected string
	}{
		{
			name: "basic error",
			input: &ParseError{
				LineNumber: 1,
				Message:    "invalid format",
				Line:       "header1,header2",
			},
			expected: `parse error on line 1: invalid format (line: "header1,header2")`,
		},
		{
			name: "empty line",
			input: &ParseError{
				LineNumber: 5,
				Message:    "empty line",
				Line:       "",
			},
			expected: `parse error on line 5: empty line (line: "")`,
		},
		{
			name: "special characters in line",
			input: &ParseError{
				LineNumber: 10,
				Message:    "invalid character",
				Line:       "data\twith\ttabs",
			},
			expected: `parse error on line 10: invalid character (line: "data\twith\ttabs")`,
		},
		{
			name: "large line number",
			input: &ParseError{
				LineNumber: 99999,
				Message:    "file too large",
				Line:       "very long line...",
			},
			expected: `parse error on line 99999: file too large (line: "very long line...")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.input.Error()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "uppercase null",
			input:    "NULL",
			expected: nil,
		},
		{
			name:     "lowercase null",
			input:    "null",
			expected: nil,
		},
		{
			name:     "mixed case null",
			input:    "Null",
			expected: nil,
		},
		{
			name:     "random case null",
			input:    "nUlL",
			expected: nil,
		},
		{
			name:     "true boolean",
			input:    "true",
			expected: true,
		},
		{
			name:     "false boolean",
			input:    "false",
			expected: false,
		},
		{
			name:     "mixed case true",
			input:    "TrUe",
			expected: true,
		},
		{
			name:     "mixed case false",
			input:    "FaLsE",
			expected: false,
		},
		{
			name:     "positive integer",
			input:    "123",
			expected: int64(123),
		},
		{
			name:     "negative integer",
			input:    "-456",
			expected: int64(-456),
		},
		{
			name:     "zero",
			input:    "0",
			expected: int64(0),
		},
		{
			name:     "positive float",
			input:    "123.456",
			expected: 123.456,
		},
		{
			name:     "negative float",
			input:    "-789.123",
			expected: -789.123,
		},
		{
			name:     "scientific notation",
			input:    "1.23e4",
			expected: 12300.0,
		},
		{
			name:     "float with leading zero",
			input:    "0.123",
			expected: 0.123,
		},
		{
			name:     "integer as float",
			input:    "123.0",
			expected: 123.0,
		},
		{
			name:     "invalid number",
			input:    "123abc",
			expected: "123abc",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace string",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "regular string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with numbers",
			input:    "123 main st",
			expected: "123 main st",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := InferType(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInferPrintType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase null",
			input:    "NULL",
			expected: "null",
		},
		{
			name:     "lowercase null",
			input:    "null",
			expected: "null",
		},
		{
			name:     "mixed case null",
			input:    "Null",
			expected: "null",
		},
		{
			name:     "random case null",
			input:    "nUlL",
			expected: "null",
		},
		{
			name:     "true boolean",
			input:    "true",
			expected: "true",
		},
		{
			name:     "false boolean",
			input:    "false",
			expected: "false",
		},
		{
			name:     "mixed case true",
			input:    "TrUe",
			expected: "true",
		},
		{
			name:     "mixed case false",
			input:    "FaLsE",
			expected: "false",
		},
		{
			name:     "positive integer",
			input:    "123",
			expected: "123",
		},
		{
			name:     "negative integer",
			input:    "-456",
			expected: "-456",
		},
		{
			name:     "zero",
			input:    "0",
			expected: "0",
		},
		{
			name:     "positive float",
			input:    "123.456",
			expected: "123.456",
		},
		{
			name:     "negative float",
			input:    "-789.123",
			expected: "-789.123",
		},
		{
			name:     "scientific notation",
			input:    "1.23e4",
			expected: "1.23e4",
		},
		{
			name:     "float with leading zero",
			input:    "0.123",
			expected: "0.123",
		},
		{
			name:     "integer as float",
			input:    "123.0",
			expected: "123.0",
		},
		{
			name:     "invalid number",
			input:    "123abc",
			expected: "123abc",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace string",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "regular string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with numbers",
			input:    "123 main st",
			expected: "123 main st",
		},
		{
			name:     "string with leading/trailing spaces",
			input:    "  test  ",
			expected: "  test  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := InferPrintType(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// TestDetectTableFormatByData tests format detection by content analysis
func TestDetectTableFormatByData(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		shouldErr bool
	}{
		{
			name:      "HTML table",
			input:     "<table><tr><td>test</td></tr></table>",
			expected:  TableFormatHTML,
			shouldErr: false,
		},
		{
			name:      "Markdown header",
			input:     "# Header\nSome text",
			expected:  TableFormatMarkdown,
			shouldErr: false,
		},
		{
			name:      "LaTeX table",
			input:     "\\begin{tabular}{cc}\nA & B \\\\\n\\end{tabular}",
			expected:  TableFormatLatex,
			shouldErr: false,
		},
		{
			name:      "MediaWiki table",
			input:     "{|\n| A || B\n|}",
			expected:  TableFormatMediaWiki,
			shouldErr: false,
		},
		{
			name:      "CSV",
			input:     "name,age\nAlice,30",
			expected:  TableFormatCSV,
			shouldErr: false,
		},
		{
			name:      "SQL INSERT",
			input:     "INSERT INTO table VALUES (1);",
			expected:  TableFormatSQL,
			shouldErr: false,
		},
		{
			name:      "JSON array",
			input:     `[{"name":"Alice"}]`,
			expected:  TableFormatJSON,
			shouldErr: false,
		},
		{
			name:      "JSON object",
			input:     `{"name":"Alice"}`,
			expected:  TableFormatJSON,
			shouldErr: false,
		},
		{
			name:      "XML",
			input:     "<root><item>test</item></root>",
			expected:  TableFormatXML,
			shouldErr: false,
		},
		{
			name:      "Unknown format",
			input:     "just plain text",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Empty input",
			input:     "",
			expected:  "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			format, err := DetectTableFormatByData(reader)

			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, format)
			}
		})
	}
}

// TestDetectTableFormatByExtension comprehensive tests
func TestDetectTableFormatByExtensionComprehensive(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"data.csv", TableFormatCSV},
		{"data.json", TableFormatJSON},
		{"data.md", TableFormatMarkdown},
		{"data.markdown", TableFormatMarkdown},
		{"data.xlsx", TableFormatExcel},
		{"data.xls", TableFormatExcel},
		{"data.html", TableFormatHTML},
		{"data.xml", TableFormatXML},
		{"data.sql", TableFormatSQL},
		{"data.tex", TableFormatLatex},
		{"data.wiki", TableFormatMediaWiki},
		{"data.txt", ""},
		{"data.unknown", ""},
		{"", ""},
		{"path/to/file.csv", TableFormatCSV},
		{"file.with.dots.csv", TableFormatCSV},
		{"UPPERCASE.CSV", TableFormatCSV},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := DetectTableFormatByExtension(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTableStruct tests the Table struct
func TestTableStruct(t *testing.T) {
	table := &Table{
		Headers: []string{"name", "age"},
		Rows: [][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	}

	assert.Len(t, table.Headers, 2)
	assert.Len(t, table.Rows, 2)
	assert.Equal(t, "Alice", table.Rows[0][0])
}

// TestParseErrorFormatting tests parse error formatting
func TestParseErrorFormatting(t *testing.T) {
	err := &ParseError{
		LineNumber: 42,
		Message:    "test error",
		Line:       "bad line",
	}

	formatted := err.Error()
	assert.Contains(t, formatted, "42")
	assert.Contains(t, formatted, "test error")
	assert.Contains(t, formatted, "bad line")
}

// TestFormatConstants tests format constants
func TestFormatConstants(t *testing.T) {
	assert.Equal(t, "excel", TableFormatExcel)
	assert.Equal(t, "csv", TableFormatCSV)
	assert.Equal(t, "markdown", TableFormatMarkdown)
	assert.Equal(t, "html", TableFormatHTML)
	assert.Equal(t, "mediawiki", TableFormatMediaWiki)
	assert.Equal(t, "latex", TableFormatLatex)
	assert.Equal(t, "json", TableFormatJSON)
	assert.Equal(t, "xml", TableFormatXML)
	assert.Equal(t, "ascii", TableFormatASCII)
	assert.Equal(t, "sql", TableFormatSQL)
	assert.Equal(t, "tracwiki", TableFormatTracWiki)
	assert.Equal(t, "twiki", TableFormatTWiki)
}

// TestInferTypeEdgeCases tests edge cases for type inference
func TestInferTypeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"very large integer", "999999999999999999999999", float64(1e+24)}, // Too large for int64, becomes float
		{"negative zero", "-0", int64(0)},
		{"positive zero", "+0", int64(0)},
		{"scientific large", "1e100", 1e100},
		{"scientific negative", "-1.5e-3", -0.0015},
		{"hex number", "0x10", "0x10"}, // Not supported, stays string
		{"octal number", "0123", int64(123)},
		{"binary number", "0b1010", "0b1010"}, // Not supported
		{"infinity", "Infinity", math.Inf(1)},
		{"null with spaces", "  null  ", nil},
		{"true with spaces", "  true  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InferType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Special test for NaN since NaN != NaN in Go
	t.Run("NaN", func(t *testing.T) {
		result := InferType("NaN")
		resultFloat, ok := result.(float64)
		assert.True(t, ok, "NaN should be returned as float64")
		assert.True(t, math.IsNaN(resultFloat), "Result should be NaN")
	})
}

// TestExtensionsMap tests the extensions mapping
func TestExtensionsMap(t *testing.T) {
	// This test verifies the extensions map is properly initialized
	// by checking a few key mappings
	assert.Equal(t, TableFormatExcel, DetectTableFormatByExtension("file.xlsx"))
	assert.Equal(t, TableFormatCSV, DetectTableFormatByExtension("file.csv"))
	assert.Equal(t, TableFormatJSON, DetectTableFormatByExtension("file.json"))
}

// TestFormatDetectionOrder tests that format detection follows correct priority
func TestFormatDetectionOrder(t *testing.T) {
	// Test that HTML is detected before Markdown (important for priority)
	htmlContent := "<table><tr><td>test</td></tr></table>"
	reader := strings.NewReader(htmlContent)
	format, err := DetectTableFormatByData(reader)
	assert.NoError(t, err)
	assert.Equal(t, TableFormatHTML, format)

	// Test that LaTeX is detected before MediaWiki
	latexContent := "\\begin{tabular}{c}\ntest\n\\end{tabular}"
	reader = strings.NewReader(latexContent)
	format, err = DetectTableFormatByData(reader)
	assert.NoError(t, err)
	assert.Equal(t, TableFormatLatex, format)
}

// TestReaderErrors tests error handling in DetectTableFormatByData
func TestReaderErrors(t *testing.T) {
	// Create a reader that always fails
	errorReader := &failingReader{}
	_, err := DetectTableFormatByData(errorReader)
	assert.Error(t, err)
}

// failingReader is a helper that always returns an error
type failingReader struct{}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
