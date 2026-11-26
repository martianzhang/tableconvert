package common

import "fmt"

// SelectHeaderRows applies --skip-rows and --header-row to a raw string grid.
// skipRows drops physical rows before the header; headerRow is 1-based within
// the remaining rows (0 means the first remaining row).
func SelectHeaderRows(rows [][]string, skipRows, headerRow int) ([]string, [][]string, error) {
	if skipRows < 0 {
		return nil, nil, fmt.Errorf("--skip-rows must be >= 0")
	}
	if headerRow < 0 {
		return nil, nil, fmt.Errorf("--header-row must be >= 1")
	}
	if skipRows > len(rows) {
		return nil, nil, fmt.Errorf("--skip-rows=%d exceeds input rows (%d)", skipRows, len(rows))
	}
	rows = rows[skipRows:]
	if len(rows) == 0 {
		return nil, nil, nil
	}
	index := 0
	if headerRow > 0 {
		index = headerRow - 1
	}
	if index >= len(rows) {
		return nil, nil, fmt.Errorf("--header-row=%d exceeds remaining rows (%d)", headerRow, len(rows))
	}
	headers := append([]string(nil), rows[index]...)
	data := rows[index+1:]
	return headers, data, nil
}
