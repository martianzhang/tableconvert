package excel

import (
	"os"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	// Marshal
	cfg := &common.Config{
		To:     "xlsx",
		Result: "test_output.xlsx",
	}

	table := &common.Table{
		Headers: []string{"Header1", "Header2"},
		Rows: [][]string{
			{"Data1", "Data2"},
			{"Data3", "Data4"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Unmarshal
	table2 := &common.Table{}
	cfg2 := &common.Config{
		From: "xlsx",
		File: cfg.Result,
	}
	err = Unmarshal(cfg2, table2)
	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)

	// Clean up
	_ = os.Remove(cfg.Result)
}

func TestUnmarshalWithFirstColumnHeader(t *testing.T) {
	// Create test file
	testFile := "test_first_col.xlsx"
	defer os.Remove(testFile)

	// First create a file with first column as headers
	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
	}

	table := &common.Table{
		Headers: []string{"", "Name", "Age", "City"},
		Rows: [][]string{
			{"Row1", "Alice", "30", "NYC"},
			{"Row2", "Bob", "25", "LA"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Now read it back with first-column-header option
	cfg2 := &common.Config{
		From: "xlsx",
		File: testFile,
		Extension: map[string]string{
			"first-column-header": "true",
		},
	}

	table2 := &common.Table{}
	err = Unmarshal(cfg2, table2)
	assert.NoError(t, err)

	// With first-column-header, the first column becomes headers
	// and the remaining columns become data
	expectedHeaders := []string{"", "Row1", "Row2"}
	expectedRows := [][]string{
		{"Name", "Alice", "Bob"},
		{"Age", "30", "25"},
		{"City", "NYC", "LA"},
	}

	assert.Equal(t, expectedHeaders, table2.Headers)
	assert.Equal(t, expectedRows, table2.Rows)
}

func TestUnmarshalEmptyFile(t *testing.T) {
	// Create an empty Excel file
	testFile := "test_empty.xlsx"
	defer os.Remove(testFile)

	// Create empty file
	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
	}

	table := &common.Table{
		Headers: []string{},
		Rows:    [][]string{},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Try to read it back
	cfg2 := &common.Config{
		From: "xlsx",
		File: testFile,
	}

	table2 := &common.Table{}
	err = Unmarshal(cfg2, table2)
	assert.NoError(t, err)

	// Should have nil headers and rows (or empty)
	// Excel files with no data return nil slices
	assert.True(t, table2.Headers == nil || len(table2.Headers) == 0)
	assert.True(t, table2.Rows == nil || len(table2.Rows) == 0)
}

func TestUnmarshalFileWithMissingCells(t *testing.T) {
	// Create test file with variable row lengths
	testFile := "test_missing_cells.xlsx"
	defer os.Remove(testFile)

	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
	}

	table := &common.Table{
		Headers: []string{"A", "B", "C"},
		Rows: [][]string{
			{"1", "2", "3"},
			{"4", "5"}, // Missing third cell
			{"6", "7", "8"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Read it back
	cfg2 := &common.Config{
		From: "xlsx",
		File: testFile,
	}

	table2 := &common.Table{}
	err = Unmarshal(cfg2, table2)
	assert.NoError(t, err)

	// Should fill missing cells with empty strings
	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, [][]string{
		{"1", "2", "3"},
		{"4", "5", ""},
		{"6", "7", "8"},
	}, table2.Rows)
}

func TestMarshalWithAutoWidth(t *testing.T) {
	testFile := "test_auto_width.xlsx"
	defer os.Remove(testFile)

	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
		Extension: map[string]string{
			"auto-width": "true",
		},
	}

	table := &common.Table{
		Headers: []string{"VeryLongHeaderName", "Short"},
		Rows: [][]string{
			{"Short", "A"},
			{"VeryLongDataValue", "B"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(testFile)
	assert.NoError(t, err)
}

func TestMarshalWithTextFormat(t *testing.T) {
	testFile := "test_text_format.xlsx"
	defer os.Remove(testFile)

	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
		Extension: map[string]string{
			"text-format": "true",
		},
	}

	table := &common.Table{
		Headers: []string{"Number", "String"},
		Rows: [][]string{
			{"12345678901234567890", "text"}, // Long number that should be text
			{"00123", "text2"},               // Leading zeros
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(testFile)
	assert.NoError(t, err)
}

func TestMarshalWithTextFormatDisabled(t *testing.T) {
	testFile := "test_no_text_format.xlsx"
	defer os.Remove(testFile)

	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
		Extension: map[string]string{
			"text-format": "false",
		},
	}

	table := &common.Table{
		Headers: []string{"Number"},
		Rows: [][]string{
			{"123"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(testFile)
	assert.NoError(t, err)
}

func TestMarshalWithCustomSheetName(t *testing.T) {
	testFile := "test_sheet_name.xlsx"
	defer os.Remove(testFile)

	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
		Extension: map[string]string{
			"sheet-name": "CustomSheet",
		},
	}

	table := &common.Table{
		Headers: []string{"A"},
		Rows: [][]string{
			{"1"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(testFile)
	assert.NoError(t, err)
}

func TestUnmarshalWithEmptyRows(t *testing.T) {
	testFile := "test_empty_rows.xlsx"
	defer os.Remove(testFile)

	// Create file with empty rows
	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
	}

	table := &common.Table{
		Headers: []string{"A", "B"},
		Rows: [][]string{
			{"1", "2"},
			{"", ""},
			{"3", "4"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Read it back
	cfg2 := &common.Config{
		From: "xlsx",
		File: testFile,
	}

	table2 := &common.Table{}
	err = Unmarshal(cfg2, table2)
	assert.NoError(t, err)

	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)
}

func TestUnmarshalFileWithoutExtension(t *testing.T) {
	// Test error handling for file without extension
	cfg := &common.Config{
		From: "xlsx",
		File: "nonexistent_file",
	}

	table := &common.Table{}
	err := Unmarshal(cfg, table)
	assert.Error(t, err)
}

func TestMarshalWithEmptyTable(t *testing.T) {
	testFile := "test_empty_table.xlsx"
	defer os.Remove(testFile)

	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
	}

	table := &common.Table{
		Headers: []string{},
		Rows:    [][]string{},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(testFile)
	assert.NoError(t, err)
}

func TestUnmarshalWithVariableColumnCount(t *testing.T) {
	testFile := "test_variable_cols.xlsx"
	defer os.Remove(testFile)

	// Create a file where some rows have different column counts
	cfg := &common.Config{
		To:     "xlsx",
		Result: testFile,
	}

	// Use Excel directly to create variable columns
	// Since our Marshal function pads rows, we need to test the Unmarshal
	// behavior when reading files with variable row lengths

	// First, create a normal file
	table := &common.Table{
		Headers: []string{"A", "B", "C"},
		Rows: [][]string{
			{"1", "2", "3"},
			{"4", "5", "6"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Read it back
	cfg2 := &common.Config{
		From: "xlsx",
		File: testFile,
	}

	table2 := &common.Table{}
	err = Unmarshal(cfg2, table2)
	assert.NoError(t, err)

	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)
}
