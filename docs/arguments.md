# extended arguments

Each file type may has its own extended arguments, for its own style. Following are the arguments for each file type.

## actionscript

None

## ascii

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| style             | box     | box, plus(+), dot(·), bubble(◌) | Table Style |

## asciidoc

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| last-row-footer   | true    | true, false    | Use last row as footers   |
| minify            | false   | true, false    | Minify AsciiDoc table     |
| title             |         |                | Table Title               |

## asp

None

## avro

None

## bbcode

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| minify            | false   | true, false    | Minify BBCode table       |

## csv

| Argument            | Default | Allowed Values              | Description                  |
|---------------------|---------|-----------------------------|------------------------------|
| first-column-header | false   | true, false                 | Use first column as headers  |
| transpose           | false   | true, false            | Transpose table columns with rows |
| bom                 | false   |                             | Add Byte Order Mark          |
| delimiter           | ,(COMMA)| COMMA, TAB, SEMICOLON, PIPE, SLASH, HASH | Value Delimiter |

## dax

None

## excel

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-column-header | false | true, false  | Use first column as headers |
| transpose    | false   | true, false | Transpose table columns with rows |
| sheet-name        | Sheet1  |                | Excel Sheet Name          |
| auto-width        | false   | true, false    | Auto Width                |
| text-format       | true    | true, false    | force text format         |

## firebase

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| parsing-json      | false   | true, false    | Parsing JSON              |

## html

| Argument            | Default | Allowed Values | Description                 |
|---------------------|---------|----------------|-----------------------------|
| first-column-header | false   | true, false    | Use first column as headers |
| div                 | false   | true, false    | Convert into div table      |
| escape              | true    | true, false    | Escape HTML table           |
| minify              | false   | true, false    | Minify HTML table           |
| thead               | false   | true, false    | Include thead and tbody tags|

## ini

None

## jira

| Argument            | Default | Allowed Values | Description                 |
|---------------------|---------|----------------|-----------------------------|
| escape              | true    | true, false    | Escape Jira table           |
| first-row-header    | false   | true, false    | Use first row as headers    |

## jpeg

None

## json

| Argument          | Default | Allowed Values | Description               |
| ------------------|---------|----------------|---------------------------|
| format            | object  | object, 2d, column, keyed | JSON Format    |
| minify            | false   | true, false    | Minify JSON               |
| parsing-json      | false   | true, false    | Parsing JSON              |
| wrapper           | false   | true, false    | Wrap with 'data'          |

## jsonlines

| Argument          | Default | Allowed Values | Description               |
| ------------------|---------|----------------|---------------------------|
| format            | object  | object, array  | JSONLines Format          |
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

## magic

| Argument          | Default | Allowed Values | Description               |
| ------------------|---------|----------------|---------------------------|
| builtin           | TEMPLATE-CUSTOM | TEMPLATE-CUSTOM, TEMPLATE-SQL, TEMPLATE-DOUBLE-QUOTE-CELL, TEMPLATE-SINGLE-QUOTE-CELL, TEMPLATE-DOUBLE-QUOTE-COlUMN, TEMPLATE-RENAME-REPLACE, TEMPLATE-RENAME-COLUMN, TEMPLATE-NEW-EMPTY-FILE, TEMPLATE-NEW-FILE, TEMPLATE-PHPARR, TEMPLATE-C, TEMPLATE-ROWS, TEMPLATE-CSV, TEMPLATE-OPTION, TEMPLATE-JIRA, TEMPLATE-JSON, TEMPLATE-LDIF, TEMPLATE-URL | Built-in Magic |
| footer            |         |                | Footer Magic              |
| header            |         |                | Header Magic              |
| rows              |         |                | Rows Magic                |

## markdown

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| align             | l       | l, c, r        | Text Alignment            |
| bold-header       | false   | true, false    | Table Header Bold         |
| bold-first-column | false   | true, false    | Bold first column         |
| escape            | false   | true, false    | Escape Markdown table     |
| pretty            | true    | true, false    | Pretty-print Markdown     |

## matlab

None

## mediawiki

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-row-header  | false   | true, false    | Use first row as headers  |
| minify            | false   | true, false    | Minify MediaWiki table    |
| sort              | false   | true, false    | Make table sortable in Wikipedia |

## pandasdataframe

None

## pdf

None

## php

None

## png

None

## protobuf

None

## qlik

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| table             |         |                | Table Name                |

## rdataframe

None

## rdf

None

## restructuredtext

None

## ruby

None

## sql

| Argument          | Default | Allowed Values | Description                     |
|-------------------|---------|----------------|---------------------------------|
| one-insert        | false   | true, false    | Insert multiple rows at once    |
| replace           | false   | true, false    | Use REPLACE instead of INSERT   |
| dialect           | mysql   | none, mysql, oracle, mssql, postgresql | identity escape SQL Dialect, none for no escape |
| table             |         |                | Table Name                      |

## textile

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| escape            | true    | true, false    | Escape Textile table      |
| first-row-header  | false   | true, false    | Use first row as headers  |
| thead             | true    | true, false    | Include thead and tbody tags |

## toml

None

## tracwiki

None

## xml

| Argument            | Default | Allowed Values | Description                |
|---------------------|---------|----------------|----------------------------|
| minify              | false   | true, false    | Minify XML                 |
| root-element        | dataset | string         | Root Element Tag           |
| row-element         | record  | string         | Row Element Tag            |
| declaration         | true    | true, false    | Include XML Declaration    |

## yaml

None
