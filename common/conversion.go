package common

import (
	"fmt"
	"io"
	"reflect"
)

// UnmarshalFunc is a function type that parses input data into a Table
type UnmarshalFunc func(cfg *Config, table *Table) error

// MarshalFunc is a function type that converts a Table to output format
type MarshalFunc func(cfg *Config, table *Table) error

// FormatRegistry holds format-specific unmarshal and marshal functions
type FormatRegistry struct {
	UnmarshalMap map[string]UnmarshalFunc
	MarshalMap   map[string]MarshalFunc
}

// NewFormatRegistry creates a new format registry
func NewFormatRegistry() *FormatRegistry {
	return &FormatRegistry{
		UnmarshalMap: make(map[string]UnmarshalFunc),
		MarshalMap:   make(map[string]MarshalFunc),
	}
}

// RegisterFormat adds a format's unmarshal and marshal functions to the registry
func (fr *FormatRegistry) RegisterFormat(name string, unmarshal UnmarshalFunc, marshal MarshalFunc) {
	fr.UnmarshalMap[name] = unmarshal
	fr.MarshalMap[name] = marshal
}

// RegisterFormatAlias adds an alias for an existing format
func (fr *FormatRegistry) RegisterFormatAlias(alias, format string) {
	if unmarshal, ok := fr.UnmarshalMap[format]; ok {
		fr.UnmarshalMap[alias] = unmarshal
	}
	if marshal, ok := fr.MarshalMap[format]; ok {
		fr.MarshalMap[alias] = marshal
	}
}

// RegisterWriteOnlyFormat registers a format that only supports writing (no Unmarshal)
func (fr *FormatRegistry) RegisterWriteOnlyFormat(name string, marshal MarshalFunc) {
	fr.UnmarshalMap[name] = nil
	fr.MarshalMap[name] = marshal
}

// GetUnmarshalFunc returns the unmarshal function for a format
func (fr *FormatRegistry) GetUnmarshalFunc(format string) (UnmarshalFunc, bool) {
	fn, ok := fr.UnmarshalMap[format]
	return fn, ok
}

// GetMarshalFunc returns the marshal function for a format
func (fr *FormatRegistry) GetMarshalFunc(format string) (MarshalFunc, bool) {
	fn, ok := fr.MarshalMap[format]
	return fn, ok
}

// ConversionError represents a conversion error with context
type ConversionError struct {
	Stage  string // "unmarshal" or "marshal"
	Format string
	Err    error
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("error %s format %s: %v", e.Stage, e.Format, e.Err)
}

func (e *ConversionError) Unwrap() error {
	return e.Err
}

// PerformConversionWithRegistry performs conversion using a registry of format functions
func PerformConversionWithRegistry(registry *FormatRegistry, cfg *Config) error {
	// Parse input
	var table Table

	unmarshalFn, ok := registry.GetUnmarshalFunc(cfg.From)
	if !ok {
		return fmt.Errorf("unsupported `--from` format: %s", cfg.From)
	}

	if unmarshalFn == nil {
		return fmt.Errorf("format %s does not support reading (unmarshal)", cfg.From)
	}

	if err := unmarshalFn(cfg, &table); err != nil {
		return &ConversionError{Stage: "unmarshal", Format: cfg.From, Err: err}
	}

	// Apply transformations
	cfg.ApplyTransformations(&table)

	// Generate output
	marshalFn, ok := registry.GetMarshalFunc(cfg.To)
	if !ok {
		return fmt.Errorf("unsupported `--to` format: %s", cfg.To)
	}

	if marshalFn == nil {
		return fmt.Errorf("format %s does not support writing (marshal)", cfg.To)
	}

	if err := marshalFn(cfg, &table); err != nil {
		return &ConversionError{Stage: "marshal", Format: cfg.To, Err: err}
	}

	return nil
}

// ValidateIO validates that reader and writer are properly set
func ValidateIO(cfg *Config) error {
	if cfg.Reader == nil {
		return fmt.Errorf("reader is not set")
	}
	if cfg.Writer == nil {
		return fmt.Errorf("writer is not set")
	}
	return nil
}

// CopyReaderToWriter is a utility function that copies data from reader to writer
// Useful for passthrough scenarios or testing
func CopyReaderToWriter(reader io.Reader, writer io.Writer) error {
	// Check for nil using reflection to handle typed nil pointers
	// This handles the case where a typed nil pointer (e.g., *bytes.Buffer(nil))
	// is passed to an interface parameter
	if reader == nil {
		return fmt.Errorf("reader and writer must not be nil")
	}
	if readerVal := reflect.ValueOf(reader); readerVal.IsValid() && readerVal.Kind() == reflect.Ptr && readerVal.IsNil() {
		return fmt.Errorf("reader and writer must not be nil")
	}

	if writer == nil {
		return fmt.Errorf("reader and writer must not be nil")
	}
	if writerVal := reflect.ValueOf(writer); writerVal.IsValid() && writerVal.Kind() == reflect.Ptr && writerVal.IsNil() {
		return fmt.Errorf("reader and writer must not be nil")
	}

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if _, writeErr := writer.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
