# tableconvert

Offline table convert tool. **Not production ready.**

## Usage Example

```bash
# Convert from CSV to JSON
tableconvert --from=csv --to=json --file=input.csv --result=output.json
# Read from stdin and write to stdout
cat input.csv | tableconvert --from=csv --to=json
# Convert from MySQL to Markdown using template
tableconvert --from=mysql --to=template --file=input.csv --tempalte=markdown.tmpl
```

Simple usage please refer to [Usage](https://github.com/martianzhang/tableconvert/blob/main/common/usage.txt).

Each format or file type has its own arguments, please refer to the [arguments.md](https://github.com/martianzhang/tableconvert/blob/main/docs/arguments.md) for more details.

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