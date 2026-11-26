package common

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

//go:embed usage.txt
var usageText string

func Usage() {
	fmt.Fprint(os.Stderr, usageText)
}

// ShowFormatsHelp displays help for all supported formats and their parameters
func ShowFormatsHelp() {
	fmt.Fprintln(os.Stderr, "Supported Formats and Their Parameters:")
	fmt.Fprintln(os.Stderr, "========================================")
	fmt.Fprintln(os.Stderr, "")

	// Get all formats in a consistent order (including aliases)
	formats := []string{"ascii", "csv", "tsv", "excel", "xlsx", "html", "json", "jsonl", "jsonlines", "latex", "markdown", "md", "mediawiki", "mysql", "sql", "tmpl", "template", "twiki", "tracwiki", "xml"}

	for _, format := range formats {
		params := GetFormatParams(format)
		if len(params) > 0 {
			fmt.Fprintf(os.Stderr, "## %s\n\n", format)

			// Calculate column widths for this format
			maxParamLen := len("Parameter")
			maxDefaultLen := len("Default")
			maxValuesLen := len("Allowed Values")
			maxDescLen := len("Description")

			for _, p := range params {
				if len(p.Name) > maxParamLen {
					maxParamLen = len(p.Name)
				}
				if len(p.DefaultValue) > maxDefaultLen {
					maxDefaultLen = len(p.DefaultValue)
				}
				if len(p.AllowedValues) > maxValuesLen {
					maxValuesLen = len(p.AllowedValues)
				}
				if len(p.Description) > maxDescLen {
					maxDescLen = len(p.Description)
				}
			}

			// Print header
			fmt.Fprintf(os.Stderr, "%-*s  %-*s  %-*s  %-*s\n",
				maxParamLen, "Parameter",
				maxDefaultLen, "Default",
				maxValuesLen, "Allowed Values",
				maxDescLen, "Description")

			// Print separator
			fmt.Fprintf(os.Stderr, "%s  %s  %s  %s\n",
				strings.Repeat("-", maxParamLen),
				strings.Repeat("-", maxDefaultLen),
				strings.Repeat("-", maxValuesLen),
				strings.Repeat("-", maxDescLen))

			// Print rows
			for _, param := range params {
				fmt.Fprintf(os.Stderr, "%-*s  %-*s  %-*s  %-*s\n",
					maxParamLen, param.Name,
					maxDefaultLen, param.DefaultValue,
					maxValuesLen, param.AllowedValues,
					maxDescLen, param.Description)
			}
			fmt.Fprintln(os.Stderr, "")
		}
	}

	// Show global transformation parameters
	fmt.Fprintln(os.Stderr, "## Global Transformation Parameters")
	fmt.Fprintln(os.Stderr, "These parameters work with all formats:")
	fmt.Fprintln(os.Stderr, "")

	// Calculate column widths for global params
	maxParamLen := len("Parameter")
	maxDefaultLen := len("Default")
	maxDescLen := len("Description")

	for _, p := range GlobalTransformParams {
		if len(p.Name) > maxParamLen {
			maxParamLen = len(p.Name)
		}
		if len(p.DefaultValue) > maxDefaultLen {
			maxDefaultLen = len(p.DefaultValue)
		}
		if len(p.Description) > maxDescLen {
			maxDescLen = len(p.Description)
		}
	}

	// Print header
	fmt.Fprintf(os.Stderr, "%-*s  %-*s  %-*s\n",
		maxParamLen, "Parameter",
		maxDefaultLen, "Default",
		maxDescLen, "Description")

	// Print separator
	fmt.Fprintf(os.Stderr, "%s  %s  %s\n",
		strings.Repeat("-", maxParamLen),
		strings.Repeat("-", maxDefaultLen),
		strings.Repeat("-", maxDescLen))

	// Print rows
	for _, param := range GlobalTransformParams {
		fmt.Fprintf(os.Stderr, "%-*s  %-*s  %-*s\n",
			maxParamLen, param.Name,
			maxDefaultLen, param.DefaultValue,
			maxDescLen, param.Description)
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage examples:")
	fmt.Fprintln(os.Stderr, "  tableconvert --from=csv --to=markdown --align=l,c,r --bold-header")
	fmt.Fprintln(os.Stderr, "  tableconvert --from=mysql --to=json --format=2d --minify")
}

// ShowFormatHelp displays help for a specific format
func ShowFormatHelp(format string) {
	// Normalize format name
	format = strings.ToLower(format)

	// Handle aliases
	switch format {
	case "md":
		format = "markdown"
	case "xlsx":
		format = "excel"
	case "jsonlines", "jsonl":
		format = "jsonl"
	case "tracwiki":
		format = "twiki"
	case "template":
		format = "tmpl"
	}

	params := GetFormatParams(format)

	if len(params) == 0 {
		fmt.Fprintf(os.Stderr, "Unknown format: %s\n\n", format)
		fmt.Fprintln(os.Stderr, "Use --help-formats to see all supported formats.")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Format-Specific Parameters for %s:\n", strings.ToUpper(format))
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 50))
	fmt.Fprintln(os.Stderr, "")

	// Calculate column widths for pretty alignment
	maxParamLen := len("Parameter")
	maxDefaultLen := len("Default")
	maxValuesLen := len("Allowed Values")
	maxDescLen := len("Description")

	for _, p := range params {
		if len(p.Name) > maxParamLen {
			maxParamLen = len(p.Name)
		}
		if len(p.DefaultValue) > maxDefaultLen {
			maxDefaultLen = len(p.DefaultValue)
		}
		if len(p.AllowedValues) > maxValuesLen {
			maxValuesLen = len(p.AllowedValues)
		}
		if len(p.Description) > maxDescLen {
			maxDescLen = len(p.Description)
		}
	}

	// Print header
	fmt.Fprintf(os.Stderr, "%-*s  %-*s  %-*s  %-*s\n",
		maxParamLen, "Parameter",
		maxDefaultLen, "Default",
		maxValuesLen, "Allowed Values",
		maxDescLen, "Description")

	// Print separator
	fmt.Fprintf(os.Stderr, "%s  %s  %s  %s\n",
		strings.Repeat("-", maxParamLen),
		strings.Repeat("-", maxDefaultLen),
		strings.Repeat("-", maxValuesLen),
		strings.Repeat("-", maxDescLen))

	// Print rows
	for _, param := range params {
		fmt.Fprintf(os.Stderr, "%-*s  %-*s  %-*s  %-*s\n",
			maxParamLen, param.Name,
			maxDefaultLen, param.DefaultValue,
			maxValuesLen, param.AllowedValues,
			maxDescLen, param.Description)
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Global Transformation Parameters (also available):")
	fmt.Fprintln(os.Stderr, "  --transpose, --delete-empty, --deduplicate, --uppercase, --lowercase, --capitalize")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "Usage Example:\n")
	fmt.Fprintf(os.Stderr, "  tableconvert --from=%s --to=csv", format)
	for i, param := range params {
		if i < 2 { // Show first 2 params as examples
			fmt.Fprintf(os.Stderr, " --%s=%s", param.Name, param.DefaultValue)
		}
	}
	fmt.Fprintln(os.Stderr, "")
}

