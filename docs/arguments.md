# Format-Specific Arguments

Each format has its own extension parameters for custom styling. Following are the parameters for each format.

**Note:** Global transformation parameters work with all formats:
- `--transpose` - Transpose the table (swap rows and columns)
- `--delete-empty` - Remove empty rows from the table
- `--deduplicate` - Remove duplicate rows
- `--uppercase` - Convert all text to UPPERCASE
- `--lowercase` - Convert all text to lowercase
- `--capitalize` - Capitalize the first letter of each cell

## ascii

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| style             | box     | box, plus(+), dot(·), bubble(◌) | Table Style |


## csv

| Argument            | Default | Allowed Values              | Description                  |
|---------------------|---------|-----------------------------|------------------------------|
| first-column-header | false   | true, false                 | Use first column as headers  |
| bom                 | false   |                             | Add Byte Order Mark          |
| delimiter           | ,       | COMMA, TAB, SEMICOLON, PIPE, SLASH, HASH | Value Delimiter |

## excel (or xlsx)

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-column-header | false | true, false  | Use first column as headers |
| sheet-name        | Sheet1  |                | Excel Sheet Name          |
| auto-width        | false   | true, false    | Auto Width                |
| text-format       | true    | true, false    | force text format         |

## html

| Argument            | Default | Allowed Values | Description                 |
|---------------------|---------|----------------|-----------------------------|
| first-column-header | false   | true, false    | Use first column as headers |
| div                 | false   | true, false    | Convert into div table      |
| minify              | false   | true, false    | Minify HTML table           |
| thead               | false   | true, false    | Include thead and tbody tags|

## json

| Argument          | Default | Allowed Values | Description               |
| ------------------|---------|----------------|---------------------------|
| format            | object  | object, 2d, column, keyed | JSON Format    |
| minify            | false   | true, false    | Minify JSON               |
| parsing-json      | false   | true, false    | Parsing JSON              |

## jsonl (or jsonlines)

| Argument          | Default | Allowed Values | Description               |
| ------------------|---------|----------------|---------------------------|
| parsing-json      | false   | true, false    | Parsing JSON              |

## latex

| Argument          | Default | Allowed Values | Description               |
| ------------------|---------|----------------|---------------------------|
| bold-first-column | false   | true, false    | Bold first column         |
| bold-first-row    | false   | true, false    | Bold first row            |
| borders           | 1111,1111 | 1111,1111, 1101,1101, 0000,1101, 1111,0100, 0000,0100, 0000,0000 | Table Border |
| caption           |         |                | Table Caption             |
| escape            | true    | true, false    | Escape LaTeX table        |
| ht                | false   | true, false    | Place here or top of page |
| label             |         |                | Table Label               |
| location          | above   |  above, below  | Caption Location          |
| mwe               | false   | true, false    | Minimal working example   |
| table-align       | centering | centering, raggedleft, raggedright  | Table Alignment |
| text-align        | l       | l, c, r        | Text Alignment            |

## markdown (or md)

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| align             | l       | l, c, r        | Text Alignment, columns seperate by comma |
| bold-header       | false   | true, false    | Table Header Bold         |
| bold-first-column | false   | true, false    | Bold first column         |
| escape            | true    | true, false    | Escape Markdown table     |
| pretty            | true    | true, false    | Pretty-print Markdown     |

## mediawiki

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-row-header  | false   | true, false    | Use first row as headers  |
| minify            | false   | true, false    | Minify MediaWiki table    |
| sort              | false   | true, false    | Make table sortable in Wikipedia |

## mysql

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| style             | box     | box            | MySQL table style (box format) |

## tmpl (or template)

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| template          |         |                | Template file path        |

## sql

| Argument          | Default | Allowed Values | Description                     |
|-------------------|---------|----------------|---------------------------------|
| one-insert        | false   | true, false    | Insert multiple rows at once    |
| replace           | false   | true, false    | Use REPLACE instead of INSERT   |
| dialect           | mysql   | none, mysql, oracle, mssql, postgresql | identity escape SQL Dialect, none for no escape |
| table             |         |                | Table Name                      |

## twiki (or tracwiki)

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-row-header  | false   | true, false    | Use first row as headers  |

## xml

| Argument            | Default | Allowed Values | Description                |
|---------------------|---------|----------------|----------------------------|
| minify              | false   | true, false    | Minify XML                 |
| root-element        | dataset | string         | Root Element Tag           |
| row-element         | record  | string         | Row Element Tag            |
| declaration         | true    | true, false    | Include XML Declaration    |
