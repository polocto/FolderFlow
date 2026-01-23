// Copyright (c) 2026 Paul Sade.
//
// This file is part of the FolderFlow project.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3,
// as published by the Free Software Foundation (see the LICENSE file).
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU General Public License for more details.

package filter

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"gopkg.in/yaml.v3"
)

// CustomFilter is an example custom filter.
type RegexFilter struct {
	Patterns   []string `yaml:"patterns"`
	compiledRe []*regexp.Regexp
}

func (f *RegexFilter) Match(ctx filter.Context) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context is nil")
	}
	basename := ctx.Info().Name()
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
	filter.RegisterFilter("regex", func() filter.Filter {
		return &RegexFilter{}
	})
}
