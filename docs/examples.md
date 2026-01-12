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

## 🤖 AI Assistant Integration (MCP Mode)

### Using tableconvert with Claude Code

tableconvert can run as an MCP server for AI assistants:

```bash
# Start MCP server (stdio transport)
tableconvert --mcp

# Add to Claude Code (Linux/macOS)
claude mcp add tableconvert -- /path/to/tableconvert --mcp

# Add to Claude Code (Windows)
claude mcp add tableconvert -- "C:\\path\\to\\tableconvert.exe" --mcp
```

### MCP Tool Usage Examples

Once configured, AI assistants can use these tools:

**Convert data with natural language:**
```
"Convert this CSV to JSON: name,age\nAlice,30\nBob,25"
```

**Get format information:**
```
"What options does the markdown format support?"
```

**Complex transformations:**
```
"Convert my Excel file to markdown with bold headers and center alignment"
```

### Programmatic MCP Usage

```bash
# Start server for external clients
tableconvert --mcp --verbose

# Test MCP server
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"convert_table","arguments":{"from":"csv","to":"markdown","input":"name,age\\nAlice,30\\nBob,25"}}}' | tableconvert --mcp
```

## 🔧 Debugging & Troubleshooting

### Verbose Mode for Debugging

**Basic verbose output:**
```bash
tableconvert input.csv output.json --verbose
```

**Verbose with batch processing:**
```bash
tableconvert --batch="data/*.csv" --to=json --verbose
```

**Verbose to log file:**
```bash
tableconvert --batch="data/*.csv" --to=json --verbose 2>&1 | tee conversion.log
```

**Verbose output shows:**
- Format auto-detection results
- Parsing progress
- Extension parameters being used
- Error details with line numbers
- File processing status

### Common Error Scenarios

**Parse errors with verbose:**
```bash
# MySQL format issues
tableconvert --from=mysql --to=markdown --verbose < bad_mysql.txt

# CSV with malformed quotes
tableconvert --from=csv --to=json --verbose < messy.csv
```

**What verbose reveals:**
```
Processing: stdin -> stdout
Detected format: mysql
Parsing line 1: +----------+--------------+
Parsing line 2: | Field    | Type         |
...
Error at line 15: Expected closing border, got: "| id | int(11) |"
```

### Testing Conversions

**Test with stdin/stdout first:**
```bash
echo "name,age\nAlice,30\nBob,25" | tableconvert --from=csv --to=json --verbose
```

**Test format detection:**
```bash
# Let tableconvert auto-detect
echo "name,age\nAlice,30" | tableconvert --to=json --verbose
```

**Test specific format parameters:**
```bash
# Check what parameters are available
tableconvert --help-format=markdown

# Test with parameters
echo "a,b,c\n1,2,3" | tableconvert --to=markdown --bold-header --align=l,c,r --verbose
```

## 🚀 Performance & Large Files

### Batch Processing with Progress

**Verbose batch with summary:**
```bash
tableconvert --batch="data/*.csv" --to=json --output-dir=results --verbose
```

**Expected verbose output:**
```
Processing: data/sales.csv -> results/sales.json
  ✓ Success
Processing: data/users.csv -> results/users.json
  ✓ Success
Processed: 2 files, 2 succeeded, 0 failed
```

### Recursive Directory Processing

**Deep batch conversion:**
```bash
# Process all CSV files in directory tree
tableconvert --batch="**/*.csv" --to=json --recursive --verbose

# With output structure preserved
tableconvert --batch="**/data/*.csv" --to=json --recursive --output-dir=converted
```

### Error Handling in Batch

**Partial failure handling:**
```bash
# Some files may fail, but others continue
tableconvert --batch="*.csv" --to=json --verbose
# Output shows which files succeeded/failed
```

**Clean up failed outputs:**
```bash
# tableconvert automatically removes failed output files
# Use verbose to see cleanup actions
tableconvert --batch="*.csv" --to=json --verbose
```

## 💡 Advanced Real-World Patterns

### CI/CD Integration

**GitHub Actions:**
```yaml
- name: Convert CSV to JSON
  run: tableconvert data.csv data.json --verbose

- name: Batch convert reports
  run: tableconvert --batch="reports/*.csv" --to=json --output-dir=dist
```

**Pre-commit hook:**
```bash
#!/bin/bash
# Verify all CSV files can be converted to JSON
for file in $(find . -name "*.csv"); do
  echo "Checking $file..."
  if ! tableconvert "$file" /tmp/test.json --verbose; then
    echo "Failed to convert $file"
    exit 1
  fi
done
```

### Data Pipeline Script

