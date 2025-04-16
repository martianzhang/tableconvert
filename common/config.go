package common

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	From      string
	To        string
	File      string
	Reader    io.Reader
	Result    string
	Writer    io.Writer
	Verbose   bool
	Extension map[string]string
}

func Usage() {
	fmt.Println(`Usage: tableconvert [OPTIONS]

Convert between different table formats (MySQL, Markdown, CSV, JSON, Excel, etc.)

Options:
  --from|-f={FORMAT}     Source format (e.g. mysql, csv, json, xlsx)
  --to|-t={FORMAT}       Target format (e.g. mysql, csv, json, xlsx)
  --file={PATH}          Input file path (or use stdin if not specified)
  --result|-r={PATH}     Output file path (or use stdout if not specified)
  --verbose|-v           Enable verbose output
  -h|--help              Show this help message

Examples:
  tableconvert --from=csv --to=json --file=input.csv --result=output.json
  cat input.csv | tableconvert --from=csv --to=json

Extension Arguments:
  For eash format there are many extension config, please refer to:
  https://github.com/martianzhang/tableconvert/blob/main/docs/arguments.md`)
}

// ParseConfig parses arguments in format "--key=value" or "--key value" and returns a key-value map
func ParseConfig(args []string) (Config, error) {
	var err error
	configs := make(map[string]string)
	for i, arg := range args {
		// Only process arguments starting with "--"
		if strings.HasPrefix(arg, "-") {
			// Remove prefix "--"
			trimmed := strings.TrimLeft(arg, "-")
			var key, value string
			// If contains "=", parse directly
			if idx := strings.Index(trimmed, "="); idx != -1 {
				key = trimmed[:idx]
				value = trimmed[idx+1:]
			} else {
				// Otherwise, check if next argument exists and doesn't start with "-", use as value if true
				key = trimmed
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					value = args[i+1]
					i++ // This will skip one element in next iteration
				}
			}
			configs[key] = value
		}
	}

	// Process known parameters
	var cfg = Config{Extension: make(map[string]string)}
	for k, v := range configs {
		switch k {
		case "from", "f":
			cfg.From = v
		case "to", "t":
			cfg.To = v
		case "file":
			cfg.File = v
		case "result", "r":
			cfg.Result = v
		case "verbose", "v":
			if v == "" || strings.ToLower(v) == "true" {
				cfg.Verbose = true
			} else {
				cfg.Verbose = false
			}
		case "h", "help":
			Usage()
			os.Exit(0)
		default:
			cfg.Extension[k] = v
		}
	}

	// Determine input target (Reader)
	if cfg.File != "" {
		// Check if file exists
		if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
			return cfg, fmt.Errorf("file does not exist: %s", cfg.File)
		}

		// Try to open file
		cfg.Reader, err = os.Open(cfg.File)
		if err != nil {
			return cfg, fmt.Errorf("failed to open file: %v", err)
		}
	} else {
		cfg.Reader = os.Stdin
	}

	// Auto detect `--from` and `--to` format
	if cfg.From == "" && cfg.File != "" {
		if ext := DetectTableFormatByExtension(cfg.File); ext != "" {
			cfg.From = ext
		} else {
			format, err := DetectTableFormatByData(cfg.Reader)
			if err != nil {
				cfg.From = format
			}
		}
	}
	if cfg.To == "" && cfg.Result != "" {
		if ext := DetectTableFormatByExtension(cfg.Result); ext != "" {
			cfg.To = ext
		}
	}

	// Check if required parameters are provided
	if cfg.From == "" || cfg.To == "" {
		return cfg, fmt.Errorf("must provide -f|--from and -t|--to parameters")
	}

	// Determine output destination (Writer)
	if cfg.Result != "" {
		file, err := os.Create(cfg.Result)
		if err != nil {
			return cfg, err
		}
		cfg.Writer = file // Set output to file
	} else {
		cfg.Writer = os.Stdout // Set output to standard output
	}

	return cfg, nil
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
