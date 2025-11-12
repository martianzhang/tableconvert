package xml

import (
	"encoding/xml"
	"fmt"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	reader := cfg.Reader
	if reader == nil {
		return fmt.Errorf("reader is nil, please provide a valid reader")
	}

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
		return err
	}

	// Extract headers and row data
	if len(root.Rows) > 0 {
		// Extract headers from the first row
		for _, field := range root.Rows[0].Fields {
			table.Headers = append(table.Headers, field.XMLName.Local)
		}

		// Extract all row data
		for _, row := range root.Rows {
			var dataRow []string
			for _, field := range row.Fields {
				dataRow = append(dataRow, field.Value)
			}
			table.Rows = append(table.Rows, dataRow)
		}
	}

	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	// Get the configuration for minify
	minify := cfg.GetExtensionBool("minify", false)

	// Get the configuration for root-element and row-element
	rootElement := cfg.GetExtensionString("root-element", "dataset")
	rowElement := cfg.GetExtensionString("row-element", "record")

	// Use the Encoder from encoding/xml
	xmlEncoder := xml.NewEncoder(cfg.Writer)
	if minify {
		xmlEncoder.Indent("", "")
	} else {
		xmlEncoder.Indent("", "  ")
	}

	// Write the XML declaration
	if cfg.GetExtensionBool("declaration", false) {
		if _, err := fmt.Fprintln(cfg.Writer, `<?xml version="1.0" encoding="UTF-8" ?>`); err != nil {
			return err
		}
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

	for _, row := range table.Rows {
		r := Row{
			XMLName: xml.Name{Local: rowElement},
		}
		for i, cell := range row {
			header := table.Headers[i]
			if header == "" {
				header = "NULL"
			}
			r.Cells = append(r.Cells, Cell{
				XMLName: xml.Name{Local: header},
				Value:   cell,
			})
		}
		root.Rows = append(root.Rows, r)
	}

	// Encode and write XML
	if err := xmlEncoder.Encode(root); err != nil {
		return err
	}

	return nil
}
