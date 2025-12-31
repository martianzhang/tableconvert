package common

// FormatParam represents a single format-specific parameter
type FormatParam struct {
	Name          string
	DefaultValue  string
	AllowedValues string
	Description   string
}

// FormatParamsRegistry maps format names to their supported extension parameters
var FormatParamsRegistry = map[string][]FormatParam{
	"ascii": {
		{Name: "style", DefaultValue: "box", AllowedValues: "box, plus(+), dot(·), bubble(◌)", Description: "Table Style"},
	},
	"csv": {
		{Name: "first-column-header", DefaultValue: "false", AllowedValues: "true, false", Description: "Use first column as headers"},
		{Name: "bom", DefaultValue: "false", AllowedValues: "", Description: "Add Byte Order Mark"},
		{Name: "delimiter", DefaultValue: ",", AllowedValues: "COMMA, TAB, SEMICOLON, PIPE, SLASH, HASH", Description: "Value Delimiter"},
	},
	"excel": {
		{Name: "first-column-header", DefaultValue: "false", AllowedValues: "true, false", Description: "Use first column as headers"},
		{Name: "sheet-name", DefaultValue: "Sheet1", AllowedValues: "", Description: "Excel Sheet Name"},
		{Name: "auto-width", DefaultValue: "false", AllowedValues: "true, false", Description: "Auto Width"},
		{Name: "text-format", DefaultValue: "true", AllowedValues: "true, false", Description: "force text format"},
	},
	"html": {
		{Name: "first-column-header", DefaultValue: "false", AllowedValues: "true, false", Description: "Use first column as headers"},
		{Name: "div", DefaultValue: "false", AllowedValues: "true, false", Description: "Convert into div table"},
		{Name: "minify", DefaultValue: "false", AllowedValues: "true, false", Description: "Minify HTML table"},
		{Name: "thead", DefaultValue: "false", AllowedValues: "true, false", Description: "Include thead and tbody tags"},
	},
	"json": {
		{Name: "format", DefaultValue: "object", AllowedValues: "object, 2d, column, keyed", Description: "JSON Format"},
		{Name: "minify", DefaultValue: "false", AllowedValues: "true, false", Description: "Minify JSON"},
		{Name: "parsing-json", DefaultValue: "false", AllowedValues: "true, false", Description: "Parsing JSON"},
	},
	"jsonl": {
		{Name: "parsing-json", DefaultValue: "false", AllowedValues: "true, false", Description: "Parsing JSON"},
	},
	"latex": {
		{Name: "bold-first-column", DefaultValue: "false", AllowedValues: "true, false", Description: "Bold first column"},
		{Name: "bold-first-row", DefaultValue: "false", AllowedValues: "true, false", Description: "Bold first row"},
		{Name: "borders", DefaultValue: "1111,1111", AllowedValues: "1111,1111, 1101,1101, 0000,1101, 1111,0100, 0000,0100, 0000,0000", Description: "Table Border"},
		{Name: "caption", DefaultValue: "", AllowedValues: "", Description: "Table Caption"},
		{Name: "escape", DefaultValue: "true", AllowedValues: "true, false", Description: "Escape LaTeX table"},
		{Name: "ht", DefaultValue: "false", AllowedValues: "true, false", Description: "Place here or top of page"},
		{Name: "label", DefaultValue: "", AllowedValues: "", Description: "Table Label"},
		{Name: "location", DefaultValue: "above", AllowedValues: "above, below", Description: "Caption Location"},
		{Name: "mwe", DefaultValue: "false", AllowedValues: "true, false", Description: "Minimal working example"},
		{Name: "table-align", DefaultValue: "centering", AllowedValues: "centering, raggedleft, raggedright", Description: "Table Alignment"},
		{Name: "text-align", DefaultValue: "l", AllowedValues: "l, c, r", Description: "Text Alignment"},
	},
	"markdown": {
		{Name: "align", DefaultValue: "l", AllowedValues: "l, c, r", Description: "Text Alignment, columns seperate by comma"},
		{Name: "bold-header", DefaultValue: "false", AllowedValues: "true, false", Description: "Table Header Bold"},
		{Name: "bold-first-column", DefaultValue: "false", AllowedValues: "true, false", Description: "Bold first column"},
		{Name: "escape", DefaultValue: "false", AllowedValues: "true, false", Description: "Escape Markdown table"},
		{Name: "pretty", DefaultValue: "true", AllowedValues: "true, false", Description: "Pretty-print Markdown"},
	},
	"mediawiki": {
		{Name: "first-row-header", DefaultValue: "false", AllowedValues: "true, false", Description: "Use first row as headers"},
		{Name: "minify", DefaultValue: "false", AllowedValues: "true, false", Description: "Minify MediaWiki table"},
		{Name: "sort", DefaultValue: "false", AllowedValues: "true, false", Description: "Make table sortable in Wikipedia"},
	},
	"tmpl": {
		{Name: "template", DefaultValue: "", AllowedValues: "", Description: "Template file path"},
	},
	"sql": {
		{Name: "one-insert", DefaultValue: "false", AllowedValues: "true, false", Description: "Insert multiple rows at once"},
		{Name: "replace", DefaultValue: "false", AllowedValues: "true, false", Description: "Use REPLACE instead of INSERT"},
		{Name: "dialect", DefaultValue: "mysql", AllowedValues: "none, mysql, oracle, mssql, postgresql", Description: "identity escape SQL Dialect, none for no escape"},
		{Name: "table", DefaultValue: "", AllowedValues: "", Description: "Table Name"},
	},
	"xml": {
		{Name: "minify", DefaultValue: "false", AllowedValues: "true, false", Description: "Minify XML"},
		{Name: "root-element", DefaultValue: "dataset", AllowedValues: "string", Description: "Root Element Tag"},
		{Name: "row-element", DefaultValue: "record", AllowedValues: "string", Description: "Row Element Tag"},
		{Name: "declaration", DefaultValue: "true", AllowedValues: "true, false", Description: "Include XML Declaration"},
	},
}

// GlobalTransformParams are transformation parameters available for all formats
var GlobalTransformParams = []FormatParam{
	{Name: "transpose", DefaultValue: "false", AllowedValues: "true, false", Description: "Transpose the table (swap rows and columns)"},
	{Name: "delete-empty", DefaultValue: "false", AllowedValues: "true, false", Description: "Remove empty rows from the table"},
	{Name: "deduplicate", DefaultValue: "false", AllowedValues: "true, false", Description: "Remove duplicate rows"},
	{Name: "uppercase", DefaultValue: "false", AllowedValues: "true, false", Description: "Convert all text to UPPERCASE"},
	{Name: "lowercase", DefaultValue: "false", AllowedValues: "true, false", Description: "Convert all text to lowercase"},
	{Name: "capitalize", DefaultValue: "false", AllowedValues: "true, false", Description: "Capitalize the first letter of each cell"},
}

// GetFormatParams returns the parameters for a specific format
func GetFormatParams(format string) []FormatParam {
	return FormatParamsRegistry[format]
}

// GetAllFormats returns a list of all supported formats
func GetAllFormats() []string {
	formats := make([]string, 0, len(FormatParamsRegistry))
	for format := range FormatParamsRegistry {
		formats = append(formats, format)
	}
	return formats
}

// FormatExists checks if a format is supported
func FormatExists(format string) bool {
	_, exists := FormatParamsRegistry[format]
	return exists
}
