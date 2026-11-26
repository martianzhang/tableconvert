package common

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// ApplyInputLimit enforces --max-bytes before parsing. It reads at most
// limit+1 bytes, so it works for stdin and files without buffering the whole
// input. Zero or an unset value means unlimited.
func (c *Config) ApplyInputLimit() error {
	if c == nil || c.Reader == nil {
		return nil
	}
	limit, set, err := c.optionalNonNegativeInt("max-bytes")
	if err != nil {
		return err
	}
	if !set || limit == 0 {
		return nil
	}
	if c.File != "" {
		if info, err := os.Stat(c.File); err == nil && info.Size() > int64(limit) {
			return fmt.Errorf("input size %d exceeds --max-bytes=%d", info.Size(), limit)
		}
	}

	buffer := make([]byte, limit+1)
	n, err := io.ReadFull(c.Reader, buffer)
	if err == nil {
		return fmt.Errorf("input exceeds --max-bytes=%d", limit)
	}
	if err != io.ErrUnexpectedEOF && err != io.EOF {
		return err
	}
	c.Reader = bytes.NewReader(buffer[:n])
	return nil
}

// EnforceRowLimit returns an error when the parsed table exceeds --max-rows.
// Zero or an unset value means unlimited.
func (c *Config) EnforceRowLimit(table *Table) error {
	if table == nil {
		return nil
	}
	limit, set, err := c.optionalNonNegativeInt("max-rows")
	if err != nil {
		return err
	}
	if set && limit > 0 && len(table.Rows) > limit {
		return fmt.Errorf("input has %d rows, exceeding --max-rows=%d", len(table.Rows), limit)
	}
	return nil
}
