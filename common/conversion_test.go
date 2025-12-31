package common

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormatRegistry(t *testing.T) {
	registry := NewFormatRegistry()
	assert.NotNil(t, registry)
	assert.NotNil(t, registry.UnmarshalMap)
	assert.NotNil(t, registry.MarshalMap)
	assert.Equal(t, 0, len(registry.UnmarshalMap))
	assert.Equal(t, 0, len(registry.MarshalMap))
}

func TestRegisterFormat(t *testing.T) {
	registry := NewFormatRegistry()

	unmarshalCalled := false
	marshalCalled := false

	unmarshal := func(cfg *Config, table *Table) error {
		unmarshalCalled = true
		return nil
	}
	marshal := func(cfg *Config, table *Table) error {
		marshalCalled = true
		return nil
	}

	registry.RegisterFormat("test", unmarshal, marshal)

	// Verify it's registered
	assert.Equal(t, 1, len(registry.UnmarshalMap))
	assert.Equal(t, 1, len(registry.MarshalMap))

	// Test calling the registered functions
	unmarshalFn, ok := registry.GetUnmarshalFunc("test")
	assert.True(t, ok)
	assert.NotNil(t, unmarshalFn)

	marshalFn, ok := registry.GetMarshalFunc("test")
	assert.True(t, ok)
	assert.NotNil(t, marshalFn)

	// Call them to verify they're the right functions
	cfg := &Config{}
	table := &Table{}
	unmarshalFn(cfg, table)
	marshalFn(cfg, table)

	assert.True(t, unmarshalCalled)
	assert.True(t, marshalCalled)
}

func TestRegisterFormatAlias(t *testing.T) {
	registry := NewFormatRegistry()

	unmarshal := func(cfg *Config, table *Table) error { return nil }
	marshal := func(cfg *Config, table *Table) error { return nil }

	registry.RegisterFormat("original", unmarshal, marshal)
	registry.RegisterFormatAlias("alias", "original")

	// Verify alias works
	unmarshalFn, ok := registry.GetUnmarshalFunc("alias")
	assert.True(t, ok)
	assert.NotNil(t, unmarshalFn)

	marshalFn, ok := registry.GetMarshalFunc("alias")
	assert.True(t, ok)
	assert.NotNil(t, marshalFn)

	// Verify original still works
	unmarshalFn, ok = registry.GetUnmarshalFunc("original")
	assert.True(t, ok)
	assert.NotNil(t, unmarshalFn)
}

func TestRegisterFormatAliasNonExistent(t *testing.T) {
	registry := NewFormatRegistry()

	// Register alias for non-existent format
	registry.RegisterFormatAlias("alias", "nonexistent")

	// Alias should not be registered
	_, ok := registry.GetUnmarshalFunc("alias")
	assert.False(t, ok)

	_, ok = registry.GetMarshalFunc("alias")
	assert.False(t, ok)
}

func TestRegisterWriteOnlyFormat(t *testing.T) {
	registry := NewFormatRegistry()

	marshal := func(cfg *Config, table *Table) error { return nil }
	registry.RegisterWriteOnlyFormat("writeonly", marshal)

	// Unmarshal should be nil
	unmarshalFn, ok := registry.GetUnmarshalFunc("writeonly")
	assert.True(t, ok)
	assert.Nil(t, unmarshalFn)

	// Marshal should exist
	marshalFn, ok := registry.GetMarshalFunc("writeonly")
	assert.True(t, ok)
	assert.NotNil(t, marshalFn)
}

func TestGetUnmarshalFunc(t *testing.T) {
	registry := NewFormatRegistry()

	// Non-existent format
	_, ok := registry.GetUnmarshalFunc("nonexistent")
	assert.False(t, ok)

	// Register a format
	unmarshal := func(cfg *Config, table *Table) error { return nil }
	registry.RegisterFormat("test", unmarshal, nil)

	// Should exist now
	fn, ok := registry.GetUnmarshalFunc("test")
	assert.True(t, ok)
	assert.NotNil(t, fn)
}

