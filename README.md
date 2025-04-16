# tableconvert

Offline table convert tool. **Not production ready.**

## Usage

```txt
Usage: tableconvert [OPTIONS]

Convert between different table formats (MySQL, Markdown, CSV, JSON, Excel, etc.)

Options:
  --from|-f={FORMAT}     Source format (e.g. mysql, csv, json, xlsx)
  --to|-t={FORMAT}       Target format (e.g. mysql, csv, json, xlsx)
  --file={PATH}          Input file path (or use stdin if not specified)
  --result|-r={PATH}     Output file path (or use stdout if not specified)
  --verbose|-v           Enable verbose output
  -h|--help              Show this help message

Examples:
  tableconvert --from=csv --to=json --file=input.csv --result=output.json
  cat input.csv | tableconvert --from=csv --to=json
```

Each format or file type has its own arguments, please refer to the [arguments.md](https://github.com/martianzhang/tableconvert/blob/master/arguments.md) for more details.

## Support source

- [x] Excel
- [x] CSV
- [x] XML
- [x] HTML
- [x] Markdown
- [x] JSON
- [x] SQL
- [x] MySQL
- [x] LaTeX
- [x] MediaWiki
- [x] TWiki/TracWiki

## Support destination

- [ ] actionscript
- [x] ascii
- [ ] asciidoc
- [ ] asp
- [ ] avro
- [ ] bbcode
- [x] csv
- [ ] dax
- [x] excel
- [ ] firebase
- [x] html
- [ ] ini
- [ ] jira
- [ ] jpeg
- [x] json
- [ ] jsonlines
- [x] latex
- [ ] magic
- [x] markdown
- [ ] matlab
- [x] mediawiki
- [ ] pandasdataframe
- [ ] pdf
- [ ] php
- [ ] png
- [ ] protobuf
- [ ] qlik
- [ ] rdataframe
- [ ] rdf
- [ ] restructuredtext
- [ ] ruby
- [x] sql
- [ ] textile
- [ ] toml
- [x] tracwiki
- [x] xml
- [ ] yaml

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