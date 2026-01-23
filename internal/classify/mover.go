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

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
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

func resolveConflict(src, dst filehandler.Context, onConflict string) (destPath string, action MoveAction, err error) {
	destPath = dst.Path()
	switch onConflict {
	case "skip":
		action = MoveSkipped
	case "overwrite":
		action = MoveOverwritten
	case "rename": // rename

		if ok, err := filehandler.Equal(src, dst); err != nil {
			return "", MoveFailed, fmt.Errorf("failed to compare files for equality : source=%s dest=%s err=%w", src.Path(), dst.Path(), err)
		} else if ok {
			slog.Warn("Source and destination files are identical, skipping move", "source", src.Path(), "dest", dst.Path())
			action = MoveSkippedIdentical
		} else {
			destPath = filehandler.GetUniquePath(dst.Path())
			slog.Warn("Renaming destination file to avoid conflict", "originalDest", dst.Path(), "newDest", destPath)
			action = MoveRenamed
		}
	default:
		return "", MoveFailed, fmt.Errorf("unknown conflict resolution strategy: %s", onConflict)
	}

	return destPath, action, nil
}

func moveFile(file filehandler.Context, destPath, onConflict string, dryRun bool) (MoveAction, filehandler.Context, error) {
	action := MoveMoved

	if dst, err := filehandler.NewContextFile(destPath); err == nil {
		slog.Debug("Conflic found resolving it")
		if destPath, action, err = resolveConflict(file, dst, onConflict); err != nil {
			return action, nil, fmt.Errorf("failed to resolve conflict at %s", destPath)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return MoveFailed, nil, err
	}
	if dryRun {
		slog.Debug("Dry run enabled, not moving file", "source", file.Path(), "dest", destPath)
		return action, nil, nil
	}

	if action == MoveSkipped {
		return action, nil, nil
	}
	srcPath := file.Path()
	newFile, err := executeMove(file, destPath)
	if err != nil {
		return MoveSkipped, newFile, err
	}
	slog.Debug("File moved with success", "src", srcPath, "dst", destPath, "filepath", file.Path())
	return action, newFile, nil
}

func executeMove(file filehandler.Context, dst string) (newFile filehandler.Context, err error) {
	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return nil, err
	}

	// Tentative rapide et atomique
	if newFile, err = filehandler.Replace(file, dst); err == nil {
		return newFile, nil
	}

	return newFile, nil
}
