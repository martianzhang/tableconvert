package html

import (
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

// Marshal converts a table structure to HTML format and writes it to the config's Writer.
// It generates a properly formatted HTML table with headers and data rows.
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

	writer := cfg.Writer

	// Parse configuration options
	useDiv := false
	if _, ok := cfg.Extension["div"]; ok {
		useDiv = true
	}

	shouldEscape := true
	if v, ok := cfg.Extension["escape"]; ok && strings.ToLower(v) == "false" {
		shouldEscape = false
	}

	minify := false
	if _, ok := cfg.Extension["minify"]; ok {
		minify = true
	}

	useThead := false
	if _, ok := cfg.Extension["thead"]; ok {
		useThead = true
	}

	nl := "\n"
	indent := "  "
	if minify {
		nl = ""
		indent = ""
	}

	// Helper function to escape HTML if needed
	escape := func(s string) string {
		if shouldEscape {
			return html.EscapeString(s)
		}
		return s
	}

	if useDiv {
		// Write div-based table
		writer.Write([]byte(fmt.Sprintf(`<div class="table">%s`, nl)))

		if useThead {
			writer.Write([]byte(fmt.Sprintf("%s<div class=\"thead\">%s", indent, nl)))
		}

		// Header row
		writer.Write([]byte(fmt.Sprintf(`%s<div class="tr">%s`, indent, nl)))
		for _, header := range table.Headers {
			writer.Write([]byte(fmt.Sprintf(`%s%s<div class="th">%s</div>%s`, indent, indent, escape(header), nl)))
		}
		writer.Write([]byte(fmt.Sprintf(`%s</div>%s`, indent, nl)))

		if useThead {
			writer.Write([]byte(fmt.Sprintf("%s</div>%s", indent, nl)))
			writer.Write([]byte(fmt.Sprintf("%s<div class=\"tbody\">%s", indent, nl)))
		}

		// Data rows
		for _, row := range table.Rows {
			writer.Write([]byte(fmt.Sprintf(`%s<div class="tr">%s`, indent, nl)))
			for _, cell := range row {
				cellClass := "td"
				writer.Write([]byte(fmt.Sprintf(`%s%s<div class="%s">%s</div>%s`, indent, indent, cellClass, escape(cell), nl)))
			}
			writer.Write([]byte(fmt.Sprintf(`%s</div>%s`, indent, nl)))
		}
		writer.Write([]byte("</div>"))
	} else {
		// Write traditional table
		writer.Write([]byte(fmt.Sprintf("<table>%s", nl)))

		if useThead {
			writer.Write([]byte(fmt.Sprintf("%s<thead>%s", indent, nl)))
		}

		// Header row
		writer.Write([]byte(fmt.Sprintf("%s<tr>%s", indent, nl)))
		for _, header := range table.Headers {
			writer.Write([]byte(fmt.Sprintf("%s%s<th>%s</th>%s", indent, indent, escape(header), nl)))
		}
		writer.Write([]byte(fmt.Sprintf("%s</tr>%s", indent, nl)))

		if useThead {
			writer.Write([]byte(fmt.Sprintf("%s</thead>%s%s<tbody>%s", indent, nl, indent, nl)))
		}

		// Data rows
		for _, row := range table.Rows {
			writer.Write([]byte(fmt.Sprintf("%s<tr>%s", indent, nl)))
			for _, cell := range row {
				tag := "td"
				writer.Write([]byte(fmt.Sprintf("%s%s<%s>%s</%s>%s", indent, indent, tag, escape(cell), tag, nl)))
			}
			writer.Write([]byte(fmt.Sprintf("%s</tr>%s", indent, nl)))
		}

		if useThead {
			writer.Write([]byte(fmt.Sprintf("%s</tbody>%s", indent, nl)))
		}

		writer.Write([]byte("</table>"))
	}

	return nil
}

// Unmarshal parses HTML table content from an io.Reader and populates the given Table struct.
// It extracts headers from <th> tags and data from <td> tags.
func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Unmarshal: output table cannot be nil")
	}

	// Reset the table fields to ensure clean population
	table.Headers = nil
	table.Rows = nil

	// Read all content
	content, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return fmt.Errorf("Unmarshal: failed to read input: %w", err)
	}

	htmlContent := string(content)

	// Extract the table content
	tableRegex := regexp.MustCompile(`(?s)<table.*?>(.*?)</table>`)
	tableMatch := tableRegex.FindStringSubmatch(htmlContent)
	if len(tableMatch) < 2 {
		return fmt.Errorf("Unmarshal: no table found in HTML content")
	}

	tableContent := tableMatch[1]

	// Extract header row
	headerRegex := regexp.MustCompile(`(?s)<tr.*?>(.*?)</tr>`)
	headerMatch := headerRegex.FindStringSubmatch(tableContent)
	if len(headerMatch) < 2 {
		return fmt.Errorf("Unmarshal: no header row found in table")
	}

	headerContent := headerMatch[1]

	// Extract header cells
	thRegex := regexp.MustCompile(`<th.*?>(.*?)</th>`)
	thMatches := thRegex.FindAllStringSubmatch(headerContent, -1)

	if len(thMatches) == 0 {
		return fmt.Errorf("Unmarshal: no header cells found in header row")
	}

	// Extract header values
	for _, match := range thMatches {
		if len(match) >= 2 {
			table.Headers = append(table.Headers, strings.TrimSpace(match[1]))
		}
	}

	// Extract data rows
	rowRegex := regexp.MustCompile(`(?s)<tr.*?>(.*?)</tr>`)
	rowMatches := rowRegex.FindAllStringSubmatch(tableContent, -1)

	// Skip the first row (header)
	for i := 1; i < len(rowMatches); i++ {
		if len(rowMatches[i]) < 2 {
			continue
		}

		rowContent := rowMatches[i][1]

		// Extract data cells
		tdRegex := regexp.MustCompile(`<td.*?>(.*?)</td>`)
		tdMatches := tdRegex.FindAllStringSubmatch(rowContent, -1)

		if len(tdMatches) > 0 {
			var row []string
			for _, match := range tdMatches {
				if len(match) >= 2 {
					row = append(row, strings.TrimSpace(match[1]))
				}
			}

			if len(row) > 0 {
				table.Rows = append(table.Rows, row)
			}
		}
	}

	return nil
}
