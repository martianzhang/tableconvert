# Quick Reference Card

Fast lookup for common tableconvert commands.

## 🚀 Basic Commands

```bash
# Auto-detect formats
tableconvert input.csv output.json

# Explicit formats
tableconvert --from=csv --to=markdown input.csv output.md

# Short flags
tableconvert -i input.csv -o output.json

# Stdin/stdout
echo "a,b\n1,2" | tableconvert --from=csv --to=json
```

## 📊 Format Conversions

| From | To | Command |
|------|----|---------|
| CSV | JSON | `tableconvert data.csv data.json` |
| CSV | Markdown | `tableconvert data.csv data.md --bold-header` |
| JSON | CSV | `tableconvert data.json data.csv` |
| MySQL | Markdown | `mysql -t -e "QUERY" \| tableconvert --from=mysql --to=markdown` |
| Excel | HTML | `tableconvert data.xlsx data.html --thead` |
| CSV | SQL | `tableconvert data.csv data.sql --table=users` |
| CSV | LaTeX | `tableconvert data.csv data.tex --caption="Table"` |
| CSV | XML | `tableconvert data.csv data.xml` |
| CSV | MediaWiki | `tableconvert data.csv data.wiki --sort` |

## 🎨 Format-Specific Options

### Markdown
```bash
--bold-header          # Bold headers
--align=l,c,r          # Column alignment
--escape=false         # Disable escaping
--pretty=false         # Compact format
```

### JSON
```bash
--format=object        # Array of objects (default)
--format=2d            # 2D array
--format=column        # Column-oriented
--minify               # Minified output
```

### CSV
```bash
--delimiter=TAB        # Tab delimiter
--bom                  # Add BOM
--first-column-header  # First column as headers
```

### Excel
```bash
--auto-width           # Auto column width
--sheet-name="Name"    # Custom sheet name
--first-column-header  # First column as headers
```

### HTML
```bash
--div                  # Div wrapper
--minify               # Minify output
--thead                # Include thead/tbody
```

### SQL
```bash
--table=users          # Table name (required)
--one-insert           # Single INSERT with multiple rows
--dialect=mysql        # SQL dialect (mysql, postgresql, mssql, oracle)
--replace              # Use REPLACE instead of INSERT
```

### LaTeX
```bash
--caption="Text"       # Table caption
--text-align=c         # Text alignment (l, c, r)
--bold-first-row       # Bold headers
--table-align=centering # Table alignment
--mwe                  # Minimal working example
```

### XML
```bash
--minify               # Minify output
--root-element=data    # Root element name
--row-element=record   # Row element name
```

### ASCII
```bash
--style=box            # Box style (default)
--style=plus           # Plus style
--style=dot            # Dot style
--style=bubble         # Bubble style
```

## 🔧 Data Transformations

```bash
--transpose            # Swap rows and columns
--delete-empty         # Remove empty rows
--deduplicate          # Remove duplicates
--uppercase            # Convert to UPPERCASE
--lowercase            # Convert to lowercase
--capitalize           # Capitalize first letters
```

**Combine transformations:**
```bash
tableconvert data.csv output.md --transpose --capitalize --bold-header
```

## ⚡ Batch Processing

```bash
# Basic batch
tableconvert --batch="*.csv" --to=json

# Recursive
tableconvert --batch="**/*.csv" --to=json --recursive

# With output directory
tableconvert --batch="data/*.csv" --to=json --output-dir=results

# Verbose
tableconvert --batch="*.csv" --to=json --verbose

# With transformations
tableconvert --batch="*.csv" --to=md --bold-header --delete-empty
```

## 🆘 Help Commands

```bash
# General help
tableconvert --help

# All formats
tableconvert --help-formats

# Format-specific help
tableconvert --help-format=markdown
tableconvert --help-format=json
tableconvert --help-format=latex
tableconvert --help-format=sql
```

## 🔍 Auto-Detection Reference

| Extension | Format |
|-----------|--------|
| `.csv` | csv |
| `.json` | json |
| `.jsonl` | jsonl |
| `.md`, `.markdown` | markdown |
| `.xlsx`, `.xls` | excel |
| `.html`, `.htm` | html |
| `.xml` | xml |
| `.sql` | sql |
| `.tex`, `.latex` | latex |
| `.wiki` | mediawiki |
| `.tmpl`, `.template` | tmpl |
| `.twiki` | twiki |

**Note:** `.txt` files require explicit format specification.

## 💡 Common Patterns

### Database to Documentation
```bash
# Schema to Markdown
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown > schema.md

# Query results to JSON
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=json > data.json
```

### Data Cleaning Pipeline
```bash
# Clean and convert
tableconvert raw.csv clean.json --delete-empty --deduplicate --uppercase
```

### Report Generation
```bash
# CSV to formatted Markdown
tableconvert data.csv report.md --bold-header --align=l,c,r

# CSV to LaTeX for papers
tableconvert results.csv table.tex --caption="Results" --text-align=c
```

### API Data Preparation
```bash
# CSV to JSON for API
tableconvert api_data.csv api_data.json --format=object --minify
```

## 📋 Quick Examples

```bash
# 1. CSV → JSON
tableconvert data.csv data.json

# 2. CSV → Markdown with styling
tableconvert data.csv data.md --bold-header --align=l,c,r

# 3. MySQL → Markdown
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown

# 4. Batch convert
tableconvert --batch="data/*.csv" --to=json --output-dir=json_files

# 5. SQL generation
tableconvert data.csv data.sql --table=users --one-insert

# 6. LaTeX table
tableconvert data.csv table.tex --caption="Data" --text-align=c

# 7. Data transformation
tableconvert data.csv output.md --transpose --capitalize

# 8. Excel to HTML
tableconvert data.xlsx data.html --thead --minify
```

## 🎯 Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| Format not detected | Use `--from` and `--to` explicitly |
| Parse error | Use `--verbose` to see details |
| Batch not working | Quote pattern: `--batch="*.csv"` |
| Permission denied | Check file/directory permissions |
| Special characters | Use `--escape=false` or check encoding |

## 📂 File Structure

```
tableconvert/
├── README.md              # Main documentation
├── docs/
│   ├── arguments.md       # Complete parameter reference
│   ├── examples.md        # Practical examples
│   ├── troubleshooting.md # Problem solving
│   └── quick-reference.md # This file
└── common/
    └── usage.txt          # CLI help text
```

## 🚀 Getting Started Checklist

- [ ] Build: `make build`
- [ ] Test: `echo "a,b\n1,2" | tableconvert --from=csv --to=json`
- [ ] Check formats: `tableconvert --help-formats`
- [ ] Try basic conversion: `tableconvert input.csv output.json`
- [ ] Explore: `tableconvert --help-format=markdown`

---

**Tip:** Use `--verbose` for debugging and `--help-format=<format>` for format-specific options.