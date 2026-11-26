package html

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

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

	writer := cfg.Writer

	// Choose the appropriate template based on configuration
	var templateStr string
	if minify {
		if useDiv {
			templateStr = minifyDivTemplate
		} else {
			templateStr = minfyTableTemplate
		}
	} else {
		if useDiv {
			templateStr = divTemplate
		} else {
			templateStr = tableTemplate
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
		UseThead: useThead,
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

	reader := cfg.Reader
	if rawURL := cfg.GetExtensionString("url", ""); rawURL != "" {
		body, err := fetchHTML(rawURL)
		if err != nil {
			return err
		}
		defer body.Close()
		reader = body
	}

	// Parse the HTML
	doc, err := html.Parse(reader)
	if err != nil {
		return fmt.Errorf("Unmarshal: failed to parse HTML: %w", err)
	}

	tables := collectTables(doc)
	if len(tables) == 0 {
		if divTable := findDivTable(doc); divTable != nil {
			return parseDivTable(divTable, table)
		}
		return fmt.Errorf("Unmarshal: no table found in HTML content")
	}

	tableNumber := 1
	if raw := strings.TrimSpace(cfg.GetExtensionString("table-index", "")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return fmt.Errorf("Unmarshal: invalid --table-index %q (use a 1-based table number)", raw)
		}
		tableNumber = parsed
	}
	if tableNumber > len(tables) {
		return fmt.Errorf("Unmarshal: --table-index=%d out of range (%d tables found)", tableNumber, len(tables))
	}
	tableNode := tables[tableNumber-1]

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

func nodeClass(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return attr.Val
		}
	}
	return ""
}

func findDivTable(root *html.Node) *html.Node {
	var found *html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if found != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "div" && strings.Contains(nodeClass(n), "table") {
			found = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)
	return found
}

func parseDivTable(tableNode *html.Node, table *common.Table) error {
	var rows [][]string
	var headers []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "div" && strings.Contains(nodeClass(c), "tr") {
				var row []string
				hasHeader := false
				for cell := c.FirstChild; cell != nil; cell = cell.NextSibling {
					if cell.Type != html.ElementNode || cell.Data != "div" {
						continue
					}
					class := nodeClass(cell)
					if strings.Contains(class, "th") {
						hasHeader = true
						row = append(row, strings.TrimSpace(nodeText(cell)))
					} else if strings.Contains(class, "td") {
						row = append(row, strings.TrimSpace(nodeText(cell)))
					}
				}
				if len(row) > 0 {
					if hasHeader && len(headers) == 0 {
						headers = row
					} else {
						rows = append(rows, row)
					}
				}
				continue
			}
			if c.Type == html.ElementNode && c.Data == "div" && strings.Contains(nodeClass(c), "table") {
				continue
			}
			walk(c)
		}
	}
	walk(tableNode)

	if len(headers) == 0 {
		if len(rows) == 0 {
			return fmt.Errorf("Unmarshal: no rows found in div table")
		}
		headers = rows[0]
		rows = rows[1:]
	}
	return setRectangularTable(table, headers, rows)
}

func setRectangularTable(table *common.Table, headers []string, rows [][]string) error {
	maxCols := len(headers)
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	for len(headers) < maxCols {
		headers = append(headers, "")
	}
	outRows := make([][]string, len(rows))
	for i, row := range rows {
		outRows[i] = make([]string, maxCols)
		copy(outRows[i], row)
	}
	table.Headers = headers
	table.Rows = outRows
	return nil
}

// isFirstColumnHeader detects if the table uses first column as header.
func isFirstColumnHeader(tableNode *html.Node) bool {
	// A thead or a th in the first row means the table has a normal row header.
	rows := collectRows(tableNode)
	if len(rows) == 0 {
		return false
	}
	if hasTraditionalHeader(tableNode) {
		return false
	}

	// Check if first column looks like headers. Requiring an actual uppercase
	// letter (or colon) avoids treating numeric IDs, CJK text, or empty cells as
	// headers just because they contain no lowercase letters.
	headerLikeCount := 0
	for _, row := range rows {
		if len(row) > 0 && looksLikeHeaderCell(row[0]) {
			headerLikeCount++
		}
	}
	return headerLikeCount > len(rows)/2
}

func hasTraditionalHeader(tableNode *html.Node) bool {
	if findFirstElement(tableNode, "thead") != nil {
		return true
	}
	firstRow := findFirstElement(tableNode, "tr")
	if firstRow == nil {
		return false
	}
	for c := firstRow.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "th" {
			return true
		}
	}
	return false
}

func looksLikeHeaderCell(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if strings.Contains(s, ":") {
		return true
	}
	hasUpper := false
	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
		if unicode.IsUpper(r) {
			hasUpper = true
		}
	}
	return hasUpper
}

// findFirstElement returns the first descendant element with the given tag name.
func findFirstElement(root *html.Node, tag string) *html.Node {
	var found *html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if found != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == tag {
			found = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)
	return found
}

