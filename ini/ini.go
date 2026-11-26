package ini

import (
	"fmt"
	"io"
	"strings"

	"github.com/martianzhang/tableconvert/common"

	iniLib "gopkg.in/ini.v1"
)

const (
	sectionHeader = "section"
	keyHeader     = "key"
	valueHeader   = "value"
)

func Unmarshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("ini: target table is nil")
	}
	data, err := io.ReadAll(cfg.Reader)
	if err != nil {
		return err
	}
	file, err := iniLib.Load(data)
	if err != nil {
		return fmt.Errorf("ini: %w", err)
	}
	table.Headers = []string{sectionHeader, keyHeader, valueHeader}
	table.Rows = nil
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			table.Rows = append(table.Rows, []string{section.Name(), key.Name(), key.Value()})
		}
	}
	return nil
}

func Marshal(cfg *common.Config, table *common.Table) error {
	if table == nil {
		return fmt.Errorf("ini: input table is nil")
	}
	sectionIndex, keyIndex, valueIndex := 0, 1, 2
	for i, header := range table.Headers {
		switch strings.ToLower(header) {
		case sectionHeader:
			sectionIndex = i
		case keyHeader:
			keyIndex = i
		case valueHeader:
			valueIndex = i
		}
	}
	if len(table.Headers) < 3 {
		return fmt.Errorf("ini: output requires section, key, and value columns")
	}

	file := iniLib.Empty()
	for _, row := range table.Rows {
		sectionName := ""
		key := ""
		value := ""
		if sectionIndex < len(row) {
			sectionName = row[sectionIndex]
		}
		if keyIndex < len(row) {
			key = row[keyIndex]
		}
		if valueIndex < len(row) {
			value = row[valueIndex]
		}
		if key == "" {
			continue
		}
		section, err := file.NewSection(sectionName)
		if err != nil {
			return err
		}
		if _, err := section.NewKey(key, value); err != nil {
			return err
		}
	}
	_, err := file.WriteTo(cfg.Writer)
	return err
}
