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

package config

import (
	"testing"

	_ "github.com/polocto/FolderFlow/internal/filter"
	_ "github.com/polocto/FolderFlow/internal/strategy"
	"gopkg.in/yaml.v3"
)

func TestConfigValidation(t *testing.T) {
	data := `
source_dirs: []
dest_dirs: []
`

	var cfg Config
	err := yaml.Unmarshal([]byte(data), &cfg)

	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}

func TestConfigValidation_OK(t *testing.T) {
	data := `
source_dirs:
  - /tmp
dest_dirs:
  - name: out
    path: /dest
    strategy:
      name: dirchain
`

	var cfg Config
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
