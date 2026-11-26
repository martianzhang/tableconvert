package common

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var textFormats = map[string]bool{
	"ascii":     true,
	"csv":       true,
	"tsv":       true,
	"html":      true,
	"json":      true,
	"jsonl":     true,
	"jsonlines": true,
	"latex":     true,
	"markdown":  true,
	"md":        true,
	"mediawiki": true,
	"mysql":     true,
	"sql":       true,
	"tmpl":      true,
	"template":  true,
	"tracwiki":  true,
	"twiki":     true,
	"xml":       true,
	"yaml":      true,
	"toml":      true,
	"ini":       true,
	"bbcode":    true,
	"jira":      true,
	"asciidoc":  true,
	"rst":       true,
	"textile":   true,
	"php":       true,
	"ruby":      true,
	"pandas":    true,
}

// IsTextFormat reports whether format is decoded/encoded as text. Excel is the
// only format in the current registry that is binary.
func IsTextFormat(format string) bool {
	return textFormats[strings.ToLower(format)]
}

// DecodeReader wraps r with the requested input encoding. The default is
// UTF-8 with an optional BOM. "auto" recognizes UTF-8/UTF-16 BOMs and otherwise
// requires explicit --encoding to avoid guessing between GBK/Big5.
func DecodeReader(r io.Reader, name string) (io.Reader, error) {
	if r == nil {
		return nil, fmt.Errorf("encoding: reader is nil")
	}
	name = normalizeEncodingName(name)
	switch name {
	case "", "utf-8", "utf8", "utf-8-bom":
		return StripUTF8BOM(r), nil

	case "auto":
		return decodeAuto(r)

	case "utf-16", "utf-16le", "utf-16be":
		return decodeUTF16(r, name)

	case "gbk":
		return transform.NewReader(r, simplifiedchinese.GBK.NewDecoder()), nil
	case "gb18030":
		return transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder()), nil
	case "big5":
		return transform.NewReader(r, traditionalchinese.Big5.NewDecoder()), nil
	case "shift_jis", "shift-jis", "sjis":
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder()), nil
	case "latin1", "iso-8859-1":
		return transform.NewReader(r, charmap.ISO8859_1.NewDecoder()), nil
	default:
		return nil, fmt.Errorf("unsupported --encoding %q (use utf-8, utf-16le, utf-16be, gbk, gb18030, big5, shift_jis, latin1, or auto)", name)
	}
}

// EncodeWriter wraps w with the requested output encoding. Closing the
// returned writer flushes the transformer; it does not close the underlying
// writer.
func EncodeWriter(w io.Writer, name string) (io.WriteCloser, error) {
	if w == nil {
		return nil, fmt.Errorf("encoding: writer is nil")
	}
	name = normalizeEncodingName(name)
	var encoder *encoding.Encoder
	switch name {
	case "", "utf-8", "utf8", "utf-8-bom":
		return nopWriteCloser{Writer: w}, nil
	case "utf-16", "utf-16le":
		encoder = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder()
	case "utf-16be":
		encoder = unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewEncoder()
	case "gbk":
		encoder = simplifiedchinese.GBK.NewEncoder()
	case "gb18030":
		encoder = simplifiedchinese.GB18030.NewEncoder()
	case "big5":
		encoder = traditionalchinese.Big5.NewEncoder()
	case "shift_jis", "shift-jis", "sjis":
		encoder = japanese.ShiftJIS.NewEncoder()
	case "latin1", "iso-8859-1":
		encoder = charmap.ISO8859_1.NewEncoder()
	case "auto":
		return nil, fmt.Errorf("--output-encoding=auto is not supported; choose an explicit encoding")
	default:
		return nil, fmt.Errorf("unsupported --output-encoding %q (use utf-8, utf-16le, utf-16be, gbk, gb18030, big5, shift_jis, or latin1)", name)
	}
	return transform.NewWriter(w, encoder), nil
}

func normalizeEncodingName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func decodeAuto(r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)
	head, err := br.Peek(3)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(head) >= 3 && head[0] == 0xEF && head[1] == 0xBB && head[2] == 0xBF {
		_, _ = br.Discard(3)
		return br, nil
	}
	if len(head) >= 2 && head[0] == 0xFF && head[1] == 0xFE {
		_, _ = br.Discard(2)
		return transform.NewReader(br, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()), nil
	}
	if len(head) >= 2 && head[0] == 0xFE && head[1] == 0xFF {
		_, _ = br.Discard(2)
		return transform.NewReader(br, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()), nil
	}
	if utf8.Valid(head) {
		return br, nil
	}
	return nil, fmt.Errorf("cannot auto-detect input encoding; use --encoding=gbk/gb18030/big5/utf-16le")
}

func decodeUTF16(r io.Reader, name string) (io.Reader, error) {
	br := bufio.NewReader(r)
	head, err := br.Peek(2)
	if err != nil && err != io.EOF {
		return nil, err
	}
	endian := unicode.LittleEndian
	if name == "utf-16be" {
		endian = unicode.BigEndian
	}
	if len(head) >= 2 {
		if head[0] == 0xFF && head[1] == 0xFE {
			endian = unicode.LittleEndian
			_, _ = br.Discard(2)
		} else if head[0] == 0xFE && head[1] == 0xFF {
			endian = unicode.BigEndian
			_, _ = br.Discard(2)
		}
	}
	return transform.NewReader(br, unicode.UTF16(endian, unicode.IgnoreBOM).NewDecoder()), nil
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
