package classify

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/internal/stats"
	"github.com/polocto/FolderFlow/pkg/concurrency"
)

type Classifier struct {
	cfg    config.Config
	stats  *stats.Stats
	dryRun bool
}

func NewClassifier(cfg config.Config, s *stats.Stats, dryRun bool) (*Classifier, error) {
	return &Classifier{
		cfg:    cfg,
		stats:  s,
		dryRun: dryRun,
	}, nil
}

func (c *Classifier) Classify() error {
	slog.Info("Starting classification",
		"sources", len(c.cfg.SourceDirs),
		"destinations", len(c.cfg.DestDirs),
		"workers", c.cfg.MaxWorkers,
	)

	c.stats.StartRun()
	defer func() {
		c.stats.EndRun()
		slog.Info("Classification completed", "Stats", c.stats.String())
	}()
	if len(c.cfg.SourceDirs) == 0 {
		slog.Error("No source directories configured, nothing to classify")
		return fmt.Errorf("no source directories configured")
	}

	if len(c.cfg.DestDirs) == 0 {
		slog.Error("No destination directories configured, nothing to classify")
		return fmt.Errorf("no destination directories configured")
	}
	if c.cfg.Regroup != nil {
		slog.Info("All moved file will be regrouped", "regroupDir", c.cfg.Regroup.Path, "mode", c.cfg.Regroup.Mode)
	}
	for _, sourceDir := range c.cfg.SourceDirs {
		if sourceDir == "" {
			slog.Warn("Skipping empty source directory")
			continue
		}
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			slog.Warn("Source directory does not exist, skipping", "sourceDir", sourceDir)
			continue
		}

		if c.cfg.Regroup != nil && sourceDir == c.cfg.Regroup.Path {
			slog.Warn("Source directory is the same as regroup path, skipping to avoid conflicts", "sourceDir", sourceDir, "regroupPath", c.cfg.Regroup.Path)
			continue
		}

		if err := c.processSourceDir(sourceDir); err != nil {
			slog.Error("Failed to process source directory", "sourceDir", sourceDir, "err", err, "stats", c.stats.String())
			continue
		}
	}
	return nil
}

func (c *Classifier) processSourceDir(sourceDir string) error {
	defer c.stats.Time(&c.stats.Timing.Walk)()
	wp := concurrency.NewWorkerPool(c.cfg.MaxWorkers)
	skipDir := map[string]struct{}{
		".git":         {},
		"node_modules": {},
	}

	slog.Info("Skipping directories", "directories skipped", skipDir)
	err := filepath.WalkDir(sourceDir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("WalkDir error", "path", filePath, "err", err)
			c.stats.Error(err)
			return err
		}
		// Ignorer les répertoires (sauf si on veut les traiter)
		if d.IsDir() {
			// Exemple : sauter les répertoires comme .git
			if _, ok := skipDir[d.Name()]; ok {
				return fs.SkipDir
			}

			return nil // Continuer sans descendre dans le répertoire
		}
		info, err := d.Info()
		if err != nil {
			slog.Error("Failed to stat file", "path", filePath, "err", err)
			c.stats.Error(err)
			return err
		}
		c.stats.FileSeen(info.Size())
		wp.Add()

		go func(sourceDir, filePath string, info fs.FileInfo) {
			defer wp.Done()
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("panic in worker: %v", r)
					wp.ReportError(err)
					c.stats.Error(err)
				}
			}()
			if err := c.processFile(sourceDir, filePath, info); err != nil {
				wp.ReportError(err)
			}
		}(sourceDir, filePath, info)

		return nil
	})
	if err != nil {
		slog.Error("Error walking source directory", "sourceDir", sourceDir, "err", err)
		return err
	}

	if err := wp.Wait(); err != nil {
		slog.Error(
			"Errors occurred while processing source directory",
			"sourceDir", sourceDir,
			"error", err,
		)

		return err
	}
	return nil
}

func (c *Classifier) processFile(sourceDir, filePath string, info fs.FileInfo) error {
	defer c.stats.Time(&c.stats.Timing.Classify)()
	for _, dest := range c.cfg.DestDirs {
		if sourceDir == dest.Path {
			slog.Warn("Source and destination paths overlap, skipping to avoid conflicts", "sourceDir", sourceDir, "destPath", dest.Path)
			continue
		}
		// Check if file matches all filters for this DestDir
		if ok, err := matchFile(filePath, info, dest.Filters); err != nil {
			c.stats.Error(err)
			return err
		} else if !ok {
			continue
		}
		// File matched all filters for this DestDir
		c.stats.FileMatched()

		destFile, err := destPath(sourceDir, dest.Path, filePath, info, dest.Strategy)
		if err != nil {
			c.stats.Error(err)
			return err
		}

		// Move the file using the destination
		var action MoveAction
		if action, err = moveFile(filePath, destFile, dest.OnConflict, c.dryRun); err != nil {
			c.stats.Error(err)
			return err
		}
		// Succès : moved
		switch action {
		case MoveCopy:
			c.stats.FileCopied(info.Size())
		case MoveRenamed:
			c.stats.FileRenamed(info.Size())
		case MoveOverwritten:
			c.stats.FileOverwrtitten(info.Size())
		case MoveSkipped:
			c.stats.FileSkipped()
		default:
			c.stats.FileMoved(info.Size())
		}

		// Handle regrouping
		if c.cfg.Regroup == nil || c.cfg.Regroup.Path == "" {
			return nil
		}
		if err := c.regroupFile(sourceDir, filePath, destFile, info); err != nil {
			c.stats.Error(err)
			return err
		}
		return nil
	}
	c.stats.FileSkipped()
	return nil
}
