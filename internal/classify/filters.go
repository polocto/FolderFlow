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

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	internalfilter "github.com/polocto/FolderFlow/internal/filter"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

// matchFile checks if a file matches all the rules in DestDir.
func matchFile(file filehandler.Context, filters []filter.Filter) (bool, error) {
	// If no filters are provided, match all files
	if len(filters) == 0 {
		return true, nil
	}

	ctx, err := internalfilter.NewContextFilter(file)
	if err != nil {
		return false, err
	}

	// Run all filters
	for _, f := range filters {
		matched, err := f.Match(ctx)
		if err != nil {
			slog.Error("Filter error", "filter", f.Selector(), "path", file.Path(), "err", err)
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	slog.Debug("File matched", "path", file.Path(), "filers", filters)
	return true, nil
}

func (c *Classifier) runFilters(path filehandler.Context, filters []filter.Filter) (bool, error) {
	var ok bool
	err := c.safeRun("filters", func() error {
		var err error
		ok, err = matchFile(path, filters)
		return err
	})
	return ok, err
}
