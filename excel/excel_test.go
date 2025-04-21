package excel

import (
	"os"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	// Marshal
	cfg := &common.Config{
		To:     "xlsx",
		Result: "test_output.xlsx",
	}

	table := &common.Table{
		Headers: []string{"Header1", "Header2"},
		Rows: [][]string{
			{"Data1", "Data2"},
			{"Data3", "Data4"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Unmarshal
	table2 := &common.Table{}
	cfg2 := &common.Config{
		From: "xlsx",
		File: cfg.Result,
	}
	err = Unmarshal(cfg2, table2)
	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)

	// Clean up
	_ = os.Remove(cfg.Result)
}
