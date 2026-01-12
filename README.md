# tableconvert

🚀 **A powerful offline table format converter written in Go** - Convert between MySQL, Markdown, CSV, JSON, Excel, SQL, LaTeX, and more.

`tableconvert` is a command-line tool designed for converting between different table formats. It's perfect for developers, data analysts, and anyone who needs to transform tabular data between different systems.

## ✨ Key Features

- **🔄 Multi-format Support**: 13+ formats including MySQL, Markdown, CSV, JSON, Excel, HTML, XML, SQL, LaTeX, MediaWiki, TWiki, and custom templates
- **🔒 Offline & Secure**: All processing happens locally - your data never leaves your machine
- **⚡ Batch Processing**: Convert multiple files at once with glob patterns
- **🔧 Data Transformations**: Built-in support for transpose, deduplication, case conversion, and more
- **🎨 Format-Specific Styling**: Fine-grained control over output formatting for each format
- **🤖 MCP Server**: Built-in Model Context Protocol server for AI assistants
- **🌍 Cross-Platform**: Works on Windows, Linux, and macOS

## 🚀 Quick Start

### Installation

```bash
# Option 1: Download pre-built binaries (recommended)
# Visit: https://github.com/martianzhang/tableconvert/releases
# Or use wget (Linux/macOS):
wget https://github.com/martianzhang/tableconvert/releases/latest/download/tableconvert-linux-amd64
chmod +x tableconvert-linux-amd64
sudo mv tableconvert-linux-amd64 /usr/local/bin/tableconvert

# Option 2: Build from source (requires Go 1.23+)
git clone https://github.com/martianzhang/tableconvert.git
cd tableconvert
make build

# Binary will be available at:
# - Linux/macOS: ./bin/tableconvert
# - Windows: bin\tableconvert.exe

# Option 3: Go install
go install github.com/martianzhang/tableconvert@latest
```

### Basic Usage

```bash
# Auto-detect formats from file extensions
tableconvert input.csv output.json

# Explicit format specification
tableconvert --from=csv --to=markdown input.csv output.md

# Short flags
tableconvert -i input.csv -o output.json

# Read from stdin, write to stdout
echo "name,age\nAlice,30\nBob,25" | tableconvert --from=csv --to=json
```

### Common Scenarios

```bash
# Convert MySQL query results to Markdown (use -t for table format)
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown > users_data.md

# Convert MySQL table schema to Markdown
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown > users_schema.md

# Convert CSV to Excel with styling
tableconvert data.csv report.xlsx --auto-width --bold-header

# Batch convert all CSV files to JSON
tableconvert --batch="data/*.csv" --to=json --output-dir=json_output

# Transform data: transpose and capitalize
tableconvert input.csv output.md --transpose --capitalize

# Generate SQL INSERT statements from CSV
tableconvert data.csv data.sql --table=users --dialect=mysql

# Create LaTeX table for academic paper
tableconvert data.csv table.tex --bold-header --text-align=c
```

### Format-Specific Examples

```bash
# Markdown with custom alignment and bold headers
tableconvert data.csv output.md --align=l,c,r --bold-header

# JSON with specific format (object, 2d array, column-oriented)
tableconvert data.csv output.json --format=object --minify

# HTML table with div wrapper and minification
tableconvert data.csv output.html --div --minify --thead

# SQL with multiple rows in one INSERT
tableconvert data.csv output.sql --one-insert --table=products

# Custom template output
tableconvert data.csv output.php --template=php_array.tmpl
```

## 📖 Command Reference

**Basic Options:**
- `--from|-f={FORMAT}` - Source format (required if not auto-detected)
- `--to|-t={FORMAT}` - Target format (required if not auto-detected)
- `--file|--input|-i={PATH}` - Input file path (or stdin if omitted)
- `--result|--output|-o={PATH}` - Output file path (or stdout if omitted)

**Quick Options:**
- `-v, --verbose` - Show detailed processing information
- `-h, --help` - Show help message
- `--help-formats` - List all supported formats
- `--help-format={FORMAT}` - Show format-specific parameters
- `--mcp` - Run as MCP server for AI assistants

