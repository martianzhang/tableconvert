# Format-Specific Arguments Reference

This document provides a complete reference for all format-specific parameters supported by tableconvert.

## 📋 Global Transformations

These transformations work with **all formats**:

| Parameter | Description | Example |
|-----------|-------------|---------|
| `--transpose` | Swap rows and columns | `tableconvert data.csv output.md --transpose` |
| `--delete-empty` | Remove empty rows | `tableconvert data.csv output.md --delete-empty` |
| `--deduplicate` | Remove duplicate rows | `tableconvert data.csv output.md --deduplicate` |
| `--uppercase` | Convert all text to UPPERCASE | `tableconvert data.csv output.md --uppercase` |
| `--lowercase` | Convert all text to lowercase | `tableconvert data.csv output.md --lowercase` |
| `--capitalize` | Capitalize first letter of each cell | `tableconvert data.csv output.md --capitalize` |
| `--strict-columns` | Fail on ragged rows or empty/duplicate headers instead of normalizing | `tableconvert data.csv output.json --strict-columns` |

By default, short rows are padded, long rows extend the table with `Column N` headers, and duplicate/empty headers are renamed (`a`, `a_2`, `Column 3`). Use `--strict-columns` to fail instead.

---

## 🎨 Format Parameters

### ASCII Tables

**Usage:** `tableconvert data.csv output.ascii --style=box`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `style` | `box` | `box`, `plus`, `dot`, `bubble` | Visual table style |

**Examples:**
```bash
# Box style (default)
tableconvert data.csv output.ascii --style=box

# Plus style
tableconvert data.csv output.ascii --style=plus

# Dot style
tableconvert data.csv output.ascii --style=dot
```

---

### CSV / TSV (Delimited Values)

**Usage:** `tableconvert data.json output.csv --delimiter=TAB --bom`

TSV is available as `--from=tsv` / `--to=tsv`, or automatically from `.tsv` / `.tab` extensions. CSV/TSV default to comma/tab respectively; `--delimiter` overrides both.

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `first-column-header` | `false` | `true`, `false` | Use first column as headers |
| `bom` | `false` | `true`, `false` | Add Byte Order Mark (BOM) |
| `delimiter` | `,` | `COMMA`, `TAB`, `SEMICOLON`, `PIPE`, `SLASH`, `HASH` | Value delimiter |

**Examples:**
```bash
# Tab-separated values
tableconvert data.json output.csv --delimiter=TAB

# With BOM for Excel compatibility
tableconvert data.json output.csv --bom

# Semicolon delimiter (common in Europe)
tableconvert data.json output.csv --delimiter=SEMICOLON
```

---

### Excel (XLSX)

**Usage:** `tableconvert data.csv output.xlsx --auto-width`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `first-column-header` | `false` | `true`, `false` | Use first column as headers |
| `sheet-name` | `Sheet1` | Any string | Excel sheet name (input/output) |
| `sheet-index` | `` | 1-based number | Read the Nth sheet when `sheet-name` is empty |
| `expand-merged` | `true` | `true`, `false` | Expand merged cells on input |
| `auto-width` | `false` | `true`, `false` | Auto-adjust column widths |
| `text-format` | `true` | `true`, `false` | Force text format for all cells (intentionally kept) |

**Input format:** Excel input/output supports `.xlsx` only. Legacy `.xls` files must be saved as `.xlsx` or exported to CSV first; `tableconvert` returns a clear error instead of attempting to parse them.

**List sheets:**

```bash
tableconvert --list-sheets report.xlsx
# Sheet1
# Sales Data
```

**Examples:**
```bash
# Auto-width columns
tableconvert data.csv output.xlsx --auto-width

# Custom sheet name
tableconvert data.csv output.xlsx --sheet-name="Sales Data"

# Headers from first column
tableconvert data.csv output.xlsx --first-column-header=true

# Read the second worksheet by index
tableconvert --sheet-index=2 report.xlsx output.csv
```

---

### HTML

**Usage:** `tableconvert data.csv output.html --div --minify --thead`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `first-column-header` | `false` | `true`, `false` | Use first column as headers |
| `div` | `false` | `true`, `false` | Wrap in `<div>` instead of `<table>` |
| `minify` | `false` | `true`, `false` | Minify HTML output |
| `thead` | `false` | `true`, `false` | Include `<thead>` and `<tbody>` tags |

