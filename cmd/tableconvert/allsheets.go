package main

import (
	"fmt"
	"os"

	"github.com/martianzhang/tableconvert/common"
	"github.com/martianzhang/tableconvert/excel"
)

// performAllSheets handles --all-sheets for Excel input. Excel output writes
// one worksheet per source sheet; text output merges sheets with a leading
// "Sheet" column.
func performAllSheets(cfg *common.Config) error {
	if cfg.From != "excel" && cfg.From != "xlsx" {
		return fmt.Errorf("--all-sheets requires Excel input")
	}
	if cfg.GetExtensionString("sheet-name", "") != "" || cfg.GetExtensionString("sheet-index", "") != "" {
		return fmt.Errorf("--all-sheets cannot be combined with --sheet-name or --sheet-index")
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	sheets, err := excel.UnmarshalAll(cfg)
	if err != nil {
		return err
	}
	if cfg.DryRun {
		fmt.Fprintf(os.Stderr, "=== DRY RUN MODE (all sheets) ===\n")
		for _, sheet := range sheets {
			fmt.Fprintf(os.Stderr, "Sheet %q: %d rows, %d columns\n", sheet.Name, len(sheet.Table.Rows), len(sheet.Table.Headers))
		}
		return nil
	}

	if cfg.To == "excel" || cfg.To == "xlsx" {
		if err := excel.MarshalAll(cfg, sheets); err != nil {
			return err
		}
		return cfg.CommitOutput()
	}

	table := excel.MergeSheets(sheets)
	if err := cfg.EnsureTextIO(); err != nil {
		return err
	}
	if err := common.NormalizeTable(table, cfg.GetExtensionBool("strict-columns", false)); err != nil {
		return &common.ConversionError{Stage: "normalize", Format: cfg.From, Err: err}
	}
	if err := cfg.ApplyTransformations(table); err != nil {
		return &common.ConversionError{Stage: "transform", Format: cfg.From, Err: err}
	}

	marshalFn, ok := formatRegistry.GetMarshalFunc(cfg.To)
	if !ok {
		return fmt.Errorf("unsupported `--to` format: %s", cfg.To)
	}
	if marshalFn == nil {
		return fmt.Errorf("format %s does not support writing (marshal)", cfg.To)
	}
	if err := marshalFn(cfg, table); err != nil {
		return &common.ConversionError{Stage: "marshal", Format: cfg.To, Err: err}
	}
	return cfg.CommitOutput()
}
