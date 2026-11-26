package common

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var bordersPattern = regexp.MustCompile(`^[01]{4},[01]{4}$`)

// Validate checks option names and enum values after --from/--to are known.
// Unknown options are errors by default so typos cannot silently produce
// plausible-looking but wrong output. --allow-unknown is an escape hatch for
// experiments and forward compatibility.
func (c *Config) Validate() error {
	if c == nil {
		return nil
	}
	known := make(map[string]bool)
	for _, param := range GlobalTransformParams {
		known[param.Name] = true
	}
	for _, format := range []string{c.From, c.To} {
		for _, param := range FormatParamsRegistry[format] {
			known[param.Name] = true
		}
	}
	known["encoding"] = true
	known["output-encoding"] = true

	for key, value := range c.Extension {
		if !known[key] {
			if c.AllowUnknown {
				continue
			}
			return unknownOptionError(key, known)
		}
		if isBooleanParam(key) {
			if value != "" && !isBooleanLiteral(value) {
				return fmt.Errorf("invalid boolean value %q for --%s (use true/false, yes/no, or 1/0)", value, key)
			}
			continue
		}
		if err := validateOptionValue(key, value); err != nil {
			return err
		}
	}
	return nil
}

func unknownOptionError(key string, known map[string]bool) error {
	message := fmt.Sprintf("unknown option --%s", key)
	if suggestion := closestOption(key, known); suggestion != "" {
		message += fmt.Sprintf(" (did you mean --%s?)", suggestion)
	}
	return fmt.Errorf("%s\nRun 'tableconvert --help-formats' for supported options", message)
}

func isBooleanParam(key string) bool {
	for _, param := range GlobalTransformParams {
		if param.Name == key {
			return isBooleanParamSpec(param)
		}
	}
	for _, params := range FormatParamsRegistry {
		for _, param := range params {
			if param.Name == key {
				return isBooleanParamSpec(param)
			}
		}
	}
	return false
}

func isBooleanParamSpec(param FormatParam) bool {
	return param.DefaultValue == "true" || param.DefaultValue == "false" ||
		strings.Contains(param.AllowedValues, "true, false")
}