// collectRows collects all rows from the table (excluding nested tables) and
// expands colspan/rowspan so every logical column is represented.
func collectRows(tableNode *html.Node) [][]string {
	var rowNodes []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "table" {
				continue // don't merge rows from a nested table
			}
			if c.Type == html.ElementNode && c.Data == "tr" {
				rowNodes = append(rowNodes, c)
				continue
			}
			walk(c)
		}
	}
	walk(tableNode)

	pending := make(map[int]rowSpan)
	rows := make([][]string, 0, len(rowNodes))
	for _, rowNode := range rowNodes {
		var cells []*html.Node
		for cell := rowNode.FirstChild; cell != nil; cell = cell.NextSibling {
			if cell.Type == html.ElementNode && (cell.Data == "td" || cell.Data == "th") {
				cells = append(cells, cell)
			}
		}

		row := make([]string, 0)
		cellIndex := 0
		for col := 0; ; col++ {
			if active, ok := pending[col]; ok {
				ensureColumns(&row, col)
				row[col] = active.value
				if active.remaining <= 1 {
					delete(pending, col)
				} else {
					active.remaining--
					pending[col] = active
				}
				continue
			}
			if cellIndex >= len(cells) {
				if !hasPendingAtOrAfter(pending, col) {
					break
				}
				continue
			}
			cell := cells[cellIndex]
			cellIndex++
			value := strings.TrimSpace(nodeText(cell))
			colspan := spanValue(cell, "colspan")
			rowspan := spanValue(cell, "rowspan")
			for offset := 0; offset < colspan; offset++ {
				ensureColumns(&row, col+offset)
				row[col+offset] = value
				if rowspan > 1 {
					pending[col+offset] = rowSpan{value: value, remaining: rowspan - 1}
				}
			}
			col += colspan - 1
		}
		if len(row) > 0 {
			rows = append(rows, row)
		}
	}
	return rows
}

func collectTables(root *html.Node) []*html.Node {
	var tables []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "table" {
				tables = append(tables, c)
				continue
			}
			walk(c)
		}
	}
	walk(root)
	return tables
}

func ensureColumns(row *[]string, index int) {
	for len(*row) <= index {
		*row = append(*row, "")
	}
}

type rowSpan struct {
	value     string
	remaining int
}

func hasPendingAtOrAfter(pending map[int]rowSpan, col int) bool {
	for key := range pending {
		if key >= col {
			return true
		}
	}
	return false
}

func spanValue(node *html.Node, key string) int {
	for _, attr := range node.Attr {
		if strings.EqualFold(attr.Key, key) {
			if value, err := strconv.Atoi(attr.Val); err == nil && value > 0 {
				return value
			}
		}
	}
	return 1
}

func fetchHTML(rawURL string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: failed to fetch %s: %w", rawURL, err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		response.Body.Close()
		return nil, fmt.Errorf("Unmarshal: failed to fetch %s: HTTP %s", rawURL, response.Status)
	}
	return struct {
		io.Reader
		io.Closer
	}{Reader: io.LimitReader(response.Body, 10<<20), Closer: response.Body}, nil
}

// nodeText returns the visible text content of a node, including text nested
// inside formatting elements. <br> is rendered as a newline.
func nodeText(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			b.WriteString(n.Data)
		case html.ElementNode:
			if n.Data == "br" {
				b.WriteString("\n")
				return
			}
			if n.Data == "script" || n.Data == "style" {
				return
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				walk(c)
			}
		}
	}
	walk(n)
	return b.String()
}

// parseFirstRowAsHeader parses a table whose first row is the header. Tables
// without <th>/<thead> are supported too; the first row is still the header.
func parseFirstRowAsHeader(tableNode *html.Node, table *common.Table) error {
	rows := collectRows(tableNode)
	if len(rows) == 0 {
		return fmt.Errorf("Unmarshal: no rows found in table")
	}

	headers := append([]string(nil), rows[0]...)
	dataRows := rows[1:]

	// Keep the table rectangular. Extra cells in a data row create additional
	// empty-named columns instead of being silently dropped; short rows are
	// padded with empty strings.
	maxCols := len(headers)
	for _, row := range dataRows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	for len(headers) < maxCols {
		headers = append(headers, "")
	}
	for i, row := range dataRows {
		if len(row) < maxCols {
			dataRows[i] = append(row, make([]string, maxCols-len(row))...)
		}
	}

	table.Headers = headers
	table.Rows = dataRows
	return nil
}

// parseFirstColumnAsHeader parses table with first column as header
func parseFirstColumnAsHeader(tableNode *html.Node, table *common.Table) error {
	rows := collectRows(tableNode)
	if len(rows) == 0 {
		return fmt.Errorf("Unmarshal: no rows found in table")
	}

	headers := make([]string, len(rows))
	maxDataCols := 0
	for i, row := range rows {
		if len(row) > 0 {
			headers[i] = row[0]
			if len(row)-1 > maxDataCols {
				maxDataCols = len(row) - 1
			}
		}
	}

	if len(headers) == 0 {
		return fmt.Errorf("Unmarshal: no header cells found in first column")
	}

	table.Headers = headers
	table.Rows = make([][]string, len(rows))
	for i, row := range rows {
		data := make([]string, maxDataCols)
		if len(row) > 1 {
			copy(data, row[1:])
		}
		table.Rows[i] = data
	}

	return nil
}
