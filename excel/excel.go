package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"github.com/mattn/go-runewidth"
	"github.com/xuri/excelize/v2"
)

// SheetTable pairs a worksheet name with its parsed table.
type SheetTable struct {
	Name  string
	Table common.Table
}

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if cfg == nil {
		return fmt.Errorf("Unmarshal: config cannot be nil")
	}
	if table == nil {
		return fmt.Errorf("Unmarshal: output table cannot be nil")
	}
	if err := validateLegacyXLS(cfg); err != nil {
		return err
	}

	f, err := openWorkbook(cfg)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName, err := selectSheet(f, cfg)
	if err != nil {
		return err
	}
	result, err := readSheet(f, sheetName, cfg)
	if err != nil {
		return err
	}
	*table = *result
	return nil
}

// UnmarshalAll reads every worksheet. It is used by --all-sheets.
func UnmarshalAll(cfg *common.Config) ([]SheetTable, error) {
	if cfg == nil {
		return nil, fmt.Errorf("UnmarshalAll: config cannot be nil")
	}
	if err := validateLegacyXLS(cfg); err != nil {
		return nil, err
	}
	f, err := openWorkbook(cfg)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names := f.GetSheetList()
	if len(names) == 0 {
		return nil, fmt.Errorf("empty Excel file: no sheets found")
	}
	sheets := make([]SheetTable, 0, len(names))
	for _, name := range names {
		table, err := readSheet(f, name, cfg)
		if err != nil {
			return nil, fmt.Errorf("sheet %q: %w", name, err)
		}
		sheets = append(sheets, SheetTable{Name: name, Table: *table})
	}
	return sheets, nil
}

func readSheet(f *excelize.File, sheetName string, cfg *common.Config) (*common.Table, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	expanded, err := expandMergedCells(f, sheetName, rows, cfg.GetExtensionBool("expand-merged", true))
	if err != nil {
		return nil, err
	}
	rows = expanded

	table := &common.Table{}
	skipRows := cfg.GetExtensionInt("skip-rows", 0)
	headerRow := cfg.GetExtensionInt("header-row", 0)

	useFirstColAsHeader := cfg.GetExtensionBool("first-column-header", false)
	if useFirstColAsHeader {
		if headerRow > 0 {
			return nil, fmt.Errorf("--header-row cannot be combined with --first-column-header")
		}
		if skipRows < 0 || skipRows > len(rows) {
			return nil, fmt.Errorf("--skip-rows=%d exceeds input rows (%d)", skipRows, len(rows))
		}
		rows = rows[skipRows:]
		if len(rows) > 0 {
			maxDataCols := 0
			for _, row := range rows {
				if len(row) > 0 {
					table.Headers = append(table.Headers, row[0])
					if len(row)-1 > maxDataCols {
						maxDataCols = len(row) - 1
					}
				} else {
					table.Headers = append(table.Headers, "")
				}
			}

			table.Rows = make([][]string, maxDataCols)
			for colIdx := range table.Rows {
				table.Rows[colIdx] = make([]string, len(table.Headers))
				for headerIdx := range table.Headers {
					if headerIdx < len(rows) && colIdx+1 < len(rows[headerIdx]) {
						table.Rows[colIdx][headerIdx] = rows[headerIdx][colIdx+1]
					}
				}
			}
		}
		return table, nil
	}

	headers, data, err := common.SelectHeaderRows(rows, skipRows, headerRow)
	if err != nil {
		return nil, err
	}
	table.Headers = headers
	for _, row := range data {
		if len(row) < len(table.Headers) {
			row = append(row, make([]string, len(table.Headers)-len(row))...)
		}
		table.Rows = append(table.Rows, row)
	}
	return table, nil
}

// ListSheets returns the worksheet names in workbook order.
func ListSheets(cfg *common.Config) ([]string, error) {
	if cfg == nil {
		return nil, fmt.Errorf("ListSheets: config cannot be nil")
	}
	if err := validateLegacyXLS(cfg); err != nil {
		return nil, err
	}
	f, err := openWorkbook(cfg)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.GetSheetList(), nil
}

func validateLegacyXLS(cfg *common.Config) error {
	if cfg.File != "" && strings.EqualFold(filepath.Ext(cfg.File), ".xls") {
		return fmt.Errorf("legacy .xls input is not supported; save as .xlsx or export CSV")
	}
	return nil
}

