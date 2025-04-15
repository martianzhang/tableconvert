package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	reader, ok := cfg.Reader.(io.Reader)
	if !ok {
		return fmt.Errorf("writer is not an io.Reader, please provide a valid reader")
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

	// Get the root element and row element names from configuration
	rootElement := cfg.Extension["root-element"]
	if rootElement == "" {
		rootElement = "dataset"
	}
	rowElement := cfg.Extension["row-element"]
	if rowElement == "" {
		rowElement = "record"
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
	var minify bool
	if v, ok := cfg.Extension["minify"]; ok && strings.ToLower(v) != "false" {
		minify = true
	}

	// Get the configuration for root-element and row-element
	rootElement := cfg.Extension["root-element"]
	if rootElement == "" {
		rootElement = "dataset"
	}
	rowElement := cfg.Extension["row-element"]
	if rowElement == "" {
		rowElement = "record"
	}

	writer, ok := cfg.Writer.(io.Writer)
	if !ok {
		return fmt.Errorf("writer is not an io.Writer")
	}

	// Use the Encoder from encoding/xml
	xmlEncoder := xml.NewEncoder(writer)
	if minify {
		xmlEncoder.Indent("", "")
	} else {
		xmlEncoder.Indent("", "  ")
	}

	// Write the XML declaration
	if v, ok := cfg.Extension["declaration"]; ok && strings.ToLower(v) != "false" {
		if _, err := writer.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>` + "\n")); err != nil {
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
