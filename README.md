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
# Convert from CSV to JSON
tableconvert --from=csv --to=json --file=input.csv --result=output.json
# Read from stdin and write to stdout
cat input.csv | tableconvert --from=csv --to=json
# Convert from MySQL to Markdown using template
tableconvert --from=mysql --to=template --file=input.csv --tempalte=markdown.tmpl
```

More usage info please refer to [Usage](https://github.com/martianzhang/tableconvert/blob/main/common/usage.txt).

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