// canonicalFlagName maps short/long aliases to one key so repeated aliases are
// resolved by argument order rather than Go map iteration order.
func canonicalFlagName(key string) string {
	switch key {
	case "f":
		return "from"
	case "t":
		return "to"
	case "file", "input", "i":
		return "file"
	case "result", "output", "o":
		return "result"
	case "v":
		return "verbose"
	case "b":
		return "batch"
	case "r":
		return "recursive"
	case "dir":
		return "output-dir"
	case "dryrun", "preview":
		return "dry-run"
	case "h":
		return "help"
	default:
		return key
	}
}

func isBooleanFlag(key string) bool {
	switch canonicalFlagName(key) {
	case "verbose", "recursive", "dry-run", "mcp", "help", "help-formats", "list-sheets", "allow-unknown":
		return true
	}
	for _, params := range FormatParamsRegistry {
		for _, param := range params {
			if param.Name == key {
				return param.DefaultValue == "true" || param.DefaultValue == "false" ||
					strings.Contains(param.AllowedValues, "true, false")
			}
		}
	}
	for _, param := range GlobalTransformParams {
		if param.Name == key {
			return param.DefaultValue == "true" || param.DefaultValue == "false"
		}
	}
	return false
}

func isRepeatableFlag(key string) bool {
	switch key {
	case "replace", "filter":
		return true
	default:
		return false
	}
}

func isBooleanLiteral(value string) bool {
	switch strings.ToLower(value) {
	case "true", "false", "yes", "no", "y", "n", "1", "0":
		return true
	default:
		return false
	}
}

