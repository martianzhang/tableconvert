package common

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ReplaceCells applies one or more regexp replacements to cell values. The
// optional column name/index limits replacements to one column.
func ReplaceCells(table *Table, specs []string, column string) error {
	if table == nil || len(specs) == 0 {
		return nil
	}
	type replacement struct {
		re          *regexp.Regexp
		replacement string
	}
	replacements := make([]replacement, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "=>", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --replace %q (use regex=>replacement)", spec)
		}
		re, err := regexp.Compile(parts[0])
		if err != nil {
			return fmt.Errorf("invalid --replace regexp %q: %w", parts[0], err)
		}
		replacements = append(replacements, replacement{re: re, replacement: parts[1]})
	}

	columnIndex := -1
	if strings.TrimSpace(column) != "" {
		index, err := resolveColumn(table.Headers, column)
		if err != nil {
			return err
		}
		columnIndex = index
	}

	for i := range table.Rows {
		for j := range table.Rows[i] {
			if columnIndex >= 0 && j != columnIndex {
				continue
			}
			value := table.Rows[i][j]
			for _, replacement := range replacements {
				value = replacement.re.ReplaceAllString(value, replacement.replacement)
			}
			table.Rows[i][j] = value
		}
	}
	table.RawRows = nil
	table.Raw = nil
	return nil
}

// FilterRows keeps rows whose selected columns match every filter. Filters use
// col=regexp, where col is a header name or 1-based index.
func FilterRows(table *Table, specs []string) error {
	if table == nil || len(specs) == 0 {
		return nil
	}
	type filter struct {
		index int
		re    *regexp.Regexp
	}
	filters := make([]filter, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --filter %q (use column=regexp)", spec)
		}
		index, err := resolveColumn(table.Headers, parts[0])
		if err != nil {
			return err
		}
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return fmt.Errorf("invalid --filter regexp %q: %w", parts[1], err)
		}
		filters = append(filters, filter{index: index, re: re})
	}

	oldRowCount := len(table.Rows)
	kept := make([]int, 0, len(table.Rows))
	rows := make([][]string, 0, len(table.Rows))
	for i, row := range table.Rows {
		match := true
		for _, filter := range filters {
			value := ""
			if filter.index < len(row) {
				value = row[filter.index]
			}
			if !filter.re.MatchString(value) {
				match = false
				break
			}
		}
		if match {
			kept = append(kept, i)
			rows = append(rows, row)
		}
	}
	table.Rows = rows
	filterRawRows(table, kept, oldRowCount)
	return nil
}

// SelectColumns keeps and reorders columns using names, 1-based indexes, or
// ranges such as 1,3-5,name.
func SelectColumns(table *Table, spec string) error {
	if table == nil || strings.TrimSpace(spec) == "" {
		return nil
	}
	indices := make([]int, 0)
	seen := make(map[int]bool)
	add := func(index int) error {
		if index < 0 || index >= len(table.Headers) {
			return fmt.Errorf("column index %d out of range (1-%d)", index+1, len(table.Headers))
		}
		if seen[index] {
			return fmt.Errorf("duplicate column %q in --columns", table.Headers[index])
		}
		seen[index] = true
		indices = append(indices, index)
		return nil
	}

	for _, token := range splitSpec(spec) {
		if match := regexp.MustCompile(`^(\d+)-(\d+)$`).FindStringSubmatch(token); match != nil {
			start, _ := strconv.Atoi(match[1])
			end, _ := strconv.Atoi(match[2])
			if start > end {
				start, end = end, start
			}
			for i := start; i <= end; i++ {
				if err := add(i - 1); err != nil {
					return err
				}
			}
			continue
		}
		index, err := resolveColumn(table.Headers, token)
		if err != nil {
			return err
		}
		if err := add(index); err != nil {
			return err
		}
	}

	headers := make([]string, len(indices))
	for i, index := range indices {
		headers[i] = table.Headers[index]
	}
	rows := make([][]string, len(table.Rows))
	for i, row := range table.Rows {
		rows[i] = make([]string, len(indices))
		for j, index := range indices {
			if index < len(row) {
				rows[i][j] = row[index]
			}
		}
	}
	var raw [][]interface{}
	if table.Raw != nil && len(table.Raw) == len(table.Rows) {
		raw = make([][]interface{}, len(table.Raw))
		for i, row := range table.Raw {
			raw[i] = make([]interface{}, len(indices))
			for j, index := range indices {
				if index < len(row) {
					raw[i][j] = row[index]
				}
			}
		}
	}
	table.Headers = headers
	table.Rows = rows
	table.RawRows = nil
	table.Raw = raw
	return nil
}

