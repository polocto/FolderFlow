// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


package classify

import (
	"io/fs"
	"log/slog"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

// matchFile checks if a file matches all the rules in DestDir.
func matchFile(path string, info fs.FileInfo, filters []filter.Filter) (bool, error) {
	// If no filters are provided, match all files
	if len(filters) == 0 {
		return true, nil
	}

	// Run all filters
	for _, f := range filters {
		matched, err := f.Match(path, info)
		if err != nil {
			slog.Error("Filter error", "filter", f.Selector(), "path", path, "err", err)
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	slog.Debug("File matched", "path", path, "filers", filters)
	return true, nil
}
