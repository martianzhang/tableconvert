package tmpl

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/martianzhang/tableconvert/common"
)

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: target table pointer cannot be nil")
	}

	templateFile := cfg.GetExtensionString("template", "")
	if templateFile == "" {
		return fmt.Errorf("template file path is required (use --template=<file>)")
	}

	// Read file content into templateStr
	templateStr, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file %q: %w", templateFile, err)
	}

	// Create template object with helper functions
	tmpl, err := template.New("table").Funcs(template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"Lower": func(s string) string {
			return strings.ToLower(s)
		},
		"Capitalize": func(s string) string {
			return strings.Title(s)
		},
		"Sub": func(a, b int) int {
			return a - b
		},
		"Quote":                      strconv.Quote,
		"CSVForceQuote":              common.CSVForceQuote,
		"CSVQuoteEscape":             common.CSVQuoteEscape,
		"HtmlEscape":                 common.HtmlEscape,
		"MarkdownEscape":             common.MarkdownEscape,
		"LaTeXEscape":                common.LaTeXEscape,
		"SQLValueEscape":             common.SQLValueEscape,
		"SQLIdentifierEscape":        common.SQLIdentifierEscape,
		"OracleIdentifierEscape":     common.OracleIdentifierEscape,
		"MssqlIdentifierEscape":      common.MssqlIdentifierEscape,
		"PostgreSQLIdentifierEscape": common.PostgreSQLIdentifierEscape,
	}).Parse(string(templateStr))

	if err != nil {
		return fmt.Errorf("failed to parse template %q: %w", templateFile, err)
	}

	if err := tmpl.Execute(cfg.Writer, table); err != nil {
		return fmt.Errorf("failed to execute template %q: %w", templateFile, err)
	}

	return nil
}
