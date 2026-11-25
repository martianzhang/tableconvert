package excel

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"github.com/mattn/go-runewidth"
	"github.com/xuri/excelize/v2"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if cfg == nil {
		return fmt.Errorf("Unmarshal: config cannot be nil")
	}
	if table == nil {
		return fmt.Errorf("Unmarshal: output table cannot be nil")
	}

	// Prefer a real file path when available (normal CLI use), but fall back to
	// the configured Reader so batch mode and in-memory callers work.
	var f *excelize.File
	var err error
	if cfg.File != "" {
		f, err = excelize.OpenFile(cfg.File)
	} else if cfg.Reader != nil {
		f, err = excelize.OpenReader(cfg.Reader)
	} else {
		return fmt.Errorf("Unmarshal: Excel input requires File or Reader")
	}
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the requested sheet, defaulting to the first one.
	sheetName := cfg.GetExtensionString("sheet-name", "")
	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}
	if sheetName == "" {
		return fmt.Errorf("empty Excel file: no sheets found")
	}
	sheetIndex, err := f.GetSheetIndex(sheetName)
	if err != nil {
		return err
	}
	if sheetIndex < 0 {
		return fmt.Errorf("Excel sheet %q not found", sheetName)
	}

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	useFirstColAsHeader := cfg.GetExtensionBool("first-column-header", false)
	if useFirstColAsHeader {
		if len(rows) > 0 {
			// 1. Process headers and calculate max number of data columns
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

			// 2. Initialize result rows and populate data
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
	} else {
		// Process headers
		if len(rows) > 0 {
			table.Headers = rows[0]
		}

		// Process data rows
		for i := 1; i < len(rows); i++ {
			row := rows[i]
			// Ensure each row has the same number of columns as headers
			if len(row) < len(table.Headers) {
				// Fill missing cells with empty strings
				row = append(row, make([]string, len(table.Headers)-len(row))...)
			}
			table.Rows = append(table.Rows, row)
		}
	}

	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if cfg == nil {
		return fmt.Errorf("Marshal: config cannot be nil")
	}
	if table == nil {
		return fmt.Errorf("Marshal: input table pointer cannot be nil")
	}
	if cfg.Result != "" && strings.EqualFold(filepath.Ext(cfg.Result), ".xls") {
		return fmt.Errorf("legacy .xls output is not supported; use .xlsx")
	}

	f := excelize.NewFile()

	// Sheet Name: reuse/rename the default sheet so we do not leave an empty
	// "Sheet1" behind when a custom name is requested.
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

	// Auto-width configuration
	autoWidth := cfg.GetExtensionBool("auto-width", false)

	// Text format configuration
	textFormat := cfg.GetExtensionBool("text-format", true)

	// Write headers
	for colIndex, header := range table.Headers {
		cell, err := excelize.CoordinatesToCellName(colIndex+1, 1)
		if err != nil {
			return err
		}
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return err
		}
	}

	// Write data
	// Create text format style once if needed
	var textStyleID int
	if textFormat {
		// Set cell format to text - style 49 is built-in text format
		var err error
		textStyleID, err = f.NewStyle(&excelize.Style{
			NumFmt: 49, // 49 is the built-in text format
		})
		if err != nil {
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

	// Auto adjust column widths if enabled
	if autoWidth {
		for colIndex, header := range table.Headers {
			colName, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				return err
			}

			maxWidth := runewidth.StringWidth(header)
			for _, row := range table.Rows {
				if colIndex < len(row) {
					if w := runewidth.StringWidth(row[colIndex]); w > maxWidth {
						maxWidth = w
					}
				}
			}
			// Keep columns readable but bounded, like a spreadsheet's auto-fit.
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

	// Write to the configured Writer when present. This covers stdout, batch
	// mode, MCP buffers, and dry-run (io.Discard). Only fall back to Result for
	// direct API callers that did not provide a Writer.
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
