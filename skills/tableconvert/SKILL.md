---
name: tableconvert
description: Convert table data between different formats (MySQL, CSV, JSON, Markdown, HTML, SQL, Excel, XML, LaTeX, MediaWiki, TWiki, templates)
---

# tableconvert

Convert table data between different formats using the tableconvert command line tool.

## Quick Start

Use `tableconvert` to transform table data. **Format is determined by file extension:**

```
tableconvert --from=markdown --to=csv --file=input.md --result=output.csv
```

Or with stdin/stdout:
```
cat input.md | tableconvert --from=markdown --to=csv
```

**Common file extensions:**
- `.md` = markdown
- `.csv` = csv
- `.json` = json
- `.sql` = sql
- `.html` = html
- `.xml` = xml
- `.xlsx` = excel
- `.latex` = latex
- `.tmpl` = template
- `.mediawiki` = mediawiki
- `.twiki` = twiki

## Commands

### `tableconvert`

Convert table data from one format to another.

**Parameters:**
- `--from={FORMAT}` or `-f={FORMAT}`: Source format (required)
- `--to={FORMAT}` or `-t={FORMAT}`: Target format (required)
- `--file={PATH}`: Input file path (optional, uses stdin if omitted)
- `--result={PATH}` or `-r={PATH}`: Output file path (optional, uses stdout if omitted)
- `--verbose` or `-v`: Enable verbose output

**Global Transformations:**
- `--transpose`: Swap rows and columns
- `--delete-empty`: Remove empty rows
- `--deduplicate`: Remove duplicate rows
- `--uppercase`: Convert all text to UPPERCASE
- `--lowercase`: Convert all text to lowercase
- `--capitalize`: Capitalize first letter of each cell

**Example:**
```
tableconvert --from=csv --to=markdown --file=input.csv --align=l,c,r --bold-header
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
- `--transpose`: Swap rows and columns
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
- `--transpose`: Swap rows and columns
- `--div`: Use div instead of table
- `--minify`: Minify HTML
- `--thead`: Include thead/tbody tags

### Excel
- `--first-column-header`: Use first column as headers
- `--transpose`: Swap rows and columns
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

### MySQL to Markdown
```
tableconvert --from=mysql --to=markdown --file=mysql.txt --bold-header
```

### CSV to JSON with options
```
tableconvert --from=csv --to=json --file=input.csv --minify --format=object
```

### SQL to CSV with transformation
```
tableconvert --from=sql --to=csv --file=input.sql --table=users --capitalize
```

### Get format info
```
tableconvert --help-format=markdown
```

## Notes

- **Format selection**: Use the format name that matches your file extension (e.g., `--from=markdown` for `.md` files, `--from=mysql` for MySQL query output format)
- Input data must be properly formatted for the source format
- Special characters are automatically escaped where needed
- UTF-8 characters are handled correctly
- Use `tableconvert --help-formats` to see all available options for a format
- All options from `docs/arguments.md` are supported

## Common Patterns

**Converting database schema:**
```
tableconvert --from=mysql --to=markdown --file=schema.txt --bold-header --escape
```

**Creating SQL inserts from data:**
```
tableconvert --from=csv --to=sql --file=data.csv --table=users --one-insert --dialect=mysql
```

**Generating HTML tables:**
```
tableconvert --from=json --to=html --file=data.json --thead --minify=false
```

**Excel to Markdown for docs:**
```
tableconvert --from=excel --to=markdown --file=data.xlsx --bold-header --align=l,c,r
```
