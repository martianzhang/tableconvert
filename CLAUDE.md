# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`tableconvert` is a Go CLI tool for converting between different table formats. It reads tabular data from one format and outputs it in another format, supporting formats like MySQL query output, CSV, JSON, Markdown, HTML, Excel, XML, SQL, LaTeX, MediaWiki, TWiki, and custom templates.

## Architecture

### Core Design Pattern

The codebase follows a **plugin architecture** with a central dispatcher pattern:

1. **Common Package** (`common/`): Shared types and utilities
   - `Config`: Holds configuration, reader/writer, and extension parameters
   - `Table`: The core data structure (headers + rows)
   - Format detection functions (by extension and content analysis)
   - Type inference and escape functions

2. **Format Packages** (`ascii/`, `csv/`, `excel/`, `html/`, `json/`, `jsonl/`, `latex/`, `markdown/`, `mediawiki/`, `mysql/`, `sql/`, `tmpl/`, `twiki/`, `xml/`): Each implements:
   - `Unmarshal(*Config, *Table) error`: Parse input into Table
   - `Marshal(*Config, *Table) error`: Convert Table to output format

3. **Main Entry Point** (`cmd/tableconvert/main.go`):
   - Parses CLI arguments via `common.ParseConfig()`
   - Dispatches to appropriate `Unmarshal()` based on `--from`
   - Dispatches to appropriate `Marshal()` based on `--to`

### Key Design Decisions

- **Interface-based**: Each format is self-contained with Unmarshal/Marshal functions
- **Config-driven**: Extension parameters passed via `Config.Extension` map
- **UTF-8 aware**: Uses `runewidth` package for proper character width handling
- **Error handling**: Returns `ParseError` with line numbers for debugging

## Common Development Tasks

### Build and Run

```bash
# Build the CLI
make build

# Run tests
make test

# Run specific package tests
go test ./markdown/
go test ./mysql/

# Run with verbose output
go test -v ./...
```

### Adding a New Format

1. Create new package: `mkdir newformat && touch newformat/newformat.go`
2. Implement `Unmarshal(cfg *common.Config, table *common.Table) error`
3. Implement `Marshal(cfg *common.Config, table *common.Table) error`
4. Register in `cmd/tableconvert/main.go`:
   - Add import: `"github.com/martianzhang/tableconvert/newformat"`
   - Add case to reader switch (line ~43-72)
   - Add case to writer switch (line ~79-107)
5. Add tests in `newformat/newformat_test.go`
6. Update `docs/arguments.md` if format supports extension parameters

### Testing Strategy

- Each format has its own `_test.go` file
- Tests typically cover:
  - Basic Unmarshal/Marshal round-trips
  - Edge cases (empty tables, special characters, UTF-8)
  - Error conditions (malformed input)
- Use `common.ParseConfig()` in tests to simulate CLI args

### Extension Parameters

Format-specific options are passed via `Config.Extension` map:
- `cfg.GetExtensionBool(key, default)` - for boolean flags
- `cfg.GetExtensionString(key, default)` - for string values
- `cfg.GetExtensionInt(key, default)` - for integer values

See `docs/arguments.md` for all format-specific parameters.

## Important Files

- `cmd/tableconvert/main.go` - CLI entry point and format dispatcher
- `common/types.go` - Core `Table` struct and type definitions
- `common/config.go` - Config parsing and format detection
- `common/escape.go` - String escaping utilities for various formats
- `common/usage.txt` - Embedded CLI help text
- `docs/arguments.md` - Extension parameter reference

## Format-Specific Notes

### MySQL Format
- Uses state machine parser for robust handling of multi-byte characters
- Handles lines split across buffer reads
- Supports missing bottom border (tolerant parsing)

### Markdown Format
- Supports pretty-printing with column width calculation
- Handles UTF-8 characters via `runewidth` package
- Supports alignment (l/c/r), bold headers, and escaping

### Template Format
- Uses Go's `text/template` package
- Provides helper functions: `Upper`, `Lower`, `Capitalize`, `Sub`, `Quote`, and various escape functions

### JSON Format
- Multiple output modes: object, 2d array, column-oriented, keyed
- Type inference converts strings to bool/int/float/null where appropriate

## Configuration Flow

```
CLI Args → ParseConfig() → Config{From, To, Reader, Writer, Extension}
    ↓
Reader → Unmarshal() → Table{Headers, Rows}
    ↓
Table → Marshal() → Writer
```

## Dependencies

- `github.com/xuri/excelize/v2` - Excel file handling
- `vitess.io/vitess` - SQL parsing (for SQL format)
- `github.com/mattn/go-runewidth` - UTF-8 width handling
- `github.com/stretchr/testify` - Testing utilities

## Notes

- The project is marked as "In progress, Not production ready" in README
- Format detection order in `common/types.go:78-106` is important for auto-detection
