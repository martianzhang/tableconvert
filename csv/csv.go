package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/martianzhang/tableconvert/common"
)

func Unmarshal(reader io.Reader, table *common.Table) error {
	csvReader := csv.NewReader(reader)

	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return &common.ParseError{
			LineNumber: 0,
			Message:    "empty CSV file",
			Line:       "",
		}
	}

	table.Headers = records[0]
	table.Rows = records[1:]

	return nil
}

func Marshal(table *common.Table, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)

	// CSV header
	if err := csvWriter.Write(table.Headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// CSV records
	for _, row := range table.Rows {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
