package core

import (
	"io/fs"
	"log/slog"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

func Classify(cfg config.Config, dryRun bool) error {
	slog.Info("Starting classification")
	for _, sourceDir := range cfg.SourceDirs {
		if err := processSourceDir(sourceDir, cfg.DestDirs, cfg.Regroup, dryRun); err != nil {
			return err
		}
	}
	return nil
}

// processSourceDir processes a single source directory according to the provided destination directories and regroup settings.
func processSourceDir(sourceDir string, destDirs map[string]config.DestDir, regroup config.Regroup, dryRun bool) error {
	slog.Info("Processing source directory", "sourceDir", sourceDir)

	// Load and configure filters as needed
	// Example: load an extension filter
	// extFilter := core.LoadExtensionsFilter()
	// extFilter.SetConfig(map[string]interface{}{"extensions": []string{".jpg", ".png"}})
	// filters = append(filters, extFilter)
	// Implementation of processing logic goes here
	return filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			slog.Warn("Skipping file due to error", "path", path, "err", err)
			return nil
		}
		slog.Info("Visiting file", "path", path)

		// Placeholder for actual classification logic

		// Match files against DestDir rules
		for destName, dest := range destDirs {
			filters, _, err := dest.LoadPlugins()
			if err != nil {
				return err
			}

			if ok, err := matchFile(path, info, filters); err != nil {
				return err
			} else if !ok {
				slog.Debug("File did not match", "path", path, "dest", destName)
				continue
			}
			slog.Debug("File matched", "path", path, "dest", destName)
			// Build destination path
			// _, err := filepath.Rel(sourceDir, filepath.Dir(path))
			// if err != nil {
			// 	return err
			// }

			// var destPath string

			// slog.Debug("Creating directory", "path", destPath)
			// if !dryRun {
			// 	if err := os.MkdirAll(destPath, 0755); err != nil {
			// 		return err
			// 	}
			// }

			// destFile := filepath.Join(destPath, info.Name())
			// slog.Debug("Would move", "source", path, "dest", destFile)

			// if !dryRun {
			// 	if err := os.Rename(path, destFile); err != nil {
			// 		if err := utils.CopyFile(path, destFile); err != nil {
			// 			return err
			// 		}
			// 		os.Remove(path)
			// 	}
			// }

			// Handle regrouping
			// if regroup.Path != "" {
			// 	if err := handleRegroup(path, destFile, regroup, dryRun); err != nil {
			// 		return err
			// 	}
			// }
		}

		return nil
	})
}

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
			slog.Warn("Filter error", "filter", f.Selector(), "path", path, "err", err)
			return false, nil
		}
		if !matched {
			slog.Debug("File rejected by filter", "filter", f.Selector(), "path", path)
			return false, nil
		}
	}

	return true, nil
}
