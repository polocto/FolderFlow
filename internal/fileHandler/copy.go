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

package filehandler

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func CopyFileAtomic(file Context, dst string) error {
	if file == nil {
		return fmt.Errorf("cannot copy: %w", ErrContextIsNil)
	}
	in, err := os.Open(file.Path())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := in.Close(); cerr != nil {
			slog.Warn("failed to close source file", "path", file.Path(), "error", cerr)
		}
	}()

	tmpFile, err := os.CreateTemp(filepath.Dir(dst), filepath.Base(dst)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	// Ensure temp file is closed & deleted on any error
	defer func() {
		if cerr := tmpFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if err != nil {
			if rmErr := os.Remove(tmpPath); rmErr != nil {
				slog.Warn("failed to delete temp file after error", "path", tmpPath, "error", rmErr)
			}
		}
	}()

	if _, err := io.Copy(tmpFile, in); err != nil {
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		return err
	}

	if err := replaceFile(tmpPath, dst); err != nil {
		if rmErr := os.Remove(tmpPath); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after replace failure",
				"path", tmpPath,
				"error", rmErr,
			)
		}
		return err
	}

	return nil
}
