package common

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigMissingRequiredParameters(t *testing.T) {
	// Arrange
	args := []string{"--file", "input.txt"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NotNil(t, err, "Expected error for missing required parameters")
	assert.Equal(t, "", cfg.From, "From parameter should be empty")
	assert.Equal(t, "", cfg.To, "To parameter should be empty")
	assert.Equal(t, "input.txt", cfg.File, "File parameter should be set")
}

func TestParseConfigWithEqualSign(t *testing.T) {
	// Arrange
	args := []string{"--from=csv", "--to=json"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "csv", cfg.From)
	assert.Equal(t, "json", cfg.To)
	assert.False(t, cfg.Verbose)
	assert.Empty(t, cfg.File)
	assert.Empty(t, cfg.Result)
	assert.NotNil(t, cfg.Reader)
	assert.NotNil(t, cfg.Writer)
	assert.Empty(t, cfg.Extension)
}

func TestParseConfigWithSpaceFormat(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--to", "json"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "csv", cfg.From)
	assert.Equal(t, "json", cfg.To)
	assert.False(t, cfg.Verbose)
	assert.Empty(t, cfg.File)
	assert.Empty(t, cfg.Result)
	assert.NotNil(t, cfg.Reader)
	assert.NotNil(t, cfg.Writer)
	assert.Empty(t, cfg.Extension)
}

func TestParseConfigShortFormParameters(t *testing.T) {
	// Arrange
	args := []string{"-f", "csv", "-t", "json"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, "csv", config.From, "From parameter should be 'csv'")
	assert.Equal(t, "json", config.To, "To parameter should be 'json'")
	assert.Empty(t, config.File, "File parameter should be empty")
	assert.Empty(t, config.Result, "Result parameter should be empty")
	assert.False(t, config.Verbose, "Verbose should be false by default")
	assert.NotNil(t, config.Reader, "Reader should not be nil")
	assert.NotNil(t, config.Writer, "Writer should not be nil")
	assert.Empty(t, config.Extension, "Extension map should be empty")
}

func TestParseConfigDuplicateParameters(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--from", "json", "--to", "markdown"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "json", config.From, "Expected last 'from' value to be used")
	assert.Equal(t, "markdown", config.To, "Expected 'to' parameter to be set correctly")
}

func TestParseConfigMultipleUnknownParameters(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--to", "json", "--param1", "val1", "--param2", "val2"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "csv", config.From)
	assert.Equal(t, "json", config.To)
	assert.Equal(t, 2, len(config.Extension))
	assert.Equal(t, "val1", config.Extension["param1"])
	assert.Equal(t, "val2", config.Extension["param2"])
}

func TestParseConfigWithEmptyInput(t *testing.T) {
	// Arrange
	emptyArgs := []string{}

	// Act
	cfg, err := ParseConfig(emptyArgs)

	// Assert
	assert.NotNil(t, err, "Expected error for empty input")
	assert.Contains(t, err.Error(), "Missing required parameters", "Error should mention missing parameters")

	// Verify default values in config
	assert.Empty(t, cfg.From, "From should be empty")
	assert.Empty(t, cfg.To, "To should be empty")
	assert.Empty(t, cfg.File, "File should be empty")
	assert.Empty(t, cfg.Result, "Result should be empty")
	assert.False(t, cfg.Verbose, "Verbose should be false by default")
	assert.NotNil(t, cfg.Extension, "Extension map should be initialized")
	assert.Empty(t, cfg.Extension, "Extension map should be empty")
}

func TestParseConfigUnknownParameters(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--to", "json", "--unknown", "value"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "csv", config.From)
	assert.Equal(t, "json", config.To)
	assert.Contains(t, config.Extension, "unknown")
	assert.Equal(t, "value", config.Extension["unknown"])
}

func TestParseConfigMixedParameterFormats(t *testing.T) {
	// Arrange
	args := []string{"--from=csv", "-t", "json", "--verbose"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "csv", cfg.From, "From parameter should be 'csv'")
	assert.Equal(t, "json", cfg.To, "To parameter should be 'json'")
	assert.True(t, cfg.Verbose, "Verbose should be true")

	// Verify default values
	assert.Empty(t, cfg.File, "File should be empty")
	assert.Empty(t, cfg.Result, "Result should be empty")
	assert.NotNil(t, cfg.Extension, "Extension map should be initialized")
	assert.Empty(t, cfg.Extension, "Extension map should be empty")
}

func TestParseConfigEmptyVerboseFlag(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--to", "json", "--verbose"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.True(t, config.Verbose)
	assert.Equal(t, "csv", config.From)
	assert.Equal(t, "json", config.To)
}

func TestParseConfigVerboseFlagFalse(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--to", "json", "--verbose=false"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "csv", config.From)
	assert.Equal(t, "json", config.To)
	assert.False(t, config.Verbose)
}

func TestParseConfigVerboseFlagTrue(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--to", "json", "--verbose=true"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.True(t, config.Verbose)
	assert.Equal(t, "csv", config.From)
	assert.Equal(t, "json", config.To)
}

func TestParseConfigWithResultFileOutput(t *testing.T) {
	// Setup
	testArgs := []string{"--from", "csv", "--to", "json", "--result", "output.txt"}

	// Clean up any existing test file before and after the test
	defer os.Remove("output.txt")

	// Execute
	cfg, err := ParseConfig(testArgs)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "csv", cfg.From)
	assert.Equal(t, "json", cfg.To)
	assert.Equal(t, "output.txt", cfg.Result)

	// Verify that Writer is not nil and is a *os.File
	fileWriter, ok := cfg.Writer.(*os.File)
	assert.True(t, ok, "Writer should be a *os.File")
	assert.NotNil(t, fileWriter)

	// Clean up: close the file
	if fileWriter != nil {
		fileWriter.Close()
	}
}

