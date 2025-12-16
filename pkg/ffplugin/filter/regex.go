package filter

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"regexp"
)

// CustomFilter is an example custom filter.
type RegexFilter struct {
	Patterns   []string `yaml:"regex"`
	compiledRe []*regexp.Regexp
}

func (f *RegexFilter) Match(path string, info fs.FileInfo) (bool, error) {
	basename := filepath.Base(path)
	for i, re := range f.compiledRe {
		if re.MatchString(basename) {
			slog.Debug("Match found", "basename", basename, "pattern", f.Patterns[i])
			return true, nil
		}
	}
	slog.Debug("No match ", "basename", basename, "patterns", f.Patterns)
	return false, nil
}

func (f *RegexFilter) Selector() string {
	return "regex"
}

func (f *RegexFilter) LoadConfig(config map[string]interface{}) error {
	if patterns, ok := config["patterns"].([]string); ok {
		if len(patterns) == 0 {
			slog.Error("Patterns list is empty", "config", config)
			return fmt.Errorf("'patterns' config cannot be empty")
		}
		f.Patterns = patterns
	} else {
		slog.Error("Failed to load patterns", "config", config)
		return fmt.Errorf("invalid or missing 'patterns' config")
	}
	compiled := make([]*regexp.Regexp, len(f.Patterns))
	for i, pat := range f.Patterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			return fmt.Errorf("invalid patterns %q: %v", pat, err)
		}
		compiled[i] = re
	}
	f.compiledRe = compiled
	slog.Debug("Loading regex was successful", "patterns", f.Patterns)
	return nil
}

func init() {
	RegisterFilter("regex", func() Filter {
		return &RegexFilter{}
	})
}
