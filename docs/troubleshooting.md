# Troubleshooting Guide

This guide helps you solve common issues when using tableconvert.

## 🔍 Common Issues and Solutions

### 1. Format Detection Issues

**Problem:** "Format not detected" or "Unknown format" error

**Solutions:**
```bash
# Explicitly specify formats
tableconvert --from=csv --to=json input.csv output.json

# Check supported formats
tableconvert --help-formats

# Verify file extension
tableconvert --help-format=markdown  # Check format details
```

**Common mistakes:**
- ❌ `tableconvert data.txt output.json` (txt not auto-detected)
- ✅ `tableconvert --from=csv --to=json data.txt output.json`

### 2. Invalid Table Format

**Problem:** "Parse error" or "Invalid table format"

**Solutions:**
```bash
# Use verbose mode to see detailed errors
tableconvert input.csv output.json --verbose

# Check input format
head input.csv

# Test with simple data first
echo "a,b,c\n1,2,3" | tableconvert --from=csv --to=json
```

**For MySQL format specifically:**
- Ensure complete box borders: `+----+----+`
- Check for multi-byte characters
- Verify line endings (Unix vs Windows)

### 2a. MySQL Format Compatibility Issues

**Problem:** "Parse error" or "Invalid table format" when using MySQL commands

**Root Cause:** `tableconvert` only supports MySQL **box format** output (like `mysql -t`), NOT:
- `mysqldump` output (SQL INSERT/CREATE statements)
- Raw SQL queries without `-t` flag
- Multi-statement SQL scripts

**Solutions:**
```bash
# ✅ CORRECT - Use mysql -t for query results
mysql -t -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown

# ✅ CORRECT - Use mysql -t for schema
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown

# ❌ WRONG - mysqldump not supported
mysqldump --no-data users | tableconvert --from=mysql --to=markdown

# ❌ WRONG - Missing -t flag
mysql -e "SELECT * FROM users" | tableconvert --from=mysql --to=markdown
```

**Alternatives for Unsupported Formats:**
```bash
# For raw SQL (Requires additional processing tools)
tableconvert dump.sql output.md --from=sql --to=markdown

# For schema documentation: Use DESCRIBE
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown > schema.md
```

**Verify MySQL Output Format:**
```bash
# This should show box borders
mysql -t -e "SELECT 1, 2"

# Expected output:
# +---+---+
# | 1 | 2 |
# +---+---+
# | 1 | 2 |
# +---+---+
```

**Quick Test:**
```bash
# Test if your MySQL output is compatible
mysql -t -e "SELECT 'test' as col1, 123 as col2" | tableconvert --from=mysql --to=json
# Should work

# This will fail
mysql -e "SELECT 'test' as col1, 123 as col2" | tableconvert --from=mysql --to=json
# Will show parse error
```

### 3. Batch Mode Issues

**Problem:** Batch processing not working or no files processed

**Solutions:**
```bash
# Quote glob patterns
tableconvert --batch="*.csv" --to=json

# Use recursive search for subdirectories
tableconvert --batch="**/*.csv" --to=json --recursive

# Check pattern matches
ls data/*.csv  # Verify files exist

# Use verbose mode
tableconvert --batch="data/*.csv" --to=json --verbose
```

**Common mistakes:**
- ❌ `tableconvert --batch=*.csv --to=json` (unquoted glob)
- ✅ `tableconvert --batch="*.csv" --to=json`

### 4. File Permission Issues

**Problem:** "Permission denied" or "Cannot create file"

**Solutions:**
```bash
# Check file permissions
ls -la output_directory/

# Create output directory first
mkdir -p output_dir
tableconvert --batch="data/*.csv" --to=json --output-dir=output_dir

# Use current directory with write permissions
tableconvert input.csv ./output.json
```

### 5. Encoding Issues

**Problem:** Special characters, UTF-8 issues, or garbled output

