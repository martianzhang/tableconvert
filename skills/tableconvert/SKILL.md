---
name: tableconvert
description: Convert table data between different formats (MySQL, CSV, JSON, Markdown, HTML, SQL, Excel, XML, LaTeX, MediaWiki, TWiki, templates)
---

# tableconvert

Convert table data between different formats using the tableconvert command line tool.

## Quick Start

**Multiple ways to use tableconvert:**

### 1. Auto-detect formats from file extensions (simplest):
```
tableconvert input.csv output.json
```

### 2. Short flags:
```
tableconvert -i input.csv -o output.json
```

### 3. Space-separated format values:
```
tableconvert --from csv --to json --file input.csv
```

### 4. Traditional with equals signs:
```
tableconvert --from=markdown --to=csv --file=input.md --result=output.csv
```

### 5. stdin/stdout:
```
cat input.md | tableconvert --from=markdown --to=csv
```

### 6. With format-specific options:
```
tableconvert input.csv output.md --bold-header --align=l,c,r
```

### 7. With transformations:
```
tableconvert input.csv output.json --transpose --capitalize
```

**Common file extensions:**
- `.md`, `.markdown` = markdown
- `.csv` = csv
- `.json` = json
- `.jsonl`, `.jsonlines` = jsonl
- `.sql` = sql
- `.html`, `.htm` = html
- `.xml` = xml
- `.xlsx`, `.xls` = excel
- `.tex`, `.latex` = latex
- `.tmpl`, `.template` = template
- `.wiki` = mediawiki
- `.twiki` = twiki
- `.txt` = **must specify format explicitly**

## Commands

### `tableconvert`

Convert table data from one format to another.

**Basic Options:**
- `--from={FORMAT}` or `-f={FORMAT}`: Source format (auto-detected from file if omitted)
- `--to={FORMAT}` or `-t={FORMAT}`: Target format (auto-detected from file if omitted)
- `--file={PATH}` or `--input={PATH}` or `-i={PATH}`: Input file path (optional, uses stdin if omitted)
- `--result={PATH}` or `--output={PATH}` or `-o={PATH}`: Output file path (optional, uses stdout if omitted)
- `--verbose` or `-v`: Enable verbose output

**Quick Options:**
- `-h, --help`: Show help message
- `--help-formats`: Show all supported formats and their parameters
- `--help-format={FORMAT}`: Show parameters for a specific format
- `--mcp`: Run as MCP (Model Context Protocol) server

**Global Transformations:**
- `--transpose`: Swap rows and columns
- `--delete-empty`: Remove empty rows
- `--deduplicate`: Remove duplicate rows
- `--uppercase`: Convert all text to UPPERCASE
- `--lowercase`: Convert all text to lowercase
- `--capitalize`: Capitalize first letter of each cell

**Examples:**
```
# Auto-detect (simplest)
tableconvert input.csv output.json

# Short flags
tableconvert -i input.csv -o output.json

# Space-separated
tableconvert --from csv --to json --file input.csv

# With options
tableconvert --from=csv --to=markdown --file=input.csv --align=l,c,r --bold-header

# With transformations
tableconvert input.csv output.json --transpose --capitalize
```

### `tableconvert --help-formats`

Get information about supported formats and their parameters.

**Example:**
```
tableconvert --help-formats
```

Or for a specific format:
```
tableconvert --help-format=markdown
```

## Supported Formats

| Format | Description | Common Use |
|--------|-------------|------------|
| **mysql** | MySQL query output | Database schema |
| **csv** | Comma-separated values | Spreadsheet data |
| **json** | JSON (object, 2d, column, keyed modes) | API data |
| **jsonl** | JSON Lines | Streaming data |
| **markdown** | Markdown tables | Documentation |
| **html** | HTML tables | Web pages |
| **sql** | SQL INSERT statements | Database imports |
| **excel** | Excel files | Spreadsheet files |
| **xml** | XML format | Structured data |
| **latex** | LaTeX tables | Academic papers |
| **mediawiki** | MediaWiki tables | Wiki content |
| **twiki** | TWiki/TracWiki format | Wiki content |
| **template** | Custom templates | Custom output |

## Format-Specific Options

Format-specific options are passed as flags without values (true) or with values:

