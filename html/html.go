package html

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"golang.org/x/net/html"
)

const (
	divTemplate = `<div class="table">
  <div class="tr">{{range .Headers}}<div class="th">{{.}}</div>{{end}}</div>{{range .Rows}}
  <div class="tr">{{range .}}<div class="td">{{.}}</div>{{end}}</div>{{end}}
</div>`
	minifyDivTemplate = `<div class="table"><div class="tr">{{range .Headers}}<div class="th">{{.}}</div>{{end}}</div>{{range .Rows}}<div class="tr">{{range .}}<div class="td">{{.}}</div>{{end}}</div>{{end}}</div>`

	tableTemplate = `<table>
{{if .UseThead}}  <thead>
{{end}}  <tr>{{range .Headers}}<th>{{.}}</th>{{end}}</tr>{{if .UseThead}}
  </thead>
  <tbody>{{end}}
{{range .Rows}}  <tr>{{range .}}<td>{{.}}</td>{{end}}</tr>
{{end}}{{if .UseThead}}  </tbody>
{{end}}</table>`

	minfyTableTemplate = `{{if .UseThead}}<table><thead><tr>{{range .Headers}}<th>{{.}}</th>{{end}}</tr></thead><tbody>{{range .Rows}}<tr>{{range .}}<td>{{.}}</td>{{end}}</tr>{{end}}</tbody></table>{{else}}<table><tr>{{range .Headers}}<th>{{.}}</th>{{end}}</tr>{{range .Rows}}<tr>{{range .}}<td>{{.}}</td>{{end}}</tr>{{end}}</table>{{end}}`

	// Transposed templates
	transposedDivTemplate = `<div class="table">
{{range $i, $header := .Headers}}  <div class="tr"><div class="th">{{$header}}</div>{{range $.Rows}}<div class="td">{{index . $i}}</div>{{end}}</div>
{{end}}</div>`
	transposedMinifyDivTemplate = `<div class="table">{{range $i, $header := .Headers}}<div class="tr"><div class="th">{{$header}}</div>{{range $.Rows}}<div class="td">{{index . $i}}</div>{{end}}</div>{{end}}</div>`

	transposedTableTemplate = `<table>
{{range $i, $header := .Headers}}  <tr><th>{{$header}}</th>{{range $.Rows}}<td>{{index . $i}}</td>{{end}}</tr>
{{end}}</table>`
	transposedMinifyTableTemplate = `<table>{{range $i, $header := .Headers}}<tr><th>{{$header}}</th>{{range $.Rows}}<td>{{index . $i}}</td>{{end}}</tr>{{end}}</table>`
)

// Marshal converts a table structure to HTML format and writes it to the config's Writer.
func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}

	// Check if table has headers
	columnCounts := len(table.Headers)
	if columnCounts == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}

	// Validate row lengths
	for i, row := range table.Rows {
		if len(row) != columnCounts {
			return fmt.Errorf("Marshal: row %d has %d columns, but table has %d", i, len(row), columnCounts)
		}
	}

	// Parse configuration options
	useDiv := cfg.GetExtensionBool("div", false)
	minify := cfg.GetExtensionBool("minify", false)
	useThead := cfg.GetExtensionBool("thead", false)
	transpose := cfg.GetExtensionBool("transpose", false)

	writer := cfg.Writer

	// Choose the appropriate template based on configuration
	var templateStr string
	if minify {
		if useDiv {
			if transpose {
				templateStr = transposedMinifyDivTemplate
			} else {
				templateStr = minifyDivTemplate
			}
		} else {
			if transpose {
				templateStr = transposedMinifyTableTemplate
			} else {
				templateStr = minfyTableTemplate
			}
		}
	} else {
		if useDiv {
			if transpose {
				templateStr = transposedDivTemplate
			} else {
				templateStr = divTemplate
			}
		} else {
			if transpose {
				templateStr = transposedTableTemplate
			} else {
				templateStr = tableTemplate
			}
		}
	}

	tmpl, err := template.New("table").Parse(templateStr)
	if err != nil {
		return err
	}

	context := struct {
		UseThead bool
		Headers  []string
		Rows     [][]string
	}{
		UseThead: useThead && !transpose, // Disable thead when transposed
		Headers:  table.Headers,
		Rows:     table.Rows,
	}

	return tmpl.Execute(writer, context)
}

// Unmarshal parses HTML table content using the html parser
func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Unmarshal: output table cannot be nil")
	}

	// Reset the table
	table.Headers = nil
	table.Rows = nil

	// Parse the HTML
	doc, err := html.Parse(cfg.Reader)
	if err != nil {
		return fmt.Errorf("Unmarshal: failed to parse HTML: %w", err)
	}

	// Find the first table
	var tableNode *html.Node
	var findTable func(*html.Node)
	findTable = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tableNode = n
			return
		}
		for c := n.FirstChild; c != nil && tableNode == nil; c = c.NextSibling {
			findTable(c)
		}
	}
	findTable(doc)

	if tableNode == nil {
		return fmt.Errorf("Unmarshal: no table found in HTML content")
	}

	// Check if first-column-header is explicitly set in config
	if cfg.GetExtensionBool("first-column-header", false) {
		return parseFirstColumnAsHeader(tableNode, table)
	}

	// Auto-detect header type
	if isFirstColumnHeader(tableNode) {
		return parseFirstColumnAsHeader(tableNode, table)
	}
	return parseFirstRowAsHeader(tableNode, table)
}

