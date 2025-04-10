package main

import (
	"fmt"
	"os"

	"github.com/martianzhang/tableconvert/ascii"
	"github.com/martianzhang/tableconvert/common"
	"github.com/martianzhang/tableconvert/csv"
	"github.com/martianzhang/tableconvert/markdown"
	"github.com/martianzhang/tableconvert/mysql"
)

func main() {
	// Parse config
	args := os.Args[1:]
	cfg, err := common.ParseConfig(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parameter parsing error: %v\n", err)
		os.Exit(1)
	}

	if cfg.Verbose {
		// Here we can process the read data according to business needs, this is just a print example
		fmt.Printf("# From: %s\n", cfg.From)
		fmt.Printf("# To: %s\n", cfg.To)
		fmt.Printf("# Extra Configs: %v\n", cfg.Others)
	}

	// Reader
	var table common.Table
	switch cfg.From {
	case "markdown":
		err = markdown.Unmarshal(cfg.Reader, &table)
	case "ascii":
		err = ascii.Unmarshal(cfg.Reader, &table)
	case "mysql":
		err = mysql.Unmarshal(cfg.Reader, &table)
	case "csv":
		err = csv.Unmarshal(cfg.Reader, &table)
	default:
		panic("Unsupported format")
	}
	if err != nil {
		panic(err)
	}

	// Writer
	switch cfg.To {
	case "markdown":
		err = markdown.Marshal(&table, cfg.Writer)
	case "ascii":
		err = ascii.Marshal(&table, cfg.Writer)
	case "mysql":
		err = mysql.Marshal(&table, cfg.Writer)
	case "csv":
		err = csv.Marshal(&table, cfg.Writer)
	default:
		panic("Unsupported format")
	}
	if err != nil {
		panic(err)
	}
}
