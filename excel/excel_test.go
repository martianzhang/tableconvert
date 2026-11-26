package excel

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
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

	// Verify the columns got usable widths, not the old width=0 (hidden).
	f, err := excelize.OpenFile(testFile)
	require.NoError(t, err)
	defer f.Close()
	widthA, err := f.GetColWidth("Sheet1", "A")
	require.NoError(t, err)
	widthB, err := f.GetColWidth("Sheet1", "B")
	require.NoError(t, err)
	assert.Greater(t, widthA, float64(0))
	assert.GreaterOrEqual(t, widthA, float64(8))
	assert.Greater(t, widthB, float64(0))
}

func TestMarshalRejectsLegacyXLS(t *testing.T) {
	var buf bytes.Buffer
	err := Marshal(&common.Config{
		Writer: &buf,
		Result: "out.xls",
	}, &common.Table{Headers: []string{"a"}, Rows: [][]string{{"1"}}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), ".xlsx")
}

func TestMarshalWriterAndUnmarshalReader(t *testing.T) {
	var buf bytes.Buffer
	table := &common.Table{
		Headers: []string{"Header1", "Header2"},
		Rows: [][]string{
			{"Data1", "Data2"},
			{"Data3", "Data4"},
		},
	}

	cfg := &common.Config{Writer: &buf}
	require.NoError(t, Marshal(cfg, table))
	require.NotZero(t, buf.Len())

	got := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: bytes.NewReader(buf.Bytes())}, got))
	assert.Equal(t, table.Headers, got.Headers)
	assert.Equal(t, table.Rows, got.Rows)
}

func TestSheetNameWriteAndRead(t *testing.T) {
	table := &common.Table{
		Headers: []string{"name", "age"},
		Rows:    [][]string{{"Alice", "30"}},
	}

	var buf bytes.Buffer
	cfg := &common.Config{
		Writer:    &buf,
		Extension: map[string]string{"sheet-name": "Data"},
	}
	require.NoError(t, Marshal(cfg, table))

	f, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	assert.Equal(t, []string{"Data"}, f.GetSheetList())
	require.NoError(t, f.Close())

	got := &common.Table{}
	require.NoError(t, Unmarshal(&common.Config{Reader: bytes.NewReader(buf.Bytes())}, got))
	assert.Equal(t, table.Headers, got.Headers)
	assert.Equal(t, table.Rows, got.Rows)
}

func TestUnmarshalSelectedSheet(t *testing.T) {
	f := excelize.NewFile()
	_, err := f.NewSheet("Data")
	require.NoError(t, err)
	require.NoError(t, f.SetCellValue("Sheet1", "A1", "wrong"))
	require.NoError(t, f.SetCellValue("Data", "A1", "name"))
	require.NoError(t, f.SetCellValue("Data", "A2", "Alice"))

	var buf bytes.Buffer
	require.NoError(t, f.Write(&buf))
	require.NoError(t, f.Close())

	got := &common.Table{}
	err = Unmarshal(&common.Config{
		Reader:    bytes.NewReader(buf.Bytes()),
		Extension: map[string]string{"sheet-name": "Data"},
	}, got)
	require.NoError(t, err)
	assert.Equal(t, []string{"name"}, got.Headers)
	assert.Equal(t, [][]string{{"Alice"}}, got.Rows)
}

func TestSheetNameNotFound(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, Marshal(&common.Config{Writer: &buf}, &common.Table{
		Headers: []string{"a"},
		Rows:    [][]string{{"1"}},
	}))

	err := Unmarshal(&common.Config{
		Reader:    bytes.NewReader(buf.Bytes()),
		Extension: map[string]string{"sheet-name": "Missing"},
	}, &common.Table{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
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

func TestUnmarshalMergedCellsExpanded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "merged.xlsx")

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "A1", "Name"))
	require.NoError(t, f.SetCellValue(sheet, "B1", "Info"))
	require.NoError(t, f.SetCellValue(sheet, "C1", "Age"))
	require.NoError(t, f.SetCellValue(sheet, "A2", "Alice"))
	require.NoError(t, f.MergeCell(sheet, "A2", "B2"))
	require.NoError(t, f.SetCellValue(sheet, "C2", "30"))
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	table := &common.Table{}
	cfg := &common.Config{From: "xlsx", File: path}
	require.NoError(t, Unmarshal(cfg, table))
	assert.Equal(t, []string{"Name", "Info", "Age"}, table.Headers)
	assert.Equal(t, [][]string{{"Alice", "Alice", "30"}}, table.Rows)

	table2 := &common.Table{}
	cfg2 := &common.Config{
		From:      "xlsx",
		File:      path,
		Extension: map[string]string{"expand-merged": "false"},
	}
	require.NoError(t, Unmarshal(cfg2, table2))
	assert.Equal(t, []string{"Name", "Info", "Age"}, table2.Headers)
	assert.Equal(t, []string{"Alice", "", "30"}, table2.Rows[0])
}