**Examples:**
```bash
# Minified HTML table
tableconvert data.csv output.html --minify

# Div-based table
tableconvert data.csv output.html --div

# With proper thead/tbody structure
tableconvert data.csv output.html --thead
```

---

### JSON

**Usage:** `tableconvert data.csv output.json --format=object --minify`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `format` | `object` | `object`, `2d`, `column`, `column-array`, `keyed` | JSON output format |
| `indent` | `2` | `0-8` | Indent size; `0` minifies |
| `root` | `` | Any string | Wrap the output in a root object, e.g. `{"data":...}` |
| `minify` | `false` | `true`, `false` | Minify JSON output (overrides `indent`) |
| `parse-types` | `false` | `true`, `false` | Convert values to numbers/bools/null on output (alias: `parsing-json`) |
| `parsing-json` | `false` | `true`, `false` | Alias for `parse-types` |
| `json-path` | `` | Dot/bracket path | Select the row array from a wrapper object (e.g. `data.items`) |
| `flatten-arrays` | `json` | `json`, `index` | Serialize arrays into one cell, or expand to `tags[0]`, `tags[1]` |
| `sort-headers` | `false` | `true`, `false` | Sort headers alphabetically instead of first-seen order |

**JSON Format Options:**
- `object`: Array of objects (recommended)
- `2d`: 2D array (rows and columns)
- `column`: Column-oriented object: `{"a":[...],"b":[...]}`
- `column-array`: Legacy column array: `[{"a":[...]},{"b":[...]}]`
- `keyed`: Keyed by first column

**Nested JSON:**
- Nested objects become dot paths: `{"user":{"name":"Alice"}}` → `user.name`.
- Arrays are serialized as compact JSON by default: `["a","b"]`.
- Use `--flatten-arrays=index` to expand arrays into `tags[0]`, `tags[1]`.
- A top-level object is treated as one row; use `--json-path` for wrapper objects.

**Type preservation:**
- JSON → JSON/JSONL preserves nested values, `null`, booleans, and the original number text when no destructive transformation is applied.
- For CSV/Excel/SQL sources, use `--parse-types` to emit numbers, booleans, and nulls instead of strings. Leading-zero values and long IDs stay strings.

**Examples:**
```bash
# Array of objects (default)
tableconvert data.csv output.json --format=object

# 2D array format
tableconvert data.csv output.json --format=2d

# Minified output
tableconvert data.csv output.json --minify

# Column-oriented object
tableconvert data.csv output.json --format=column

# Root object and 4-space indentation
tableconvert data.csv output.json --root=data --indent=4
```

---

### JSONL (JSON Lines)

**Usage:** `tableconvert data.csv output.jsonl --parsing-json`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `parse-types` | `false` | `true`, `false` | Convert values to numbers/bools/null on output (alias: `parsing-json`) |
| `parsing-json` | `false` | `true`, `false` | Alias for `parse-types` |
| `flatten-arrays` | `json` | `json`, `index` | Serialize arrays into one cell, or expand to indexed columns |
| `sort-headers` | `false` | `true`, `false` | Sort headers alphabetically instead of first-seen order |

**Example:**
```bash
tableconvert data.csv output.jsonl
```

---

### LaTeX

**Usage:** `tableconvert data.csv output.tex --bold-header --text-align=c`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `bold-first-column` | `false` | `true`, `false` | Bold first column |
| `bold-first-row` | `false` | `true`, `false` | Bold first row (headers) |
| `borders` | `1111,1111` | Various border patterns | Table border style |
| `caption` | `` | Any string | Table caption |
| `escape` | `true` | `true`, `false` | Escape LaTeX special characters |
| `ht` | `false` | `true`, `false` | Place here or top of page |
| `label` | `` | Any string | Table label for referencing |
| `location` | `above` | `above`, `below` | Caption location |
| `mwe` | `false` | `true`, `false` | Minimal working example |
| `table-align` | `centering` | `centering`, `raggedleft`, `raggedright` | Table alignment |
| `text-align` | `l` | `l`, `c`, `r` | Text alignment |

