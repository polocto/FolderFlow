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

package strategy

import (
	"fmt"
	"io/fs"
	"log/slog"
)

// Strategy defines the interface for file organization strategies.
// FinalDirPath should ONLY compute the destination path and MUST NOT modify the filesystem.
type Strategy interface {
	// FinalDirPath computes the final directory path for a file based on the strategy.
	FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error)
	// Selector returns a unique identifier for the strategy (e.g., "date", "dirchain")
	Selector() string
	// LoadConfig allows the strategy to be configured from the YAML config
	LoadConfig(config map[string]interface{}) error
}

var strategyRegistry = make(map[string]func() Strategy)

func RegisterStrategy(name string, factory func() Strategy) {
	strategyRegistry[name] = factory
	slog.Debug("Create a new strategy", "name", name)
}

func NewStrategy(name string) (Strategy, error) {
	factory, ok := strategyRegistry[name]
	slog.Debug("Retrieving strategy", "name", name, "found", ok)
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
	return factory(), nil
}
