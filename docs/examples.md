# Practical Examples

This document provides real-world examples of tableconvert usage for common scenarios.

## 📊 Data Format Conversions

### CSV to JSON
```bash
# Basic conversion
tableconvert data.csv data.json

# Array of objects (default)
tableconvert data.csv data.json --format=object

# 2D array format
tableconvert data.csv data.json --format=2d

# Minified JSON
tableconvert data.csv data.json --minify

# With type inference (strings to numbers/booleans)
tableconvert data.csv data.json --format=object
```

**Input CSV:**
```csv
name,age,active
Alice,30,true
Bob,25,false
```

**Output JSON:**
```json
[
  {"name": "Alice", "age": 30, "active": true},
  {"name": "Bob", "age": 25, "active": false}
]
```

### JSON to CSV
```bash
# Basic conversion
tableconvert data.json data.csv

# Tab-separated values
tableconvert data.json data.csv --delimiter=TAB

# With BOM for Excel
tableconvert data.json data.csv --bom

# Semicolon delimiter (European format)
tableconvert data.json data.csv --delimiter=SEMICOLON
```

### MySQL to Markdown
```bash
# Convert table schema (use -t for box format)
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown > schema.md

# With bold headers and custom alignment
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown --bold-header --align=l,c,r,c,c

# Note: tableconvert only supports MySQL box format output (mysql -t)
# It does NOT support mysqldump output or raw SQL statements
```

**MySQL Output:**
```
+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| id       | int(11)      | NO   | PRI | NULL    | auto_increment |
| username | varchar(50)  | NO   |     | NULL    |                |
| email    | varchar(100) | YES  |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
```

**Markdown Output:**
```markdown
| Field    | Type         | Null | Key | Default | Extra          |
|----------|--------------|------|-----|---------|----------------|
| id       | int(11)      | NO   | PRI | NULL    | auto_increment |
| username | varchar(50)  | NO   |     | NULL    |                |
| email    | varchar(100) | YES  |     | NULL    |                |
```

### Excel to HTML
```bash
# Convert Excel to HTML table
tableconvert data.xlsx data.html --thead

# Minified HTML for web
tableconvert data.xlsx data.html --minify --thead

# Div-based table
tableconvert data.xlsx data.html --div
```

## 🔧 Data Transformations

### Transpose Table
```bash
# Swap rows and columns
tableconvert data.csv output.md --transpose

# Useful for converting wide tables to tall tables
tableconvert sales_data.csv rotated.md --transpose
```

**Before (CSV):**
```csv
Month,Sales,Expenses,Profit
Jan,1000,800,200
Feb,1200,900,300
```

**After (Transposed Markdown):**
```markdown
| | Month | Sales | Expenses | Profit |
|---|---|---|---|---|
| Jan | 1000 | 800 | 200 |
| Feb | 1200 | 900 | 300 |
```

### Clean Data
```bash
# Remove empty rows
tableconvert data.csv clean.csv --delete-empty

# Remove duplicates
tableconvert data.csv unique.csv --deduplicate

# Combine transformations
tableconvert data.csv clean.csv --delete-empty --deduplicate --uppercase
```

### Case Conversion
```bash
# All uppercase
tableconvert data.csv output.md --uppercase

# All lowercase
tableconvert data.csv output.md --lowercase

# Capitalize first letters
tableconvert data.csv output.md --capitalize
```

## 🗄️ Database Work

### Generate SQL INSERT Statements
```bash
# Basic INSERT statements
tableconvert data.csv data.sql --table=users

# Single INSERT with multiple rows
tableconvert data.csv data.sql --table=users --one-insert

# REPLACE instead of INSERT
tableconvert data.csv data.sql --table=users --replace

# PostgreSQL dialect
tableconvert data.csv data.sql --table=users --dialect=postgresql --one-insert

# Oracle dialect
tableconvert data.csv data.sql --table=users --dialect=oracle
```

**Input CSV:**
```csv
id,name,email
1,Alice,alice@example.com
2,Bob,bob@example.com
```

**Output SQL (one-insert):**
```sql
INSERT INTO `users` (`id`, `name`, `email`) VALUES
(1, 'Alice', 'alice@example.com'),
(2, 'Bob', 'bob@example.com');
```

