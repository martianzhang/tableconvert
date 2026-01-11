# tableconvert

Offline table convert tool. **In progress, Not production ready.**

`tableconvert` command line tool, which is writen in Go. It's designed for converting between different table formats, such as MySQL Query Output, Excel/CSV, SQL queries, and potentially other tabular data structures.

## Key Features

- **Multi-format Conversion**: Easily transform data between different formats, such as MySQL Query Output, Markdown, Excel/CSV, SQL queries, and potentially other tabular data structures.
- **Offline available**: No need to connect to the internet to convert between formats. Your data is safe and secure. 
- **Template Support**: Support template output, which can be used to generate any other formats, like `PHP`, `Python` data structs, etc.
- **Cross-platform**: Support Windows, Linux, and Mac OS.

## Usage Example

```bash
# Basic conversion with flags
tableconvert --from=csv --to=json --file=input.csv --result=output.json

# Short flags
tableconvert -i input.csv -o output.json

# Auto-detect formats from file extensions
tableconvert input.csv output.json

# Read from stdin and write to stdout
cat input.csv | tableconvert --from=csv --to=json

# With format-specific options
tableconvert input.csv output.md --bold-header --align=l,c,r

# Transpose and capitalize
tableconvert input.csv output.md --transpose --capitalize

# Convert from MySQL to Markdown using template
tableconvert --from=mysql --to=template --file=input.csv --template=markdown.tmpl

# Batch mode: convert all CSV files in a directory
tableconvert --batch="data/*.csv" --to=json

# Batch mode with recursive search
tableconvert --batch="data/**/*.csv" --to=json --recursive

# Batch mode with output directory
tableconvert --batch="input/*.csv" --to=json --output-dir=output

# Batch mode with short flags
tableconvert -b "*.csv" -t json -r -v
```

More usage info please refer to [Usage](https://github.com/martianzhang/tableconvert/blob/main/common/usage.txt).

### Command Line Options

**Basic Options:**
- `--from|-f={FORMAT}` - Source format (e.g. mysql, csv, json, xlsx)
- `--to|-t={FORMAT}` - Target format (e.g. mysql, csv, json, xlsx)
- `--file|--input|-i={PATH}` - Input file path (or use stdin if not specified)
- `--result|--output|-o={PATH}` - Output file path (or use stdout if not specified)

**Quick Options:**
- `-v, --verbose` - Enable verbose output
- `--mcp` - Run as MCP (Model Context Protocol) server
- `-h, --help` - Show help message
- `--help-formats` - Show all supported formats and their parameters
- `--help-format={FORMAT}` - Show parameters for a specific format

**Batch Processing Options:**
- `--batch|-b={PATTERN}` - Process multiple files matching a pattern (e.g., `*.csv`, `data/*.json`)
- `--recursive|-r` - Enable recursive directory traversal for batch mode
- `--output-dir|--dir={PATH}` - Specify output directory for batch results (default: same as input)

**Data Transformation Options:**
- `--transpose` - Transpose the table (swap rows and columns)
- `--delete-empty` - Remove empty rows from the table
- `--deduplicate` - Remove duplicate rows
- `--uppercase` - Convert all text to UPPERCASE
- `--lowercase` - Convert all text to lowercase
- `--capitalize` - Capitalize the first letter of each cell

**Auto-Detection:**
When `--from` or `--to` are not specified, tableconvert will attempt to detect the format from file extensions:
- `.csv` → csv
- `.json` → json
- `.jsonl` → jsonl
- `.md`, `.markdown` → markdown
- `.xlsx`, `.xls` → excel
- `.html`, `.htm` → html
- `.xml` → xml
- `.sql` → sql
- `.tex`, `.latex` → latex
- `.wiki` → mediawiki
- `.tmpl`, `.template` → tmpl
- `.txt` → (not auto-detected, must specify)

Each format or file type has its own arguments, please refer to the [arguments.md](https://github.com/martianzhang/tableconvert/blob/main/docs/arguments.md) for more details.

## MCP (Model Context Protocol) Usage

`tableconvert` provides MCP stdio tools for AI assistants like Claude Code.

### Add to Claude Code

Add this to your Claude Code settings:

```bash
claude mcp add tableconvert -- /path/to/tableconvert --mcp
```

On Windows:
```bash
claude mcp add tableconvert -- "C:\path\to\tableconvert.exe" --mcp
```

Or add directly to your config file.

```json
{
  "mcpServers": {
    "tableconvert": {
      "command": "/path/to/tableconvert",
      "args": ["--mcp"]
    }
  }
}
```

### Available Tools

- **`convert_table`**: Convert table data between formats
- **`get_formats`**: Get information about supported formats and their parameters

## Support Format

- [x] Excel
- [x] CSV
- [x] XML
- [x] HTML
- [x] Markdown
- [x] JSON
- [x] JSONL
- [x] SQL
- [x] MySQL
- [x] LaTeX
- [x] MediaWiki
- [x] TWiki/TracWiki
- [x] User Define template Output

### MySQL Query Output Example

```txt
+----------+--------------+------+-----+---------+----------------+
| FIELD    | TYPE         | NULL | KEY | DEFAULT | EXTRA          |
+----------+--------------+------+-----+---------+----------------+
| user_id  | smallint(5)  | NO   | PRI | NULL    | auto_increment |
| username | varchar(10)  | NO   |     | NULL    |                |
| password | varchar(100) | NO   |     |         |                |
+----------+--------------+------+-----+---------+----------------+
```

## CSV Format Table Example

```csv
FIELD,TYPE,NULL,KEY,DEFAULT,EXTRA
user_id,smallint(5),NO,PRI,NULL,auto_increment
username,varchar(10),NO,,NULL,
password,varchar(100),NO,,,
```

## Markdown Format Table Example

```md
| FIELD    | TYPE         | NULL | KEY | DEFAULT | EXTRA          |
|----------|--------------|------|-----|---------|----------------|
| user_id  | smallint(5)  | NO   | PRI | NULL    | auto_increment |
| username | varchar(10)  | NO   |     | NULL    |                |
| password | varchar(100) | NO   |     |         |                |
```

## INSERT SQL Example

```sql
INSERT INTO `{table_name}` (`FIELD`, `TYPE`, `NULL`, `KEY`, `DEFAULT`, `EXTRA`) VALUES ('user_id', 'smallint(5)', 'NO', 'PRI', NULL, 'auto_increment');
INSERT INTO `{table_name}` (`FIELD`, `TYPE`, `NULL`, `KEY`, `DEFAULT`, `EXTRA`) VALUES ('username', 'varchar(10)', 'NO', '', NULL, '');
INSERT INTO `{table_name}` (`FIELD`, `TYPE`, `NULL`, `KEY`, `DEFAULT`, `EXTRA`) VALUES ('password', 'varchar(100)', 'NO', '', '', '');
```

## Reference

* [online tableconvert](https://tableconvert.com/)
* [ascii-tables](https://github.com/ozh/ascii-tables)
* [tablewriter](https://github.com/olekukonko/tablewriter)
* [csvq](https://github.com/mithrandie/csvq)

## Dependency

* [excelize](https://github.com/xuri/excelize)
* [sqlparser](https://vitess.io/vitess)

## License

[Apache License 2.0](https://github.com/martianzhang/tableconvert/blob/main/LICENSE)