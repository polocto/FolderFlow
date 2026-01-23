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
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"

	"github.com/polocto/FolderFlow/pkg/concurrency"
)

func (c *Classifier) processSourceDir(sourceDir string) error {
	defer c.stats.Time(&c.stats.Timing.Walk)()

	wp := concurrency.NewWorkerPool(c.cfg.MaxWorkers)

	err := filepath.WalkDir(sourceDir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkDir error : path=%s err=%w", filePath, err)
		}

		if d.IsDir() {
			return nil // Continuer until processing a file
		} else {
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("unable to read a file info: path=%s err=%w", filePath, err)
			}
			c.stats.FileSeen(info.Size())
		}
		wp.Add()

		go func(sourceDir, filePath string) {
			defer wp.Done()
			if err := c.safeRun("processFile", func() error { return c.processFile(sourceDir, filePath) }); err != nil {
				wp.ReportError(err)
			}
		}(sourceDir, filePath)

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
