package jira

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("jira: input table is nil")
	}
	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()
	writeRow := func(prefix string, cells []string) {
		fmt.Fprint(writer, prefix)
		for _, cell := range cells {
			fmt.Fprintf(writer, "%s|", escape(cell))
		}
		fmt.Fprintln(writer)
	}
	writeRow("||", table.Headers)
	for _, row := range table.Rows {
		writeRow("|", row)
	}
	return writer.Flush()
}

func escape(value string) string {
	return strings.ReplaceAll(value, "|", `\|`)
}
