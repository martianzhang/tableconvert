package twiki

import (
	"bytes"
	"os"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	table := &common.Table{
		Headers: []string{"Header"},
		Rows: [][]string{
			{"R1"},
			{"R2"},
		},
	}

	// Act
	err := Marshal(cfg, table)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "|=Header=|\n|R1|\n|R2|\n", buf.String())
}

func TestUnmarshal(t *testing.T) {
	// Arrange
	input := "|=Header1=|=Header2=|\n|Data1|Data2|\n|Data3|Data4|"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	// Act
	err := Unmarshal(cfg, table)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Data1", "Data2"},
		{"Data3", "Data4"},
	}, table.Rows)
}

func TestMarshalWithPipeCharacters(t *testing.T) {
	// Test that pipe characters in content are properly escaped
	var buf bytes.Buffer
	cfg := &common.Config{
		Writer: &buf,
	}

	table := &common.Table{
		Headers: []string{"Name|Age", "City"},
		Rows: [][]string{
			{"Alice|30", "New York"},
			{"Bob|25", "Los Angeles"},
		},
	}

	err := Marshal(cfg, table)
	assert.NoError(t, err)

	// Pipe characters should be escaped with backslash
	// TWiki format: |=Header1=|=Header2=|
	expected := "|=Name\\|Age=|=City=|\n|Alice\\|30|New York|\n|Bob\\|25|Los Angeles|\n"
	assert.Equal(t, expected, buf.String())
}

func TestUnmarshalWithEmptyCells(t *testing.T) {
	// Test handling of empty cells
	input := "|=Header1=|=Header2=|=Header3=|\n||Data2||\n|Data1||Data3|"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2", "Header3"}, table.Headers)
	assert.Equal(t, [][]string{
		{"", "Data2", ""},
		{"Data1", "", "Data3"},
	}, table.Rows)
}

func TestUnmarshalWithExtraSpaces(t *testing.T) {
	// Test handling of extra spaces
	input := "|=  Header1  =|=Header2=|\n|  Data1  |  Data2  |"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Data1", "Data2"},
	}, table.Rows)
}

func TestUnmarshalWithEmptyLines(t *testing.T) {
	// Test handling of empty lines
	input := "\n\n|=Header1=|=Header2=|\n\n|Data1|Data2|\n\n\n|Data3|Data4|\n\n"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Data1", "Data2"},
		{"Data3", "Data4"},
	}, table.Rows)
}

func TestUnmarshalWithNonTableContent(t *testing.T) {
	// Test handling of non-table content before and after table
	input := "Some text before\n|=Header1=|=Header2=|\n|Data1|Data2|\nSome text after\nMore text"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2"}, table.Headers)
	assert.Equal(t, [][]string{
		{"Data1", "Data2"},
	}, table.Rows)
}

func TestUnmarshalWithUTF8Characters(t *testing.T) {
	// Test handling of UTF-8 characters
	input := "|=测试=|=🎉=|=中文=|\n|测试|emoji|中文|"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"测试", "🎉", "中文"}, table.Headers)
	assert.Equal(t, [][]string{
		{"测试", "emoji", "中文"},
	}, table.Rows)
}