// ParseConfig parses arguments in format "--key=value" or "--key value" and returns a key-value map
func ParseConfig(args []string) (Config, error) {
	configs := make(map[string]string)
	var positionalArgs []string

	// Use traditional for loop to allow skipping elements with i++
	for i := 0; i < len(args); i++ {
		arg := args[i]
		// Check if argument is a flag (starts with -)
		if strings.HasPrefix(arg, "-") {
			// Remove prefix "--" or "-"
			trimmed := strings.TrimLeft(arg, "-")
			var key, value string
			// If contains "=", parse directly
			if idx := strings.Index(trimmed, "="); idx != -1 {
				key = trimmed[:idx]
				value = trimmed[idx+1:]
			} else {
				// Flag-style: boolean flags do not consume a following
				// positional argument; non-boolean flags take the next token.
				key = trimmed
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					if !isBooleanFlag(key) || isBooleanLiteral(args[i+1]) {
						value = args[i+1]
						i++ // Skip the next element since we consumed it as the value
					}
				}
			}
			canonical := canonicalFlagName(key)
			if isRepeatableFlag(canonical) {
				if previous, ok := configs[canonical]; ok {
					configs[canonical] = previous + "\x00" + value
				} else {
					configs[canonical] = value
				}
				continue
			}
			configs[canonical] = value
		} else {
			// Positional argument (file path)
			positionalArgs = append(positionalArgs, arg)
		}
	}

	// Process known parameters. Keys are canonicalized above, so aliases such
	// as --from/-f and --file/-i cannot race each other in map iteration order.
	var cfg = Config{Extension: make(map[string]string)}
	for k, v := range configs {
		switch k {
		case "from":
			cfg.From = v
		case "to":
			cfg.To = v
		case "file":
			cfg.File = v
		case "result":
			cfg.Result = v
		case "verbose":
			parsed, err := parseBoolStrict(v, true)
			if err != nil {
				return cfg, fmt.Errorf("--verbose: %w", err)
			}
			cfg.Verbose = parsed
		case "batch":
			cfg.Batch = v
		case "recursive":
			parsed, err := parseBoolStrict(v, true)
			if err != nil {
				return cfg, fmt.Errorf("--recursive: %w", err)
			}
			cfg.Recursive = parsed
		case "output-dir":
			cfg.OutputDir = v
		case "dry-run":
			parsed, err := parseBoolStrict(v, true)
			if err != nil {
				return cfg, fmt.Errorf("--dry-run: %w", err)
			}
			cfg.DryRun = parsed
		case "help":
			Usage()
			os.Exit(0)
		case "help-formats":
			ShowFormatsHelp()
			os.Exit(0)
		case "help-format":
			ShowFormatHelp(v)
			os.Exit(0)
		case "mcp":
			parsed, err := parseBoolStrict(v, true)
			if err != nil {
				return cfg, fmt.Errorf("--mcp: %w", err)
			}
			cfg.MCPMode = parsed
		case "list-sheets":
			parsed, err := parseBoolStrict(v, true)
			if err != nil {
				return cfg, fmt.Errorf("--list-sheets: %w", err)
			}
			cfg.ListSheets = parsed
		case "allow-unknown":
			parsed, err := parseBoolStrict(v, true)
			if err != nil {
				return cfg, fmt.Errorf("--allow-unknown: %w", err)
			}
			cfg.AllowUnknown = parsed
		default:
			cfg.Extension[k] = v
		}
	}

	// Handle positional arguments (auto-detect input/output files)
	// Pattern: tableconvert input output
	if len(positionalArgs) > 0 {
		if cfg.File == "" {
			cfg.File = positionalArgs[0]
		}
		if len(positionalArgs) > 1 && cfg.Result == "" {
			cfg.Result = positionalArgs[1]
		}
	}

	// If MCP mode is enabled, skip from/to validation
	if cfg.MCPMode {
		return cfg, nil
	}

	// Batch mode handling
	if cfg.Batch != "" {
		// Batch mode: skip single file validation
		// Validate that to format is provided
		if cfg.To == "" {
			return cfg, fmt.Errorf("batch mode requires --to format\n\nExample:\n  tableconvert --batch=\"*.csv\" --to=json")
		}
		if !FormatExists(cfg.To) {
			return cfg, fmt.Errorf("unsupported output format: %s\n\nSupported formats: %v", cfg.To, getSupportedFormats())
		}
		// Batch mode doesn't use Reader/Writer for single file
		return cfg, nil
	}

	// Auto-detect formats from file extensions if not specified
	if cfg.From == "" && cfg.File != "" {
		cfg.From = detectFormatFromExtension(cfg.File)
		if cfg.From != "" && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "# Auto-detected input format: %s\n", cfg.From)
		}
	}
	if cfg.To == "" && cfg.Result != "" {
		cfg.To = detectFormatFromExtension(cfg.Result)
		if cfg.To != "" && cfg.Verbose {
			fmt.Fprintf(os.Stderr, "# Auto-detected output format: %s\n", cfg.To)
		}
	}

	if cfg.ListSheets {
		if cfg.From != "excel" && cfg.From != "xlsx" {
			return cfg, fmt.Errorf("--list-sheets requires an Excel input (.xlsx)")
		}
		if cfg.File == "" {
			return cfg, fmt.Errorf("--list-sheets requires an input file")
		}
		return cfg, nil
	}

	// If the extension is unknown (or input is stdin), inspect the content.
	if cfg.From == "" {
		if cfg.File == "" && stdinIsTerminal() {
			return cfg, newParseError(cfg.From, cfg.To, cfg.File, cfg.Result)
		}
		data, err := readInputForDetection(&cfg)
		if err != nil {
			return cfg, err
		}
		detected, err := DetectTableFormatByData(bytes.NewReader(data))
		if err != nil {
			if cfg.File == "" && cfg.To == "" {
				return cfg, newParseError(cfg.From, cfg.To, cfg.File, cfg.Result)
			}
			return cfg, fmt.Errorf("cannot auto-detect input format%s: %w\nUse --from=<format> to specify it explicitly", inputDescription(cfg.File), err)
		}
		cfg.From = detected
		cfg.BufferedInput = data
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "# Auto-detected input format from content: %s\n", detected)
		}
	}

	// Check if required parameters are provided
	if cfg.From == "" || cfg.To == "" {
		return cfg, newParseError(cfg.From, cfg.To, cfg.File, cfg.Result)
	}

	// Validate formats are supported
	if !FormatExists(cfg.From) {
		return cfg, fmt.Errorf("unsupported input format: %s\n\nSupported formats: %v\nRun 'tableconvert --help-formats' for details", cfg.From, getSupportedFormats())
	}
	if !FormatExists(cfg.To) {
		return cfg, fmt.Errorf("unsupported output format: %s\n\nSupported formats: %v\nRun 'tableconvert --help-formats' for details", cfg.To, getSupportedFormats())
	}

	// Determine input target (Reader)
	if cfg.BufferedInput != nil {
		cfg.Reader = bytes.NewReader(cfg.BufferedInput)
	} else if cfg.File != "" {
		// Check if file exists
		if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
			return cfg, fmt.Errorf("file does not exist: %s", cfg.File)
		}
		// Open file for reading
		file, err := os.Open(cfg.File)
		if err != nil {
			return cfg, fmt.Errorf("failed to open file: %w", err)
		}
		cfg.Reader = file
		cfg.InputFile = file
	} else {
		// Use stdin
		cfg.Reader = os.Stdin
	}

	// Determine output destination (Writer)
	// Skip file creation in dry-run mode
	if cfg.Result != "" {
		if cfg.DryRun {
			// In dry-run mode, use a discard writer since we won't write output
			cfg.Writer = io.Discard
		} else {
			// Write to a temporary file in the destination directory first. The
			// final file is only replaced after a successful conversion, so a
			// failure (or input == output) cannot destroy existing data.
			dir := filepath.Dir(cfg.Result)
			base := filepath.Base(cfg.Result)
			tmp, err := os.CreateTemp(dir, "."+base+".tmp-*")
			if err != nil {
				_ = cfg.CloseInput()
				return cfg, err
			}
			if err := tmp.Chmod(0644); err != nil {
				tmp.Close()
				os.Remove(tmp.Name())
				_ = cfg.CloseInput()
				return cfg, err
			}
			cfg.Writer = tmp
			cfg.OutputFile = tmp
			cfg.TempResult = tmp.Name()
		}
	} else {
		cfg.Writer = os.Stdout
	}

	return cfg, nil
}

