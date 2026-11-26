# Documentation Index

Welcome to the tableconvert documentation! This index will help you find what you need quickly.

## 🚀 Quick Start

### New to tableconvert?
1. **[README.md](../README.md)** - Overview, features, and quick start
2. **[Quick Reference](quick-reference.md)** - Common commands at a glance
3. **[Practical Examples](examples.md)** - Real-world usage scenarios

### Need to convert something right now?
```bash
# Basic conversion
tableconvert input.csv output.json

# Get help for any format
tableconvert --help-format=markdown
```

## 📚 Documentation by Topic

### Usage & Commands
| Topic | File | What you'll learn |
|-------|------|-------------------|
| **All Commands** | [quick-reference.md](quick-reference.md) | Complete command cheat sheet |
| **Format Options** | [arguments.md](arguments.md) | Every parameter for every format |
| **Examples** | [examples.md](examples.md) | Real-world scenarios |
| **CLI Help** | [../common/usage.txt](../common/usage.txt) | Built-in help text |

### Formats & Conversions
| Format | Documentation | Common Uses |
|--------|---------------|-------------|
| **CSV** | [arguments.md#csv](arguments.md#csv) | Data import/export |
| **TSV** | [arguments.md#csv--tsv-delimited-values](arguments.md#csv--tsv-delimited-values) | Tab-separated data |
| **JSON** | [arguments.md#json](arguments.md#json) | APIs, web apps |
| **Markdown** | [arguments.md#markdown](arguments.md#markdown) | Documentation |
| **Excel** | [arguments.md#excel](arguments.md#excel) | Spreadsheets |
| **HTML** | [arguments.md#html](arguments.md#html) | Web pages |
| **SQL** | [arguments.md#sql](arguments.md#sql) | Database operations |
| **LaTeX** | [arguments.md#latex](arguments.md#latex) | Academic papers |
| **MySQL** | [arguments.md#mysql](arguments.md#mysql) | Database output |
| **XML** | [arguments.md#xml](arguments.md#xml) | Data interchange |
| **MediaWiki** | [wiki.md](wiki.md) | Wikipedia/Wikis |
| **TWiki** | [wiki.md](wiki.md) | Corporate wikis |
| **YAML** | [arguments.md#yaml](arguments.md#yaml) | Configuration/data |
| **TOML** | [arguments.md#toml](arguments.md#toml) | Configuration/data |
| **INI** | [arguments.md#ini](arguments.md#ini) | Configuration/data |
| **BBCode / Jira / AsciiDoc / reST / Textile** | [arguments.md#documentation-formats](arguments.md#documentation-formats) | Docs and wikis |
| **PHP / Ruby / Pandas** | [arguments.md#code-formats](arguments.md#code-formats) | Code literals |
| **Template** | [arguments.md#template](arguments.md#template) | Custom formats |

### Problem Solving
| Issue | Solution |
|-------|----------|
| **Format not detected** | [Troubleshooting - Format Detection](troubleshooting.md#1-format-detection-issues) |
| **Parse errors** | [Troubleshooting - Invalid Format](troubleshooting.md#2-invalid-table-format) |
| **Batch mode issues** | [Troubleshooting - Batch Mode](troubleshooting.md#3-batch-mode-issues) |
| **Encoding problems** | [Troubleshooting - Encoding](troubleshooting.md#5-encoding-issues) |
| **Performance issues** | [Troubleshooting - Memory Issues](troubleshooting.md#6-memory-issues-with-large-files) |
| **Known limitations** | [limitations.md](limitations.md) | Current boundaries and workarounds |

## 🎯 Common Tasks

### I want to...
- **Convert CSV to JSON** → `tableconvert data.csv data.json`
- **Convert MySQL to Markdown** → `mysql -t -e "QUERY" | tableconvert --from=mysql --to=markdown`
- **Batch convert files** → `tableconvert --batch="*.csv" --to=json`
- **Generate SQL INSERTs** → `tableconvert data.csv data.sql --table=users`
- **Create LaTeX tables** → `tableconvert data.csv table.tex --caption="Results"`
- **Transform data** → `tableconvert data.csv output.md --transpose --capitalize`
- **Find all format options** → `tableconvert --help-format=<format>`
- **Debug issues** → Use `--verbose` flag

### Quick Examples
```bash
# CSV → JSON
tableconvert data.csv data.json

# CSV → Markdown (styled)
tableconvert data.csv data.md --bold-header --align=l,c,r

# MySQL → Markdown
mysql -t -e "DESCRIBE users" | tableconvert --from=mysql --to=markdown

# Batch conversion
tableconvert --batch="data/*.csv" --to=json --output-dir=results

# Data cleaning
tableconvert data.csv clean.csv --delete-empty --deduplicate

# SQL generation
tableconvert data.csv data.sql --table=users --one-insert
```

## 🔍 Find Information

### By Use Case
- **Data Processing** → [examples.md](examples.md#-data-transformations)
- **Database Work** → [examples.md](examples.md#-database-work)
- **Documentation** → [examples.md](examples.md#-reporting--documentation)
- **Web Development** → [examples.md](examples.md#-web-development-html)
- **Academic Writing** → [examples.md](examples.md#-academic-papers-latex)

### By Format
- **CSV/JSON/XML** → [arguments.md](arguments.md)
- **Excel** → [arguments.md#excel](arguments.md#excel)
- **SQL** → [arguments.md#sql](arguments.md#sql)
- **LaTeX** → [arguments.md#latex](arguments.md#latex) + [latex.md](latex.md)
- **Wikis** → [wiki.md](wiki.md)
- **Templates** → [arguments.md#template](arguments.md#template)

### By Feature
- **Batch Processing** → [examples.md](examples.md#-batch-processing)
- **Data Transformations** → [arguments.md#global-transformations](arguments.md#-global-transformations)
- **MCP Integration** → [README.md#-mcp-model-context-protocol-integration](../README.md#-mcp-model-context-protocol-integration)
- **Auto-Detection** → [README.md#auto-detection](../README.md#auto-detection)

## 🛠️ Getting Help

### Built-in Help
```bash
# General help
tableconvert --help

# All formats
tableconvert --help-formats

# Format-specific help
tableconvert --help-format=markdown
tableconvert --help-format=json
tableconvert --help-format=latex
```

### Debugging
```bash
# Verbose mode
tableconvert input.csv output.json --verbose

# Test with simple data
echo "a,b\n1,2" | tableconvert --from=csv --to=json
```

### Additional Resources
- **[Troubleshooting Guide](troubleshooting.md)** - Common issues and solutions
- **[GitHub Issues](https://github.com/martianzhang/tableconvert/issues)** - Report bugs
- **[CLAUDE.md](../CLAUDE.md)** - Developer documentation

## 📖 Reading Order for Beginners

### Step 1: Understanding (5 minutes)
1. Read [README.md](../README.md) overview
2. Look at [Quick Reference](quick-reference.md)

### Step 2: First Conversion (2 minutes)
```bash
# Create test file
echo "name,age\nAlice,30\nBob,25" > test.csv

# Convert it
tableconvert test.csv test.json

# Check result
cat test.json
```

### Step 3: Explore Formats (10 minutes)
1. Try different formats: `tableconvert test.csv test.md`
2. Add styling: `tableconvert test.csv test.md --bold-header`
3. See all options: `tableconvert --help-format=markdown`

### Step 4: Learn Advanced Features
1. **Transformations**: [arguments.md#global-transformations](arguments.md#-global-transformations)
2. **Batch Processing**: [examples.md#-batch-processing](examples.md#-batch-processing)
3. **Real Examples**: [examples.md](examples.md)

### Step 5: Solve Problems
- Check [Troubleshooting](troubleshooting.md) for common issues
- Use `--verbose` for debugging
- Ask for help on GitHub

## 📋 Quick Commands Reference

```bash
# Installation
git clone https://github.com/martianzhang/tableconvert.git
cd tableconvert
make build

# Basic usage
tableconvert input.csv output.json
tableconvert --from=csv --to=markdown input.csv output.md

# Help
tableconvert --help
tableconvert --help-formats
tableconvert --help-format=markdown

# Batch
tableconvert --batch="*.csv" --to=json --output-dir=results

# Transformations
tableconvert data.csv output.md --transpose --capitalize --bold-header

# Debug
tableconvert input.csv output.json --verbose
```

---

## 🤝 Contributing

Found a typo? Missing example? Want to add documentation for a new feature?

1. Edit the appropriate markdown file
2. Add clear examples
3. Test your examples
4. Submit a pull request

See [CLAUDE.md](../CLAUDE.md) for development guidelines.

## 📞 Contact

- **Issues**: [GitHub Issues](https://github.com/martianzhang/tableconvert/issues)
- **Documentation**: This directory (`docs/`)
- **Main Project**: [README.md](../README.md)