**Robust conversion script:**
```bash
#!/bin/bash
set -e

INPUT_DIR="raw_data"
OUTPUT_DIR="processed"
LOG_FILE="conversion.log"

echo "Starting batch conversion..." | tee $LOG_FILE

tableconvert --batch="$INPUT_DIR/*.csv" \
  --to=json \
  --output-dir=$OUTPUT_DIR \
  --verbose 2>&1 | tee -a $LOG_FILE

# Check exit code
if [ $? -eq 0 ]; then
  echo "✅ All conversions successful" | tee -a $LOG_FILE
else
  echo "❌ Some conversions failed - check $LOG_FILE" >&2
  exit 1
fi
```

### Data Cleaning Pipeline

**Multi-step transformation:**
```bash
# Clean, deduplicate, then convert
tableconvert raw.csv cleaned.csv \
  --delete-empty \
  --deduplicate \
  --uppercase \
  --verbose

tableconvert cleaned.csv final.json \
  --format=object \
  --verbose
```

**One-liner pipeline:**
```bash
tableconvert raw.csv cleaned.json \
  --delete-empty --deduplicate --uppercase \
  --format=object --verbose
```

### Database Documentation

**Generate docs from multiple tables:**
```bash
#!/bin/bash
DB="mydb"
OUTPUT="docs"

mkdir -p $OUTPUT

# Get all tables
TABLES=$(mysql -NBe "SHOW TABLES" $DB)

for table in $TABLES; do
  echo "Documenting $table..."
  mysql -t -e "DESCRIBE $table" $DB | \
    tableconvert --from=mysql --to=markdown \
      --bold-header --align=l,c,c,c,c,c \
      > "$OUTPUT/${table}_schema.md"

  mysql -t -e "SELECT * FROM $table LIMIT 10" $DB | \
    tableconvert --from=mysql --to=markdown \
      --bold-header --align=l,c,r \
      > "$OUTPUT/${table}_sample.md"
done

echo "Documentation generated in $OUTPUT/"
```

### Performance Monitoring

**Time large conversions:**
```bash
# Time a batch operation
time tableconvert --batch="large_data/*.csv" --to=json --output-dir=results --verbose
```

**Monitor memory usage:**
```bash
# Use /usr/bin/time if available
/usr/bin/time -v tableconvert --batch="*.csv" --to=json 2>&1 | grep -E "(User|System|Elapsed|Maximum resident)"
```

## 🎯 Quick Reference: Debugging Commands

### Format Detection Issues
```bash
# See what format is detected
echo "name,age\nAlice,30" | tableconvert --to=json --verbose

# Force specific format
tableconvert --from=csv --to=json --verbose < data.txt
```

### Parsing Errors
```bash
# MySQL box format issues
tableconvert --from=mysql --to=markdown --verbose < bad_output.txt

# CSV quote issues
tableconvert --from=csv --to=json --verbose < messy.csv
```

### Batch Failures
```bash
# Process with verbose to see individual file status
tableconvert --batch="*.csv" --to=json --verbose

# Check which files failed
tableconvert --batch="*.csv" --to=json --verbose 2>&1 | grep "Failed"
```

### MCP Server Issues
```bash
# Test MCP server directly
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | tableconvert --mcp

# Run with verbose to see communication
tableconvert --mcp --verbose
```

## 📊 Performance Tips

1. **Use verbose mode to identify bottlenecks**
   ```bash
   tableconvert --batch="*.csv" --to=json --verbose
   ```

2. **Process in chunks for very large datasets**
   ```bash
   # Split then batch
   split -l 10000 large.csv chunk_
   tableconvert --batch="chunk_*" --to=json --output-dir=results
   ```

3. **Use appropriate JSON format for your use case**
   ```bash
   # object (default) - best for APIs
   # 2d - best for spreadsheets
   # column - best for analytics
   # keyed - best for lookups
   tableconvert data.csv data.json --format=object --verbose
   ```

4. **Monitor conversion progress**
   ```bash
   # Count files first
   ls -1 data/*.csv | wc -l

   # Then process with verbose
   tableconvert --batch="data/*.csv" --to=json --verbose
   ```

## 🛡️ Safety Best Practices

### Always Test First
```bash
# Test single file
tableconvert test.csv test.json --verbose

# Test with sample
head -5 data.csv | tableconvert --from=csv --to=json --verbose
```

### Backup Before Batch
```bash
# Create backup
cp -r data data_backup

# Process backup
tableconvert --batch="data_backup/*.csv" --to=json --output-dir=results --verbose
```

### Verify Output
```bash
# Check conversion worked
tableconvert results/file.json results_check.csv --verbose

# Compare row counts
wc -l data_backup/*.csv
wc -l results/*.json
```

### Use Dry-Run Concept (via verbose)
```bash
# See what would be processed without actually converting
tableconvert --batch="*.csv" --to=json --verbose 2>&1 | grep "Processing:"
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

9. **Use MCP mode for AI-assisted conversions**
   ```bash
  tableconvert --mcp
   # Then ask your AI: "Convert CSV to JSON with type inference"
   ```

10. **Debug with verbose when things go wrong**
    ```bash
    tableconvert input.csv output.json --verbose
    # Look for: format detection, parsing progress, error details
    ```
