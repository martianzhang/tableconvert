package common

import (
	"fmt"
	"io"
	"os"
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
		default:
			cfg.Extension[k] = v
		}
	}

	// Check if required parameters are provided
	if cfg.From == "" || cfg.To == "" {
		return cfg, fmt.Errorf("must provide -f|--from and -t|--to parameters")
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