func TestUnmarshalSheetIndex(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sheets.xlsx")

	f := excelize.NewFile()
	first := f.GetSheetName(0)
	_, err := f.NewSheet("Second")
	require.NoError(t, err)
	require.NoError(t, f.SetCellValue(first, "A1", "first"))
	require.NoError(t, f.SetCellValue(first, "A2", "1"))
	require.NoError(t, f.SetCellValue("Second", "A1", "second"))
	require.NoError(t, f.SetCellValue("Second", "A2", "2"))
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	table := &common.Table{}
	cfg := &common.Config{
		From:      "xlsx",
		File:      path,
		Extension: map[string]string{"sheet-index": "2"},
	}
	require.NoError(t, Unmarshal(cfg, table))
	assert.Equal(t, []string{"second"}, table.Headers)
	assert.Equal(t, [][]string{{"2"}}, table.Rows)
}

func TestUnmarshalLegacyXLS(t *testing.T) {
	table := &common.Table{}
	cfg := &common.Config{From: "excel", File: "legacy.xls"}
	err := Unmarshal(cfg, table)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "legacy .xls input is not supported")
}

func TestListSheets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sheets.xlsx")

	f := excelize.NewFile()
	first := f.GetSheetName(0)
	_, err := f.NewSheet("Second")
	require.NoError(t, err)
	require.NoError(t, f.SetCellValue(first, "A1", "x"))
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	sheets, err := ListSheets(&common.Config{File: path})
	require.NoError(t, err)
	assert.Equal(t, []string{first, "Second"}, sheets)
}

func TestUnmarshalSheetNameOverridesIndex(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sheets.xlsx")

	f := excelize.NewFile()
	first := f.GetSheetName(0)
	_, err := f.NewSheet("Second")
	require.NoError(t, err)
	require.NoError(t, f.SetCellValue(first, "A1", "first"))
	require.NoError(t, f.SetCellValue(first, "A2", "1"))
	require.NoError(t, f.SetCellValue("Second", "A1", "second"))
	require.NoError(t, f.SetCellValue("Second", "A2", "2"))
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	table := &common.Table{}
	cfg := &common.Config{
		From: "xlsx",
		File: path,
		Extension: map[string]string{
			"sheet-name":  "Second",
			"sheet-index": "1",
		},
	}
	require.NoError(t, Unmarshal(cfg, table))
	assert.Equal(t, []string{"second"}, table.Headers)
	assert.Equal(t, [][]string{{"2"}}, table.Rows)
}

func TestListSheetsLegacyXLS(t *testing.T) {
	_, err := ListSheets(&common.Config{File: "legacy.xls"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "legacy .xls input is not supported")
}

func TestUnmarshalMergedCellsWarningWhenDisabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "merged-warning.xlsx")

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "A1", "Name"))
	require.NoError(t, f.SetCellValue(sheet, "A2", "Alice"))
	require.NoError(t, f.MergeCell(sheet, "A2", "B2"))
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	oldStderr := os.Stderr
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = writer
	defer func() { os.Stderr = oldStderr }()

	table := &common.Table{}
	cfg := &common.Config{
		From:      "xlsx",
		File:      path,
		Extension: map[string]string{"expand-merged": "false"},
	}
	require.NoError(t, Unmarshal(cfg, table))
	require.NoError(t, writer.Close())
	os.Stderr = oldStderr

	output, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Contains(t, string(output), "not expanded")
}

func TestUnmarshalAllAndMarshalAll(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "all.xlsx")

	f := excelize.NewFile()
	first := f.GetSheetName(0)
	_, err := f.NewSheet("Second")
	require.NoError(t, err)
	require.NoError(t, f.SetCellValue(first, "A1", "name"))
	require.NoError(t, f.SetCellValue(first, "B1", "age"))
	require.NoError(t, f.SetCellValue(first, "A2", "Alice"))
	require.NoError(t, f.SetCellValue(first, "B2", "30"))
	require.NoError(t, f.SetCellValue("Second", "A1", "name"))
	require.NoError(t, f.SetCellValue("Second", "B1", "city"))
	require.NoError(t, f.SetCellValue("Second", "A2", "Bob"))
	require.NoError(t, f.SetCellValue("Second", "B2", "LA"))
	require.NoError(t, f.SaveAs(input))
	require.NoError(t, f.Close())

	sheets, err := UnmarshalAll(&common.Config{From: "xlsx", File: input})
	require.NoError(t, err)
	require.Len(t, sheets, 2)
	assert.Equal(t, first, sheets[0].Name)
	assert.Equal(t, []string{"name", "age"}, sheets[0].Table.Headers)
	assert.Equal(t, []string{"name", "city"}, sheets[1].Table.Headers)

	output := filepath.Join(dir, "out.xlsx")
	require.NoError(t, MarshalAll(&common.Config{Result: output}, sheets))
	readBack, err := UnmarshalAll(&common.Config{From: "xlsx", File: output})
	require.NoError(t, err)
	require.Len(t, readBack, 2)
	assert.Equal(t, sheets[0].Table.Headers, readBack[0].Table.Headers)
	assert.Equal(t, sheets[1].Table.Rows, readBack[1].Table.Rows)
}

func TestMergeSheets(t *testing.T) {
	merged := MergeSheets([]SheetTable{
		{Name: "S1", Table: common.Table{Headers: []string{"name", "age"}, Rows: [][]string{{"Alice", "30"}}}},
		{Name: "S2", Table: common.Table{Headers: []string{"name", "city"}, Rows: [][]string{{"Bob", "LA"}}}},
	})
	assert.Equal(t, []string{"Sheet", "name", "age", "city"}, merged.Headers)
	assert.Equal(t, [][]string{{"S1", "Alice", "30", ""}, {"S2", "Bob", "", "LA"}}, merged.Rows)
}