**Solutions:**
```bash
# Ensure UTF-8 input
file input.csv  # Check encoding

# For CSV with BOM
tableconvert input.csv output.csv --bom

# Disable escaping if needed
tableconvert input.csv output.md --escape=false

# Check terminal encoding
echo $LANG  # Should be UTF-8
```

### 6. Memory Issues with Large Files

**Problem:** Slow performance or memory errors with large files

**Solutions:**
```bash
# Use streaming for large files
cat large.csv | tableconvert --from=csv --to=json > output.json

# Process in chunks if possible
head -n 1000 large.csv > chunk1.csv
tableconvert chunk1.csv chunk1.json

# Use JSONL for very large datasets
tableconvert large.csv large.jsonl  # Line-by-line format
```

### 7. Template Issues

**Problem:** Template not found or template syntax errors

**Solutions:**
```bash
# Specify full path to template
tableconvert data.csv output.php --template=/full/path/to/template.tmpl

# Check template file exists
ls template.tmpl

# Verify template syntax
# Template should use Go template syntax: {{.FieldName}}
```

### 8. Excel Issues

**Problem:** Cannot read or write Excel files

**Solutions:**
```bash
# Verify file format
file data.xlsx  # Should show "Microsoft Excel"

# Try opening in Excel first to verify file integrity

# For reading: ensure first row is headers if needed
tableconvert data.xlsx output.csv --first-column-header=true

# For writing: check output directory permissions
tableconvert data.csv output.xlsx --auto-width
```

### 9. SQL Generation Issues

**Problem:** SQL syntax errors or missing table name

**Solutions:**
```bash
# Always specify table name
tableconvert data.csv output.sql --table=users

# Choose correct dialect
tableconvert data.csv output.sql --table=users --dialect=postgresql

# For multiple rows in one INSERT
tableconvert data.csv output.sql --table=users --one-insert

# Check output before running
cat output.sql
```

### 10. LaTeX Compilation Issues

**Problem:** Generated LaTeX doesn't compile

**Solutions:**
```bash
# Use minimal working example
tableconvert data.csv output.tex --mwe

# Check for special characters that need escaping
tableconvert data.csv output.tex --escape=true

# Add proper caption and label
tableconvert data.csv output.tex --caption="My Table" --label=tab:mytable

# Verify output
cat output.tex
```

## 🛠️ Diagnostic Commands

### Check Installation
```bash
# Verify binary exists
which tableconvert

# Check version (if available)
tableconvert --version 2>/dev/null || echo "No version info"

# Test basic functionality
echo "a,b\n1,2" | tableconvert --from=csv --to=json
```

### Debug Input Issues
```bash
# Show first few lines
head input.csv

# Check file size
ls -lh input.csv

# Verify encoding
file input.csv

# Count lines
wc -l input.csv
```

### Test Conversions
```bash
# Simple test
echo "name,age\nAlice,30" | tableconvert --from=csv --to=json

# Verbose test
echo "name,age\nAlice,30" | tableconvert --from=csv --to=json --verbose

# Test with file
echo "name,age\nAlice,30" > test.csv
tableconvert test.csv test.json --verbose
rm test.csv test.json
```

## 📋 Error Messages Reference

### "Format not detected"
**Cause:** File extension not recognized or missing
**Solution:** Use `--from` and `--to` explicitly

### "Parse error on line X"
**Cause:** Invalid table format in input
**Solution:** Check input format, use `--verbose` for details

### "Parse error" or "Invalid table format" (MySQL)
**Cause:** Input is not MySQL box format (missing `-t` flag or using mysqldump)
**Solution:** Use `mysql -t -e "QUERY"` instead of `mysql -e "QUERY"` or `mysqldump`

### "Unknown format: XYZ"
**Cause:** Format not supported
**Solution:** Run `tableconvert --help-formats` to see supported formats

### "Permission denied"
**Cause:** Cannot read input or write output
**Solution:** Check file permissions and directory access