### MySQL query result output as table box
```bash
# Generate Markdown documentation from query
mysql -t -e "select * from  users" | tableconvert --from=mysql --to=markdown > users.md

# Batch document all tables
for table in $(mysql -NBe "SHOW TABLES" | tail -n +2); do
  mysql -t -e "select * from $table" | tableconvert --from=mysql --to=markdown > "${table}.md"
done
```

## 📑 Reporting & Documentation

### Academic Papers (LaTeX)
```bash
# Basic LaTeX table
tableconvert results.csv table.tex

# With caption and centered alignment
tableconvert results.csv table.tex --caption="Experiment Results" --table-align=centering --text-align=c

# Bold headers and specific borders
tableconvert results.csv table.tex --bold-first-row --borders=1111,1111

# Minimal working example
tableconvert results.csv table.tex --mwe --caption="Results" --label=tab:results
```

### Web Development (HTML)
```bash
# Clean HTML table
tableconvert data.csv table.html --thead

# Minified for production
tableconvert data.csv table.html --thead --minify

# Div-based for CSS styling
tableconvert data.csv table.html --div --minify
```

### GitHub Documentation (Markdown)
```bash
# Nice looking table
tableconvert data.csv README.md --bold-header --align=l,c,r

# Without escaping (if content is safe)
tableconvert data.csv README.md --bold-header --escape=false

# Compact format
tableconvert data.csv README.md --pretty=false
```

## 🔄 Batch Processing

### Convert Multiple Files
```bash
# All CSV files in directory
tableconvert --batch="data/*.csv" --to=json

# Recursive search
tableconvert --batch="**/*.csv" --to=json --recursive

# With output directory
tableconvert --batch="input/*.csv" --to=json --output-dir=output

# Verbose mode for progress
tableconvert --batch="data/*.csv" --to=json --verbose

# Different format per file type
tableconvert --batch="data/*.csv" --to=json --output-dir=json_files
tableconvert --batch="data/*.xlsx" --to=csv --output-dir=csv_files
```

### Transform All Files
```bash
# Add transformations to batch
tableconvert --batch="data/*.csv" --to=md --transpose --bold-header

# Clean and convert
tableconvert --batch="data/*.csv" --to=json --delete-empty --deduplicate
```

## 🎨 Format-Specific Examples

### Markdown with Custom Styling
```bash
# Left, center, right alignment
tableconvert data.csv output.md --align=l,c,r --bold-header

# All columns centered
tableconvert data.csv output.md --align=c,c,c --bold-header

# Bold first column
tableconvert data.csv output.md --bold-first-column

# Compact without pretty printing
tableconvert data.csv output.md --pretty=false
```

### JSON with Different Formats
```bash
# Array of objects (default)
tableconvert data.csv output.json --format=object

# 2D array
tableconvert data.csv output.json --format=2d

# Column-oriented
tableconvert data.csv output.json --format=column

# Keyed by first column
tableconvert data.csv output.json --format=keyed

# Minified
tableconvert data.csv output.json --format=object --minify
```

### Excel with Styling
```bash
# Auto-width columns
tableconvert data.csv output.xlsx --auto-width

# Headers from first row
tableconvert data.csv output.xlsx --first-column-header=true

# Custom sheet name
tableconvert data.csv output.xlsx --sheet-name="Sales Report"

# Combined options
tableconvert data.csv output.xlsx --auto-width --first-column-header=true --sheet-name="Data"
```

### SQL with Various Options
```bash
# Multiple dialects
tableconvert data.csv output.sql --table=users --dialect=mysql
tableconvert data.csv output.sql --table=users --dialect=postgresql
tableconvert data.csv output.sql --table=users --dialect=mssql
tableconvert data.csv output.sql --table=users --dialect=oracle

# REPLACE instead of INSERT
tableconvert data.csv output.sql --table=users --replace

# No escaping
tableconvert data.csv output.sql --table=users --dialect=none
```

### LaTeX Advanced
```bash
# Full-featured table
tableconvert data.csv output.tex \
  --caption="Performance Metrics" \
  --label=tab:metrics \
  --table-align=centering \
  --text-align=c \
  --bold-first-row \
  --location=above

# Minimal example
tableconvert data.csv output.tex --mwe
```

