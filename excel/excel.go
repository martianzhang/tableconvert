package excel

import (
	"fmt"

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
		return fmt.Errorf("empty Excel file: no sheets found")
	}

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	// Process rows
	useFirstColAsHeader := cfg.GetExtensionBool("first-column-header", false)
	if useFirstColAsHeader {
		if len(rows) > 0 {
			// Process headers - handle empty rows by using empty string
			for _, row := range rows {
				if len(row) > 0 {
					table.Headers = append(table.Headers, row[0])
				} else {
					table.Headers = append(table.Headers, "")
				}
			}

			// Process data rows
			for i := range rows {
				if len(rows[i]) > 1 {
					table.Rows = append(table.Rows, rows[i][1:])
				} else {
					// If row is empty or has only 1 column, add empty slice
					table.Rows = append(table.Rows, []string{})
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
	f := excelize.NewFile()

	// Sheet Name
	sheetName := cfg.GetExtensionString("sheet-name", "Sheet1")
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Auto-width configuration
	autoWidth := cfg.GetExtensionBool("auto-width", false)

	// Text format configuration
	textFormat := cfg.GetExtensionBool("text-format", true)

	// Transpose configuration
	transpose := cfg.GetExtensionBool("transpose", false)

	if transpose {
		// Write headers
		for colIndex := range table.Headers {
			cell, err := excelize.CoordinatesToCellName(1, colIndex+1)
			if err != nil {
				return err
			}
			if err := f.SetCellValue(sheetName, cell, table.Headers[colIndex]); err != nil {
				return err
			}
		}

		// Write data
		for rowIndex, row := range table.Rows {
			for colIndex, value := range row {
				cell, err := excelize.CoordinatesToCellName(rowIndex+2, colIndex+1)
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
	} else {
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
	}

	// Auto adjust column widths if enabled
	if autoWidth {
		if transpose {
			for colIndex := range table.Rows {
				colName, err := excelize.ColumnNumberToName(colIndex + 2)
				if err != nil {
					return err
				}
				if err := f.SetColWidth(sheetName, colName, colName, 0); err != nil {
					return err
				}
			}
		} else {
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
	}

	return f.SaveAs(cfg.Result)
}
