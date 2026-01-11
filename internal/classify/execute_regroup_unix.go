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

//go:build !windows

package classify

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/fsutil"
)

func executeRegroup(source, target, mode string) error {
	// Actual regrouping logic (symlink, hardlink, copy) would be implemented here
	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	switch mode {
	case "symlink":
		rel := source
		if r, err := filepath.Rel(filepath.Dir(target), source); err == nil {
			rel = r
		}

		if err := os.Symlink(rel, target); err != nil {

			if errors.Is(err, fs.ErrExist) {
				slog.Debug("Regroup target already exists, skipping", "path", target)
				return nil
			}
			return err
		}
	case "hardlink":
		if err := os.Link(source, target); err != nil {
			if errors.Is(err, fs.ErrExist) {
				slog.Debug("Regroup target already exists, skipping", "path", target)
				return nil
			}
			if !fsutil.IsCrossDeviceError(err) {
				return err
			}
			slog.Warn(
				"Hardlink failed (cross-device), falling back to copy",
				"source", source,
				"dest", target,
			)
			if err := fsutil.CopyFileAtomic(source, target); err != nil {
				return err
			}
		}
	case "copy":
		if err := fsutil.CopyFileAtomic(source, target); err != nil {
			return err
		}
	default:
		return ErrInvalidRegroupMode(mode)
	}
	slog.Debug(
		"File regrouped successfully",
		"source", source,
		"target", target,
		"mode", mode,
	)

	return nil
}
