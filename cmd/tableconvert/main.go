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
		fmt.Printf("# Extra Configs: %v\n", cfg.Extension)
	}

	// Reader
	var table common.Table
	switch cfg.From {
	case "markdown":
		err = markdown.Unmarshal(&cfg, &table)
	case "ascii":
		err = ascii.Unmarshal(&cfg, &table)
	case "mysql":
		err = mysql.Unmarshal(&cfg, &table)
	case "csv":
		err = csv.Unmarshal(&cfg, &table)
	default:
		panic("Unsupported format")
	}
	if err != nil {
		panic(err)
	}

	// Writer
	switch cfg.To {
	case "markdown":
		err = markdown.Marshal(&cfg, &table)
	case "ascii":
		err = ascii.Marshal(&cfg, &table)
	case "mysql":
		err = mysql.Marshal(&cfg, &table)
	case "csv":
		err = csv.Marshal(&cfg, &table)
	default:
		panic("Unsupported format")
	}
	if err != nil {
		panic(err)
	}
}
