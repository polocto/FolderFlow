// Copyright (c) 2026 Paul Sade.
//
// This file is part of the FoderFlow project.
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
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CustomFilter is an example custom filter.
type ExtensionFilter struct {
	Extensions []string `yaml:"extensions"`
}

func (f *ExtensionFilter) Match(path string, info fs.FileInfo) (bool, error) {
	ext := strings.ToLower(filepath.Ext(path))
	for _, allowedExt := range f.Extensions {
		if ext == strings.ToLower(allowedExt) {
			slog.Debug("Match found", "file's path", path, "matched extension", allowedExt)
			return true, nil
		}
	}
	return false, nil
}

func (f *ExtensionFilter) Selector() string {
	return "extensions"
}

func (f *ExtensionFilter) LoadConfig(config map[string]interface{}) error {
	var cfg struct {
		Extensions []string `yaml:"extensions"`
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	if len(cfg.Extensions) == 0 {
		return fmt.Errorf("invalid or missing 'extensions' config")
	}

	f.Extensions = cfg.Extensions

	slog.Debug("Loading extensions was successful", "extensions", f.Extensions)
	return nil
}

func init() {
	RegisterFilter("extensions", func() Filter {
		return &ExtensionFilter{}
	})
}
