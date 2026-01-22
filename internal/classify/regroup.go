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
	"os"
	"path/filepath"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
)

func execute(source filehandler.Context, target, mode string) (file filehandler.Context, err error) {
	if source == nil {
		return nil, fmt.Errorf("failed to regroup non existing file: %w", filehandler.ErrContextIsNil)
	}
	// Actual regrouping logic (symlink, hardlink, copy) would be implemented here
	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return nil, err
	}
	switch mode {
	case "symlink":
		file, err = filehandler.Symlink(source, target)
	case "hardlink":
		file, err = filehandler.Hardlink(source, target)
	case "copy":
		file, err = filehandler.CopyFileAtomic(source, target)
	default:
		return nil, ErrInvalidRegroupMode(mode)
	}
	slog.Debug(
		"File regrouped successfully",
		"source", source.Path(),
		"target", file.Path(),
		"mode", mode,
	)

	return file, err

}
