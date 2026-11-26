package common

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func encodeWith(t *testing.T, enc encoding.Encoding, value string) []byte {
	t.Helper()
	var buf bytes.Buffer
	writer := transform.NewWriter(&buf, enc.NewEncoder())
	_, err := writer.Write([]byte(value))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return buf.Bytes()
}

func TestDecodeReaderGBK(t *testing.T) {
	data := encodeWith(t, simplifiedchinese.GBK, "张三,30\n")
	reader, err := DecodeReader(bytes.NewReader(data), "gbk")
	require.NoError(t, err)
	got, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "张三,30\n", string(got))
}

func TestDecodeReaderUTF16LEBOM(t *testing.T) {
	data := encodeWith(t, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM), "name,age\n张三,30\n")
	reader, err := DecodeReader(bytes.NewReader(data), "utf-16le")
	require.NoError(t, err)
	got, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "name,age\n张三,30\n", string(got))
}

func TestDecodeReaderAutoUTF8(t *testing.T) {
	reader, err := DecodeReader(strings.NewReader("a,b\n1,2\n"), "auto")
	require.NoError(t, err)
	got, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "a,b\n1,2\n", string(got))
}

func TestDecodeReaderAutoUnknown(t *testing.T) {
	_, err := DecodeReader(bytes.NewReader([]byte{0x81, 0x40}), "auto")
	require.Error(t, err)
}

func TestEncodeWriterGBK(t *testing.T) {
	var buf bytes.Buffer
	writer, err := EncodeWriter(&buf, "gbk")
	require.NoError(t, err)
	_, err = writer.Write([]byte("张三"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	reader := transform.NewReader(bytes.NewReader(buf.Bytes()), simplifiedchinese.GBK.NewDecoder())
	got, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "张三", string(got))
}

func TestEncodingErrors(t *testing.T) {
	_, err := DecodeReader(strings.NewReader("x"), "utf-32")
	assert.Error(t, err)
	_, err = EncodeWriter(io.Discard, "auto")
	assert.Error(t, err)
}

func TestIsTextFormat(t *testing.T) {
	assert.True(t, IsTextFormat("csv"))
	assert.True(t, IsTextFormat("json"))
	assert.False(t, IsTextFormat("excel"))
	assert.False(t, IsTextFormat("xlsx"))
}
