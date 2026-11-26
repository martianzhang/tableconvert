package rst

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	"github.com/mattn/go-runewidth"
)

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("rst: input table is nil")
	}
	columns := len(table.Headers)
	widths := make([]int, columns)
	for i, header := range table.Headers {
		widths[i] = runewidth.StringWidth(header)
	}
	for _, row := range table.Rows {
		for i := 0; i < columns && i < len(row); i++ {
			if width := runewidth.StringWidth(row[i]); width > widths[i] {
				widths[i] = width
			}
		}
	}

	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()
	border := func(char string) {
		fmt.Fprint(writer, "+")
		for _, width := range widths {
			fmt.Fprint(writer, strings.Repeat(char, width+2), "+")
		}
		fmt.Fprintln(writer)
	}
	writeRow := func(cells []string) {
		fmt.Fprint(writer, "|")
		for i := 0; i < columns; i++ {
			value := ""
			if i < len(cells) {
				value = cells[i]
			}
			padding := widths[i] - runewidth.StringWidth(value)
			fmt.Fprintf(writer, " %s%s |", value, strings.Repeat(" ", padding))
		}
		fmt.Fprintln(writer)
	}

	border("-")
	writeRow(table.Headers)
	border("=")
	for _, row := range table.Rows {
		writeRow(row)
	}
	border("-")
	return writer.Flush()
}
