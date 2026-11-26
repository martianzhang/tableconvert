package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteEmptyRowsKeepsRaw(t *testing.T) {
	table := &Table{
		Headers: []string{"a"},
		Rows:    [][]string{{"1"}, {""}, {"3"}},
		RawRows: []interface{}{"row1", "row2", "row3"},
		Raw:     [][]interface{}{{int64(1)}, {nil}, {int64(3)}},
	}
	DeleteEmptyRows(table)
	assert.Equal(t, [][]string{{"1"}, {"3"}}, table.Rows)
	assert.Equal(t, []interface{}{"row1", "row3"}, table.RawRows)
	assert.Equal(t, [][]interface{}{{int64(1)}, {int64(3)}}, table.Raw)
}

func TestDeduplicateRowsKeepsRaw(t *testing.T) {
	table := &Table{
		Headers: []string{"a"},
		Rows:    [][]string{{"1"}, {"1"}, {"2"}},
		RawRows: []interface{}{"row1", "row2", "row3"},
		Raw:     [][]interface{}{{int64(1)}, {int64(1)}, {int64(2)}},
	}
	DeduplicateRows(table)
	assert.Equal(t, [][]string{{"1"}, {"2"}}, table.Rows)
	assert.Equal(t, []interface{}{"row1", "row3"}, table.RawRows)
	assert.Equal(t, [][]interface{}{{int64(1)}, {int64(2)}}, table.Raw)
}

func TestCaseTransformsDropRaw(t *testing.T) {
	table := &Table{
		Headers: []string{"a"},
		Rows:    [][]string{{"x"}},
		RawRows: []interface{}{"row1"},
		Raw:     [][]interface{}{{int64(1)}},
	}
	Uppercase(table)
	assert.Nil(t, table.RawRows)
	assert.Nil(t, table.Raw)
}

func TestTransposeDropsRaw(t *testing.T) {
	table := &Table{
		Headers: []string{"a"},
		Rows:    [][]string{{"1"}},
		RawRows: []interface{}{"row1"},
		Raw:     [][]interface{}{{int64(1)}},
	}
	Transpose(table)
	assert.Nil(t, table.RawRows)
	assert.Nil(t, table.Raw)
}
