# extended arguments

Each file type may has its own extended arguments, for its own style. Following are the arguments for each file type.

## actionscript

None

## ascii

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| align             | l       | l, c, r        | Text Alignment            |
| comment           | ∅       | ∅, //, #       | Comment style             |
| force-separate    | false   | true, false    | Force separate lines      |
| style             | 1       | 1, 2, 3, 4     | Plain-Text Table Style    |

### ascii table style 1

```txt

```

## asciidoc

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-row-header  | true    | true, false    | Use first row as headers  |
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

| Argument          | Default | Allowed Values              | Description                  |
|-------------------|---------|-----------------------------|------------------------------|
| bom               | false   |                             | Add Byte Order Mark          |
| delimiter         | ,(COMMA)| COMMA, TAB, SEMICOLON, PIPE, SLASH, HASH | Value Delimiter |

## dax

None

## excel

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| auto-width        | false   | true, false    | Auto Width                |
| sheet-name        | Sheet1  |                | Excel Sheet Name          |
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
| first-column-header | false   | true, false    | Use first column as headers |
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
| first-row-header  | true    | true, false    | Use first row as headers  |
| pretty            | true    | true, false    | Pretty-print Markdown     |
| simple            | false   | true, false    | Use simple Markdown table |

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
| create            | false   | true, false    | Generate CREATE TABLE statement |
| drop              | false   | true, false    | Drop table if exists            |
| one-insert        | false   | true, false    | Insert multiple rows at once    |
| quote             | ∅       | ∅, `, [], ", ' | Use Quotes                      |
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

| Argument            | Default | Allowed Values | Description                |
|---------------------|---------|----------------|----------------------------|
| first-column-header | false   | true, false    | Use first column as headers|
| first-row-header    | true    | true, false    | Use first row as headers   |

## xml

| Argument            | Default | Allowed Values | Description                |
|---------------------|---------|----------------|----------------------------|
| escape              | true    | true, false    | Escape XML                 |
| minify              | false   | true, false    | Minify XML                 |
| root-element        | table   | table, row     | Root Element               |
| row-element         | row     | row, cell      | Row Element                |

## yaml

None