func TestParseConfigWithNonExistentFile(t *testing.T) {
	// Arrange
	args := []string{"--from", "csv", "--to", "json", "--file", "nonexistent.txt"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
	assert.Equal(t, "nonexistent.txt", cfg.File)
	assert.Equal(t, "csv", cfg.From)
	assert.Equal(t, "json", cfg.To)
}

func TestParseConfigWithValidFileInput(t *testing.T) {
	// Setup: Create a temporary test file
	tmpFile, err := os.CreateTemp("", "existing_file.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after test
	defer tmpFile.Close()

	// Test input
	args := []string{"--from", "csv", "--to", "json", "--file", tmpFile.Name()}

	// Execute
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "csv", cfg.From)
	assert.Equal(t, "json", cfg.To)
	assert.Equal(t, tmpFile.Name(), cfg.File)
	assert.NotNil(t, cfg.Reader)

	// Verify that the Reader is actually a *os.File
	_, ok := cfg.Reader.(*os.File)
	assert.True(t, ok, "Reader should be a *os.File")

	// Clean up: Close the reader if it's different from the temporary file
	if reader, ok := cfg.Reader.(*os.File); ok && reader != tmpFile {
		reader.Close()
	}
}

func TestGetProjectRootPath(t *testing.T) {
	path, err := GetProjectRootPath()
	assert.Nil(t, err)
	assert.Equal(t, true, strings.Contains(path, "tableconvert"))
}

func TestGetExtensionBool(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		key          string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "nil extension map returns default",
			config:       &Config{Extension: nil},
			key:          "test",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "key not found returns default",
			config:       &Config{Extension: map[string]string{}},
			key:          "missing",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "true value",
			config:       &Config{Extension: map[string]string{"test": "true"}},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "yes value",
			config:       &Config{Extension: map[string]string{"test": "yes"}},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "y value",
			config:       &Config{Extension: map[string]string{"test": "y"}},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "1 value",
			config:       &Config{Extension: map[string]string{"test": "1"}},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "empty string value returns true",
			config:       &Config{Extension: map[string]string{"test": ""}},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "false value",
			config:       &Config{Extension: map[string]string{"test": "false"}},
			key:          "test",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "no value",
			config:       &Config{Extension: map[string]string{"test": "no"}},
			key:          "test",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "n value",
			config:       &Config{Extension: map[string]string{"test": "n"}},
			key:          "test",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "0 value",
			config:       &Config{Extension: map[string]string{"test": "0"}},
			key:          "test",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "unknown value returns default",
			config:       &Config{Extension: map[string]string{"test": "unknown"}},
			key:          "test",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "whitespace handling",
			config:       &Config{Extension: map[string]string{"test": "  true  "}},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "case insensitive - TRUE",
			config:       &Config{Extension: map[string]string{"test": "TRUE"}},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "case insensitive - FALSE",
			config:       &Config{Extension: map[string]string{"test": "FALSE"}},
			key:          "test",
			defaultValue: true,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetExtensionBool(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetExtensionString(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "nil extension map returns default",
			config:       &Config{Extension: nil},
			key:          "test",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "key not found returns default",
			config:       &Config{Extension: map[string]string{}},
			key:          "missing",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "key found returns value",
			config:       &Config{Extension: map[string]string{"test": "value"}},
			key:          "test",
			defaultValue: "default",
			expected:     "value",
		},
		{
			name:         "empty string value",
			config:       &Config{Extension: map[string]string{"test": ""}},
			key:          "test",
			defaultValue: "default",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetExtensionString(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetExtensionInt(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "nil extension map returns default",
			config:       &Config{Extension: nil},
			key:          "test",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "key not found returns default",
			config:       &Config{Extension: map[string]string{}},
			key:          "missing",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "valid integer string",
			config:       &Config{Extension: map[string]string{"test": "123"}},
			key:          "test",
			defaultValue: 42,
			expected:     123,
		},
		{
			name:         "negative integer",
			config:       &Config{Extension: map[string]string{"test": "-5"}},
			key:          "test",
			defaultValue: 42,
			expected:     -5,
		},
		{
			name:         "invalid integer returns default",
			config:       &Config{Extension: map[string]string{"test": "abc"}},
			key:          "test",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "empty string returns default",
			config:       &Config{Extension: map[string]string{"test": ""}},
			key:          "test",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "float string returns default",
			config:       &Config{Extension: map[string]string{"test": "3.14"}},
			key:          "test",
			defaultValue: 42,
			expected:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetExtensionInt(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseConfigMCPMode(t *testing.T) {
	// Test MCP mode enabled
	args := []string{"--mcp"}
	cfg, err := ParseConfig(args)
	assert.NoError(t, err)
	assert.True(t, cfg.MCPMode)

	// Test MCP mode with value
	args = []string{"--mcp=true"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err)
	assert.True(t, cfg.MCPMode)

	// Test MCP mode disabled - normal validation applies
	args = []string{"--mcp=false", "--from", "csv", "--to", "json"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err) // Valid formats provided, should succeed
	assert.False(t, cfg.MCPMode)
	assert.Equal(t, "csv", cfg.From)
	assert.Equal(t, "json", cfg.To)

	// Test MCP mode skips validation
	args = []string{"--mcp", "--from", "csv"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err)
	assert.True(t, cfg.MCPMode)
	assert.Equal(t, "csv", cfg.From)
}

func TestParseConfigHelpFlags(t *testing.T) {
	// These tests will call os.Exit(0), so we can't test them directly
	// But we can verify the parsing logic doesn't fail before exit
	// For now, we'll skip these as they require mocking os.Exit
	t.Skip("Help flags call os.Exit and cannot be tested directly")
}

func TestParseConfigVerboseFlagVariations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"verbose without value", []string{"--mcp", "--verbose"}, true},
		{"verbose=true", []string{"--mcp", "--verbose=true"}, true},
		{"verbose=false", []string{"--mcp", "--verbose=false"}, false},
		{"verbose=True", []string{"--mcp", "--verbose=True"}, true},
		{"verbose=FALSE", []string{"--mcp", "--verbose=FALSE"}, false},
		{"verbose=yes", []string{"--mcp", "--verbose=yes"}, true},
		{"verbose=no", []string{"--mcp", "--verbose=no"}, false},
		{"verbose=1", []string{"--mcp", "--verbose=1"}, true},
		{"verbose=0", []string{"--mcp", "--verbose=0"}, false},
		{"verbose=y", []string{"--mcp", "--verbose=y"}, true},
		{"verbose=n", []string{"--mcp", "--verbose=n"}, false},
		{"verbose=unknown", []string{"--mcp", "--verbose=unknown"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ParseConfig(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cfg.Verbose)
		})
	}
}

func TestParseConfigMCPFlagVariations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
		needsIO  bool // whether from/to are needed for validation
	}{
		{"mcp without value", []string{"--mcp"}, true, false},
		{"mcp=true", []string{"--mcp=true"}, true, false},
		{"mcp=false", []string{"--mcp=false", "--from", "csv", "--to", "json"}, false, true},
		{"mcp=True", []string{"--mcp=True"}, true, false},
		{"mcp=FALSE", []string{"--mcp=FALSE", "--from", "csv", "--to", "json"}, false, true},
		{"mcp=yes", []string{"--mcp=yes"}, true, false},
		{"mcp=no", []string{"--mcp=no", "--from", "csv", "--to", "json"}, false, true},
		{"mcp=1", []string{"--mcp=1"}, true, false},
		{"mcp=0", []string{"--mcp=0", "--from", "csv", "--to", "json"}, false, true},
		{"mcp=y", []string{"--mcp=y"}, true, false},
		{"mcp=n", []string{"--mcp=n", "--from", "csv", "--to", "json"}, false, true},
		{"mcp=unknown", []string{"--mcp=unknown", "--from", "csv", "--to", "json"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ParseConfig(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cfg.MCPMode)
		})
	}
}

func TestParseConfigMixedAndEdgeCases(t *testing.T) {
	// Test mixed parameter formats
	args := []string{"--from=csv", "-t", "json", "--verbose", "--param1", "val1", "--param2=val2"}
	cfg, err := ParseConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, "csv", cfg.From)
	assert.Equal(t, "json", cfg.To)
	assert.True(t, cfg.Verbose)
	assert.Equal(t, "val1", cfg.Extension["param1"])
	assert.Equal(t, "val2", cfg.Extension["param2"])

	// Test parameters with special characters
	args = []string{"--from", "csv", "--to", "json", "--delimiter", ";", "--table", "users_data"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, ";", cfg.Extension["delimiter"])
	assert.Equal(t, "users_data", cfg.Extension["table"])

	// Test multiple unknown parameters
	args = []string{"--from", "csv", "--to", "json", "--param1", "1", "--param2", "2", "--param3", "3"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, "1", cfg.Extension["param1"])
	assert.Equal(t, "2", cfg.Extension["param2"])
	assert.Equal(t, "3", cfg.Extension["param3"])
}

func TestParseConfigBatchMode(t *testing.T) {
	// Test basic batch mode
	args := []string{"--batch=*.csv", "--to=json"}
	cfg, err := ParseConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, "*.csv", cfg.Batch)
	assert.Equal(t, "json", cfg.To)
	assert.False(t, cfg.Recursive)
	assert.Empty(t, cfg.OutputDir)

	// Test batch mode with recursive
	args = []string{"--batch=dir/**/*.csv", "--to=json", "--recursive"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, "dir/**/*.csv", cfg.Batch)
	assert.Equal(t, "json", cfg.To)
	assert.True(t, cfg.Recursive)

	// Test batch mode with output directory
	args = []string{"--batch=*.csv", "--to=json", "--output-dir=output"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, "*.csv", cfg.Batch)
	assert.Equal(t, "json", cfg.To)
	assert.Equal(t, "output", cfg.OutputDir)

	// Test batch mode with short flags
	args = []string{"-b", "*.csv", "-t", "json", "-r"}
	cfg, err = ParseConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, "*.csv", cfg.Batch)
	assert.Equal(t, "json", cfg.To)
	assert.True(t, cfg.Recursive)

	// Test batch mode requires --to
	args = []string{"--batch=*.csv"}
	cfg, err = ParseConfig(args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "batch mode requires --to")

	// Test batch mode with invalid format
	args = []string{"--batch=*.csv", "--to=invalid"}
	cfg, err = ParseConfig(args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported output format")
}

