package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/martianzhang/tableconvert/common"
)

type xmlNode struct {
	Name     xml.Name
	Text     string
	Children []xmlNode
}

func (n xmlNode) fullText() string {
	return n.Text
}

func readXMLNode(decoder *xml.Decoder, start xml.StartElement) (xmlNode, error) {
	node := xmlNode{Name: start.Name}
	var text strings.Builder
	for {
		token, err := decoder.Token()
		if err != nil {
			return node, err
		}
		switch t := token.(type) {
		case xml.CharData:
			text.Write([]byte(t))
		case xml.StartElement:
			child, err := readXMLNode(decoder, t)
			if err != nil {
				return node, err
			}
			node.Children = append(node.Children, child)
			text.WriteString(child.fullText())
		case xml.EndElement:
			node.Text = text.String()
			return node, nil
		}
	}
}

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

	decoder := xml.NewDecoder(reader)
	var root *xmlNode
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to decode XML: %w", err)
		}
		if start, ok := token.(xml.StartElement); ok {
			node, err := readXMLNode(decoder, start)
			if err != nil {
				return fmt.Errorf("failed to decode XML: %w", err)
			}
			root = &node
			break
		}
	}

	if root == nil {
		return nil // Empty XML is a valid empty table.
	}

	rowElement := cfg.GetExtensionString("row-element", "")
	records, err := selectXMLRecords(root, rowElement)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}

	headerIndex := make(map[string]int)
	for _, record := range records {
		for _, field := range recordFields(record) {
			name := field.Name.Local
			if _, exists := headerIndex[name]; !exists {
				headerIndex[name] = len(table.Headers)
				table.Headers = append(table.Headers, name)
			}
		}
	}

	for _, record := range records {
		row := make([]string, len(table.Headers))
		for _, field := range recordFields(record) {
			index := headerIndex[field.Name.Local]
			if row[index] == "" {
				row[index] = strings.TrimSpace(field.fullText())
			}
		}
		table.Rows = append(table.Rows, row)
	}
	return nil
}

// recordFields returns a record's field elements. A scalar record (for example
// <record>42</record>) is represented as a synthetic "value" column.
func recordFields(record xmlNode) []xmlNode {
	if len(record.Children) > 0 {
		return record.Children
	}
	if strings.TrimSpace(record.Text) != "" {
		return []xmlNode{{Name: xml.Name{Local: "value"}, Text: record.Text}}
	}
	return nil
}

func selectXMLRecords(root *xmlNode, explicitRowElement string) ([]xmlNode, error) {
	if len(root.Children) == 0 {
		return nil, nil // <root></root> is an empty table.
	}

	if explicitRowElement != "" {
		var records []xmlNode
		for _, child := range root.Children {
			if child.Name.Local == explicitRowElement {
				records = append(records, child)
			}
		}
		if len(records) == 0 {
			return nil, fmt.Errorf("no <%s> elements found under <%s>; check --row-element", explicitRowElement, root.Name.Local)
		}
		return records, nil
	}

	type group struct {
		name  string
		nodes []xmlNode
	}
	groups := make([]group, 0)
	groupIndex := make(map[string]int)
	for _, child := range root.Children {
		if len(child.Children) == 0 {
			continue
		}
		if index, ok := groupIndex[child.Name.Local]; ok {
			groups[index].nodes = append(groups[index].nodes, child)
			continue
		}
		groupIndex[child.Name.Local] = len(groups)
		groups = append(groups, group{name: child.Name.Local, nodes: []xmlNode{child}})
	}

	if len(groups) == 1 {
		return groups[0].nodes, nil
	}
	if len(groups) > 1 {
		repeated := make([]group, 0)
		for _, entry := range groups {
			if len(entry.nodes) > 1 {
				repeated = append(repeated, entry)
			}
		}
		if len(repeated) == 1 {
			return repeated[0].nodes, nil
		}
		if len(repeated) > 1 {
			names := make([]string, 0, len(repeated))
			for _, entry := range repeated {
				names = append(names, entry.name)
			}
			return nil, fmt.Errorf("multiple repeated row elements under <%s> (%s); use --row-element", root.Name.Local, strings.Join(names, ", "))
		}
		// Prefer conventional record element names when metadata and a record
		// list coexist (e.g. <meta>...</meta><record>...</record>).
		for _, preferred := range []string{"record", "row", "item", "entry"} {
			for _, entry := range groups {
				if entry.name == preferred {
					return entry.nodes, nil
				}
			}
		}
		// Several single child elements are fields of one root record.
		return []xmlNode{*root}, nil
	}

	// No child has nested elements. If they share one name, each child is a
	// value row; otherwise the root itself is a single record.
	allSame := true
	firstName := root.Children[0].Name.Local
	for _, child := range root.Children[1:] {
		if child.Name.Local != firstName {
			allSame = false
			break
		}
	}
	if allSame && len(root.Children) > 1 {
		records := make([]xmlNode, len(root.Children))
		for i, child := range root.Children {
			records[i] = xmlNode{
				Name: child.Name,
				Children: []xmlNode{{
					Name: xml.Name{Local: "value"},
					Text: child.Text,
				}},
			}
		}
		return records, nil
	}
	return []xmlNode{*root}, nil
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

	for i, char := range name {
		if i == 0 {
			if !(unicode.IsLetter(char) || char == '_') {
				return false
			}
			continue
		}
		if !(unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsMark(char) ||
			char == '_' || char == '-' || char == '.') {
			return false
		}
	}

	return true
}
