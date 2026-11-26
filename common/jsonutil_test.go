package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONPathValueNested(t *testing.T) {
	root, err := DecodeOrderedJSON([]byte(`{"data":{"items":[{"name":"a"},{"name":"b"}]}}`))
	require.NoError(t, err)

	value, err := JSONPathValue(root, "data.items[1].name")
	require.NoError(t, err)
	assert.Equal(t, "b", value)
}

func TestJSONPathValueErrors(t *testing.T) {
	root, err := DecodeOrderedJSON([]byte(`{"data":{"items":[1]}}`))
	require.NoError(t, err)

	_, err = JSONPathValue(root, "data.missing")
	assert.Error(t, err)
	_, err = JSONPathValue(root, "data.items[2]")
	assert.Error(t, err)
	_, err = JSONPathValue(root, "data.items.name")
	assert.Error(t, err)
}

func TestFlattenJSONValueSpecialKeys(t *testing.T) {
	root, err := DecodeOrderedJSON([]byte(`{"a.b":1,"名字":"张三","tags":["x"]}`))
	require.NoError(t, err)

	record := FlattenJSONValue(root, JSONFlattenOptions{})
	require.NotNil(t, record)
	assert.Equal(t, []string{`["a.b"]`, "名字", "tags"}, record.Keys)
	assert.Equal(t, "1", record.Values[`["a.b"]`])
	assert.Equal(t, "张三", record.Values["名字"])
	assert.Equal(t, `["x"]`, record.Values["tags"])
}

func TestDecodeOrderedJSONTrailingData(t *testing.T) {
	_, err := DecodeOrderedJSON([]byte(`{"a":1} {"b":2}`))
	assert.Error(t, err)
}

func TestGetParseTypesAliasPrecedence(t *testing.T) {
	cfg := &Config{Extension: map[string]string{"parsing-json": "true"}}
	assert.True(t, cfg.GetParseTypes())

	cfg.Extension["parse-types"] = "false"
	assert.False(t, cfg.GetParseTypes())
}
