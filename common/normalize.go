package common

import (
	"fmt"
	"os"
)

// NormalizeTable makes the table rectangular and gives every header a unique,
// non-empty name. It is intentionally called after Unmarshal so every format
// gets the same behavior instead of each parser inventing its own padding.
//
// In strict mode the first inconsistency is returned as an error; otherwise
// rows are padded/extended and headers are renamed with a warning.
func NormalizeTable(table *Table, strict bool) error {
	if table == nil {
		return nil
	}

	if strict {
		return validateStrictTable(table)
	}

	width := len(table.Headers)
	if width == 0 && len(table.Rows) > 0 {
		for _, row := range table.Rows {
			if len(row) > width {
				width = len(row)
			}
		}
	}

	extended := false
	for _, row := range table.Rows {
		if len(row) > width {
			width = len(row)
			extended = true
		}
	}
	for len(table.Headers) < width {
		table.Headers = append(table.Headers, fmt.Sprintf("Column %d", len(table.Headers)+1))
	}

	for i, row := range table.Rows {
		if len(row) < width {
			row = append(row, make([]string, width-len(row))...)
			table.Rows[i] = row
		}
	}

	renamed := normalizeHeaders(table.Headers)
	normalizeRaw(table, width)
	if extended || renamed > 0 {
		// The original JSON rows no longer match the normalized header names.
		table.RawRows = nil
	}

	if extended {
		fmt.Fprintf(os.Stderr, "warning: normalized ragged rows to %d columns\n", width)
	}
	if renamed > 0 {
		fmt.Fprintf(os.Stderr, "warning: renamed %d empty or duplicate header(s)\n", renamed)
	}
	return nil
}

func validateStrictTable(table *Table) error {
	if len(table.Headers) == 0 && len(table.Rows) > 0 {
		return fmt.Errorf("table has %d data row(s) but no headers (use --strict-columns=false to generate headers)", len(table.Rows))
	}
	width := len(table.Headers)
	for i, row := range table.Rows {
		if len(row) != width {
			return fmt.Errorf("row %d has %d columns, but table has %d (use --strict-columns=false to pad or extend)", i+1, len(row), width)
		}
	}
	seen := make(map[string]bool, len(table.Headers))
	for i, header := range table.Headers {
		if header == "" {
			return fmt.Errorf("header %d is empty (use --strict-columns=false to rename it)", i+1)
		}
		if seen[header] {
			return fmt.Errorf("duplicate header %q at column %d (use --strict-columns=false to disambiguate)", header, i+1)
		}
		seen[header] = true
	}
	return nil
}

// normalizeHeaders replaces empty names and disambiguates duplicates. It
// returns the number of changed headers.
func normalizeHeaders(headers []string) int {
	used := make(map[string]bool, len(headers))
	changed := 0
	for i := range headers {
		base := headers[i]
		if base == "" {
			base = fmt.Sprintf("Column %d", i+1)
			changed++
		}
		candidate := base
		if used[candidate] {
			changed++
			for suffix := 2; ; suffix++ {
				candidate = fmt.Sprintf("%s_%d", base, suffix)
				if !used[candidate] {
					break
				}
			}
		}
		headers[i] = candidate
		used[candidate] = true
	}
	return changed
}

func normalizeRaw(table *Table, width int) {
	if table.Raw == nil {
		return
	}
	if len(table.Raw) != len(table.Rows) {
		table.Raw = nil
		return
	}
	for i, row := range table.Raw {
		if len(row) > width {
			table.Raw[i] = row[:width]
			continue
		}
		if len(row) < width {
			padded := make([]interface{}, width)
			copy(padded, row)
			for j := len(row); j < width; j++ {
				padded[j] = ""
			}
			table.Raw[i] = padded
		}
	}
}