func openWorkbook(cfg *common.Config) (*excelize.File, error) {
	if cfg.File != "" {
		return excelize.OpenFile(cfg.File)
	}
	if cfg.Reader != nil {
		return excelize.OpenReader(cfg.Reader)
	}
	return nil, fmt.Errorf("Excel input requires File or Reader")
}

func selectSheet(f *excelize.File, cfg *common.Config) (string, error) {
	sheetName := cfg.GetExtensionString("sheet-name", "")
	if sheetName != "" {
		index, err := f.GetSheetIndex(sheetName)
		if err != nil {
			return "", err
		}
		if index < 0 {
			return "", fmt.Errorf("Excel sheet %q not found", sheetName)
		}
		return sheetName, nil
	}

	if rawIndex := cfg.GetExtensionString("sheet-index", ""); rawIndex != "" {
		index, err := strconv.Atoi(rawIndex)
		if err != nil {
			return "", fmt.Errorf("invalid --sheet-index %q: expected a number", rawIndex)
		}
		sheets := f.GetSheetList()
		if index < 1 || index > len(sheets) {
			return "", fmt.Errorf("--sheet-index %d out of range (workbook has %d sheets)", index, len(sheets))
		}
		return sheets[index-1], nil
	}

	sheetName = f.GetSheetName(0)
	if sheetName == "" {
		return "", fmt.Errorf("empty Excel file: no sheets found")
	}
	return sheetName, nil
}

func expandMergedCells(f *excelize.File, sheetName string, rows [][]string, expand bool) ([][]string, error) {
	merges, err := f.GetMergeCells(sheetName)
	if err != nil {
		return nil, err
	}
	if len(merges) == 0 {
		return rows, nil
	}
	if !expand {
		fmt.Fprintf(os.Stderr, "warning: %d merged cell range(s) in sheet %q are not expanded; values may be empty (use --expand-merged=true)\n", len(merges), sheetName)
		return rows, nil
	}

	maxRow := len(rows)
	maxCol := 0
	for _, row := range rows {
		if len(row) > maxCol {
			maxCol = len(row)
		}
	}

	type rangeCoords struct {
		startRow, startCol, endRow, endCol int
		value                              string
	}
	ranges := make([]rangeCoords, 0, len(merges))
	for _, merge := range merges {
		startCol, startRow, err := excelize.CellNameToCoordinates(merge.GetStartAxis())
		if err != nil {
			return nil, err
		}
		endCol, endRow, err := excelize.CellNameToCoordinates(merge.GetEndAxis())
		if err != nil {
			return nil, err
		}
		value := merge.GetCellValue()
		if value != "" {
			if endRow > maxRow {
				maxRow = endRow
			}
			if endCol > maxCol {
				maxCol = endCol
			}
		}
		ranges = append(ranges, rangeCoords{startRow, startCol, endRow, endCol, value})
	}

	if maxRow == 0 || maxCol == 0 {
		return rows, nil
	}
	matrix := make([][]string, maxRow)
	for i := range matrix {
		matrix[i] = make([]string, maxCol)
	}
	for i, row := range rows {
		if i >= len(matrix) {
			break
		}
		copy(matrix[i], row)
	}
	for _, r := range ranges {
		if r.value == "" {
			continue
		}
		for row := r.startRow; row <= r.endRow && row <= maxRow; row++ {
			for col := r.startCol; col <= r.endCol && col <= maxCol; col++ {
				matrix[row-1][col-1] = r.value
			}
		}
	}
	return matrix, nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if cfg == nil {
		return fmt.Errorf("Marshal: config cannot be nil")
	}
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}
	if err := validateLegacyXLSOutput(cfg); err != nil {
		return err
	}

	f := excelize.NewFile()
	sheetName := cfg.GetExtensionString("sheet-name", "Sheet1")
	defaultSheet := f.GetSheetName(0)
	if defaultSheet == "" {
		return fmt.Errorf("failed to create default Excel sheet")
	}
	if sheetName != defaultSheet {
		if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
			return err
		}
	}
	f.SetActiveSheet(0)
	if err := writeSheet(f, sheetName, table, cfg); err != nil {
		return err
	}
	return writeWorkbook(f, cfg)
}

