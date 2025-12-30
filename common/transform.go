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

	newRowCount := len(originalHeaders)
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
		newRows[i][0] = originalHeaders[i] // First column is original header names
		for j := 0; j < len(originalRows); j++ {
			newRows[i][j+1] = originalRows[j][i]
		}
	}

	table.Headers = newHeaders
	table.Rows = newRows
}

// DeleteEmptyRows removes rows that are completely empty (all cells are empty strings)
func DeleteEmptyRows(table *Table) {
	if table == nil {
		return
	}

	nonEmptyRows := make([][]string, 0, len(table.Rows))
	for _, row := range table.Rows {
		isEmpty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			nonEmptyRows = append(nonEmptyRows, row)
		}
	}
	table.Rows = nonEmptyRows
}

// DeduplicateRows removes duplicate rows (rows with identical values in all columns)
func DeduplicateRows(table *Table) {
	if table == nil {
		return
	}

	seen := make(map[string]bool)
	uniqueRows := make([][]string, 0, len(table.Rows))

	for _, row := range table.Rows {
		// Create a key by joining all cells with a delimiter
		// Using a delimiter that's unlikely to appear in normal data
		key := strings.Join(row, "\x00")
		if !seen[key] {
			seen[key] = true
			uniqueRows = append(uniqueRows, row)
		}
	}
	table.Rows = uniqueRows
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
}
