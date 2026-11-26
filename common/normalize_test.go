package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeTableRaggedRows(t *testing.T) {
	table := &Table{
		Headers: []string{"a", "b"},
		Rows:    [][]string{{"1"}, {"2", "3", "4"}},
	}
	require.NoError(t, NormalizeTable(table, false))
	assert.Equal(t, []string{"a", "b", "Column 3"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "", ""}, {"2", "3", "4"}}, table.Rows)
}

func TestNormalizeTableHeaders(t *testing.T) {
	table := &Table{
		Headers: []string{"a", "a", ""},
		Rows:    [][]string{{"1", "2", "3"}},
	}
	require.NoError(t, NormalizeTable(table, false))
	assert.Equal(t, []string{"a", "a_2", "Column 3"}, table.Headers)
}

func TestNormalizeTableNoHeaders(t *testing.T) {
	table := &Table{
		Rows: [][]string{{"1", "2"}, {"3", "4"}},
	}
	require.NoError(t, NormalizeTable(table, false))
	assert.Equal(t, []string{"Column 1", "Column 2"}, table.Headers)
	assert.Equal(t, [][]string{{"1", "2"}, {"3", "4"}}, table.Rows)
}

func TestNormalizeTableRaw(t *testing.T) {
	table := &Table{
		Headers: []string{"a", "b"},
		Rows:    [][]string{{"1"}, {"2", "3"}},
		Raw:     [][]interface{}{{int64(1)}, {int64(2), int64(3)}},
	}
	require.NoError(t, NormalizeTable(table, false))
	assert.Equal(t, [][]interface{}{{int64(1), ""}, {int64(2), int64(3)}}, table.Raw)
}

func TestNormalizeTableStrict(t *testing.T) {
	tests := []struct {
		name  string
		table *Table
	}{
		{"ragged row", &Table{Headers: []string{"a"}, Rows: [][]string{{"1", "2"}}}},
		{"no headers", &Table{Rows: [][]string{{"1"}}}},
		{"empty header", &Table{Headers: []string{"a", ""}, Rows: [][]string{{"1", "2"}}}},
		{"duplicate header", &Table{Headers: []string{"a", "a"}, Rows: [][]string{{"1", "2"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, NormalizeTable(tt.table, true))
		})
	}
}

func TestNormalizeTableNil(t *testing.T) {
	require.NoError(t, NormalizeTable(nil, false))
}

func TestNormalizeTableDropsRawRowsWhenHeadersChange(t *testing.T) {
	table := &Table{
		Headers: []string{"a", "a"},
		Rows:    [][]string{{"1", "2"}},
		RawRows: []interface{}{"row"},
		Raw:     [][]interface{}{{int64(1), int64(2)}},
	}
	require.NoError(t, NormalizeTable(table, false))
	assert.Nil(t, table.RawRows)
	assert.Equal(t, []string{"a", "a_2"}, table.Headers)
	assert.Equal(t, [][]interface{}{{int64(1), int64(2)}}, table.Raw)
}
