package xml

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/martianzhang/tableconvert/common"
)

// Unmarshal parses XML data from the reader and populates the table
func Unmarshal(cfg *common.Config, table *common.Table) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	reader := cfg.Reader
	if reader == nil {
		return fmt.Errorf("reader is nil, please provide a valid reader")
	}

	if table == nil {
		return fmt.Errorf("output table cannot be nil")
	}

	// Reset the table fields to ensure clean population
	table.Headers = nil
	table.Rows = nil

	// Define the structure for dynamic XML parsing
	type GenericField struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}

	type GenericRow struct {
		XMLName xml.Name
		Fields  []GenericField `xml:",any"`
	}

	// Dynamic XML parsing
	var root struct {
		XMLName xml.Name
		Rows    []GenericRow `xml:",any"`
	}

	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&root); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	// Extract headers and row data
	if len(root.Rows) == 0 {
		return nil // Empty table is valid
	}

	// Extract headers from the first row
	for _, field := range root.Rows[0].Fields {
		table.Headers = append(table.Headers, field.XMLName.Local)
	}

	// Extract all row data and validate column counts
	headerCount := len(table.Headers)
	for rowIdx, row := range root.Rows {
		var dataRow []string
		for _, field := range row.Fields {
			dataRow = append(dataRow, field.Value)
		}

		// Validate column count matches headers
		if len(dataRow) != headerCount {
			return &common.ParseError{
				LineNumber: rowIdx + 1,
				Message:    fmt.Sprintf("data row has %d columns, but header has %d", len(dataRow), headerCount),
				Line:       fmt.Sprintf("row %d", rowIdx+1),
			}
		}

		table.Rows = append(table.Rows, dataRow)
	}

	return nil
}

// Marshal converts table data to XML format
func Marshal(cfg *common.Config, table *common.Table) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}

	if cfg.Writer == nil {
		return fmt.Errorf("writer is nil, please provide a valid writer")
	}

	// Get the configuration for minify
	minify := cfg.GetExtensionBool("minify", false)

	// Get the configuration for root-element and row-element
	rootElement := cfg.GetExtensionString("root-element", "dataset")
	rowElement := cfg.GetExtensionString("row-element", "record")

	// Validate element names are valid XML element names
	if !isValidXMLElementName(rootElement) {
		return fmt.Errorf("invalid root element name: %s", rootElement)
	}
	if !isValidXMLElementName(rowElement) {
		return fmt.Errorf("invalid row element name: %s", rowElement)
	}

	// Define the element structure
	type Cell struct {
		XMLName xml.Name `xml:""`
		Value   string   `xml:",chardata"`
	}

	type Row struct {
		XMLName xml.Name `xml:""`
		Cells   []Cell   `xml:""`
	}

	type Root struct {
		XMLName xml.Name `xml:""`
		Rows    []Row    `xml:""`
	}

	root := Root{
		XMLName: xml.Name{Local: rootElement},
	}

	// Build the data structure
	for rowIdx, row := range table.Rows {
		r := Row{
			XMLName: xml.Name{Local: rowElement},
		}

		// Validate row length matches headers
		if len(row) != len(table.Headers) {
			return fmt.Errorf("row %d has %d columns but headers have %d columns",
				rowIdx+1, len(row), len(table.Headers))
		}

		for i, cell := range row {
			header := table.Headers[i]
			if header == "" {
				header = "NULL"
			}

			// Validate header is a valid XML element name
			if !isValidXMLElementName(header) {
				return fmt.Errorf("invalid XML element name for header '%s' at index %d", header, i)
			}

			r.Cells = append(r.Cells, Cell{
				XMLName: xml.Name{Local: header},
				Value:   cell,
			})
		}
		root.Rows = append(root.Rows, r)
	}

	// Handle XML declaration separately if needed
	declaration := cfg.GetExtensionBool("declaration", false)

	// Create a temporary buffer if we need to handle declaration separately
	var writer io.Writer = cfg.Writer

	// Write the XML declaration if requested
	if declaration {
		if _, err := fmt.Fprintln(writer, `<?xml version="1.0" encoding="UTF-8" ?>`); err != nil {
			return fmt.Errorf("failed to write XML declaration: %w", err)
		}
	}

	// Use the Encoder from encoding/xml
	xmlEncoder := xml.NewEncoder(writer)
	if minify {
		xmlEncoder.Indent("", "")
	} else {
		xmlEncoder.Indent("", "  ")
	}

	// Encode and write XML
	if err := xmlEncoder.Encode(root); err != nil {
		return fmt.Errorf("failed to encode XML: %w", err)
	}

	// Flush the encoder to ensure all data is written
	if err := xmlEncoder.Flush(); err != nil {
		return fmt.Errorf("failed to flush XML encoder: %w", err)
	}

	return nil
}

// isValidXMLElementName checks if a string is a valid XML element name
func isValidXMLElementName(name string) bool {
	if name == "" {
		return false
	}

	// XML element names must start with a letter or underscore
	firstChar := name[0]
	if !((firstChar >= 'a' && firstChar <= 'z') ||
		(firstChar >= 'A' && firstChar <= 'Z') ||
		firstChar == '_') {
		return false
	}

	// Check remaining characters
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-' || char == '.') {
			return false
		}
	}

	return true
}