func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func readInputForDetection(cfg *Config) ([]byte, error) {
	if cfg.File != "" {
		if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", cfg.File)
		}
		data, err := os.ReadFile(cfg.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read input file: %w", err)
		}
		return data, nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %w", err)
	}
	return data, nil
}

func inputDescription(path string) string {
	if path == "" {
		return " from stdin"
	}
	return " for " + path
}

// detectFormatFromExtension detects format from file extension
// Returns empty string if extension is not recognized
func detectFormatFromExtension(filename string) string {
	// .txt is too common to auto-detect.
	if strings.EqualFold(filepath.Ext(filename), ".txt") {
		return ""
	}
	return DetectTableFormatByExtension(filename)
}

// getSupportedFormats returns a sorted list of supported formats
func getSupportedFormats() []string {
	formats := []string{"ascii", "csv", "tsv", "excel", "html", "json", "jsonl", "latex", "markdown", "mediawiki", "mysql", "sql", "tmpl", "twiki", "xml", "xlsx", "jsonlines", "md", "template", "tracwiki"}
	return formats
}

// newParseError creates a helpful error message for missing required parameters
func newParseError(from, to, file, result string) error {
	var msg strings.Builder

	msg.WriteString("Missing required parameters: --from and --to\n\n")

	// Show what was provided
	if from == "" && to == "" {
		msg.WriteString("You need to specify both input and output formats.\n\n")
	} else if from == "" {
		msg.WriteString("Input format (--from) is missing.\n\n")
	} else if to == "" {
		msg.WriteString("Output format (--to) is missing.\n\n")
	}

	// Show auto-detection status
	if file != "" && from == "" {
		detected := detectFormatFromExtension(file)
		if detected != "" {
			msg.WriteString(fmt.Sprintf("Tip: Input file '%s' could auto-detect format: %s\n", file, detected))
			msg.WriteString("     Use: --from=" + detected + " or just rename your file\n\n")
		}
	}
	if result != "" && to == "" {
		detected := detectFormatFromExtension(result)
		if detected != "" {
			msg.WriteString(fmt.Sprintf("Tip: Output file '%s' could auto-detect format: %s\n", result, detected))
			msg.WriteString("     Use: --to=" + detected + " or just rename your file\n\n")
		}
	}

	// Show concise usage examples
	msg.WriteString("Quick examples:\n")
	msg.WriteString("  tableconvert --from=csv --to=json input.csv output.json\n")
	msg.WriteString("  tableconvert input.csv output.json  (auto-detect)\n")
	msg.WriteString("  tableconvert -i input.csv -o output.json  (short flags)\n")
	msg.WriteString("  tableconvert --batch=\"*.csv\" --to=json  (batch mode)\n\n")
	msg.WriteString("Run 'tableconvert --help' for full usage.\n")

	return fmt.Errorf("%s", msg.String())
}