func TestUnmarshalErrorCases(t *testing.T) {
	table := &common.Table{}

	t.Run("nil table", func(t *testing.T) {
		reader := bytes.NewBufferString("|=Header=|\n|Data|")
		cfg := &common.Config{Reader: reader}
		err := Unmarshal(cfg, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "output table cannot be nil")
	})

	t.Run("no header row", func(t *testing.T) {
		reader := bytes.NewBufferString("|Data1|Data2|")
		cfg := &common.Config{Reader: reader}
		err := Unmarshal(cfg, table)
		assert.Error(t, err)
		// Should get error about header cell not wrapped in '=' signs
		assert.Contains(t, err.Error(), "not wrapped in '=' signs")
	})

	t.Run("header without equals", func(t *testing.T) {
		reader := bytes.NewBufferString("|Header1|Header2|\n|Data1|Data2|")
		cfg := &common.Config{Reader: reader}
		err := Unmarshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not wrapped in '=' signs")
	})

	t.Run("mismatched column count", func(t *testing.T) {
		reader := bytes.NewBufferString("|=Header1=|=Header2=|\n|Data1|")
		cfg := &common.Config{Reader: reader}
		err := Unmarshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data row has 1 columns, but header has 2")
	})

	t.Run("empty header cells", func(t *testing.T) {
		// Test case where header line has empty cells
		// This should succeed but create empty headers
		reader := bytes.NewBufferString("|=|=|\n|Data1|Data2|")
		cfg := &common.Config{Reader: reader}
		err := Unmarshal(cfg, table)
		assert.NoError(t, err)
		// Headers should be empty strings
		assert.Equal(t, []string{"", ""}, table.Headers)
		// But data row has 2 columns while headers have 2 (empty) columns
		// So this should work
		assert.Equal(t, [][]string{{"Data1", "Data2"}}, table.Rows)
	})
}

func TestMarshalErrorCases(t *testing.T) {
	t.Run("nil table", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &common.Config{Writer: &buf}
		err := Marshal(cfg, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input table pointer cannot be nil")
	})

	t.Run("writer error", func(t *testing.T) {
		// Use a file that will cause write errors
		tmpfile, err := os.CreateTemp("", "twiki_test.*.txt")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		// Close the file to cause write errors
		tmpfile.Close()

		// Try to write to closed file
		cfg := &common.Config{Writer: tmpfile}
		table := &common.Table{
			Headers: []string{"Header"},
			Rows:    [][]string{{"Data"}},
		}

		err = Marshal(cfg, table)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write")
	})
}

func TestUnmarshalEmptyTable(t *testing.T) {
	// Test with no data
	input := "|=Header1=|=Header2=|"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Header1", "Header2"}, table.Headers)
	assert.Equal(t, 0, len(table.Rows))
}

func TestUnmarshalWithConsecutivePipes(t *testing.T) {
	// Test handling of consecutive pipes
	// Header: |=Header1=||=Header2=| has 3 parts: ["=Header1=", "", "=Header2="]
	// But empty part is skipped, so we get 2 headers: ["Header1", "Header2"]
	// Data: |Data1||Data2| has 3 parts: ["Data1", "", "Data2"]
	// So data has 3 columns but header has 2 -> error
	input := "|=Header1=||=Header2=|\n|Data1||Data2|"
	reader := bytes.NewBufferString(input)
	cfg := &common.Config{
		Reader: reader,
	}
	table := &common.Table{}

	err := Unmarshal(cfg, table)
	// This should error because header has 2 columns but data has 3
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data row has 3 columns, but header has 2")
}

func TestRoundTrip(t *testing.T) {
	// Test that we can unmarshal and then marshal back to same format
	original := "|=Name=|=Age=|=City=|\n|Alice|30|New York|\n|Bob|25|Los Angeles|"

	// Parse
	reader := bytes.NewBufferString(original)
	cfg := &common.Config{Reader: reader}
	table := &common.Table{}
	err := Unmarshal(cfg, table)
	require.NoError(t, err)

	// Marshal
	var buf bytes.Buffer
	cfg.Writer = &buf
	err = Marshal(cfg, table)
	require.NoError(t, err)

	// The output should be equivalent (though whitespace might differ)
	// Let's parse the output again and compare
	output := buf.String()
	reader2 := bytes.NewBufferString(output)
	cfg2 := &common.Config{Reader: reader2}
	table2 := &common.Table{}
	err = Unmarshal(cfg2, table2)
	require.NoError(t, err)

	assert.Equal(t, table.Headers, table2.Headers)
	assert.Equal(t, table.Rows, table2.Rows)
}
