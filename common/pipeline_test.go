package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pipelineTable() *Table {
	return &Table{
		Headers: []string{"name", "age", "city", "score"},
		Rows: [][]string{
			{"Alice", "30", "NYC", "9.5"},
			{"Bob", "25", "LA", "10"},
			{"Cara", "30", "SF", "8"},
		},
	}
}

func TestSortRows(t *testing.T) {
	table := pipelineTable()
	require.NoError(t, SortRows(table, "age:desc,city:asc"))
	assert.Equal(t, [][]string{
		{"Alice", "30", "NYC", "9.5"},
		{"Cara", "30", "SF", "8"},
		{"Bob", "25", "LA", "10"},
	}, table.Rows)
}

func TestSortRowsNumeric(t *testing.T) {
	table := pipelineTable()
	require.NoError(t, SortRows(table, "score:asc"))
	assert.Equal(t, [][]string{
		{"Cara", "30", "SF", "8"},
		{"Alice", "30", "NYC", "9.5"},
		{"Bob", "25", "LA", "10"},
	}, table.Rows)
}

func TestReplaceCells(t *testing.T) {
	table := pipelineTable()
	require.NoError(t, ReplaceCells(table, []string{`^A=>X`, `0$=>zero`}, ""))
	assert.Equal(t, "Xlice", table.Rows[0][0])
	assert.Equal(t, "NYC", table.Rows[0][2])
	require.NoError(t, ReplaceCells(table, []string{`NYC=>New York`}, "city"))
	assert.Equal(t, "New York", table.Rows[0][2])
}

func TestFilterRows(t *testing.T) {
	table := pipelineTable()
	require.NoError(t, FilterRows(table, []string{"age=^30$", "city=^N"}))
	assert.Equal(t, [][]string{{"Alice", "30", "NYC", "9.5"}}, table.Rows)

	raw := &Table{
		Headers: []string{"a"},
		Rows:    [][]string{{"keep"}, {"drop"}},
		RawRows: []interface{}{"row1", "row2"},
		Raw:     [][]interface{}{{int64(1)}, {int64(2)}},
	}
	require.NoError(t, FilterRows(raw, []string{"a=keep"}))
	assert.Equal(t, []interface{}{"row1"}, raw.RawRows)
	assert.Equal(t, [][]interface{}{{int64(1)}}, raw.Raw)
}

func TestSelectColumns(t *testing.T) {
	table := pipelineTable()
	require.NoError(t, SelectColumns(table, "city,1-2,4"))
	assert.Equal(t, []string{"city", "name", "age", "score"}, table.Headers)
	assert.Equal(t, []string{"NYC", "Alice", "30", "9.5"}, table.Rows[0])
}

func TestRenameHeaders(t *testing.T) {
	table := pipelineTable()
	require.NoError(t, RenameHeaders(table, "name=Name,2=Age"))
	assert.Equal(t, []string{"Name", "Age", "city", "score"}, table.Headers)
	assert.Error(t, RenameHeaders(table, "Name=Age"))
}

func TestLimitRows(t *testing.T) {
	table := pipelineTable()
	LimitRows(table, 2, -1)
	assert.Len(t, table.Rows, 2)
	assert.Equal(t, "Alice", table.Rows[0][0])

	table = pipelineTable()
	LimitRows(table, -1, 1)
	assert.Equal(t, [][]string{{"Cara", "30", "SF", "8"}}, table.Rows)

	table = pipelineTable()
	LimitRows(table, 2, 1)
	assert.Equal(t, [][]string{{"Bob", "25", "LA", "10"}}, table.Rows)
}

func TestTrimCells(t *testing.T) {
	table := &Table{Headers: []string{" a "}, Rows: [][]string{{" x "}}}
	TrimCells(table)
	assert.Equal(t, []string{"a"}, table.Headers)
	assert.Equal(t, [][]string{{"x"}}, table.Rows)
}

func TestSortRowsKeepsRaw(t *testing.T) {
	table := &Table{
		Headers: []string{"a"},
		Rows:    [][]string{{"2"}, {"1"}},
		RawRows: []interface{}{"row2", "row1"},
		Raw:     [][]interface{}{{int64(2)}, {int64(1)}},
	}
	require.NoError(t, SortRows(table, "a:asc"))
	assert.Equal(t, []interface{}{"row1", "row2"}, table.RawRows)
	assert.Equal(t, [][]interface{}{{int64(1)}, {int64(2)}}, table.Raw)
}

func TestApplyTransformationsPipeline(t *testing.T) {
	cfg := &Config{Extension: map[string]string{
		"sort":    "age:desc",
		"head":    "1",
		"rename":  "name=Name",
		"columns": "Name,age",
	}}
	table := pipelineTable()
	require.NoError(t, cfg.ApplyTransformations(table))
	assert.Equal(t, []string{"Name", "age"}, table.Headers)
	assert.Equal(t, [][]string{{"Alice", "30"}}, table.Rows)

	bad := &Config{Extension: map[string]string{"replace": "["}}
	assert.Error(t, bad.ApplyTransformations(pipelineTable()))
}

func TestGetExtensionList(t *testing.T) {
	cfg := &Config{Extension: map[string]string{"replace": "a=>b\x00c=>d"}}
	assert.Equal(t, []string{"a=>b", "c=>d"}, cfg.GetExtensionList("replace"))
}