// parseBoolStrict parses a boolean flag value and rejects unknown values
// instead of silently treating them as false.
func parseBoolStrict(value string, defaultValue bool) (bool, error) {
	if value == "" {
		return defaultValue, nil
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "yes", "y", "1":
		return true, nil
	case "false", "no", "n", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value %q (use true/false, yes/no, or 1/0)", value)
	}
}

// Config holds configuration, reader/writer, and extension parameters
type Config struct {
	From          string
	To            string
	File          string
	Result        string
	Verbose       bool
	MCPMode       bool
	Batch         string // Batch mode pattern (e.g., "*.csv")
	Recursive     bool   // Recursive directory traversal
	OutputDir     string // Output directory for batch mode
	DryRun        bool   // Dry run mode - preview without writing
	ListSheets    bool   // List Excel worksheet names and exit
	AllowUnknown  bool   // Do not reject unknown extension options
	Reader        io.Reader
	Writer        io.Writer
	InputFile     *os.File
	BufferedInput []byte
	Extension     map[string]string
	TempResult    string // Temporary output path used until CommitOutput
	SQLNull       bool   // Source represents SQL NULL as the string "NULL"

	OutputTransform io.WriteCloser
	OutputFile      *os.File
	ioWrapped       bool
}

// EnsureTextIO applies --encoding / --output-encoding once. Binary formats
// such as Excel are left untouched.
func (c *Config) EnsureTextIO() error {
	if c.ioWrapped {
		return nil
	}
	if c.Reader != nil && IsTextFormat(c.From) {
		reader, err := DecodeReader(c.Reader, c.GetExtensionString("encoding", ""))
		if err != nil {
			return err
		}
		c.Reader = reader
	}
	if c.Writer != nil && IsTextFormat(c.To) {
		writer, err := EncodeWriter(c.Writer, c.GetExtensionString("output-encoding", ""))
		if err != nil {
			return err
		}
		c.OutputTransform = writer
		c.Writer = writer
	}
	c.ioWrapped = true
	return nil
}

