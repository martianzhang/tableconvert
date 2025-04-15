package main

import (
	"fmt"
	"os"

	"github.com/martianzhang/tableconvert/ascii"
	"github.com/martianzhang/tableconvert/common"
	"github.com/martianzhang/tableconvert/csv"
	"github.com/martianzhang/tableconvert/excel"
	"github.com/martianzhang/tableconvert/json"
	"github.com/martianzhang/tableconvert/latex"
	"github.com/martianzhang/tableconvert/markdown"
	"github.com/martianzhang/tableconvert/mediawiki"
	"github.com/martianzhang/tableconvert/mysql"
	"github.com/martianzhang/tableconvert/sql"
	"github.com/martianzhang/tableconvert/twiki"
	"github.com/martianzhang/tableconvert/xml"
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
	case "json":
		err = json.Unmarshal(&cfg, &table)
	case "sql":
		err = sql.Unmarshal(&cfg, &table)
	case "xml":
		err = xml.Unmarshal(&cfg, &table)
	case "excel", "xlsx":
		err = excel.Unmarshal(&cfg, &table)
	case "twiki", "tracwiki":
		err = twiki.Unmarshal(&cfg, &table)
	case "mediawiki":
		err = mediawiki.Unmarshal(&cfg, &table)
	case "latex":
		err = latex.Unmarshal(&cfg, &table)
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
	case "json":
		err = json.Marshal(&cfg, &table)
	case "sql":
		err = sql.Marshal(&cfg, &table)
	case "xml":
		err = xml.Marshal(&cfg, &table)
	case "excel", "xlsx":
		err = excel.Marshal(&cfg, &table)
	case "twiki", "tracwiki":
		err = twiki.Marshal(&cfg, &table)
	case "mediawiki":
		err = mediawiki.Marshal(&cfg, &table)
	case "latex":
		err = latex.Marshal(&cfg, &table)
	default:
		panic("Unsupported format")
	}
	if err != nil {
		panic(err)
	}
}