### "Template file not found"
**Cause:** Template path incorrect or file missing
**Solution:** Use absolute path or verify file exists

### "Batch pattern no matches"
**Cause:** No files match the glob pattern
**Solution:** Verify pattern and file existence with `ls pattern`

## 🎯 Quick Fixes

### Reset and Test
```bash
# Create test data
echo "a,b,c\n1,2,3\n4,5,6" > test.csv

# Test basic conversion
tableconvert test.csv test.json

# Check output
cat test.json

# Clean up
rm test.csv test.json
```

### Verbose Debugging
```bash
# Always use verbose for issues
tableconvert input.csv output.json --verbose

# Check what was detected
tableconvert --help-formats

# Verify format parameters
tableconvert --help-format=markdown
```

### Working with Stdin/Stdout
```bash
# Test without files
echo "a,b\n1,2" | tableconvert --from=csv --to=json

# Pipe through tools
cat data.csv | tableconvert --from=csv --to=json | jq .
```

## 📞 Getting More Help

### Built-in Help
```bash
# General help
tableconvert --help

# Format list
tableconvert --help-formats

# Format-specific help
tableconvert --help-format=markdown
tableconvert --help-format=json
```

### Project Resources
- **GitHub Issues:** https://github.com/martianzhang/tableconvert/issues
- **Documentation:** Check the docs/ directory
- **Examples:** See docs/examples.md

### Information to Include in Bug Reports
1. Tableconvert version
2. Operating system
3. Input file format and sample
4. Command used
5. Error message (with `--verbose`)
6. Expected vs actual output

## 🔄 Common Workflows

### CSV → JSON → Database
```bash
# Convert CSV to SQL
tableconvert data.csv data.sql --table=users --one-insert

# Verify SQL
cat data.sql

# Run in database
mysql -u user -p < data.sql
```

### Excel → Markdown Documentation
```bash
# Convert to Markdown
tableconvert data.xlsx data.md --bold-header --align=l,c,r

# Add to README
cat data.md >> README.md
```

### Batch Clean and Convert
```bash
# Clean all CSV files
tableconvert --batch="data/*.csv" --to=csv --delete-empty --deduplicate --output-dir=clean

# Convert cleaned files
tableconvert --batch="clean/*.csv" --to=json --output-dir=json_output
```

### Database Schema to Docs
```bash
# Get schema in box format (use -t flag!)
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown > schema.md

# Or for multiple tables
for table in users products orders; do
  mysql -t -e "DESCRIBE $table" | tableconvert --from=mysql --to=markdown > "${table}_schema.md"
done

# Note: tableconvert does NOT support mysqldump output directly
# For mysqldump, you'd need to parse the SQL first or use alternative tools
```

## 💡 Prevention Tips

1. **Always test with small samples first**
2. **Use verbose mode when debugging**
3. **Quote glob patterns in batch mode**
4. **Check file permissions before batch operations**
5. **Verify input format matches expected format**
6. **Use auto-detection when file extensions are reliable**
7. **Backup data before batch transformations**
8. **Test output before processing large files**

## 🚀 Performance Tips

### Large Files
```bash
# Use JSONL for very large datasets
tableconvert large.csv large.jsonl

# Stream processing
cat large.csv | tableconvert --from=csv --to=json > large.json

# Process in chunks
split -l 10000 large.csv chunk_
for f in chunk_*; do tableconvert "$f" "${f%.csv}.json"; done
```

### Batch Operations
```bash
# Use recursive for deep directory structures
tableconvert --batch="**/*.csv" --to=json --recursive --verbose

# Monitor progress
tableconvert --batch="data/*.csv" --to=json --verbose 2>&1 | tee log.txt
```

### Memory Efficiency
```bash
# Avoid loading entire file if possible
# Use streaming for very large files
cat huge.csv | tableconvert --from=csv --to=json > huge.json

# Use JSONL for line-by-line processing
tableconvert huge.csv huge.jsonl
```