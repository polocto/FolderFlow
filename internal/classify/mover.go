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
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/fsutil"
)

type MoveAction int

const (
	MoveSkipped MoveAction = iota
	MoveMoved
	MoveOverwritten
	MoveRenamed
	MoveCopy
	MoveSkippedIdentical
	MoveFailed
)

func resolveConflict(srcPath, destPath, onConflict string) (string, MoveAction, error) {
	var action MoveAction
	switch onConflict {
	case "skip":
		action = MoveSkipped
	case "overwrite":
		action = MoveOverwritten
	case "rename": // rename
		if ok, err := fsutil.FilesEqual(srcPath, destPath); err != nil {
			return "", MoveFailed, fmt.Errorf("failed to compare files for equality : source=%s dest=%s err=%w", srcPath, destPath, err)
		} else if ok {
			slog.Warn("Source and destination files are identical, skipping move", "source", srcPath, "dest", destPath)
			action = MoveSkippedIdentical
		} else {
			temp := destPath
			destPath = fsutil.GetUniquePath(destPath)
			slog.Warn("Renaming destination file to avoid conflict", "originalDest", temp, "newDest", destPath)
			action = MoveRenamed
		}
	default:
		return "", MoveFailed, fmt.Errorf("unknown conflict resolution strategy: %s", onConflict)
	}

	return destPath, action, nil
}

func moveFile(srcPath, destPath, onConflict string, dryRun bool) (MoveAction, error) {
	action := MoveMoved

	if _, err := os.Stat(destPath); err == nil {
		slog.Debug("Conflic found resolving it")
		if destPath, action, err = resolveConflict(srcPath, destPath, onConflict); err != nil {
			return action, fmt.Errorf("failed to resolve conflict at %s", destPath)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return MoveFailed, err
	}
	if dryRun {
		slog.Debug("Dry run enabled, not moving file", "source", srcPath, "dest", destPath)
		return action, nil
	}

	if action == MoveSkipped {
		return action, nil
	}

	if err := executeMove(srcPath, destPath); err != nil {
		return MoveSkipped, err
	}
	slog.Debug("File moved with success", "src", srcPath, "dst", destPath)
	return action, nil
}

func executeMove(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	// Tentative rapide et atomique
	if err := fsutil.ReplaceFile(src, dst); err == nil {
		return nil
	} else if !fsutil.IsCrossDeviceError(err) {
		// Erreur autre que EXDEV → vraie erreur
		return fmt.Errorf("failed to rename file : src=%s dst=%s err=%w", src, dst, err)
	}
	// Fallback : copy + fsync + remove
	if err := fsutil.CopyFileAtomic(src, dst); err != nil {
		return fmt.Errorf("failed to copy file : src=%s dst=%s err=%w", src, dst, err)
	}

	return os.Remove(src)
}
