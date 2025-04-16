package sql

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"vitess.io/vitess/go/vt/sqlparser"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Unmarshal: target table pointer cannot be nil")
	}
	if cfg.Reader == nil {
		return fmt.Errorf("Unmarshal: Reader in Config cannot be nil")
	}

	// Read the entire SQL content
	sqls, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return fmt.Errorf("failed to read SQL data: %w", err)
	}

	// Create the parser instance
	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		return fmt.Errorf("failed to initialize parser: %w", err)
	}

	// Parse the SQL
	stmts, err := parser.SplitStatements(string(sqls))
	if err != nil {
		return fmt.Errorf("failed to parse SQL: %w", err)
	}

	// Handle statements
	for _, stmt := range stmts {
		switch stmt := stmt.(type) {
		case *sqlparser.Insert:
			handleInsert(stmt, table)
		default:
			return fmt.Errorf("unsupported SQL statement type %T", stmt)
		}
	}

	return nil
}

func handleInsert(insert *sqlparser.Insert, table *common.Table) {
	if len(table.Headers) == 0 {
		for _, col := range insert.Columns {
			table.Headers = append(table.Headers, col.String())
		}
	}

	rows, ok := insert.Rows.(sqlparser.Values)
	if !ok {
		// Unsupported row formats (e.g. INSERT SELECT statements)
		return
	}

	for _, row := range rows {
		var values []string
		for _, val := range row {
			switch v := val.(type) {
			case *sqlparser.Literal:
				values = append(values, strings.Trim(string(v.Val), "'"))
			case *sqlparser.NullVal:
				values = append(values, "NULL")
			default:
				values = append(values, sqlparser.String(val)) // fallback: stringify everything else
			}
		}
		table.Rows = append(table.Rows, values)
	}
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}

	columnCount := len(table.Headers)
	if columnCount == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}

	writer := cfg.Writer

	// table name
	tableName := cfg.GetExtensionString("table", "{table_name}")

	// dialect
	dialect := cfg.GetExtensionString("dialect", "mysql")

	// all-in-one
	allInOne := cfg.GetExtensionBool("one-insert", false)

	// INSERT or REPLACE
	var insert = "INSERT"
	if cfg.GetExtensionBool("replace", false) {
		insert = "REPLACE"
	}

	// SQL Prefix
	columns := make([]string, columnCount)
	for i, h := range table.Headers {
		columns[i] = escapeIdentifier(h, dialect)
	}

	// concat statements
	var stmt string
	if allInOne {
		stmt = fmt.Sprintf("%s INTO %s (%s) VALUES\n",
			insert,
			escapeIdentifier(tableName, dialect),
			strings.Join(columns, ", "),
		)
		if _, err := writer.Write([]byte(stmt)); err != nil {
			return fmt.Errorf("failed to write SQL: %w", err)
		}
	}

	// Build and write SQL INSERT statements
	for j, row := range table.Rows {
		if len(row) != columnCount {
			return fmt.Errorf("Marshal: row has %d columns, but table has %d", len(row), columnCount)
		}

		values := make([]string, columnCount)
		for i, cell := range row {
			values[i] = escapeValue(cell)
		}
		if allInOne {
			if j == 0 {
				stmt = fmt.Sprintf("(%s)", strings.Join(values, ", "))
			} else {
				stmt = fmt.Sprintf(",\n(%s)", strings.Join(values, ", "))
			}
			if _, err := writer.Write([]byte(stmt)); err != nil {
				return fmt.Errorf("failed to write SQL: %w", err)
			}
		} else {
			stmt = fmt.Sprintf("%s INTO %s (%s) VALUES (%s);\n",
				insert,
				escapeIdentifier(tableName, dialect),
				strings.Join(columns, ", "),
				strings.Join(values, ", "),
			)
			if _, err := writer.Write([]byte(stmt)); err != nil {
				return fmt.Errorf("failed to write SQL: %w", err)
			}
		}
	}
	if allInOne {
		if _, err := writer.Write([]byte(";\n")); err != nil {
			return fmt.Errorf("failed to write SQL: %w", err)
		}
	}

	return nil
}

// Escape identifier (like column and table names)
func escapeIdentifier(s string, dialect string) string {
	switch dialect {
	case "oracle":
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	case "postgres":
		return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""
	case "mssql":
		return "[" + strings.ReplaceAll(s, "[", "[]") + "]"
	case "none":
		return s
	default:
		// mysql
		return "`" + strings.ReplaceAll(s, "`", "\\`") + "`"
	}
}

// Escape values for SQL insertion (handle quotes, NULLs, etc.)
func escapeValue(s string) string {
	if strings.EqualFold(s, "NULL") {
		return "NULL"
	}

	// Double Quote string
	s = strconv.Quote(s)
	// Convert to Single Quote string
	s = s[1 : len(s)-1]
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "'", "\\'")
	return "'" + s + "'"
}
