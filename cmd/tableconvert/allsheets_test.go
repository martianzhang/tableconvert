package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestPerformAllSheetsToCSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "all.xlsx")

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
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	var buf bytes.Buffer
	cfg := &common.Config{
		From:      "xlsx",
		To:        "csv",
		File:      path,
		Writer:    &buf,
		Extension: map[string]string{"all-sheets": "true"},
	}
	require.NoError(t, performAllSheets(cfg))

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Equal(t, "Sheet,name,age,city", lines[0])
	assert.Contains(t, lines, "Sheet1,Alice,30,")
	assert.Contains(t, lines, "Second,Bob,,LA")

	cfg.Extension["sheet-name"] = "Sheet1"
	assert.Error(t, performAllSheets(cfg))
}
