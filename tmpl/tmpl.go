package tmpl

import (
	"fmt"
	"io/ioutil"
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

	// Read file content into templateStr
	templateStr, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	// Create template object with html escape function
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
		return err
	}

	return tmpl.Execute(cfg.Writer, table)
}