### Markdown
- `--align=l,c,r`: Text alignment (l, c, r) - columns separated by comma
- `--bold-header`: Make header bold
- `--bold-first-column`: Bold first column
- `--escape`: Escape Markdown characters
- `--pretty`: Pretty-print with calculated column widths (default: true)

### CSV
- `--first-column-header`: Use first column as headers
- `--bom`: Add Byte Order Mark
- `--delimiter=TAB`: Value delimiter (COMMA, TAB, SEMICOLON, PIPE, SLASH, HASH)

### JSON
- `--format=object`: Output format (object, 2d, column, keyed)
- `--minify`: Minify JSON output
- `--parsing-json`: Parse input as JSON

### SQL
- `--one-insert`: Multiple rows in one INSERT
- `--replace`: Use REPLACE instead of INSERT
- `--dialect=mysql`: SQL dialect (none, mysql, oracle, mssql, postgresql)
- `--table=tablename`: Table name

### HTML
- `--first-column-header`: Use first column as headers
- `--div`: Use div instead of table
- `--minify`: Minify HTML
- `--thead`: Include thead/tbody tags

### Excel
- `--first-column-header`: Use first column as headers
- `--sheet-name=Sheet1`: Excel sheet name
- `--auto-width`: Auto-adjust column widths
- `--text-format`: Force text format (default: true)

### LaTeX
- `--bold-first-column`: Bold first column
- `--bold-first-row`: Bold first row
- `--borders=1111,1111`: Border style
- `--caption="Table Caption"`: Table caption
- `--escape`: Escape LaTeX characters (default: true)
- `--ht`: Place here or top of page
- `--label=table:label`: Table label
- `--location=above`: Caption location (above, below)
- `--mwe`: Minimal working example
- `--table-align=centering`: Table alignment (centering, raggedleft, raggedright)
- `--text-align=l`: Text alignment (l, c, r)

### MediaWiki
- `--first-row-header`: Use first row as headers
- `--minify`: Minify MediaWiki table
- `--sort`: Make table sortable

### XML
- `--minify`: Minify XML
- `--root-element=dataset`: Root element tag
- `--row-element=record`: Row element tag
- `--declaration`: Include XML declaration (default: true)

### Template
- `--template=file.tmpl`: Template file path

### ASCII
- `--style=box`: Table style (box, plus, dot, bubble)

## Global Transformations

Apply these to any format conversion:

- `--transpose`: Swap rows and columns
- `--delete-empty`: Remove empty rows
- `--deduplicate`: Remove duplicate rows
- `--uppercase`: Convert all text to UPPERCASE
- `--lowercase`: Convert all text to lowercase
- `--capitalize`: Capitalize first letter of each cell

**Example with transformations:**
```
tableconvert --from=csv --to=json --file=input.csv --uppercase --deduplicate
```

## Usage Examples

### MySQL to Markdown (auto-detect)
```
tableconvert mysql.txt output.md --bold-header
```

### CSV to JSON with options
```
tableconvert input.csv output.json --minify --format=object
```

### SQL to CSV with transformation
```
tableconvert input.sql output.csv --table=users --capitalize
```

### Get format info
```
tableconvert --help-format=markdown
```

### Using short flags
```
tableconvert -i data.csv -o data.json -v
```

## Notes

- **Format selection**: Auto-detected from file extensions, or specify with `--from`/`--to`
- **.txt files**: Must specify format explicitly (not auto-detected)
- Input data must be properly formatted for the source format
- Special characters are automatically escaped where needed
- UTF-8 characters are handled correctly
- Use `tableconvert --help-formats` to see all available options for a format
- All options from `docs/arguments.md` are supported

## Common Patterns

**Converting database schema (auto-detect):**
```
tableconvert schema.txt output.md --bold-header --escape
```

**Creating SQL inserts from data:**
```
tableconvert data.csv output.sql --table=users --one-insert --dialect=mysql
```

**Generating HTML tables:**
```
tableconvert data.json output.html --thead --minify=false
```

**Excel to Markdown for docs:**
```
tableconvert data.xlsx output.md --bold-header --align=l,c,r
```

**Quick pipe conversion:**
```
cat data.csv | tableconvert --from csv --to json --uppercase
```