func validateOptionValue(key, value string) error {
	value = strings.TrimSpace(value)
	switch key {
	case "delimiter":
		if value == "" {
			return nil
		}
		switch strings.ToUpper(value) {
		case "COMMA", "TAB", "SEMICOLON", "PIPE", "SLASH", "HASH":
			return nil
		}
		switch value {
		case ",", "\t", ";", "|", "/", "#":
			return nil
		}
		return fmt.Errorf("invalid --delimiter %q (use COMMA, TAB, SEMICOLON, PIPE, SLASH, HASH, or the literal character)", value)

	case "format":
		if value == "" {
			return nil
		}
		switch strings.ToLower(value) {
		case "object", "2d", "column", "column-array", "keyed":
			return nil
		}
		return fmt.Errorf("invalid --format %q (use object, 2d, column, or keyed)", value)

	case "style":
		if value == "" {
			return nil
		}
		switch strings.ToLower(value) {
		case "box", "plus", "dot", "bubble":
			return nil
		}
		if utf8.RuneCountInString(value) == 1 {
			return nil
		}
		return fmt.Errorf("invalid --style %q (use box, plus, dot, bubble, or a single border character)", value)

	case "dialect":
		if value == "" {
			return nil
		}
		switch strings.ToLower(value) {
		case "none", "mysql", "oracle", "mssql", "postgres", "postgresql":
			return nil
		}
		return fmt.Errorf("invalid --dialect %q (use none, mysql, oracle, mssql, or postgresql)", value)

	case "align":
		if value == "" {
			return nil
		}
		for _, part := range strings.Split(value, ",") {
			switch strings.ToLower(strings.TrimSpace(part)) {
			case "l", "c", "r":
			default:
				return fmt.Errorf("invalid --align %q (use l, c, r, separated by commas)", value)
			}
		}
		return nil

	case "text-align":
		if value == "" {
			return nil
		}
		switch strings.ToLower(value) {
		case "l", "c", "r":
			return nil
		}
		return fmt.Errorf("invalid --text-align %q (use l, c, or r)", value)

	case "table-align":
		if value == "" {
			return nil
		}
		switch strings.ToLower(value) {
		case "centering", "raggedleft", "raggedright":
			return nil
		}
		return fmt.Errorf("invalid --table-align %q (use centering, raggedleft, or raggedright)", value)

	case "location":
		if value == "" {
			return nil
		}
		switch strings.ToLower(value) {
		case "above", "below":
			return nil
		}
		return fmt.Errorf("invalid --location %q (use above or below)", value)

	case "borders":
		if value == "" || bordersPattern.MatchString(value) {
			return nil
		}
		return fmt.Errorf("invalid --borders %q (use 4 header bits and 4 row bits, e.g. 1111,1111)", value)

	case "flatten-arrays":
		if value == "" {
			return nil
		}
		switch strings.ToLower(value) {
		case "json", "index":
			return nil
		}
		return fmt.Errorf("invalid --flatten-arrays %q (use json or index)", value)

	case "head", "tail", "skip-rows", "header-row":
		if value == "" {
			return nil
		}
		number, err := strconv.Atoi(value)
		if err != nil || number < 0 {
			return fmt.Errorf("invalid --%s %q (use a non-negative integer)", key, value)
		}
		return nil

	case "table-index":
		if value == "" {
			return nil
		}
		number, err := strconv.Atoi(value)
		if err != nil || number < 1 {
			return fmt.Errorf("invalid --table-index %q (use a 1-based positive number)", value)
		}
		return nil

	case "url":
		if value == "" {
			return nil
		}
		parsed, err := url.ParseRequestURI(value)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return fmt.Errorf("invalid --url %q (use an http(s) URL)", value)
		}
		return nil

	case "sheet-index":
		if value == "" {
			return nil
		}
		index, err := strconv.Atoi(value)
		if err != nil || index < 1 {
			return fmt.Errorf("invalid --sheet-index %q (use a 1-based positive number)", value)
		}
		return nil

	case "encoding":
		return validateInputEncoding(value)
	case "output-encoding":
		return validateOutputEncoding(value)
	}
	return nil
}

func validateInputEncoding(value string) error {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "utf-8", "utf8", "utf-8-bom", "utf-16", "utf-16le", "utf-16be",
		"gbk", "gb18030", "big5", "shift_jis", "shift-jis", "sjis", "latin1", "iso-8859-1", "auto":
		return nil
	default:
		return fmt.Errorf("invalid --encoding %q", value)
	}
}

func validateOutputEncoding(value string) error {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "utf-8", "utf8", "utf-8-bom", "utf-16", "utf-16le", "utf-16be",
		"gbk", "gb18030", "big5", "shift_jis", "shift-jis", "sjis", "latin1", "iso-8859-1":
		return nil
	case "auto":
		return fmt.Errorf("invalid --output-encoding %q (auto is only valid for input)", value)
	default:
		return fmt.Errorf("invalid --output-encoding %q", value)
	}
}

func closestOption(key string, known map[string]bool) string {
	best := ""
	bestDistance := 3 // only suggest distances 0-2 (i.e. likely typos)
	for candidate := range known {
		distance := editDistance(key, candidate)
		if distance < bestDistance {
			best = candidate
			bestDistance = distance
		}
	}
	return best
}

func editDistance(a, b string) int {
	if a == b {
		return 0
	}
	previous := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current := make([]int, len(b)+1)
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			current[j] = minInt(current[j-1]+1, previous[j]+1, previous[j-1]+cost)
		}
		previous = current
	}
	return previous[len(b)]
}

func minInt(values ...int) int {
	result := values[0]
	for _, value := range values[1:] {
		if value < result {
			result = value
		}
	}
	return result
}
