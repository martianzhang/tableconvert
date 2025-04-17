package excel

import (
	"github.com/martianzhang/tableconvert/common"

	"github.com/xuri/excelize/v2"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	// Open Excel file
	f, err := excelize.OpenFile(cfg.File)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the first sheet name
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil // Empty file
	}

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

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

	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	f := excelize.NewFile()

	// Sheet Name
	var sheetName string
	if v, ok := cfg.Extension["sheet-name"]; ok {
		sheetName = v
	} else {
		sheetName = "Sheet1"
	}

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Auto-width configuration
	autoWidth := false
	if v, ok := cfg.Extension["auto-width"]; ok && v == "true" {
		autoWidth = true
	}

	// Text format configuration
	textFormat := false
	if v, ok := cfg.Extension["text-format"]; ok && v == "true" {
		textFormat = true
	}

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
	for rowIndex, row := range table.Rows {
		for colIndex, value := range row {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				return err
			}

			if textFormat {
				// Set cell format to text
				style, err := f.NewStyle(&excelize.Style{
					NumFmt: 49, // 49 is the built-in text format
				})
				if err != nil {
					return err
				}
				if err := f.SetCellStyle(sheetName, cell, cell, style); err != nil {
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
		for colIndex := range table.Headers {
			colName, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				return err
			}
			if err := f.SetColWidth(sheetName, colName, colName, 0); err != nil {
				return err
			}
		}
	}

	return f.SaveAs(cfg.Result)
}