func (c *Config) closeOutputTransform() {
	if c.OutputTransform == nil {
		return
	}
	_ = c.OutputTransform.Close()
	c.OutputTransform = nil
}

// CloseInput closes the input file opened by ParseConfig, if any.
func (c *Config) CloseInput() error {
	if c.InputFile == nil {
		return nil
	}
	err := c.InputFile.Close()
	c.InputFile = nil
	return err
}

// CommitOutput closes the temporary output file and atomically moves it to the
// final Result path. It is a no-op when no temporary output is active.
func (c *Config) CommitOutput() error {
	if c.TempResult == "" {
		c.closeOutputTransform()
		return nil
	}
	c.closeOutputTransform()
	if c.OutputFile != nil {
		if err := c.OutputFile.Close(); err != nil {
			c.OutputFile = nil
			c.AbortOutput()
			return err
		}
		c.OutputFile = nil
	} else if f, ok := c.Writer.(*os.File); ok {
		if err := f.Close(); err != nil {
			c.AbortOutput()
			return err
		}
	}
	if err := os.Rename(c.TempResult, c.Result); err != nil {
		c.AbortOutput()
		return err
	}
	c.TempResult = ""
	return nil
}

// AbortOutput closes and removes the temporary output file, leaving any
// existing destination file untouched.
func (c *Config) AbortOutput() {
	c.closeOutputTransform()
	if c.TempResult == "" {
		return
	}
	if c.OutputFile != nil {
		_ = c.OutputFile.Close()
		c.OutputFile = nil
	} else if f, ok := c.Writer.(*os.File); ok {
		_ = f.Close()
	}
	_ = os.Remove(c.TempResult)
	c.TempResult = ""
}

// GetExtensionBool gets a boolean value from Extension with default
func (c *Config) GetExtensionBool(key string, defaultValue bool) bool {
	if c.Extension == nil {
		return defaultValue
	}
	val, ok := c.Extension[key]
	if !ok {
		return defaultValue
	}
	switch strings.TrimSpace(strings.ToLower(val)) {
	case "true", "yes", "y", "1", "":
		return true
	case "false", "no", "n", "0":
		return false
	default:
		return defaultValue
	}
}

// GetParseTypes reports whether output values should be converted to JSON
// types. --parse-types is the preferred name; --parsing-json is kept as an
// alias for compatibility.
func (c *Config) GetParseTypes() bool {
	if c.Extension == nil {
		return false
	}
	if _, ok := c.Extension["parse-types"]; ok {
		return c.GetExtensionBool("parse-types", false)
	}
	return c.GetExtensionBool("parsing-json", false)
}

// GetExtensionList returns repeated option values in command-line order.
func (c *Config) GetExtensionList(key string) []string {
	if c.Extension == nil {
		return nil
	}
	raw, ok := c.Extension[key]
	if !ok || raw == "" {
		return nil
	}
	parts := strings.Split(raw, "\x00")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}

// GetExtensionString gets a string value from Extension with default
func (c *Config) GetExtensionString(key string, defaultValue string) string {
	if c.Extension == nil {
		return defaultValue
	}
	if val, ok := c.Extension[key]; ok {
		return val
	}
	return defaultValue
}

// GetExtensionInt gets an int value from Extension with default
func (c *Config) GetExtensionInt(key string, defaultValue int) int {
	if c.Extension == nil {
		return defaultValue
	}
	if val, ok := c.Extension[key]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}

