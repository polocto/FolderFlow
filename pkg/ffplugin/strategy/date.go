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
	"log/slog"
	"path/filepath"
)

// plugin/date_strategy.go
type DateStrategy struct {
	Format string `yaml:"format"`
}

func (s *DateStrategy) Selector() string {
	return "date"
}

func (s *DateStrategy) LoadConfig(config map[string]interface{}) error {
	if format, ok := config["format"].(string); ok {
		s.Format = format
	} else {
		s.Format = "2006/01" // default
	}
	return nil
}

func (s *DateStrategy) FinalDirPath(ctx Context) (string, error) {
	yearMonth := ctx.Info().ModTime().Format(s.Format)
	finalDest := filepath.Join(ctx.DstDir(), yearMonth, ctx.Info().Name())
	return finalDest, nil
}

func init() {
	RegisterStrategy("date", func() Strategy {
		slog.Debug("Create a strategy", "name", "date")
		return &DateStrategy{}
	})
}
