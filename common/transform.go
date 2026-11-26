package common

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Transpose converts columns to rows and rows to columns.
// The first column of the new table contains the original headers.
// The first row of the new table contains "" and then "Row_1", "Row_2", etc.
func Transpose(table *Table) {
	if table == nil || len(table.Headers) == 0 {
		return
	}

	originalHeaders := table.Headers
	originalRows := table.Rows

	// New table dimensions:
	// - Headers: ["", "Row_1", "Row_2", ..., "Row_n"]
	// - Rows: Each original column becomes one row
	//   Row i: [originalHeader[i], originalRow[0][i], originalRow[1][i], ...]
	//
	// Rows may be ragged (for example CSV accepts variable field counts), so
	// transpose over the widest row and pad missing cells with empty strings.
	// This keeps the transformation total instead of panicking on a short row.
	newRowCount := len(originalHeaders)
	for _, row := range originalRows {
		if len(row) > newRowCount {
			newRowCount = len(row)
		}
	}
	newColumnCount := 1 + len(originalRows)

	// Build new headers
	newHeaders := make([]string, newColumnCount)
	newHeaders[0] = ""
	for i := 0; i < len(originalRows); i++ {
		newHeaders[i+1] = "Row_" + strconv.Itoa(i+1)
	}

	// Build new rows
	newRows := make([][]string, newRowCount)
	for i := 0; i < newRowCount; i++ {
		newRows[i] = make([]string, newColumnCount)
		if i < len(originalHeaders) {
			newRows[i][0] = originalHeaders[i] // First column is original header names
		}
		for j := 0; j < len(originalRows); j++ {
			if i < len(originalRows[j]) {
				newRows[i][j+1] = originalRows[j][i]
			}
		}
	}

	table.Headers = newHeaders
	table.Rows = newRows
	table.RawRows = nil
	table.Raw = nil
}

// DeleteEmptyRows removes rows that are completely empty (all cells are empty strings)
func DeleteEmptyRows(table *Table) {
	if table == nil {
		return
	}

	oldRowCount := len(table.Rows)
	nonEmptyRows := make([][]string, 0, len(table.Rows))
	kept := make([]int, 0, len(table.Rows))
	for rowIndex, row := range table.Rows {
		isEmpty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			nonEmptyRows = append(nonEmptyRows, row)
			kept = append(kept, rowIndex)
		}
	}
	table.Rows = nonEmptyRows
	filterRawRows(table, kept, oldRowCount)
}

// DeduplicateRows removes duplicate rows (rows with identical values in all columns)
func DeduplicateRows(table *Table) {
	if table == nil {
		return
	}

	oldRowCount := len(table.Rows)
	seen := make(map[string]bool)
	uniqueRows := make([][]string, 0, len(table.Rows))
	kept := make([]int, 0, len(table.Rows))

	for rowIndex, row := range table.Rows {
		// Length-prefix each cell so different rows cannot collide even when
		// cells contain the separator (including NUL).
		var key strings.Builder
		for _, cell := range row {
			key.WriteString(strconv.Itoa(len(cell)))
			key.WriteByte(':')
			key.WriteString(cell)
		}
		if !seen[key.String()] {
			seen[key.String()] = true
			uniqueRows = append(uniqueRows, row)
			kept = append(kept, rowIndex)
		}
	}
	table.Rows = uniqueRows
	filterRawRows(table, kept, oldRowCount)
}

// filterRawRows keeps the raw rows/values aligned when a transformation
// removes rows. The oldRowCount guard protects against callers that populated
// Raw with a different shape.
func filterRawRows(table *Table, kept []int, oldRowCount int) {
	if table.RawRows != nil {
		if len(table.RawRows) == oldRowCount {
			next := make([]interface{}, 0, len(kept))
			for _, index := range kept {
				next = append(next, table.RawRows[index])
			}
			table.RawRows = next
		} else {
			table.RawRows = nil
		}
	}
	if table.Raw != nil {
		if len(table.Raw) == oldRowCount {
			next := make([][]interface{}, 0, len(kept))
			for _, index := range kept {
				next = append(next, table.Raw[index])
			}
			table.Raw = next
		} else {
			table.Raw = nil
		}
	}
}

// Uppercase converts all cell values to uppercase
func Uppercase(table *Table) {
	if table == nil {
		return
	}

	for i := range table.Headers {
		table.Headers[i] = strings.ToUpper(table.Headers[i])
	}

	for i := range table.Rows {
		for j := range table.Rows[i] {
			table.Rows[i][j] = strings.ToUpper(table.Rows[i][j])
		}
	}
	table.RawRows = nil
	table.Raw = nil
}

// Lowercase converts all cell values to lowercase
func Lowercase(table *Table) {
	if table == nil {
		return
	}

	for i := range table.Headers {
		table.Headers[i] = strings.ToLower(table.Headers[i])
	}

	for i := range table.Rows {
		for j := range table.Rows[i] {
			table.Rows[i][j] = strings.ToLower(table.Rows[i][j])
		}
	}
	table.RawRows = nil
	table.Raw = nil
}

// Capitalize converts the first letter of each cell value to uppercase
func Capitalize(table *Table) {
	if table == nil {
		return
	}

	capitalizeFirst := func(s string) string {
		if s == "" {
			return s
		}
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return s
		}
		return string(unicode.ToUpper(r)) + s[size:]
	}

	for i := range table.Headers {
		table.Headers[i] = capitalizeFirst(table.Headers[i])
	}

	for i := range table.Rows {
		for j := range table.Rows[i] {
			table.Rows[i][j] = capitalizeFirst(table.Rows[i][j])
		}
	}
	table.RawRows = nil
	table.Raw = nil
}
