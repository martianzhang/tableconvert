# Known Limitations

This page lists current boundaries so you can choose the right workaround.

## Excel

- Input/output supports `.xlsx` only. Legacy `.xls` files must be saved as `.xlsx` or exported to CSV.
- One worksheet is read per conversion. Use `--sheet-index=N` or `--sheet-name=...`; `--list-sheets` prints the names. `--all-sheets` is not implemented yet.
- Merged cells are expanded by default. Use `--expand-merged=false` to keep raw layout; a warning is printed because only the top-left value is retained.
- Output keeps `text-format=true` by default, so numbers are stored as text in Excel. This is intentional; see the format reference.

## Encoding

- `--encoding=auto` detects UTF-8/UTF-16 BOMs and valid UTF-8. GBK/GB18030/Big5/Shift_JIS need an explicit `--encoding`.
- Excel files are binary and ignore `--encoding`.

## JSON

- JSON → JSON/JSONL preserves nested values, `null`, and number text only when no destructive transformation is applied. `transpose` and case transformations fall back to the string/flattened view.
- Nested objects are flattened to dot paths for tabular formats; arrays are compact JSON strings unless `--flatten-arrays=index` is used.

## Header options

- `--skip-rows` / `--header-row` are currently implemented for CSV, TSV, and Excel. Other formats use their native header rules.

## Data size

- CSV/JSON/SQL inputs are currently read into memory, and `Table` is fully in-memory. Very large files can exhaust memory.
- Streaming and `--max-rows` / `--max-bytes` guards are not implemented yet.

## SQL

- SQL input supports `INSERT` statements. `CREATE TABLE`, `INSERT ... SELECT`, and multi-table dumps are not parsed.
- SQL output generates `INSERT`/`REPLACE`; there is no `CREATE TABLE` generation.

## Rendering / integrations

- PDF, PNG, and JPEG output is not supported.
- MCP is text-oriented. Binary targets such as Excel are not supported over MCP.
