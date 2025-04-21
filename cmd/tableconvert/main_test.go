package main

import (
	"os"
	"testing"

	"github.com/martianzhang/tableconvert/common"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/stretchr/testify/assert"
)

var ProjectRoot string

func init() {
	var err error
	// Get project root path
	ProjectRoot, err = common.GetProjectRootPath()
	if err != nil {
		panic(err)
	}

	// Check if feature directory exists， if not create it
	if _, err := os.Stat(ProjectRoot + "feature"); os.IsNotExist(err) {
		os.Mkdir(ProjectRoot+"feature", 0755)
	}
}

func TestMain(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		args     []string
		file     string // relative to test/ directory
		expected string // relative to test/ directory
		result   string // relative to feature/ directory
	}{
		{
			name:   "mysql to markdown",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "markdown", "--key", "value"},
			file:   "mysql.txt",
			result: "mysql.md",
		},
		{
			name:   "markdown to mysql",
			args:   []string{"tableconvert", "--from", "markdown", "--to", "mysql"},
			file:   "mysql.md",
			result: "mysql.txt",
		},
		{
			name:   "mysql to csv (SEMICOLON delimiter)",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "csv", "--delimiter=SEMICOLON"},
			file:   "mysql.txt",
			result: "mysql.semicolon.csv",
		},
		{
			name:   "csv to mysql",
			args:   []string{"tableconvert", "--from", "csv", "--to", "mysql"},
			file:   "mysql.csv",
			result: "mysql.txt",
		},
		{
			name:   "mysql to json",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "json", "--parsing-json"},
			file:   "mysql.txt",
			result: "mysql.json",
		},
		{
			name:   "mysql to json (2d format, minified)",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "json", "--parsing-json", "--format=2d", "--minify"},
			file:   "mysql.txt",
			result: "mysql.2d.json",
		},
		{
			name:   "mysql to json (column format, minified)",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "json", "--parsing-json", "--format=column", "--minify"},
			file:   "mysql.txt",
			result: "mysql.column.json",
		},
		{
			name:   "mysql to sql",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "sql"},
			file:   "mysql.txt",
			result: "mysql.sql",
		},
		{
			name:   "mysql to sql (with table and dialect args)",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "sql", "--table", "tables", "--replace", "--dialect=none", "--one-insert"},
			file:   "mysql.txt",
			result: "mysql.one.sql",
		},
		{
			name:   "sql to mysql",
			args:   []string{"tableconvert", "--from", "sql", "--to", "mysql"},
			file:   "mysql.sql",
			result: "mysql.txt",
		},
		{
			name:   "replace sql to mysql",
			args:   []string{"tableconvert", "--from", "sql", "--to", "mysql"},
			file:   "mysql.replace.sql",
			result: "mysql.txt",
		},
		{
			name:   "mysql to xml",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "xml"},
			file:   "mysql.txt",
			result: "mysql.xml",
		},
		{
			name:   "xml to mysql",
			args:   []string{"tableconvert", "--from", "xml", "--to", "mysql"},
			file:   "mysql.xml",
			result: "mysql.txt",
		},
		{
			name:   "excel to mysql",
			args:   []string{"tableconvert", "--from", "xlsx", "--to", "mysql"},
			file:   "mysql.xlsx",
			result: "mysql.txt",
		},
		{
			name:   "mysql to twiki",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "twiki"},
			file:   "mysql.txt",
			result: "mysql.twiki",
		},
		{
			name:   "twiki to mysql",
			args:   []string{"tableconvert", "--from", "twiki", "--to", "mysql"},
			file:   "mysql.twiki",
			result: "mysql.txt",
		},
		{
			name:   "mysql to html",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "html"},
			file:   "mysql.txt",
			result: "mysql.html",
		},
		{
			name:   "html to mysql",
			args:   []string{"tableconvert", "--from", "html", "--to", "mysql"},
			file:   "mysql.html",
			result: "mysql.txt",
		},
		{
			name:   "mysql to mediawiki",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "mediawiki"},
			file:   "mysql.txt",
			result: "mysql.mediawiki",
		},
		{
			name:   "mediawiki to mysql",
			args:   []string{"tableconvert", "--from", "mediawiki", "--to", "mysql"},
			file:   "mysql.mediawiki",
			result: "mysql.txt",
		},
		{
			name:   "mysql to latex",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "latex"},
			file:   "mysql.txt",
			result: "mysql.latex",
		},
		{
			name:   "latex to mysql",
			args:   []string{"tableconvert", "--from", "latex", "--to", "mysql"},
			file:   "mysql.latex",
			result: "mysql.txt",
		},
		{
			name:   "mysql to ascii plus",
			args:   []string{"tableconvert", "--from", "mysql", "--to", "ascii", "--style", "plus"},
			file:   "mysql.txt",
			result: "mysql.plus.ascii",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare full paths
			inputFile := ProjectRoot + "test/" + tt.file
			resultFile := ProjectRoot + "feature/" + tt.result
			expectedFile := ProjectRoot + "test/" + tt.result

			// Set up test arguments
			args := append(tt.args, "--file", inputFile, "--result", resultFile)
			os.Args = args

			// Run main function in a separate goroutine to catch potential panics
			var mainErr interface{}
			func() {
				defer func() { mainErr = recover() }()
				main()
			}()

			// Assert no panic occurred
			assert.Nil(t, mainErr, "Main function should not panic")

			// Compare the generated file with expected file
			assertFilesEqual(t, expectedFile, resultFile)
		})
	}
}

func assertFilesEqual(t *testing.T, expectedFile, actualFile string) {
	t.Helper()

	expected, err := os.ReadFile(expectedFile)
	assert.Nil(t, err, "Should read expected file without error")

	actual, err := os.ReadFile(actualFile)
	assert.Nil(t, err, "Should read actual file without error")

	if string(expected) != string(actual) {
		edits := myers.ComputeEdits(span.URIFromPath(expectedFile), string(expected), string(actual))
		diff := gotextdiff.ToUnified(expectedFile, actualFile, string(expected), edits)
		t.Errorf("Files differ:\n%s", diff)
	}
}
