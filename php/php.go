package php

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/martianzhang/tableconvert/common"
)

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("php: input table is nil")
	}
	writer := bufio.NewWriter(cfg.Writer)
	defer writer.Flush()
	fmt.Fprintln(writer, "<?php")
	fmt.Fprintf(writer, "$headers = %s;\n", phpArray(table.Headers))
	fmt.Fprintln(writer, "$rows = [")
	for _, row := range table.Rows {
		fmt.Fprintf(writer, "    %s,\n", phpArray(row))
	}
	fmt.Fprintln(writer, "];")
	return writer.Flush()
}

func phpArray(values []string) string {
	result := "["
	for i, value := range values {
		if i > 0 {
			result += ", "
		}
		result += strconv.Quote(value)
	}
	return result + "]"
}