func GetProjectRootPath() (string, error) {
	// Get project root path
	rootPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		// Check if go.mod exists in current directory
		goModPath := filepath.Join(rootPath, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		// Move up one directory
		parent := filepath.Dir(rootPath)
		if parent == rootPath {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		rootPath = parent
	}

	return strings.TrimSuffix(rootPath, "/") + "/", nil
}

// ApplyTransformations applies data transformations in the documented order.
func (c *Config) ApplyTransformations(table *Table) error {
	if table == nil {
		return nil
	}

	if c.GetExtensionBool("transpose", false) {
		Transpose(table)
	}
	if c.GetExtensionBool("delete-empty", false) {
		DeleteEmptyRows(table)
	}
	if c.GetExtensionBool("deduplicate", false) {
		DeduplicateRows(table)
	}
	if specs := c.GetExtensionList("replace"); len(specs) > 0 {
		if err := ReplaceCells(table, specs, c.GetExtensionString("replace-column", "")); err != nil {
			return err
		}
	}
	if specs := c.GetExtensionList("filter"); len(specs) > 0 {
		if err := FilterRows(table, specs); err != nil {
			return err
		}
	}
	if spec := c.GetExtensionString("rename", ""); spec != "" {
		if err := RenameHeaders(table, spec); err != nil {
			return err
		}
	}
	if spec := c.GetExtensionString("columns", ""); spec != "" {
		if err := SelectColumns(table, spec); err != nil {
			return err
		}
	}
	if spec := c.GetExtensionString("sort", ""); spec != "" {
		if err := SortRows(table, spec); err != nil {
			return err
		}
	}

	head, headSet, err := c.optionalNonNegativeInt("head")
	if err != nil {
		return err
	}
	tail, tailSet, err := c.optionalNonNegativeInt("tail")
	if err != nil {
		return err
	}
	if headSet || tailSet {
		headValue, tailValue := -1, -1
		if headSet {
			headValue = head
		}
		if tailSet {
			tailValue = tail
		}
		LimitRows(table, headValue, tailValue)
	}
	if c.GetExtensionBool("trim", false) {
		TrimCells(table)
	}

	// Case transformations are mutually exclusive; priority: uppercase > lowercase > capitalize.
	if c.GetExtensionBool("uppercase", false) {
		Uppercase(table)
	} else if c.GetExtensionBool("lowercase", false) {
		Lowercase(table)
	} else if c.GetExtensionBool("capitalize", false) {
		Capitalize(table)
	}
	return nil
}

// BatchFile represents a single file in batch processing
type BatchFile struct {
	InputPath  string
	OutputPath string
	FromFormat string
	ToFormat   string
}

// GetBatchFiles expands the batch pattern and returns a list of files to process
func (c *Config) GetBatchFiles() ([]BatchFile, error) {
	var files []BatchFile

	// Determine base directory and pattern
	pattern := c.Batch
	var baseDir string
	var filePattern string

	// Handle glob patterns properly - don't use filepath.Dir/Base for patterns with **
	// because they treat ** as a literal directory name
	if strings.Contains(pattern, "**") {
		// For recursive patterns like "test_batch/**/*.csv"
		// Find the position of ** and split there
		idx := strings.Index(pattern, "**")
		// Everything before ** is the base directory (up to the last /)
		beforeStar := pattern[:idx]
		if strings.Contains(beforeStar, "/") {
			lastSlash := strings.LastIndex(beforeStar, "/")
			baseDir = beforeStar[:lastSlash]
			filePattern = pattern[lastSlash+1:]
		} else {
			baseDir = "."
			filePattern = pattern
		}
		// Normalize baseDir
		if baseDir == "" {
			baseDir = "."
		}
	} else if strings.Contains(pattern, "/") || strings.Contains(pattern, "\\") {
		// Simple path without ** - use standard filepath functions
		baseDir = filepath.Dir(pattern)
		filePattern = filepath.Base(pattern)
	} else {
		baseDir = "."
		filePattern = pattern
	}

	// Handle recursive mode. A ** pattern is inherently recursive even when
	// --recursive is omitted.
	var matches []string
	var err error
	if c.Recursive || strings.Contains(pattern, "**") {
		// Use Walk to find all matching files recursively
		matches, err = c.walkPattern(baseDir, filePattern)
	} else {
		// Use Glob for single directory
		fullPattern := filepath.Join(baseDir, filePattern)
		matches, err = filepath.Glob(fullPattern)
	}

	if err != nil {
		return nil, fmt.Errorf("error expanding pattern: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files found matching pattern: %s", pattern)
	}

	// Sort matches for consistent order
	sort.Strings(matches)

	// Build BatchFile list
	for _, inputPath := range matches {
		// Skip directories
		info, err := os.Stat(inputPath)
		if err != nil || info.IsDir() {
			continue
		}

		// Determine output path
		var outputPath string
		if c.OutputDir != "" {
			// Use specified output directory
			outputPath = filepath.Join(c.OutputDir, filepath.Base(inputPath))
		} else {
			// Use same directory as input
			outputPath = filepath.Join(filepath.Dir(inputPath), filepath.Base(inputPath))
		}

		// Change extension to match output format
		outputPath = c.changeExtension(outputPath, c.To)

		// Auto-detect input format from extension unless the user explicitly
		// provided --from, which must take precedence.
		fromFormat := c.From
		if fromFormat == "" {
			fromFormat = detectFormatFromExtension(inputPath)
		}
		if fromFormat == "" {
			return nil, fmt.Errorf("cannot detect format from input file: %s (use --from to specify)", inputPath)
		}

		files = append(files, BatchFile{
			InputPath:  inputPath,
			OutputPath: outputPath,
			FromFormat: fromFormat,
			ToFormat:   c.To,
		})
	}

	// Reject output collisions before touching any file. Two inputs mapping to
	// the same output would silently overwrite one result, and an output that
	// matches any input (including its own) would overwrite source data.
	inputs := make(map[string]string, len(files))
	for _, file := range files {
		inputs[canonicalPath(file.InputPath)] = file.InputPath
	}
	outputs := make(map[string]string, len(files))
	for _, file := range files {
		outKey := canonicalPath(file.OutputPath)
		if prev, ok := outputs[outKey]; ok {
			return nil, fmt.Errorf("multiple inputs map to the same output %s: %s and %s", file.OutputPath, prev, file.InputPath)
		}
		outputs[outKey] = file.InputPath
		if inputPath, ok := inputs[outKey]; ok {
			return nil, fmt.Errorf("output %s would overwrite input %s", file.OutputPath, inputPath)
		}
	}

	return files, nil
}

func canonicalPath(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	path = filepath.Clean(path)
	if runtime.GOOS == "windows" {
		path = strings.ToLower(path)
	}
	return path
}

// walkPattern recursively finds files matching the pattern
// Supports ** for recursive directory matching
func (c *Config) walkPattern(baseDir, pattern string) ([]string, error) {
	var matches []string

	// Handle ** in pattern by converting to recursive matching
	// Pattern "**/*.csv" means "match all .csv files in all subdirectories (including base)"
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden directories unless explicitly in pattern
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				// Check if pattern explicitly includes this
				if !strings.Contains(pattern, info.Name()) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Get relative path from baseDir
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		// Convert to forward slashes for consistent pattern matching
		relPath = filepath.ToSlash(relPath)

		// Check if path matches pattern. ** matches zero or more path segments.
		if matchRecursivePattern(filepath.ToSlash(pattern), relPath) {
			matches = append(matches, path)
		}
		return nil
	})

	return matches, err
}