// MarshalAll writes one worksheet per SheetTable. It is used by --all-sheets.
func MarshalAll(cfg *common.Config, sheets []SheetTable) error {
	if cfg == nil {
		return fmt.Errorf("MarshalAll: config cannot be nil")
	}
	if len(sheets) == 0 {
		return fmt.Errorf("MarshalAll: no sheets to write")
	}
	if err := validateLegacyXLSOutput(cfg); err != nil {
		return err
	}

	f := excelize.NewFile()
	defaultSheet := f.GetSheetName(0)
	if defaultSheet == "" {
		return fmt.Errorf("failed to create default Excel sheet")
	}
	for i, sheet := range sheets {
		name := sheet.Name
		if name == "" {
			name = fmt.Sprintf("Sheet%d", i+1)
		}
		if i == 0 {
			if name != defaultSheet {
				if err := f.SetSheetName(defaultSheet, name); err != nil {
					return err
				}
			}
		} else {
			if _, err := f.NewSheet(name); err != nil {
				return fmt.Errorf("sheet %q: %w", name, err)
			}
		}
		if err := writeSheet(f, name, &sheet.Table, cfg); err != nil {
			return err
		}
	}
	f.SetActiveSheet(0)
	return writeWorkbook(f, cfg)
}

func writeSheet(f *excelize.File, sheetName string, table *common.Table, cfg *common.Config) error {
	textFormat := cfg.GetExtensionBool("text-format", true)
	textStyleID := 0
	if textFormat {
		var err error
		textStyleID, err = f.NewStyle(&excelize.Style{NumFmt: 49})
		if err != nil {
			return err
		}
	}

	for colIndex, header := range table.Headers {
		cell, err := excelize.CoordinatesToCellName(colIndex+1, 1)
		if err != nil {
			return err
		}
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return err
		}
	}
	for rowIndex, row := range table.Rows {
		for colIndex, value := range row {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				return err
			}
			if textFormat {
				if err := f.SetCellStyle(sheetName, cell, cell, textStyleID); err != nil {
					return err
				}
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return err
			}
		}
	}

	if cfg.GetExtensionBool("auto-width", false) {
		for colIndex, header := range table.Headers {
			colName, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				return err
			}
			maxWidth := runewidth.StringWidth(header)
			for _, row := range table.Rows {
				if colIndex < len(row) {
					if width := runewidth.StringWidth(row[colIndex]); width > maxWidth {
						maxWidth = width
					}
				}
			}
			width := float64(maxWidth + 2)
			if width < 8 {
				width = 8
			}
			if width > 100 {
				width = 100
			}
			if err := f.SetColWidth(sheetName, colName, colName, width); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeWorkbook(f *excelize.File, cfg *common.Config) error {
	if cfg.Writer != nil {
		if err := f.Write(cfg.Writer); err != nil {
			return fmt.Errorf("failed to write workbook: %w", err)
		}
		return nil
	}
	if cfg.Result == "" {
		return fmt.Errorf("Marshal: Excel output requires Writer or Result")
	}
	return f.SaveAs(cfg.Result)
}

func validateLegacyXLSOutput(cfg *common.Config) error {
	if cfg.Result != "" && strings.EqualFold(filepath.Ext(cfg.Result), ".xls") {
		return fmt.Errorf("legacy .xls output is not supported; use .xlsx")
	}
	return nil
}

// MergeSheets combines worksheets into one table with a leading sheet column.
func MergeSheets(sheets []SheetTable) *common.Table {
	sheetColumn := "Sheet"
	for _, sheet := range sheets {
		for _, header := range sheet.Table.Headers {
			if header == sheetColumn {
				sheetColumn = "_sheet"
				break
			}
		}
	}

	headers := []string{sheetColumn}
	index := map[string]int{sheetColumn: 0}
	for _, sheet := range sheets {
		for _, header := range sheet.Table.Headers {
			if _, exists := index[header]; !exists {
				index[header] = len(headers)
				headers = append(headers, header)
			}
		}
	}

	rows := make([][]string, 0)
	for _, sheet := range sheets {
		for _, row := range sheet.Table.Rows {
			merged := make([]string, len(headers))
			merged[0] = sheet.Name
			for i, header := range sheet.Table.Headers {
				if i < len(row) {
					merged[index[header]] = row[i]
				}
			}
			rows = append(rows, merged)
		}
	}
	return &common.Table{Headers: headers, Rows: rows}
}