**Batch Processing:**
- `--batch|-b={PATTERN}` - Process multiple files (e.g., `*.csv`, `data/*.json`)
- `--recursive|-r` - Enable recursive directory search
- `--output-dir|--dir={PATH}` - Output directory (default: same as input)

**Data Transformations:**
- `--transpose` - Swap rows and columns
- `--delete-empty` - Remove empty rows
- `--deduplicate` - Remove duplicate rows
- `--uppercase` - Convert to UPPERCASE
- `--lowercase` - Convert to lowercase
- `--capitalize` - Capitalize first letter of each cell

**Auto-Detection:**
When `--from` or `--to` are omitted, formats are detected from file extensions:
- `.csv` → csv, `.json` → json, `.jsonl` → jsonl
- `.md`, `.markdown` → markdown, `.xlsx`, `.xls` → excel
- `.html`, `.htm` → html, `.xml` → xml, `.sql` → sql
- `.tex`, `.latex` → latex, `.wiki` → mediawiki
- `.tmpl`, `.template` → tmpl

**Format-Specific Parameters:**
Each format supports custom styling options. See [arguments.md](docs/arguments.md) for complete reference or use `--help-format={format}`.

## 🤖 MCP (Model Context Protocol) Integration

`tableconvert` includes a built-in MCP server for seamless integration with AI assistants like Claude Code.

### Setup with Claude Code

```bash
# Add to Claude Code (Unix/Linux/macOS)
claude mcp add tableconvert -- /path/to/tableconvert --mcp

# Windows
claude mcp add tableconvert -- "C:\\path\\to\\tableconvert.exe" --mcp
```

### MCP Tools Available

- **`convert_table`**: Convert table data between formats with full parameter support
- **`get_formats`**: Discover supported formats and their parameters

### Example MCP Usage

Once configured, you can ask your AI assistant:
- "Convert this CSV data to a Markdown table with bold headers"
- "Transform my MySQL query results to JSON format"
- "Create a LaTeX table from this data with centered alignment"

## 📊 Supported Formats

| Format | Extensions | Read | Write | Description |
|--------|------------|------|-------|-------------|
| **MySQL** | `mysql` | ✅ | ✅ | MySQL query output (box format) |
| **CSV** | `.csv` | ✅ | ✅ | Comma-separated values |
| **JSON** | `.json` | ✅ | ✅ | JavaScript Object Notation |
| **JSONL** | `.jsonl`, `.jsonlines` | ✅ | ✅ | JSON Lines format |
| **Markdown** | `.md`, `.markdown` | ✅ | ✅ | GitHub/Markdown tables |
| **Excel** | `.xlsx`, `.xls` | ✅ | ✅ | Microsoft Excel files |
| **HTML** | `.html`, `.htm` | ✅ | ✅ | HTML tables |
| **XML** | `.xml` | ✅ | ✅ | XML data format |
| **SQL** | `.sql` | ✅ | ✅ | SQL INSERT statements |
| **LaTeX** | `.tex`, `.latex` | ✅ | ✅ | LaTeX table format |
| **MediaWiki** | `.wiki` | ✅ | ✅ | MediaWiki tables |
| **TWiki** | `.twiki` | ✅ | ✅ | TWiki/TracWiki format |
| **ASCII** | - | ❌ | ✅ | ASCII art tables |
| **Template** | `.tmpl`, `.template` | ❌ | ✅ | Custom templates |

## 📚 Real-World Examples

### Database Schema Documentation
```bash
# Export table schema to Markdown
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown > schema.md

# Generate HTML documentation
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=html > table.html
```

**Note about MySQL format**: `tableconvert` only supports MySQL query output in **box format** (like `mysql -t` output). It does NOT support:
- `mysqldump` output (SQL INSERT/CREATE statements)
- Raw SQL queries without `-t` flag
- Multi-statement SQL scripts

For `mysqldump` or raw SQL, you'll need additional tools to convert to box format first, or use the `sql` format for INSERT statements.