func TestGetMarshalFunc(t *testing.T) {
	registry := NewFormatRegistry()

	// Non-existent format
	_, ok := registry.GetMarshalFunc("nonexistent")
	assert.False(t, ok)

	// Register a format
	marshal := func(cfg *Config, table *Table) error { return nil }
	registry.RegisterFormat("test", nil, marshal)

	// Should exist now
	fn, ok := registry.GetMarshalFunc("test")
	assert.True(t, ok)
	assert.NotNil(t, fn)
}

func TestConversionError(t *testing.T) {
	baseErr := fmt.Errorf("base error")
	convErr := &ConversionError{
		Stage:  "unmarshal",
		Format: "csv",
		Err:    baseErr,
	}

	assert.Equal(t, "error unmarshal format csv: base error", convErr.Error())
	assert.Equal(t, baseErr, convErr.Unwrap())
}

func TestPerformConversionWithRegistry(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*FormatRegistry)
		cfg         *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful conversion",
			setupFunc: func(r *FormatRegistry) {
				r.RegisterFormat("test", func(cfg *Config, table *Table) error {
					table.Headers = []string{"A", "B"}
					table.Rows = [][]string{{"1", "2"}}
					return nil
				}, func(cfg *Config, table *Table) error {
					return nil
				})
			},
			cfg: &Config{
				From:   "test",
				To:     "test",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: false,
		},
		{
			name: "unsupported from format",
			setupFunc: func(r *FormatRegistry) {
				// Empty registry
			},
			cfg: &Config{
				From:   "nonexistent",
				To:     "test",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: true,
			errorMsg:    "unsupported `--from` format: nonexistent",
		},
		{
			name: "unsupported to format",
			setupFunc: func(r *FormatRegistry) {
				r.RegisterFormat("test", func(cfg *Config, table *Table) error {
					return nil
				}, nil) // No marshal function
			},
			cfg: &Config{
				From:   "test",
				To:     "nonexistent",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: true,
			errorMsg:    "unsupported `--to` format: nonexistent",
		},
		{
			name: "from format with nil unmarshal",
			setupFunc: func(r *FormatRegistry) {
				r.RegisterWriteOnlyFormat("writeonly", func(cfg *Config, table *Table) error {
					return nil
				})
			},
			cfg: &Config{
				From:   "writeonly",
				To:     "test",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: true,
			errorMsg:    "format writeonly does not support reading (unmarshal)",
		},
		{
			name: "to format with nil marshal",
			setupFunc: func(r *FormatRegistry) {
				r.RegisterFormat("readonly", func(cfg *Config, table *Table) error {
					return nil
				}, nil)
			},
			cfg: &Config{
				From:   "readonly",
				To:     "readonly",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: true,
			errorMsg:    "format readonly does not support writing (marshal)",
		},
		{
			name: "unmarshal error",
			setupFunc: func(r *FormatRegistry) {
				r.RegisterFormat("test", func(cfg *Config, table *Table) error {
					return fmt.Errorf("unmarshal failed")
				}, func(cfg *Config, table *Table) error {
					return nil
				})
			},
			cfg: &Config{
				From:   "test",
				To:     "test",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: true,
			errorMsg:    "error unmarshal format test: unmarshal failed",
		},
		{
			name: "marshal error",
			setupFunc: func(r *FormatRegistry) {
				r.RegisterFormat("test", func(cfg *Config, table *Table) error {
					return nil
				}, func(cfg *Config, table *Table) error {
					return fmt.Errorf("marshal failed")
				})
			},
			cfg: &Config{
				From:   "test",
				To:     "test",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: true,
			errorMsg:    "error marshal format test: marshal failed",
		},
		{
			name: "with transformations",
			setupFunc: func(r *FormatRegistry) {
				r.RegisterFormat("test", func(cfg *Config, table *Table) error {
					table.Headers = []string{"name"}
					table.Rows = [][]string{{"alice"}}
					return nil
				}, func(cfg *Config, table *Table) error {
					// Verify transformation was applied
					if table.Headers[0] != "NAME" {
						return fmt.Errorf("expected uppercase header, got %s", table.Headers[0])
					}
					return nil
				})
			},
			cfg: &Config{
				From:   "test",
				To:     "test",
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
				Extension: map[string]string{
					"uppercase": "true",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewFormatRegistry()
			tt.setupFunc(registry)

			err := PerformConversionWithRegistry(registry, tt.cfg)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateIO(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *Config
		expectError bool
	}{
		{
			name: "both reader and writer set",
			cfg: &Config{
				Reader: strings.NewReader(""),
				Writer: &bytes.Buffer{},
			},
			expectError: false,
		},
		{
			name: "nil reader",
			cfg: &Config{
				Reader: nil,
				Writer: &bytes.Buffer{},
			},
			expectError: true,
		},
		{
			name: "nil writer",
			cfg: &Config{
				Reader: strings.NewReader(""),
				Writer: nil,
			},
			expectError: true,
		},
		{
			name: "both nil",
			cfg: &Config{
				Reader: nil,
				Writer: nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIO(tt.cfg)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCopyReaderToWriter(t *testing.T) {
	tests := []struct {
		name        string
		reader      *bytes.Buffer
		writer      *bytes.Buffer
		input       string
		expectError bool
	}{
		{
			name:        "successful copy",
			reader:      bytes.NewBufferString("Hello, World!"),
			writer:      &bytes.Buffer{},
			input:       "Hello, World!",
			expectError: false,
		},
		{
			name:        "empty input",
			reader:      bytes.NewBufferString(""),
			writer:      &bytes.Buffer{},
			input:       "",
			expectError: false,
		},
		{
			name:        "large content",
			reader:      bytes.NewBufferString(strings.Repeat("A", 10000)),
			writer:      &bytes.Buffer{},
			input:       strings.Repeat("A", 10000),
			expectError: false,
		},
		{
			name:        "nil reader",
			reader:      nil,
			writer:      &bytes.Buffer{},
			expectError: true,
		},
		{
			name:        "nil writer",
			reader:      bytes.NewBufferString("test"),
			writer:      nil,
			expectError: true,
		},
		{
			name:        "both nil",
			reader:      nil,
			writer:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CopyReaderToWriter(tt.reader, tt.writer)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.reader != nil && tt.writer != nil {
					assert.Equal(t, tt.input, tt.writer.String())
				}
			}
		})
	}
}

func TestPerformConversionWithRegistry_Advanced(t *testing.T) {
	// Test that transformations are applied correctly
	registry := NewFormatRegistry()

	unmarshal := func(cfg *Config, table *Table) error {
		table.Headers = []string{"name", "city"}
		table.Rows = [][]string{
			{"alice", "nyc"},
			{"bob", "la"},
		}
		return nil
	}

	marshal := func(cfg *Config, table *Table) error {
		// Verify uppercase transformation was applied
		if table.Headers[0] != "NAME" {
			return fmt.Errorf("expected NAME, got %s", table.Headers[0])
		}
		if table.Rows[0][0] != "ALICE" {
			return fmt.Errorf("expected ALICE, got %s", table.Rows[0][0])
		}
		return nil
	}

	registry.RegisterFormat("test", unmarshal, marshal)

	cfg := &Config{
		From:   "test",
		To:     "test",
		Reader: strings.NewReader(""),
		Writer: &bytes.Buffer{},
		Extension: map[string]string{
			"uppercase": "true",
		},
	}

	err := PerformConversionWithRegistry(registry, cfg)
	assert.NoError(t, err)
}

func TestConversionErrorWrapping(t *testing.T) {
	// Test that ConversionError properly wraps the underlying error
	baseErr := fmt.Errorf("underlying error")
	convErr := &ConversionError{
		Stage:  "marshal",
		Format: "json",
		Err:    baseErr,
	}

	// Test Error() method
	expectedMsg := "error marshal format json: underlying error"
	assert.Equal(t, expectedMsg, convErr.Error())

	// Test Unwrap() method
	assert.Equal(t, baseErr, convErr.Unwrap())

	// Test that errors.Is works correctly
	assert.ErrorIs(t, convErr, baseErr)
}
