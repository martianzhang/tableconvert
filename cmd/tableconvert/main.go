package main

import (
	"context"
	"fmt"
	"io"
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

// version is set at build time via ldflags
// Example: go build -ldflags "-X main.version=v1.0.0" ./cmd/tableconvert
var version = "dev"

var formatRegistry *common.FormatRegistry

func init() {
	// Initialize format registry
	formatRegistry = common.NewFormatRegistry()

	// Register all formats
	registerFormats()
}

func main() {
	// Check for version flag early (before ParseConfig)
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--version" {
			fmt.Printf("tableconvert version %s\n", version)
			os.Exit(0)
		}
	}

	// Parse config
	cfg, err := common.ParseConfig(args)
	if err != nil {
		// Show concise error message with helpful tips
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if MCP mode is requested
	if cfg.MCPMode {
		runMCPMode()
		return
	}

	// Check if batch mode is requested
	if cfg.Batch != "" {
		runBatchMode(&cfg)
		return
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

// runBatchMode processes multiple files in batch
func runBatchMode(cfg *common.Config) {
	// Get list of files to process
	files, err := cfg.GetBatchFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Dry-run mode for batch
	if cfg.DryRun {
		fmt.Fprintf(os.Stderr, "=== DRY RUN MODE (Batch) ===\n")
		fmt.Fprintf(os.Stderr, "Found %d files to process\n", len(files))
		for i, file := range files {
			fmt.Fprintf(os.Stderr, "  %d. %s -> %s\n", i+1, file.InputPath, file.OutputPath)
		}
		fmt.Fprintf(os.Stderr, "\n✓ Batch would process %d files (no output written)\n", len(files))
		return
	}

	// Create output directory if specified and doesn't exist
	if cfg.OutputDir != "" {
		if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Process each file
	successCount := 0
	failCount := 0

	for _, file := range files {
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Processing: %s -> %s\n", file.InputPath, file.OutputPath)
		}

		// Open input file
		inputFile, err := os.Open(file.InputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Failed to open %s: %v\n", file.InputPath, err)
			failCount++
			continue
		}

		// Create output file
		outputFile, err := os.Create(file.OutputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Failed to create %s: %v\n", file.OutputPath, err)
			inputFile.Close()
			failCount++
			continue
		}

		// Create a temporary config for this file
		fileCfg := &common.Config{
			From:      file.FromFormat,
			To:        file.ToFormat,
			Reader:    inputFile,
			Writer:    outputFile,
			Verbose:   cfg.Verbose,
			Extension: cfg.Extension,
		}

		// Perform conversion
		err = common.PerformConversionWithRegistry(formatRegistry, fileCfg)
		inputFile.Close()
		outputFile.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Conversion failed: %v\n", err)
			// Clean up failed output file
			os.Remove(file.OutputPath)
			failCount++
		} else {
			if cfg.Verbose {
				fmt.Fprintf(os.Stderr, "  ✓ Success\n")
			}
			successCount++
		}
	}

	// Print summary
	fmt.Fprintf(os.Stderr, "\nBatch processing complete: %d succeeded, %d failed\n", successCount, failCount)

	if failCount > 0 {
		os.Exit(1)
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
	// Print verbose output
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "# From: %s\n", cfg.From)
		fmt.Fprintf(os.Stderr, "# To: %s\n", cfg.To)
		fmt.Fprintf(os.Stderr, "# Extra Configs: %v\n", cfg.Extension)
	}

	// Handle dry-run mode
	if cfg.DryRun {
		return performDryRun(cfg)
	}

	// Use the registry-based conversion
	return common.PerformConversionWithRegistry(formatRegistry, cfg)
}

// performDryRun performs a dry-run conversion, showing preview without writing
func performDryRun(cfg *common.Config) error {
	// Parse input
	var table common.Table

	unmarshalFn, ok := formatRegistry.GetUnmarshalFunc(cfg.From)
	if !ok {
		return fmt.Errorf("unsupported `--from` format: %s", cfg.From)
	}

	if unmarshalFn == nil {
		return fmt.Errorf("format %s does not support reading (unmarshal)", cfg.From)
	}

	if err := unmarshalFn(cfg, &table); err != nil {
		return fmt.Errorf("error parsing input: %w", err)
	}

	// Apply transformations
	cfg.ApplyTransformations(&table)

	// Show dry-run summary
	fmt.Fprintf(os.Stderr, "=== DRY RUN MODE ===\n")
	fmt.Fprintf(os.Stderr, "Input format:  %s\n", cfg.From)
	fmt.Fprintf(os.Stderr, "Output format: %s\n", cfg.To)
	fmt.Fprintf(os.Stderr, "Rows: %d\n", len(table.Rows))
	fmt.Fprintf(os.Stderr, "Columns: %d\n", len(table.Headers))
	if len(table.Headers) > 0 {
		fmt.Fprintf(os.Stderr, "Headers: %v\n", table.Headers)
	}
	if cfg.Verbose && len(table.Rows) > 0 {
		fmt.Fprintf(os.Stderr, "First row: %v\n", table.Rows[0])
		if len(table.Rows) > 1 {
			fmt.Fprintf(os.Stderr, "Last row: %v\n", table.Rows[len(table.Rows)-1])
		}
	}

	// Try to generate output to validate it would work
	// Use a discard writer since we don't want actual output
	marshalFn, ok := formatRegistry.GetMarshalFunc(cfg.To)
	if !ok {
		return fmt.Errorf("unsupported `--to` format: %s", cfg.To)
	}

	if marshalFn == nil {
		return fmt.Errorf("format %s does not support writing (marshal)", cfg.To)
	}

	// Create a discard writer to test marshaling
	// We need to temporarily replace the writer
	originalWriter := cfg.Writer
	cfg.Writer = io.Discard

	err := marshalFn(cfg, &table)
	cfg.Writer = originalWriter

	if err != nil {
		return fmt.Errorf("error generating output: %w", err)
	}

	fmt.Fprintf(os.Stderr, "\n✓ Conversion would succeed (no output written)\n")
	return nil
}