**Examples:**
```bash
# Centered table with caption
tableconvert data.csv output.tex --caption="Experiment Results" --table-align=centering

# Bold headers and centered text
tableconvert data.csv output.tex --bold-first-row --text-align=c

# Minimal working example
tableconvert data.csv output.tex --mwe
```

---

### Markdown

**Usage:** `tableconvert data.csv output.md --bold-header --align=l,c,r`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `align` | `l` | `l`, `c`, `r` (comma-separated) | Column alignment |
| `bold-header` | `false` | `true`, `false` | Bold header row |
| `bold-first-column` | `false` | `true`, `false` | Bold first column |
| `escape` | `true` | `true`, `false` | Escape Markdown characters |
| `pretty` | `true` | `true`, `false` | Pretty-print with calculated widths |

**Examples:**
```bash
# Bold headers with custom alignment
tableconvert data.csv output.md --bold-header --align=l,c,r

# No escaping (for pre-formatted content)
tableconvert data.csv output.md --escape=false

# Compact format
tableconvert data.csv output.md --pretty=false
```

---

### MediaWiki

**Usage:** `tableconvert data.csv output.wiki --first-row-header --sort`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `first-row-header` | `false` | `true`, `false` | Use first row as headers |
| `minify` | `false` | `true`, `false` | Minify output |
| `sort` | `false` | `true`, `false` | Make table sortable in Wikipedia |

**Example:**
```bash
# Sortable MediaWiki table
tableconvert data.csv output.wiki --sort --first-row-header
```

---

### MySQL

**Usage:** `tableconvert data.csv output.ascii --style=box`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `style` | `box` | `box` | MySQL table style (box format) |

**Example:**
```bash
tableconvert data.csv output.ascii --style=box
```

**Important Notes about MySQL Format:**

**Reading MySQL (Input):**
- ✅ **Supported**: MySQL query output with box borders (from `mysql -t` command)
- ❌ **NOT Supported**:
  - `mysqldump` output (SQL INSERT/CREATE statements)
  - Raw SQL queries without `-t` flag
  - Multi-statement SQL scripts
  - MySQL export files

**Correct Usage:**
```bash
# ✅ CORRECT - Use mysql -t for box format
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown

# ✅ CORRECT - Schema description
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown

# ❌ INCORRECT - mysqldump not supported
mysqldump --no-data users | tableconvert --from=mysql --to=markdown  # FAILS

# ❌ INCORRECT - Missing -t flag
mysql -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown  # FAILS
```

**Alternatives for Unsupported Formats:**
- For `mysqldump` output: Use `sql` format for INSERT statements
- For raw SQL: Process with other tools first to convert to box format
- For schema documentation: Use `DESCRIBE` or `SHOW CREATE TABLE` with `-t` flag

---

**Writing MySQL (Output):**
- ✅ **Supported**: Generate MySQL box format tables from any input
- ✅ **Use case**: Create formatted tables for documentation

---

### Template

**Usage:** `tableconvert data.csv output.php --template=templates/php_array.tmpl`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `template` | `` | File path | Template file path |

**Template Functions Available:**
- `Upper`, `Lower`, `Capitalize`
- `Sub` (integer subtraction)
- `Substr` (substring; rune indexes, end exclusive)
- `Quote` (add quotes)
- Various escape functions

**Example:**
```bash
tableconvert data.csv output.php --template=templates/php_array.tmpl
```

---

### SQL

**Usage:** `tableconvert data.csv output.sql --table=users --dialect=mysql`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `one-insert` | `false` | `true`, `false` | Multiple rows in one INSERT |
| `parse-types` | `false` | `true`, `false` | Emit numbers/bools/null as SQL literals |
| `replace` | `false` | `true`, `false` | Use REPLACE instead of INSERT |
| `dialect` | `mysql` | `none`, `mysql`, `oracle`, `mssql`, `postgresql` | SQL dialect |
| `table` | `` | Any string | Table name |

**Examples:**
```bash
# Single INSERT with multiple rows
tableconvert data.csv output.sql --one-insert --table=users

# REPLACE instead of INSERT
tableconvert data.csv output.sql --replace --table=users

# PostgreSQL dialect
tableconvert data.csv output.sql --dialect=postgresql --table=users
```

