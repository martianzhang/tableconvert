package sql

import (
	"fmt"
	"io"
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
			if err := handleInsert(stmt, table); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported SQL statement type %T", stmt)
		}
	}

	return nil
}

// handleInsert converts INSERT statement to table rows
// This is a helper function and should only be called by Unmarshal
func handleInsert(insert *sqlparser.Insert, table *common.Table) error {
	// Check for column order mismatch
	if len(table.Headers) > 0 {
		// Verify current INSERT columns match existing headers
		if len(insert.Columns) != len(table.Headers) {
			return fmt.Errorf("column count mismatch: expected %d columns, got %d", len(table.Headers), len(insert.Columns))
		}
		// Verify column names match (order sensitive)
		for i, col := range insert.Columns {
			if col.String() != table.Headers[i] {
				return fmt.Errorf("column order mismatch at position %d: expected %s, got %s", i, table.Headers[i], col.String())
			}
		}
	} else {
		// First INSERT statement, set headers
		for _, col := range insert.Columns {
			table.Headers = append(table.Headers, col.String())
		}
	}

	rows, ok := insert.Rows.(sqlparser.Values)
	if !ok {
		// Unsupported row formats (e.g. INSERT SELECT statements)
		return fmt.Errorf("unsupported row format")
	}

	for _, row := range rows {
		var values []string
		for _, val := range row {
			switch v := val.(type) {
			case *sqlparser.Literal:
				litStr := string(v.Val)
				// Preserve escaped quotes inside literals, only trim surrounding quotes if present
				if strings.HasPrefix(litStr, "'") && strings.HasSuffix(litStr, "'") && len(litStr) > 1 {
					// Remove surrounding quotes and unescape inner quotes
					inner := litStr[1 : len(litStr)-1]
					inner = strings.ReplaceAll(inner, `\'`, `'`)
					values = append(values, inner)
				} else {
					// No surrounding quotes, use as-is
					values = append(values, litStr)
				}
			case *sqlparser.NullVal:
				values = append(values, "NULL")
			default:
				values = append(values, sqlparser.String(val)) // fallback: stringify everything else
			}
		}
		table.Rows = append(table.Rows, values)
	}
	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}

	columnCount := len(table.Headers)
	if columnCount == 0 {
		return fmt.Errorf("Marshal: table must have at least one header")
	}

	// Validate all rows have consistent column count upfront
	for _, row := range table.Rows {
		if len(row) != columnCount {
			return fmt.Errorf("Marshal: row has %d columns, but table has %d", len(row), columnCount)
		}
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
		values := make([]string, columnCount)
		for i, cell := range row {
			values[i] = common.SQLValueEscape(cell)
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
		return common.OracleIdentifierEscape(s)
	case "postgres":
		return common.PostgreSQLIdentifierEscape(s)
	case "mssql":
		return common.MssqlIdentifierEscape(s)
	case "none":
		return s
	default:
		return common.SQLIdentifierEscape(s)
	}
}
