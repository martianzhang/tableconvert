package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPServerContext holds the registry for MCP handlers
type MCPServerContext struct {
	Registry *FormatRegistry
}

// NewMCPServerContext creates a new MCP server context with a format registry
func NewMCPServerContext(registry *FormatRegistry) *MCPServerContext {
	return &MCPServerContext{
		Registry: registry,
	}
}

// ConvertTableArgs represents the arguments for the convert_table tool
type ConvertTableArgs struct {
	From            string            `json:"from" jsonschema:"source format (e.g., csv, mysql, markdown)"`
	To              string            `json:"to" jsonschema:"target format (e.g., json, csv, markdown)"`
	Input           string            `json:"input" jsonschema:"table data as a string"`
	Options         map[string]string `json:"options,omitempty" jsonschema:"format-specific options as key-value pairs"`
	Transformations map[string]bool   `json:"transformations,omitempty" jsonschema:"global transformations (transpose, delete-empty, deduplicate, uppercase, lowercase, capitalize)"`
}

// ConvertTableResult represents the result of a conversion
type ConvertTableResult struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// GetFormatsArgs represents the arguments for the get_formats tool
type GetFormatsArgs struct {
	Format string `json:"format,omitempty" jsonschema:"specific format to get help for"`
}

// GetFormatsResult represents the result of getting format information
type GetFormatsResult struct {
	Formats    map[string]string        `json:"formats"`
	Parameters map[string][]FormatParam `json:"parameters,omitempty"`
}

// MCPFormatParam represents a format-specific parameter for MCP (compatible with FormatParam)
type MCPFormatParam struct {
	Name          string `json:"name"`
	DefaultValue  string `json:"default_value"`
	AllowedValues string `json:"allowed_values"`
	Description   string `json:"description"`
}

// HandleConvertTable handles the convert_table tool call
func (s *MCPServerContext) HandleConvertTable(ctx context.Context, req *mcp.CallToolRequest, args ConvertTableArgs) (*mcp.CallToolResult, ConvertTableResult, error) {
	result := ConvertTableResult{}

	// Validate required fields
	if args.From == "" {
		return nil, ConvertTableResult{Error: "from is required"}, fmt.Errorf("from is required")
	}
	if args.To == "" {
		return nil, ConvertTableResult{Error: "to is required"}, fmt.Errorf("to is required")
	}
	if args.Input == "" {
		return nil, ConvertTableResult{Error: "input is required"}, fmt.Errorf("input is required")
	}

	// Create config
	cfg := Config{
		From:      args.From,
		To:        args.To,
		Reader:    strings.NewReader(args.Input),
		Extension: make(map[string]string),
	}

	// Add format-specific options
	for key, value := range args.Options {
		cfg.Extension[key] = value
	}

	// Add transformations
	if args.Transformations != nil {
		for key, value := range args.Transformations {
			if value {
				cfg.Extension[key] = "true"
			}
		}
	}

	// Create output buffer
	var output strings.Builder
	cfg.Writer = &output

	// Perform conversion using registry
	err := PerformConversionWithRegistry(s.Registry, &cfg)
	if err != nil {
		result.Error = err.Error()
		return nil, result, err
	}

	result.Output = output.String()
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result.Output},
		},
	}, result, nil
}

// CreateMCPServer creates the MCP server with tableconvert tools
func CreateMCPServer(registry *FormatRegistry) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "tableconvert",
		Version: "1.0.0",
		Title:   "TableConvert - Format Converter",
	}, &mcp.ServerOptions{
		Instructions: `Tableconvert MCP Server provides tools for converting between different table formats.

Available Tools:
1. convert_table - Convert table data between formats
2. get_formats - Get information about supported formats and their parameters

Supported Formats:
- csv, excel, html, json, jsonl, latex, markdown, mediawiki, mysql, sql, tmpl, twiki, xml

Format-Specific Options:
Use the options parameter to pass format-specific settings like:
- markdown: align, bold-header, bold-first-column, escape, pretty
- csv: first-column-header, bom, delimiter
- json: format, minify, parsing-json
- html: first-column-header, div, minify, thead
- excel: first-column-header, sheet-name, auto-width, text-format
- ascii: style
- latex: bold-first-column, bold-first-row, borders, caption, escape, ht, label, location, mwe, table-align, text-align
- mediawiki: first-row-header, minify, sort
- sql: one-insert, replace, dialect, table
- tmpl: template
- xml: minify, root-element, row-element, declaration

Global Transformations:
Use the transformations parameter for operations that work across all formats:
- transpose: swap rows and columns
- delete-empty: remove empty rows
- deduplicate: remove duplicate rows
- uppercase: convert all text to UPPERCASE
- lowercase: convert all text to lowercase
- capitalize: capitalize first letter of each cell

Example Usage:
{
  "from": "csv",
  "to": "markdown",
  "input": "Name,Age\\nAlice,30\\nBob,25",
  "options": {"align": "l,c", "bold-header": "true"},
  "transformations": {"uppercase": true}
}`,
	})

	// Create server context with registry
	context := NewMCPServerContext(registry)

	// Add the convert_table tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "convert_table",
		Description: "Convert table data between different formats",
	}, context.HandleConvertTable)

	// Add the get_formats tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_formats",
		Description: "Get information about supported formats and their parameters",
	}, context.HandleGetFormats)

	return server
}

// HandleGetFormats handles the get_formats tool call (static version)
func (s *MCPServerContext) HandleGetFormats(ctx context.Context, req *mcp.CallToolRequest, args GetFormatsArgs) (*mcp.CallToolResult, GetFormatsResult, error) {
	result := GetFormatsResult{
		Formats: map[string]string{
			"csv":       "Comma-Separated Values",
			"excel":     "Excel spreadsheet (XLSX)",
			"html":      "HTML table",
			"json":      "JSON (object, 2d array, column-oriented, keyed)",
			"jsonl":     "JSON Lines",
			"latex":     "LaTeX table",
			"markdown":  "Markdown table",
			"mediawiki": "MediaWiki table",
			"mysql":     "MySQL query output",
			"sql":       "SQL INSERT statements",
			"tmpl":      "Custom template",
			"twiki":     "TWiki/TracWiki table",
			"xml":       "XML",
			"ascii":     "ASCII table",
		},
	}

	// If a specific format is requested, return its parameters
	if args.Format != "" {
		format := strings.ToLower(args.Format)
		params := GetFormatParams(format)

		if len(params) == 0 {
			return nil, result, fmt.Errorf("unknown format: %s", args.Format)
		}

		result.Parameters = make(map[string][]FormatParam)
		result.Parameters[format] = make([]FormatParam, len(params))

		for i, p := range params {
			result.Parameters[format][i] = FormatParam{
				Name:          p.Name,
				DefaultValue:  p.DefaultValue,
				AllowedValues: p.AllowedValues,
				Description:   p.Description,
			}
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Available formats: %v", result.Formats)},
		},
	}, result, nil
}