## 💡 Real-World Scenarios

### Web API Data Preparation
```bash
# Convert CSV API data to JSON for frontend
tableconvert api_data.csv frontend_data.json --format=object --minify

# Transform for different endpoint
tableconvert api_data.csv users.json --format=object --transpose
```

### Data Analysis Pipeline
```bash
# Clean, transform, and export
tableconvert raw_data.csv clean_data.csv --delete-empty --deduplicate --uppercase
tableconvert clean_data.csv analysis.json --format=object

# Batch process daily reports
tableconvert --batch="reports/*.csv" --to=json --output-dir=processed --deduplicate
```

### Documentation Generation
```bash
# Generate documentation from data
tableconvert schema.csv schema.md --bold-header --align=l,c,c,c,c,c
tableconvert sample_data.csv sample.md --bold-header --align=l,c,r

# Create LaTeX appendix
tableconvert full_data.csv appendix.tex --caption="Complete Dataset" --mwe
```

### Database Migration
```bash
# Export to SQL for migration
tableconvert data.csv migration.sql --table=users --one-insert --dialect=postgresql

# Verify with MySQL format
tableconvert data.csv verification.txt --style=box
```

### Wiki Content
```bash
# MediaWiki for Wikipedia
tableconvert data.csv wiki.wiki --sort --first-row-header

# TWiki for corporate wiki
tableconvert data.csv wiki.twiki --first-row-header
```

## 🎯 Quick Reference

### Common Command Patterns
```bash
# CSV → JSON
tableconvert input.csv output.json

# CSV → Markdown
tableconvert input.csv output.md --bold-header --align=l,c,r

# JSON → CSV
tableconvert input.json output.csv --delimiter=TAB

# MySQL → Markdown (use -t for box format)
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown

# Excel → HTML
tableconvert input.xlsx output.html --thead

# Batch conversion
tableconvert --batch="*.csv" --to=json --output-dir=results

# Data cleaning
tableconvert input.csv output.csv --delete-empty --deduplicate --uppercase

# SQL generation
tableconvert data.csv output.sql --table=users --one-insert

# LaTeX for papers
tableconvert data.csv output.tex --caption="Results" --text-align=c
```

### Transformation Combinations
```bash
# Clean and format
tableconvert data.csv output.md --delete-empty --deduplicate --bold-header

# Transform and convert
tableconvert data.csv output.json --transpose --format=object

# Batch clean
tableconvert --batch="*.csv" --to=csv --delete-empty --deduplicate --output-dir=clean
```

### Format-Specific Quick Commands
```bash
# Markdown
tableconvert data.csv out.md --bold-header --align=l,c,r

# JSON
tableconvert data.csv out.json --format=object --minify

# HTML
tableconvert data.csv out.html --thead --minify

# SQL
tableconvert data.csv out.sql --table=users --one-insert

# LaTeX
tableconvert data.csv out.tex --caption="Table" --text-align=c

# Excel
tableconvert data.csv out.xlsx --auto-width

# CSV
tableconvert data.json out.csv --delimiter=TAB

# XML
tableconvert data.csv out.xml --minify

# MediaWiki
tableconvert data.csv out.wiki --sort
```

## 📝 Tips and Best Practices

1. **Always use quotes for file paths with spaces**
   ```bash
   tableconvert "my data.csv" "output file.json"
   ```

2. **Use verbose mode for debugging**
   ```bash
   tableconvert input.csv output.json --verbose
   ```

3. **Test with stdin/stdout first**
   ```bash
   echo "a,b,c\n1,2,3" | tableconvert --from=csv --to=json
   ```

4. **Check format help for specific options**
   ```bash
   tableconvert --help-format=markdown
   ```

5. **Use auto-detection when possible**
   ```bash
   tableconvert input.csv output.json  # No --from/--to needed
   ```

6. **Batch mode for multiple files**
   ```bash
   tableconvert --batch="data/*.csv" --to=json --output-dir=results
   ```

7. **Combine transformations**
   ```bash
   tableconvert data.csv output.md --transpose --capitalize --bold-header
   ```

8. **Backup before batch operations**
   ```bash
   cp -r data data_backup
   tableconvert --batch="data_backup/*.csv" --to=json
   ```