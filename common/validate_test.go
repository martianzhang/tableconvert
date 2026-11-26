package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRejectsUnknownOption(t *testing.T) {
	cfg := &Config{
		From:      "csv",
		To:        "json",
		Extension: map[string]string{"sort": "name"},
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown option --sort")
}

func TestValidateSuggestsOption(t *testing.T) {
	cfg := &Config{
		From:      "csv",
		To:        "json",
		Extension: map[string]string{"formt": "2d"},
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "did you mean --format")
}

func TestValidateAllowUnknown(t *testing.T) {
	cfg := &Config{
		From:         "csv",
		To:           "json",
		AllowUnknown: true,
		Extension:    map[string]string{"custom": "value"},
	}
	assert.NoError(t, cfg.Validate())
}

func TestValidateEnumValues(t *testing.T) {
	valid := &Config{
		From: "mysql",
		To:   "csv",
		Extension: map[string]string{
			"delimiter": "TAB",
			"style":     "plus",
		},
	}
	assert.NoError(t, valid.Validate())

	invalid := &Config{
		From:      "mysql",
		To:        "csv",
		Extension: map[string]string{"delimiter": "SPACE"},
	}
	err := invalid.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--delimiter")
}

func TestValidateBooleanValues(t *testing.T) {
	cfg := &Config{
		From:      "excel",
		To:        "csv",
		Extension: map[string]string{"expand-merged": "maybe"},
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid boolean")
}

func TestValidateJSONOptions(t *testing.T) {
	cfg := &Config{
		From: "json",
		To:   "json",
		Extension: map[string]string{
			"format":         "2d",
			"flatten-arrays": "index",
			"json-path":      "data.items",
			"parse-types":    "true",
			"sort-headers":   "false",
		},
	}
	assert.NoError(t, cfg.Validate())

	bad := &Config{
		From:      "json",
		To:        "json",
		Extension: map[string]string{"flatten-arrays": "explode"},
	}
	assert.Error(t, bad.Validate())
}

func TestValidateEncodingValues(t *testing.T) {
	cfg := &Config{
		From:      "csv",
		To:        "csv",
		Extension: map[string]string{"encoding": "gbk", "output-encoding": "gb18030"},
	}
	assert.NoError(t, cfg.Validate())

	bad := &Config{
		From:      "csv",
		To:        "csv",
		Extension: map[string]string{"encoding": "utf-32"},
	}
	assert.Error(t, bad.Validate())
}

func TestParseConfigRejectsInvalidBoolean(t *testing.T) {
	_, err := ParseConfig([]string{"--verbose=maybe", "--from", "csv", "--to", "json"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid boolean")
}

func TestParseConfigAllowUnknown(t *testing.T) {
	cfg, err := ParseConfig([]string{"--allow-unknown", "--from", "csv", "--to", "json", "--foo", "bar"})
	require.NoError(t, err)
	assert.True(t, cfg.AllowUnknown)
	assert.NoError(t, cfg.Validate())
}