func TestGetBatchFiles(t *testing.T) {
	// Create test directory structure
	testDir := t.TempDir()

	// Create test files
	os.WriteFile(testDir+"/data1.csv", []byte("a,b\n1,2"), 0644)
	os.WriteFile(testDir+"/data2.csv", []byte("c,d\n3,4"), 0644)
	os.WriteFile(testDir+"/data3.txt", []byte("e,f\n5,6"), 0644)

	// Create subdirectory
	subDir := testDir + "/sub"
	os.Mkdir(subDir, 0755)
	os.WriteFile(subDir+"/nested.csv", []byte("g,h\n7,8"), 0644)

	// Test basic pattern
	t.Run("basic pattern", func(t *testing.T) {
		cfg := &Config{
			Batch: testDir + "/*.csv",
			To:    "json",
		}
		files, err := cfg.GetBatchFiles()
		assert.NoError(t, err)
		assert.Len(t, files, 2)

		// Check that files are sorted
		assert.True(t, files[0].InputPath < files[1].InputPath)

		// Check format detection
		assert.Equal(t, "csv", files[0].FromFormat)
		assert.Equal(t, "json", files[0].ToFormat)
	})

	// Test recursive pattern
	t.Run("recursive pattern", func(t *testing.T) {
		cfg := &Config{
			Batch:     testDir + "/**/*.csv",
			To:        "json",
			Recursive: true,
		}
		files, err := cfg.GetBatchFiles()
		assert.NoError(t, err)
		assert.Len(t, files, 3) // data1.csv, data2.csv, sub/nested.csv

		// Check that nested file is included
		nestedFound := false
		for _, f := range files {
			if strings.Contains(f.InputPath, "nested.csv") {
				nestedFound = true
				assert.True(t, strings.Contains(f.OutputPath, "nested.json"))
			}
		}
		assert.True(t, nestedFound, "nested.csv should be found in recursive mode")
	})

	// Test with output directory
	t.Run("with output directory", func(t *testing.T) {
		outputDir := t.TempDir()
		cfg := &Config{
			Batch:     testDir + "/*.csv",
			To:        "json",
			OutputDir: outputDir,
		}
		files, err := cfg.GetBatchFiles()
		assert.NoError(t, err)
		assert.Len(t, files, 2)

		// All outputs should be in the output directory
		for _, f := range files {
			assert.True(t, strings.HasPrefix(f.OutputPath, outputDir))
		}
	})

	// Test no files found
	t.Run("no files found", func(t *testing.T) {
		cfg := &Config{
			Batch: testDir + "/*.xyz",
			To:    "json",
		}
		_, err := cfg.GetBatchFiles()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no files found")
	})

	// Test format auto-detection fails for .txt
	t.Run("format auto-detection fails for txt", func(t *testing.T) {
		cfg := &Config{
			Batch: testDir + "/*.txt",
			To:    "csv",
		}
		_, err := cfg.GetBatchFiles()
		// .txt files don't auto-detect, should return error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot detect format")
	})

	// Test with --from specified
	t.Run("with explicit from format", func(t *testing.T) {
		cfg := &Config{
			Batch: testDir + "/*.txt",
			From:  "csv",
			To:    "json",
		}
		files, err := cfg.GetBatchFiles()
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, "csv", files[0].FromFormat)
	})
}
