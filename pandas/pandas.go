package pandas

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/martianzhang/tableconvert/common"
)

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("pandas: input table is nil")
	}
	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()
	fmt.Fprintln(writer, "import pandas as pd")
	fmt.Fprintln(writer, "df = pd.DataFrame([")
	for _, row := range table.Rows {
		fmt.Fprintf(writer, "    %s,\n", pyArray(row))
	}
	fmt.Fprintf(writer, "], columns=%s)\n", pyArray(table.Headers))
	return writer.Flush()
}

func pyArray(values []string) string {
	result := "["
	for i, value := range values {
		if i > 0 {
			result += ", "
		}
		result += strconv.Quote(value)
	}
	return result + "]"
}
