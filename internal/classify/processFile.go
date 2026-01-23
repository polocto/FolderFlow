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
	"log/slog"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
)

func (c *Classifier) processFile(sourceDir, filePath string) error {
	defer c.stats.Time(&c.stats.Timing.Classify)()
	for _, dest := range c.cfg.DestDirs {
		if sourceDir == dest.Path {
			slog.Warn(
				"Source and destination paths overlap, skipping to avoid conflicts",
				"sourceDir",
				sourceDir,
				"destPath",
				dest.Path,
			)
			continue
		}
		file, err := filehandler.NewContextFile(filePath)
		if err != nil {
			return err
		}

		// Check if file matches all filters for this DestDir
		ok, err := c.runFilters(file, dest.Filters)
		if err != nil || !ok {
			continue
		}
		// File matched all filters for this DestDir
		c.stats.FileMatched()

		destinationPath, err := c.runStartegy(file, sourceDir, dest.Path, dest.Strategy)
		if err != nil {
			return err
		}
		regroupFile := file
		var regroupPath string
		// Handle regrouping
		if c.cfg.Regroup != nil || c.cfg.Regroup.Path != "" {
			regroupPath, err = c.runStartegy(
				file,
				sourceDir,
				c.cfg.Regroup.Path,
				c.cfg.Regroup.Strategy,
			)
			if err != nil {
				c.stats.Error(err)
				return err
			}
		}

		// Move the file using the destination
		var action MoveAction
		var copy filehandler.Context
		if action, copy, err = moveFile(file, destinationPath, dest.OnConflict, c.dryRun); err != nil {
			c.stats.Error(err)
			return err
		}
		// Succès : moved
		switch action {
		case MoveCopy:
			c.stats.FileCopied(file.Size())
		case MoveRenamed:
			c.stats.FileRenamed(file.Size())
		case MoveOverwritten:
			c.stats.FileOverwrtitten(file.Size())
		case MoveSkipped:
			c.stats.FileSkipped()
		default:
			c.stats.FileMoved(file.Size())
		}

		if copy != nil {
			regroupFile = copy
		}

		if c.cfg.Regroup != nil && c.cfg.Regroup.Path != "" {
			if _, err := execute(regroupFile, regroupPath, c.cfg.Regroup.Mode); err != nil {
				return fmt.Errorf(
					"could not regroup file: path=%q regrouppath=%q err=%w",
					file.Path(),
					regroupPath,
					err,
				)
			}
		}

		return nil

	}
	c.stats.FileSkipped()
	return nil
}
