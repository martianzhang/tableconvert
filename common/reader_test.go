package common

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStripUTF8BOM(t *testing.T) {
	withBOM := append([]byte{0xEF, 0xBB, 0xBF}, []byte("a,b\n1,2\n")...)
	got, err := io.ReadAll(StripUTF8BOM(bytes.NewReader(withBOM)))
	require.NoError(t, err)
	assert.Equal(t, "a,b\n1,2\n", string(got))

	got, err = io.ReadAll(StripUTF8BOM(strings.NewReader("a,b\n1,2\n")))
	require.NoError(t, err)
	assert.Equal(t, "a,b\n1,2\n", string(got))
}

func TestTrimUTF8BOM(t *testing.T) {
	assert.Equal(t, []byte("x"), TrimUTF8BOM([]byte("\xEF\xBB\xBFx")))
	assert.Equal(t, []byte("x"), TrimUTF8BOM([]byte("x")))
	assert.Equal(t, []byte(nil), TrimUTF8BOM(nil))
}
