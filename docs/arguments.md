# extended arguments

Each format may has its own extended arguments, for its own style. Following are the arguments for each file type.

## ascii

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| style             | box     | box, plus(+), dot(·), bubble(◌) | Table Style |


## csv

| Argument            | Default | Allowed Values              | Description                  |
|---------------------|---------|-----------------------------|------------------------------|
| first-column-header | false   | true, false                 | Use first column as headers  |
| transpose           | false   | true, false            | Transpose table columns with rows |
| bom                 | false   |                             | Add Byte Order Mark          |
| delimiter           | ,(COMMA)| COMMA, TAB, SEMICOLON, PIPE, SLASH, HASH | Value Delimiter |

## excel

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-column-header | false | true, false  | Use first column as headers |
| transpose    | false   | true, false | Transpose table columns with rows |
| sheet-name        | Sheet1  |                | Excel Sheet Name          |
| auto-width        | false   | true, false    | Auto Width                |
| text-format       | true    | true, false    | force text format         |

## html

| Argument            | Default | Allowed Values | Description                 |
|---------------------|---------|----------------|-----------------------------|
| first-column-header | false   | true, false    | Use first column as headers |
| transpose           | false   | true, false    | Transpose table columns with rows |
| div                 | false   | true, false    | Convert into div table      |
| minify              | false   | true, false    | Minify HTML table           |
| thead               | false   | true, false    | Include thead and tbody tags|

## json

| Argument          | Default | Allowed Values | Description               |
| ------------------|---------|----------------|---------------------------|
| format            | object  | object, 2d, column, keyed | JSON Format    |
| minify            | false   | true, false    | Minify JSON               |
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

## markdown

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| align             | l       | l, c, r        | Text Alignment            |
| bold-header       | false   | true, false    | Table Header Bold         |
| bold-first-column | false   | true, false    | Bold first column         |
| escape            | false   | true, false    | Escape Markdown table     |
| pretty            | true    | true, false    | Pretty-print Markdown     |

## mediawiki

| Argument          | Default | Allowed Values | Description               |
|-------------------|---------|----------------|---------------------------|
| first-row-header  | false   | true, false    | Use first row as headers  |
| minify            | false   | true, false    | Minify MediaWiki table    |
| sort              | false   | true, false    | Make table sortable in Wikipedia |

## tmpl

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

## xml

| Argument            | Default | Allowed Values | Description                |
|---------------------|---------|----------------|----------------------------|
| minify              | false   | true, false    | Minify XML                 |
| root-element        | dataset | string         | Root Element Tag           |
| row-element         | record  | string         | Row Element Tag            |
| declaration         | true    | true, false    | Include XML Declaration    |
