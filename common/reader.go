package common

import (
	"bufio"
	"bytes"
	"io"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// StripUTF8BOM returns a reader with a leading UTF-8 BOM removed. It is safe
// for readers that do not start with a BOM.
func StripUTF8BOM(r io.Reader) io.Reader {
	br := bufio.NewReader(r)
	if b, err := br.Peek(3); err == nil && bytes.Equal(b, utf8BOM) {
		_, _ = br.Discard(3)
	}
	return br
}

// TrimUTF8BOM removes a leading UTF-8 BOM from a byte slice.
func TrimUTF8BOM(data []byte) []byte {
	return bytes.TrimPrefix(data, utf8BOM)
}