// SortRows sorts rows by one or more columns: --sort=name:asc,2:desc.
func SortRows(table *Table, spec string) error {
	if table == nil || strings.TrimSpace(spec) == "" {
		return nil
	}
	type sortKey struct {
		index   int
		desc    bool
		numeric bool
	}
	keys := make([]sortKey, 0)
	for _, token := range splitSpec(spec) {
		name := token
		desc := false
		if index := strings.LastIndex(token, ":"); index >= 0 {
			switch strings.ToLower(token[index+1:]) {
			case "asc":
				name = token[:index]
			case "desc":
				name = token[:index]
				desc = true
			}
		}
		column, err := resolveColumn(table.Headers, name)
		if err != nil {
			return err
		}
		keys = append(keys, sortKey{index: column, desc: desc})
	}
	for i := range keys {
		keys[i].numeric = columnIsNumeric(table.Rows, keys[i].index)
	}

	order := make([]int, len(table.Rows))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		a := table.Rows[order[i]]
		b := table.Rows[order[j]]
		for _, key := range keys {
			av, bv := "", ""
			if key.index < len(a) {
				av = a[key.index]
			}
			if key.index < len(b) {
				bv = b[key.index]
			}
			if av == bv {
				continue
			}
			var less bool
			if key.numeric {
				af, _ := strconv.ParseFloat(strings.TrimSpace(av), 64)
				bf, _ := strconv.ParseFloat(strings.TrimSpace(bv), 64)
				if af == bf {
					continue
				}
				less = af < bf
			} else {
				less = av < bv
			}
			if key.desc {
				less = !less
			}
			return less
		}
		return false
	})
	permuteRows(table, order)
	return nil
}

// LimitRows applies --head and --tail. Head is applied first; tail then keeps
// the last N rows of that result.
func LimitRows(table *Table, head, tail int) {
	if table == nil || (head < 0 && tail < 0) {
		return
	}
	oldRowCount := len(table.Rows)
	start, end := 0, len(table.Rows)
	if head >= 0 && head < end {
		end = head
	}
	if tail >= 0 {
		start = end - tail
		if start < 0 {
			start = 0
		}
	}
	indices := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		indices = append(indices, i)
	}
	table.Rows = table.Rows[start:end]
	filterRawRows(table, indices, oldRowCount)
}

// RenameHeaders renames columns using old=new pairs. Old can be a header name
// or a 1-based index.
func RenameHeaders(table *Table, spec string) error {
	if table == nil || strings.TrimSpace(spec) == "" {
		return nil
	}
	seen := make(map[string]bool, len(table.Headers))
	for _, header := range table.Headers {
		seen[header] = true
	}
	for _, pair := range splitSpec(spec) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --rename %q (use old=new)", pair)
		}
		index, err := resolveColumn(table.Headers, parts[0])
		if err != nil {
			return err
		}
		newName := strings.TrimSpace(parts[1])
		if newName == "" {
			return fmt.Errorf("invalid --rename %q: new name is empty", pair)
		}
		oldName := table.Headers[index]
		if oldName == newName {
			continue
		}
		if seen[newName] {
			return fmt.Errorf("--rename target %q already exists", newName)
		}
		delete(seen, oldName)
		seen[newName] = true
		table.Headers[index] = newName
	}
	table.RawRows = nil
	return nil
}

// TrimCells trims leading/trailing whitespace from headers and values.
func TrimCells(table *Table) {
	if table == nil {
		return
	}
	for i := range table.Headers {
		table.Headers[i] = strings.TrimSpace(table.Headers[i])
	}
	for i := range table.Rows {
		for j := range table.Rows[i] {
			table.Rows[i][j] = strings.TrimSpace(table.Rows[i][j])
		}
	}
	table.RawRows = nil
	table.Raw = nil
}

func (c *Config) optionalNonNegativeInt(key string) (int, bool, error) {
	raw := c.GetExtensionString(key, "")
	if strings.TrimSpace(raw) == "" {
		return 0, false, nil
	}
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value < 0 {
		return 0, false, fmt.Errorf("--%s must be a non-negative integer", key)
	}
	return value, true, nil
}

func splitSpec(spec string) []string {
	parts := strings.Split(spec, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func resolveColumn(headers []string, token string) (int, error) {
	token = strings.TrimSpace(token)
	if index, err := strconv.Atoi(token); err == nil {
		if index < 1 || index > len(headers) {
			return -1, fmt.Errorf("column index %d out of range (1-%d)", index, len(headers))
		}
		return index - 1, nil
	}
	for i, header := range headers {
		if header == token {
			return i, nil
		}
	}
	return -1, fmt.Errorf("unknown column %q", token)
}

func columnIsNumeric(rows [][]string, index int) bool {
	found := false
	for _, row := range rows {
		if index >= len(row) {
			continue
		}
		value := strings.TrimSpace(row[index])
		if value == "" {
			continue
		}
		found = true
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return false
		}
	}
	return found
}

func permuteRows(table *Table, order []int) {
	rows := make([][]string, len(order))
	var rawRows []interface{}
	if table.RawRows != nil && len(table.RawRows) == len(order) {
		rawRows = make([]interface{}, len(order))
	}
	var raw [][]interface{}
	if table.Raw != nil && len(table.Raw) == len(order) {
		raw = make([][]interface{}, len(order))
	}
	for i, index := range order {
		rows[i] = table.Rows[index]
		if rawRows != nil {
			rawRows[i] = table.RawRows[index]
		}
		if raw != nil {
			raw[i] = table.Raw[index]
		}
	}
	table.Rows = rows
	table.RawRows = rawRows
	table.Raw = raw
}