---

### TWiki / TracWiki

**Usage:** `tableconvert data.csv output.twiki --first-row-header`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `first-row-header` | `false` | `true`, `false` | Use first row as headers |

**Example:**
```bash
tableconvert data.csv output.twiki --first-row-header
```

---

### XML

**Usage:** `tableconvert data.csv output.xml --minify --root-element=data`

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `minify` | `false` | `true`, `false` | Minify XML output |
| `root-element` | `dataset` | Any string | Root element tag |
| `row-element` | auto | Any string | Row element tag; auto-detects repeated child elements when omitted |
| `declaration` | `false` | `true`, `false` | Include XML declaration |

**Input behavior:**
- If `--row-element` is specified, only direct children with that name are rows; no match is an error.
- Otherwise, `tableconvert` auto-detects the repeated child element (e.g. `<item>`, `<record>`), or treats the root as one record when its children are fields.
- Header names are the union of field names; missing fields are padded.

**Example:**
```bash
# Custom element names
tableconvert data.csv output.xml --root-element=table --row-element=row

# Minified output
tableconvert data.csv output.xml --minify
```

---

## 🌐 Encoding

Text formats are decoded as UTF-8 by default (an optional UTF-8 BOM is stripped).

| Parameter | Default | Allowed Values | Description |
|-----------|---------|----------------|-------------|
| `encoding` | `utf-8` | `utf-8`, `utf-16le`, `utf-16be`, `gbk`, `gb18030`, `big5`, `shift_jis`, `latin1`, `auto` | Input text encoding |
| `output-encoding` | `utf-8` | `utf-8`, `utf-16le`, `utf-16be`, `gbk`, `gb18030`, `big5`, `shift_jis`, `latin1` | Output text encoding |

`auto` only detects UTF-8/UTF-16 BOMs; for GBK/Big5, specify the encoding explicitly. Excel files are binary and ignore these options.

```bash
# GBK CSV exported from Windows Excel
tableconvert --from=csv --to=json --encoding=gbk data.csv data.json

# Write GBK CSV for Excel
tableconvert --from=json --to=csv --output-encoding=gbk data.json data.csv
```

## 🎯 Quick Reference by Use Case

### For Documentation
```bash
# Markdown with nice formatting
tableconvert data.csv doc.md --bold-header --align=l,c,r

# HTML for web
tableconvert data.csv doc.html --thead --minify

# LaTeX for papers
tableconvert data.csv doc.tex --caption="Results" --text-align=c
```

### For Data Processing
```bash
# JSON for APIs
tableconvert data.csv data.json --format=object --minify

# CSV with specific delimiter
tableconvert data.json data.csv --delimiter=TAB

# Batch conversion
tableconvert --batch="data/*.csv" --to=json --output-dir=processed
```

### For Database Work
```bash
# SQL INSERT statements
tableconvert data.csv data.sql --table=users --one-insert

# MySQL format
tableconvert data.csv data.ascii --style=box

# Schema documentation
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown
```

### For Wikis
```bash
# MediaWiki
tableconvert data.csv wiki.wiki --sort --first-row-header

# TWiki
tableconvert data.csv wiki.twiki --first-row-header
```

---

## 🔍 Getting Help

### Show All Formats
```bash
tableconvert --help-formats
```

### Show Format-Specific Help
```bash
tableconvert --help-format=markdown
tableconvert --help-format=json
tableconvert --help-format=latex
```

### Verbose Mode for Debugging
```bash
tableconvert data.csv output.json --verbose
```

---

## 📝 Notes

- **Parameter Names**: Use lowercase with hyphens (e.g., `--bold-header`)
- **Boolean Values**: Use `true`/`false` or just the flag (e.g., `--bold-header` or `--bold-header=true`)
- **Multiple Values**: Use comma-separated for lists (e.g., `--align=l,c,r`)
- **File Paths**: Always quote paths with spaces
- **Unknown Options**: Unknown options and invalid enum values are rejected so typos cannot silently change output. Use `--allow-unknown` only for forward-compatibility experiments.
- **Auto-Detection**: When in doubt, let tableconvert detect formats from extensions