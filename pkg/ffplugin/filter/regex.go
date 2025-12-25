package filter

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

// CustomFilter is an example custom filter.
type RegexFilter struct {
	Patterns   []string `yaml:"patterns"`
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
	var cfg struct {
		Patterns []string `yaml:"patterns"`
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	if len(cfg.Patterns) == 0 {
		return fmt.Errorf("'patterns' config cannot be empty")
	}

	compiled := make([]*regexp.Regexp, len(cfg.Patterns))
	for i, pat := range cfg.Patterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			return fmt.Errorf("invalid pattern %q: %w", pat, err)
		}
		compiled[i] = re
	}

	f.Patterns = cfg.Patterns
	f.compiledRe = compiled

	slog.Debug("Loading regex was successful", "patterns", f.Patterns)
	return nil
}

func init() {
	RegisterFilter("regex", func() Filter {
		return &RegexFilter{}
	})
}
