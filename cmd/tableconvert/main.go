package main

import (
	"bytes"
	"io"
	"os"

	"github.com/martianzhang/tableconvert/common"
	"github.com/martianzhang/tableconvert/markdown"
	"github.com/martianzhang/tableconvert/mysql"
)

func main() {
	// read text from stdin pipeline
	input, err := io.ReadAll(io.TeeReader(os.Stdin, io.Discard))
	if err != nil {
		panic(err)
	}
	var table common.Table
	// 使用 bytes.NewReader 将 []byte 转换为 io.Reader
	err = mysql.Unmarshal(bytes.NewReader(input), &table)
	if err != nil {
		panic(err)
	}
	markdown.Marshal(&table, os.Stdout)
}