// isFirstColumnHeader detects if the table uses first column as header
func isFirstColumnHeader(tableNode *html.Node) bool {
	// Check for traditional header (thead or th in first row)
	hasTraditionalHeader := false
	var checkTraditionalHeader func(*html.Node) bool
	checkTraditionalHeader = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "thead" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "tr" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "th" {
					return true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if checkTraditionalHeader(c) {
				return true
			}
		}
		return false
	}
	hasTraditionalHeader = checkTraditionalHeader(tableNode)

	// If traditional header exists, it's not first-column header
	if hasTraditionalHeader {
		return false
	}

	// Check for first-column header pattern
	// A table is considered first-column-header if:
	// 1. No thead or th in first row
	// 2. First cell of each row contains header-like content
	rows := collectRows(tableNode)
	if len(rows) == 0 {
		return false
	}

	// Check if first column looks like headers (more text, different formatting, etc.)
	// This is a simple heuristic - can be enhanced based on specific needs
	headerLikeCount := 0
	for _, row := range rows {
		if len(row) > 0 {
			firstCell := row[0]
			// Simple heuristic: if first cell is all caps or contains colon, it's likely a header
			if strings.ToUpper(firstCell) == firstCell || strings.Contains(firstCell, ":") {
				headerLikeCount++
			}
		}
	}

	// If majority of first cells look like headers, assume first-column-header
	return headerLikeCount > len(rows)/2
}

// collectRows collects all rows from the table
func collectRows(tableNode *html.Node) [][]string {
	var rows [][]string
	var processRows func(*html.Node)
	processRows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var row []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
					if c.FirstChild != nil {
						row = append(row, strings.TrimSpace(c.FirstChild.Data))
					} else {
						row = append(row, "")
					}
				}
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processRows(c)
		}
	}
	processRows(tableNode)
	return rows
}

// parseFirstRowAsHeader parses table with first row as header (default behavior)
func parseFirstRowAsHeader(tableNode *html.Node, table *common.Table) error {
	// Find first header row (tr with th cells)
	var headerRow *html.Node
	var findHeaderRow func(*html.Node) bool
	findHeaderRow = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "tr" {
			// Check if this row has th cells
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "th" {
					headerRow = n
					return true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findHeaderRow(c) {
				return true
			}
		}
		return false
	}

	if !findHeaderRow(tableNode) {
		return fmt.Errorf("Unmarshal: no header row found in table")
	}

	// Extract headers from header row
	var headers []string
	for c := headerRow.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "th" {
			// Handle empty th cells
			if c.FirstChild != nil {
				headers = append(headers, strings.TrimSpace(c.FirstChild.Data))
			} else {
				headers = append(headers, "")
			}
		}
	}

	if len(headers) == 0 {
		return fmt.Errorf("Unmarshal: no header cells found in header row")
	}

	table.Headers = headers

	// Process data rows
	var processRows func(*html.Node)
	processRows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			// Skip header rows (already processed)
			isHeaderRow := false
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "th" {
					isHeaderRow = true
					break
				}
			}
			if isHeaderRow {
				return
			}

			// Process data row
			var row []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					// Handle empty td cells
					if c.FirstChild != nil {
						row = append(row, strings.TrimSpace(c.FirstChild.Data))
					} else {
						row = append(row, "")
					}
				}
			}

			// Only add row if it has the same number of columns as headers
			if len(row) == len(headers) {
				table.Rows = append(table.Rows, row)
			} else if len(row) > 0 {
				// If row has fewer columns, pad with empty strings
				for len(row) < len(headers) {
					row = append(row, "")
				}
				table.Rows = append(table.Rows, row)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processRows(c)
		}
	}

	processRows(tableNode)

	return nil
}

// parseFirstColumnAsHeader parses table with first column as header
func parseFirstColumnAsHeader(tableNode *html.Node, table *common.Table) error {
	var headers []string
	var rows [][]string

	// Process all rows
	var processRows func(*html.Node)
	processRows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var row []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
					// Handle empty cells
					if c.FirstChild != nil {
						row = append(row, strings.TrimSpace(c.FirstChild.Data))
					} else {
						row = append(row, "")
					}
				}
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processRows(c)
		}
	}

	processRows(tableNode)

	if len(rows) == 0 {
		return fmt.Errorf("Unmarshal: no rows found in table")
	}

	// Extract headers from first column
	for _, row := range rows {
		if len(row) > 0 {
			headers = append(headers, row[0])
		} else {
			headers = append(headers, "")
		}
	}

	if len(headers) == 0 {
		return fmt.Errorf("Unmarshal: no header cells found in first column")
	}

	table.Headers = headers

	// Extract data (skip first column)
	for _, row := range rows {
		if len(row) > 1 {
			table.Rows = append(table.Rows, row[1:])
		} else {
			// If row only has header column, add empty row
			table.Rows = append(table.Rows, []string{})
		}
	}

	return nil
}
