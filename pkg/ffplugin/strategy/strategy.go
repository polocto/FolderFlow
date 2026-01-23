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
	"log/slog"
	"path/filepath"
	"strings"
)

// Strategy defines the interface for file organization strategies.
// FinalDirPath should ONLY compute the destination path and MUST NOT modify the filesystem.
type Strategy interface {
	// FinalDirPath computes the final directory path for a file based on the strategy.
	FinalDirPath(ctx Context) (string, error)
	// Selector returns a unique identifier for the strategy (e.g., "date", "dirchain")
	Selector() string
	// LoadConfig allows the strategy to be configured from the YAML config
	LoadConfig(config map[string]interface{}) error
}

type safeStrategy struct {
	Strategy
}

func (s safeStrategy) FinalDirPath(ctx Context) (string, error) {
	if ctx == nil {
		return "", ErrContextIsNil
	}
	path, err := s.Strategy.FinalDirPath(ctx)
	if err != nil {
		return "", err
	}

	path = filepath.Clean(path)

	if !isSubDir(ctx.DstDir(), path) {
		return "", fmt.Errorf(
			"strategy %s returned path outside destination directory",
			s.Selector(),
		)
	}

	return path, nil
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
	return safeStrategy{factory()}, nil
}

func isSubDir(base, target string) bool {
	base = filepath.Clean(base)
	target = filepath.Clean(target)

	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}

	// Disallow:
	//   ".."
	//   "../something"
	// Allow:
	//   "."
	//   "sub"
	//   "sub/dir"
	return rel != ".." &&
		!strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
