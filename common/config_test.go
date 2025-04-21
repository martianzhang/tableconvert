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
	args := []string{"--from=en", "--to=zh"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "en", cfg.From)
	assert.Equal(t, "zh", cfg.To)
	assert.False(t, cfg.Verbose)
	assert.Empty(t, cfg.File)
	assert.Empty(t, cfg.Result)
	assert.NotNil(t, cfg.Reader)
	assert.NotNil(t, cfg.Writer)
	assert.Empty(t, cfg.Extension)
}

func TestParseConfigWithSpaceFormat(t *testing.T) {
	// Arrange
	args := []string{"--from", "en", "--to", "zh"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "en", cfg.From)
	assert.Equal(t, "zh", cfg.To)
	assert.False(t, cfg.Verbose)
	assert.Empty(t, cfg.File)
	assert.Empty(t, cfg.Result)
	assert.NotNil(t, cfg.Reader)
	assert.NotNil(t, cfg.Writer)
	assert.Empty(t, cfg.Extension)
}

func TestParseConfigShortFormParameters(t *testing.T) {
	// Arrange
	args := []string{"-f", "en", "-t", "zh"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, "en", config.From, "From parameter should be 'en'")
	assert.Equal(t, "zh", config.To, "To parameter should be 'zh'")
	assert.Empty(t, config.File, "File parameter should be empty")
	assert.Empty(t, config.Result, "Result parameter should be empty")
	assert.False(t, config.Verbose, "Verbose should be false by default")
	assert.NotNil(t, config.Reader, "Reader should not be nil")
	assert.NotNil(t, config.Writer, "Writer should not be nil")
	assert.Empty(t, config.Extension, "Extension map should be empty")
}

func TestParseConfigDuplicateParameters(t *testing.T) {
	// Arrange
	args := []string{"--from", "en", "--from", "es", "--to", "zh"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "es", config.From, "Expected last 'from' value to be used")
	assert.Equal(t, "zh", config.To, "Expected 'to' parameter to be set correctly")
}

func TestParseConfigMultipleUnknownParameters(t *testing.T) {
	// Arrange
	args := []string{"--from", "en", "--to", "zh", "--param1", "val1", "--param2", "val2"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "en", config.From)
	assert.Equal(t, "zh", config.To)
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
	assert.Equal(t, "must provide -f|--from and -t|--to parameters", err.Error(), "Expected error message about missing required parameters")

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
	args := []string{"--from", "en", "--to", "zh", "--unknown", "value"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "en", config.From)
	assert.Equal(t, "zh", config.To)
	assert.Contains(t, config.Extension, "unknown")
	assert.Equal(t, "value", config.Extension["unknown"])
}

func TestParseConfigMixedParameterFormats(t *testing.T) {
	// Arrange
	args := []string{"--from=en", "-t", "zh", "--verbose"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "en", cfg.From, "From parameter should be 'en'")
	assert.Equal(t, "zh", cfg.To, "To parameter should be 'zh'")
	assert.True(t, cfg.Verbose, "Verbose should be true")

	// Verify default values
	assert.Empty(t, cfg.File, "File should be empty")
	assert.Empty(t, cfg.Result, "Result should be empty")
	assert.NotNil(t, cfg.Extension, "Extension map should be initialized")
	assert.Empty(t, cfg.Extension, "Extension map should be empty")
}

func TestParseConfigEmptyVerboseFlag(t *testing.T) {
	// Arrange
	args := []string{"--from", "en", "--to", "zh", "--verbose"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.True(t, config.Verbose)
	assert.Equal(t, "en", config.From)
	assert.Equal(t, "zh", config.To)
}

func TestParseConfigVerboseFlagFalse(t *testing.T) {
	// Arrange
	args := []string{"--from", "en", "--to", "zh", "--verbose=false"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "en", config.From)
	assert.Equal(t, "zh", config.To)
	assert.False(t, config.Verbose)
}

func TestParseConfigVerboseFlagTrue(t *testing.T) {
	// Arrange
	args := []string{"--from", "en", "--to", "zh", "--verbose=true"}

	// Act
	config, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.True(t, config.Verbose)
	assert.Equal(t, "en", config.From)
	assert.Equal(t, "zh", config.To)
}

func TestParseConfigWithResultFileOutput(t *testing.T) {
	// Setup
	testArgs := []string{"--from", "en", "--to", "zh", "--result", "output.txt"}

	// Clean up any existing test file before and after the test
	defer os.Remove("output.txt")

	// Execute
	cfg, err := ParseConfig(testArgs)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "en", cfg.From)
	assert.Equal(t, "zh", cfg.To)
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
	args := []string{"--from", "en", "--to", "zh", "--file", "nonexistent.txt"}

	// Act
	cfg, err := ParseConfig(args)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
	assert.Equal(t, "nonexistent.txt", cfg.File)
	assert.Equal(t, "en", cfg.From)
	assert.Equal(t, "zh", cfg.To)
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
	args := []string{"--from", "en", "--to", "zh", "--file", tmpFile.Name()}

	// Execute
	cfg, err := ParseConfig(args)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "en", cfg.From)
	assert.Equal(t, "zh", cfg.To)
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
	assert.Equal(t, true, strings.HasSuffix(path, "/tableconvert/"))
}
