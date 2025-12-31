package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFormatParams(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected int // expected number of parameters
	}{
		{"markdown format", "markdown", 5},
		{"csv format", "csv", 3}, // first-column-header, bom, delimiter
		{"json format", "json", 3},
		{"latex format", "latex", 11},
		{"excel format", "excel", 4}, // first-column-header, sheet-name, auto-width, text-format
		{"html format", "html", 4},   // first-column-header, div, minify, thead
		{"ascii format", "ascii", 1},
		{"sql format", "sql", 4},
		{"xml format", "xml", 4},
		{"unknown format", "unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := GetFormatParams(tt.format)
			assert.Equal(t, tt.expected, len(params), "Expected %d parameters for %s", tt.expected, tt.format)
		})
	}
}

func TestGetAllFormats(t *testing.T) {
	formats := GetAllFormats()

	// Should contain all expected formats
	expectedFormats := []string{"ascii", "csv", "excel", "html", "json", "jsonl", "latex", "markdown", "mediawiki", "sql", "tmpl", "xml"}

	for _, expected := range expectedFormats {
		assert.Contains(t, formats, expected, "Expected format %s to be in the list", expected)
	}
}

func TestFormatExists(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{"markdown exists", "markdown", true},
		{"csv exists", "csv", true},
		{"json exists", "json", true},
		{"unknown format", "unknown", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := FormatExists(tt.format)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestGlobalTransformParams(t *testing.T) {
	assert.NotEmpty(t, GlobalTransformParams, "GlobalTransformParams should not be empty")

	// Check that expected transformation parameters exist
	paramNames := make(map[string]bool)
	for _, param := range GlobalTransformParams {
		paramNames[param.Name] = true
	}

	expectedParams := []string{"transpose", "delete-empty", "deduplicate", "uppercase", "lowercase", "capitalize"}
	for _, expected := range expectedParams {
		assert.True(t, paramNames[expected], "Expected global parameter %s to exist", expected)
	}
}

func TestFormatParamsStructure(t *testing.T) {
	// Test that format parameters have the correct structure
	markdownParams := GetFormatParams("markdown")

	assert.NotEmpty(t, markdownParams, "Markdown should have parameters")

	// Check first parameter structure
	param := markdownParams[0]
	assert.NotEmpty(t, param.Name, "Parameter should have a name")
	assert.NotEmpty(t, param.DefaultValue, "Parameter should have a default value")
	assert.NotEmpty(t, param.Description, "Parameter should have a description")
}

func TestFormatParamsRegistryComplete(t *testing.T) {
	// Verify that all formats in the registry have proper parameters
	for format, params := range FormatParamsRegistry {
		assert.NotEmpty(t, params, "Format %s should have parameters", format)

		for _, param := range params {
			assert.NotEmpty(t, param.Name, "Format %s has parameter without name", format)
			assert.NotEmpty(t, param.Description, "Format %s parameter %s has no description", format, param.Name)
		}
	}
}