### Data Pipeline
```bash
# Convert CSV to JSON for API
tableconvert data.csv data.json --format=object

# Transform for analytics
tableconvert sales.csv report.md --transpose --bold-header --align=l,c,r
```

### Reporting
```bash
# Batch process daily reports
tableconvert --batch="reports/*.csv" --to=excel --output-dir=excel_reports --bold-header

# Generate LaTeX for papers
tableconvert results.csv table.tex --caption="Experiment Results" --table-align=centering
```

## 📦 Releases

Pre-built binaries are available for Linux, macOS, and Windows on the [GitHub Releases page](https://github.com/martianzhang/tableconvert/releases).

### Release Assets

Each release includes:
- **Binaries** for all platforms (x64 and ARM64)
- **SHA256 checksums** for verification
- **Release notes** with changes and build info

### Verification

Always verify downloaded binaries:

```bash
# Download binary and checksums.txt
sha256sum -c checksums.txt

# Should output: tableconvert-linux-amd64: OK
```

### Building Releases Locally

```bash
# Build all platform binaries
make release

# Create zip archives (requires zip)
make release-zip

# Generate release notes
make release-notes
```

See [docs/release-guide.md](docs/release-guide.md) for complete release documentation.

## 🔧 Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/martianzhang/tableconvert.git
cd tableconvert

# Build
make build

# Run tests
make test

# Run specific tests
go test ./markdown/
go test ./mysql/
```

### Adding New Formats

1. Create package: `mkdir newformat && touch newformat/newformat.go`
2. Implement `Unmarshal(*common.Config, *common.Table) error`
3. Implement `Marshal(*common.Config, *common.Table) error`
4. Register in `cmd/tableconvert/main.go`
5. Add tests in `newformat/newformat_test.go`

See [CLAUDE.md](CLAUDE.md) for detailed architecture documentation.

## 📖 Documentation

### User Documentation
- **[Quick Reference](docs/quick-reference.md)** - Fast lookup for common commands
- **[Complete Parameter Reference](docs/arguments.md)** - All format-specific options
- **[Practical Examples](docs/examples.md)** - Real-world usage scenarios
- **[Troubleshooting Guide](docs/troubleshooting.md)** - Solve common issues

### Format Guides
- **[LaTeX Guide](docs/latex.md)** - LaTeX table syntax explained
- **[Wiki Formats](docs/wiki.md)** - TWiki, MediaWiki, Confluence syntax

### Developer Documentation
- **[Architecture Guide](CLAUDE.md)** - Code structure and development
- **[Contributing](#-contributing)** - How to contribute

### Quick Help
```bash
# Show all available commands
tableconvert --help

# Show all supported formats
tableconvert --help-formats

# Show format-specific parameters
tableconvert --help-format=markdown
```

## 🆘 Troubleshooting

### Common Issues

**"Format not detected"**
- Ensure file extensions are correct
- Use `--from` and `--to` to specify formats explicitly
- Run `--help-formats` to see all supported formats

**"Invalid table format"**
- Check if input follows expected format
- Use `--verbose` to see detailed parsing information
- For MySQL format, ensure box borders are complete

**"Batch mode not working"**
- Verify glob patterns are quoted properly
- Use `--recursive` for subdirectories
- Check file permissions

### Getting Help

```bash
# Show all available commands
tableconvert --help

# Show format-specific parameters
tableconvert --help-format=markdown

# List all supported formats
tableconvert --help-formats

# Verbose mode for debugging
tableconvert input.csv output.json --verbose
```

## 🤝 Contributing

Contributions are welcome! Please ensure:
- All tests pass: `make test`
- Code follows existing patterns
- Add tests for new features
- Update documentation

## 📄 License

[Apache License 2.0](LICENSE)

## 🔗 References

* [Online tableconvert](https://tableconvert.com/)
* [ascii-tables](https://github.com/ozh/ascii-tables)
* [tablewriter](https://github.com/olekukonko/tablewriter)
* [csvq](https://github.com/mithrandie/csvq)

---

## Note
* This tool is not production ready.
* When you will overwrite exists file, please check first. 
* All data processing happens locally on your machine.
