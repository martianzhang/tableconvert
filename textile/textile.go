package textile

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("textile: input table is nil")
	}
	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()
	writeRow := func(header bool, cells []string) {
		fmt.Fprint(writer, "|")
		for _, cell := range cells {
			prefix := ""
			if header {
				prefix = "_. "
			}
			fmt.Fprintf(writer, "%s%s |", prefix, escape(cell))
		}
		fmt.Fprintln(writer)
	}
	writeRow(true, table.Headers)
	for _, row := range table.Rows {
		writeRow(false, row)
	}
	return writer.Flush()
}

func escape(value string) string {
	return strings.ReplaceAll(value, "|", `\|`)
}
