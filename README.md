# tableconvert

Offline table convert tool.

## Support source

- [ ] Excel
- [ ] CSV
- [ ] XML
- [ ] HTML
- [ ] Markdown
- [ ] JSON
- [ ] SQL
- [ ] MySQL
- [ ] LaTeX
- [ ] MediaWiki

## Support destination

- [ ] actionscript
- [ ] ascii
- [ ] asciidoc
- [ ] asp
- [ ] avro
- [ ] bbcode
- [ ] csv
- [ ] dax
- [ ] excel
- [ ] firebase
- [ ] html
- [ ] ini
- [ ] jira
- [ ] jpeg
- [ ] json
- [ ] jsonlines
- [ ] latex
- [ ] magic
- [ ] markdown
- [ ] matlab
- [ ] mediawiki
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
- [ ] sql
- [ ] textile
- [ ] toml
- [ ] tracwiki
- [ ] xml
- [ ] yaml

## Usage

```bash
# read from stdin pipeline
cat {filename} | tableconvert --from {from_type} --to {to_type}

# read from file
tableconvert --from {from_type} --to {to_type} --file {filename}
```

## Reference

* [online tableconvert](https://tableconvert.com/)
* [ascii-tables](https://github.com/ozh/ascii-tables)
* [tablewriter](https://github.com/olekukonko/tablewriter)
* [csvq](https://github.com/mithrandie/csvq)

## Dependency


## License

Apache License 2.0