package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/martianzhang/tableconvert/ascii"
	"github.com/martianzhang/tableconvert/common"
	"github.com/martianzhang/tableconvert/csv"
	"github.com/martianzhang/tableconvert/excel"
	"github.com/martianzhang/tableconvert/html"
	"github.com/martianzhang/tableconvert/json"
	"github.com/martianzhang/tableconvert/jsonl"
	"github.com/martianzhang/tableconvert/latex"
	"github.com/martianzhang/tableconvert/markdown"
	"github.com/martianzhang/tableconvert/mediawiki"
	"github.com/martianzhang/tableconvert/mysql"
	"github.com/martianzhang/tableconvert/sql"
	"github.com/martianzhang/tableconvert/tmpl"
	"github.com/martianzhang/tableconvert/twiki"
	"github.com/martianzhang/tableconvert/xml"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var formatRegistry *common.FormatRegistry

func init() {
	// Initialize format registry
	formatRegistry = common.NewFormatRegistry()

	// Register all formats
	registerFormats()
}

func main() {
	// Parse config
	args := os.Args[1:]
	cfg, err := common.ParseConfig(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parameter parsing error: %v\n\n---\n", err)
		common.Usage()
		os.Exit(1)
	}

	// Check if MCP mode is requested
	if cfg.MCPMode {
		runMCPMode()
		return
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "# From: %s\n", cfg.From)
		fmt.Fprintf(os.Stderr, "# To: %s\n", cfg.To)
		fmt.Fprintf(os.Stderr, "# Extra Configs: %v\n", cfg.Extension)
	}

	// Perform conversion
	err = performConversion(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runMCPMode starts the MCP server
func runMCPMode() {
	// Create the MCP server using the common package
	server := common.CreateMCPServer(formatRegistry)

	// Run the server using stdio transport
	ctx := context.Background()
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server failed: %v", err)
	}
}

// registerFormats registers all format unmarshal/marshal functions
func registerFormats() {
	// ASCII
	formatRegistry.RegisterFormat("ascii", ascii.Unmarshal, ascii.Marshal)

	// CSV
	formatRegistry.RegisterFormat("csv", csv.Unmarshal, csv.Marshal)

	// Excel
	formatRegistry.RegisterFormat("excel", excel.Unmarshal, excel.Marshal)
	formatRegistry.RegisterFormatAlias("xlsx", "excel")

	// HTML
	formatRegistry.RegisterFormat("html", html.Unmarshal, html.Marshal)

	// JSON
	formatRegistry.RegisterFormat("json", json.Unmarshal, json.Marshal)

	// JSONL
	formatRegistry.RegisterFormat("jsonl", jsonl.Unmarshal, jsonl.Marshal)
	formatRegistry.RegisterFormatAlias("jsonlines", "jsonl")

	// LaTeX
	formatRegistry.RegisterFormat("latex", latex.Unmarshal, latex.Marshal)

	// Markdown
	formatRegistry.RegisterFormat("markdown", markdown.Unmarshal, markdown.Marshal)
	formatRegistry.RegisterFormatAlias("md", "markdown")

	// MediaWiki
	formatRegistry.RegisterFormat("mediawiki", mediawiki.Unmarshal, mediawiki.Marshal)

	// MySQL
	formatRegistry.RegisterFormat("mysql", mysql.Unmarshal, mysql.Marshal)

	// SQL
	formatRegistry.RegisterFormat("sql", sql.Unmarshal, sql.Marshal)

	// Template (write-only format)
	formatRegistry.RegisterWriteOnlyFormat("tmpl", tmpl.Marshal)
	formatRegistry.RegisterFormatAlias("template", "tmpl")

	// TWiki
	formatRegistry.RegisterFormat("twiki", twiki.Unmarshal, twiki.Marshal)
	formatRegistry.RegisterFormatAlias("tracwiki", "twiki")

	// XML
	formatRegistry.RegisterFormat("xml", xml.Unmarshal, xml.Marshal)
}

// performConversion performs the core table conversion logic
func performConversion(cfg *common.Config) error {
	// Use the registry-based conversion
	return common.PerformConversionWithRegistry(formatRegistry, cfg)
}
