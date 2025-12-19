package classify

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/pkg/concurrency"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

func Classify(cfg config.Config, dryRun bool) error {
	slog.Info("Starting classification")
	var totalStats = NewStats()
	if len(cfg.SourceDirs) == 0 {
		slog.Error("No source directories configured, nothing to classify")
		return fmt.Errorf("no source directories configured")
	}

	if len(cfg.DestDirs) == 0 {
		slog.Error("No destination directories configured, nothing to classify")
		return fmt.Errorf("no destination directories configured")
	}
	for _, sourceDir := range cfg.SourceDirs {
		if sourceDir == "" {
			slog.Warn("Skipping empty source directory")
			continue
		}
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			slog.Warn("Source directory does not exist, skipping", "sourceDir", sourceDir)
			continue
		}
		if sourceDir == cfg.Regroup.Path {
			slog.Warn("Source directory is the same as regroup path, skipping to avoid conflicts", "sourceDir", sourceDir, "regroupPath", cfg.Regroup.Path)
			continue
		}

		slog.Info("Processing source directory", "sourceDir", sourceDir)
		if srcStats, err := processSourceDir(sourceDir, cfg.DestDirs, cfg.Regroup, dryRun, cfg.MaxWorkers); err != nil {
			slog.Error("Failed to process source directory", "sourceDir", sourceDir, "err", err, "stats", srcStats.String())
			continue
		} else {
			totalStats.Merge(srcStats)
			slog.Info("Finished processing source directory", "sourceDir", sourceDir, "stats", srcStats.String())
		}
	}
	slog.Info("Classification completed", "totalStats", totalStats.String())
	return nil
}

// processSourceDir processes a single source directory according to the provided destination directories and regroup settings.
func processSourceDir(sourceDir string, destDirs map[string]config.DestDir, regroup config.Regroup, dryRun bool, maxWorkers int) (*Stats, error) {
	sourceStats := NewStats()
	// Load and configure filters as needed
	// Example: load an extension filter
	// extFilter := core.LoadExtensionsFilter()
	// extFilter.SetConfig(map[string]interface{}{"extensions": []string{".jpg", ".png"}})
	// filters = append(filters, extFilter)
	// Match files against DestDir rules
	for destName, dest := range destDirs {
		filters, strat, err := dest.LoadPlugins()
		if err != nil {
			slog.Error("Failed to load plugins for destination", "dest", destName, "err", err)
			return sourceStats, err
		}
		slog.Debug("Processing destination", "dest", destName, "strategy", strat.Selector())

		if len(filters) == 0 {
			slog.Warn("No filters configured for destination, all files will match", "dest", destName)
		}
		if strat == nil {
			slog.Error("No strategy configured for destination", "dest", destName)
			return sourceStats, fmt.Errorf("no strategy configured for destination: %s", destName)
		}
		if dest.Path == "" {
			slog.Error("No path configured for destination", "dest", destName)
			return sourceStats, fmt.Errorf("no path configured for destination: %s", destName)
		}
		if sourceDir == dest.Path {
			slog.Warn("Source directory is the same as destination path, skipping to avoid conflicts", "sourceDir", sourceDir, "destPath", dest.Path)
			continue
		}
		// Implementation of processing logic goes here
		if stats, err := processFile(sourceDir, destName, dest, filters, strat, dryRun, maxWorkers); err != nil {
			sourceStats.Merge(stats)
			slog.Error("Failed to process files for destination", "dest", destName, "err", err)
			return sourceStats, err
		} else {
			sourceStats.Merge(stats)
			slog.Debug("Finished processing destination", "dest", destName, "stats", stats.String())
		}
	}
	return sourceStats, nil
}

func processFile(sourceDir, destName string, dest config.DestDir, filters []filter.Filter, strat strategy.Strategy, dryRun bool, maxWorkers int) (*Stats, error) {
	wp := concurrency.NewWorkerPool(maxWorkers)
	s := NewStats()
	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		slog.Debug("Visiting file", "path", path)
		skipDir := []string{".git", "node_modules"}
		// Ignorer les répertoires (sauf si on veut les traiter)
		if d.IsDir() {
			// Exemple : sauter les répertoires comme .git
			if slices.Contains(skipDir, d.Name()) {
				return fs.SkipDir // Sauter ce répertoire et ses enfants
			}
			return nil // Continuer sans descendre dans le répertoire
		}

		info, err := os.Stat(path)
		if err != nil {
			slog.Error("Failed to stat file", "path", path, "err", err)
			s.RecordError(path, err.Error())
			return err
		}
		s.TotalFile(path, info.Size())
		wp.Add()
		go func(path string, info fs.FileInfo) {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("panic in goroutine: %v", r)
					wp.ReportError(err)
					s.RecordError(path, err.Error())
				}
				wp.Done()
			}()
			slog.Debug("Processing file in goroutine", "path", path)
			// Check if file matches all filters for this DestDir
			if ok, err := matchFile(path, info, filters); err != nil {
				err := fmt.Errorf("error matching file %s: %w", path, err)
				wp.ReportError(err)
				s.RecordError(path, err.Error())
				return
			} else if !ok {
				slog.Debug("File did not match", "path", path, "dest", destName)
				return
			}
			// File matched all filters for this DestDir
			slog.Debug("File matched", "path", path, "dest", destName)
			s.MatchedFile(path, info.Size())

			destFile, err := destPath(sourceDir, dest.Path, path, info, strat)
			if err != nil {
				wp.ReportError(fmt.Errorf("failed to compute destination path for file %s: %w", path, err))
				s.RecordError(path, err.Error())
				return
			}

			// Move the file using the destination
			t0 := time.Now()
			if err := moveFile(path, destFile, dest.OnConflict, dryRun); err != nil {
				wp.ReportError(fmt.Errorf("failed to move file %s to %s: %w", path, destFile, err))
				s.RecordError(path, err.Error())
				return
			}
			s.MovedFile(path, info.Size(), t0) // Succès : moved
			// // Handle regrouping
			// if regroup.Path != "" {
			// 	if err := handleRegroup(destFile, regroup, dryRun); err != nil {
			// 		return err
			// 	}
			// } else {
			// 	slog.Error("No regrouping configured", "file", destFile)
			// 	return fmt.Errorf("no regrouping configured for file: %s", regroup)
			// }
		}(path, info)

		return nil
	})
	if err != nil {
		slog.Error("Error walking source directory", "sourceDir", sourceDir, "err", err)
		return s, err
	}

	if err := wp.Wait(); err != nil {
		slog.Error("Processing completed with errors", "stats", s.String())
		return s, err
	}

	slog.Debug("Processing completed successfully", "stats", s.String())
	return s, nil
}
