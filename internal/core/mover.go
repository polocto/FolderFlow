package core

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/pkg/concurrency"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

func Classify(cfg config.Config, dryRun bool) error {
	slog.Info("Starting classification")
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
	// Match files against DestDir rules
	for destName, dest := range destDirs {
		filters, strat, err := dest.LoadPlugins()
		if err != nil {
			slog.Error("Failed to load plugins for destination", "dest", destName, "err", err)
			return err
		}
		slog.Info("Processing destination", "dest", destName, "strategy", strat.Selector())

		if len(filters) == 0 {
			slog.Warn("No filters configured for destination, all files will match", "dest", destName)
		}
		if strat == nil {
			slog.Error("No strategy configured for destination", "dest", destName)
			return fmt.Errorf("no strategy configured for destination: %s", destName)
		}
		if dest.Path == "" {
			slog.Error("No path configured for destination", "dest", destName)
			return fmt.Errorf("no path configured for destination: %s", destName)
		}
		if sourceDir == dest.Path {
			slog.Warn("Source directory is the same as destination path, skipping to avoid conflicts", "sourceDir", sourceDir, "destPath", dest.Path)
			continue
		}
		// Implementation of processing logic goes here
		if err := processFile(sourceDir, destName, dest, filters, strat, dryRun); err != nil {
			return err
		}
	}

	return nil
}

func processFile(sourceDir, destName string, dest config.DestDir, filters []filter.Filter, strat strategy.Strategy, dryRun bool) error {
	wp := concurrency.NewWorkerPool(0) // 10 workers max
	var s stats                        // to track statistics
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
		// Compter les fichiers visités
		s.totalFiles++

		info, err := os.Stat(path)
		if err != nil {
			slog.Error("Failed to stat file", "path", path, "err", err)
			s.RecordFile(false, false, true) // Enregistrer l'erreur
			return err
		}
		wp.Add()
		go func(path string, info fs.FileInfo) {
			defer func() {
				if r := recover(); r != nil {
					wp.ReportError(fmt.Errorf("panic in goroutine: %v", r))
					s.RecordFile(false, false, true) // Enregistrer l'erreur
				}
				wp.Done()
			}()
			slog.Info("Processing file in goroutine", "path", path)
			// Check if file matches all filters for this DestDir
			if ok, err := matchFile(path, info, filters); err != nil {
				wp.ReportError(fmt.Errorf("error matching file %s: %w", path, err))
				s.RecordFile(false, false, true) // Enregistrer l'erreur
				return
			} else if !ok {
				slog.Debug("File did not match", "path", path, "dest", destName)
				return
			}
			// File matched all filters for this DestDir
			slog.Info("File matched", "path", path, "dest", destName)
			s.RecordFile(true, false, err != nil) // Enregistrer le fichier comme correspondant

			destFile, err := destPath(sourceDir, dest.Path, path, info, strat)
			if err != nil {
				wp.ReportError(fmt.Errorf("failed to compute destination path for file %s: %w", path, err))
				s.RecordFile(false, false, true) // Enregistrer l'erreur
				return
			}

			// Move the file using the destination
			if err := moveFile(path, destFile, dryRun); err != nil {
				wp.ReportError(fmt.Errorf("failed to move file %s to %s: %w", path, destFile, err))
				s.RecordFile(false, false, true) // Enregistrer l'erreur
				return
			}
			s.RecordFile(false, true, false) // Succès : moved
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
		return err
	}

	if err := wp.Wait(); err != nil {
		slog.Error("Processing completed with errors", "stats", s.String())
		return err
	}

	slog.Info("Processing completed successfully", "stats", s.String())
	return nil
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
			return false, err
		}
		if !matched {
			slog.Debug("File rejected by filter", "filter", f.Selector(), "path", path)
			return false, nil
		}
	}
	slog.Debug("File matched all filters", "path", path)
	return true, nil
}

func destPath(sourceDir, destDir, filePath string, info fs.FileInfo, strat strategy.Strategy) (string, error) {
	var destPath string
	destPath, err := strat.FinalDirPath(sourceDir, destDir, filePath, info)
	if err != nil {
		slog.Error("Strategy failed to compute destination path", "strategy", strat.Selector(), "err", err)
		return "", err
	}

	// Ensure destination directory is a subdirectory of destDir
	destPathDir := filepath.Dir(destPath)
	if !IsSubDirectory(destPathDir, destDir) {
		slog.Error("Computed destination path is outside of destination directory", "computedPath", destPath, "destDir", destDir)
		return "", fmt.Errorf("computed destination path is outside of destination directory : computedPath=%s destDir=%s", destPath, destDir)
	}

	slog.Info("Creating directory", "path", destPath)

	destFile := filepath.Join(destPath, info.Name())
	slog.Debug("Would move", "source", filepath.Dir(filePath), "dest", destFile)

	return destFile, nil
}

func handleRegroup(filePath string, regroup config.Regroup, dryRun bool) error {
	slog.Info("Handling regroup", "file", filePath, "regroupPath", regroup.Path)
	// Implementation of regrouping logic goes here
	return nil
}

func moveFile(srcPath, destPath string, dryRun bool) error {
	slog.Info("Moving file", "source", srcPath, "dest", destPath)

	if dryRun {
		slog.Info("Dry run enabled, not moving file", "source", srcPath, "dest", destPath)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		slog.Error("Failed to create destination directory", "dir", filepath.Dir(destPath), "err", err)
		return err
	}

	if err := os.Rename(srcPath, destPath); err != nil {
		slog.Error("Failed to move file", "source", srcPath, "dest", destPath, "err", err)
		return err
	}

	return nil
}