func matchRecursivePattern(pattern, relPath string) bool {
	if !strings.Contains(pattern, "/") {
		matched, err := filepath.Match(pattern, filepath.Base(relPath))
		return err == nil && matched
	}

	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")
	pathSegments := strings.Split(strings.Trim(relPath, "/"), "/")
	return matchPathSegments(patternSegments, pathSegments)
}

func matchPathSegments(pattern, path []string) bool {
	if len(pattern) == 0 {
		return len(path) == 0
	}
	if pattern[0] == "**" {
		for i := 0; i <= len(path); i++ {
			if matchPathSegments(pattern[1:], path[i:]) {
				return true
			}
		}
		return false
	}
	if len(path) == 0 {
		return false
	}
	matched, err := filepath.Match(pattern[0], path[0])
	if err != nil || !matched {
		return false
	}
	return matchPathSegments(pattern[1:], path[1:])
}

// changeExtension changes the file extension to match the target format
func (c *Config) changeExtension(path, format string) string {
	ext := ""
	switch format {
	case "csv":
		ext = ".csv"
	case "tsv":
		ext = ".tsv"
	case "json":
		ext = ".json"
	case "jsonl", "jsonlines":
		ext = ".jsonl"
	case "markdown", "md":
		ext = ".md"
	case "html":
		ext = ".html"
	case "xml":
		ext = ".xml"
	case "sql":
		ext = ".sql"
	case "latex":
		ext = ".tex"
	case "excel", "xlsx":
		ext = ".xlsx"
	case "mysql":
		ext = ".txt"
	case "mediawiki":
		ext = ".wiki"
	case "twiki", "tracwiki":
		ext = ".twiki"
	case "tmpl", "template":
		ext = ".tmpl"
	case "ascii":
		ext = ".txt"
	default:
		// Keep original extension
		return path
	}

	// Replace extension
	return strings.TrimSuffix(path, filepath.Ext(path)) + ext
}
