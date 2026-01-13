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

package classify

import (
	"log/slog"

	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/internal/stats"
)

type Classifier struct {
	cfg    config.Config
	stats  *stats.Stats
	dryRun bool
}

func NewClassifier(cfg config.Config, s *stats.Stats, dryRun bool) (*Classifier, error) {
	if err := validateConfiguration(cfg); err != nil {
		return nil, err
	}

	return &Classifier{
		cfg:    cfg,
		stats:  s,
		dryRun: dryRun,
	}, nil
}

func (c *Classifier) Classify() error {
	defer func() {
		c.stats.EndRun()
		slog.Info("Classification completed", "Stats", c.stats.String())
	}()

	c.stats.StartRun()
	slog.Info("Starting classification",
		"sources", len(c.cfg.SourceDirs),
		"destinations", len(c.cfg.DestDirs),
		"workers", c.cfg.MaxWorkers,
	)

	if c.cfg.Regroup != nil {
		slog.Info("All file will be regrouped", "rgpDir", c.cfg.Regroup.Path, "mode", c.cfg.Regroup.Mode)
	}

	for _, sourceDir := range c.cfg.SourceDirs {
		if err := c.processSourceDir(sourceDir); err != nil {
			slog.Error("Failed to process source directory", "sourceDir", sourceDir, "err", err, "stats", c.stats.String())
			continue
		}
	}
	return nil
}
