package bbcode

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/martianzhang/tableconvert/common"
)

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("bbcode: input table is nil")
	}
	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()
	fmt.Fprintln(writer, "[table]")
	writeRow := func(tag string, cells []string) {
		fmt.Fprint(writer, "[tr]")
		for _, cell := range cells {
			fmt.Fprintf(writer, "[%s]%s[/%s]", tag, escape(cell), tag)
		}
		fmt.Fprintln(writer, "[/tr]")
	}
	writeRow("th", table.Headers)
	for _, row := range table.Rows {
		writeRow("td", row)
	}
	fmt.Fprintln(writer, "[/table]")
	return writer.Flush()
}

func escape(value string) string {
	value = strings.ReplaceAll(value, "[", "&#91;")
	return strings.ReplaceAll(value, "]", "&#93;")
}
