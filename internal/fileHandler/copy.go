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
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func CopyFileAtomic(file Context, dst string) (Context, error) {
	if file == nil {
		return nil, fmt.Errorf("cannot copy: %w", ErrContextIsNil)
	}

	srcHash, err := file.GetHash()
	if err != nil {
		return nil, err
	}

	in, err := os.Open(file.Path())
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := in.Close(); cerr != nil {
			slog.Warn("failed to close source file", "path", file.Path(), "error", cerr)
		}
	}()

	tmpFile, err := os.CreateTemp(filepath.Dir(dst), filepath.Base(dst)+".tmp-*")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()

	// Ensure temp file is closed & deleted on any error
	defer func() {
		if cerr := tmpFile.Close(); cerr != nil && err == nil {
			slog.Warn("failed to close temp file", "path", tmpPath, "error", cerr)
		}
		if err != nil {
			if rmErr := os.Remove(tmpPath); rmErr != nil {
				slog.Warn("failed to delete temp file after error", "path", tmpPath, "error", rmErr)
			}
		}
	}()

	h := sha256.New()

	writer := io.MultiWriter(h, tmpFile)

	if _, err := io.Copy(writer, in); err != nil {
		return nil, err
	}

	if err := tmpFile.Sync(); err != nil {
		return nil, err
	}

	var tmpHash [sha256.Size]byte
	copy(tmpHash[:], h.Sum(nil))

	if tmpHash != srcHash {
		return nil, fmt.Errorf("failed to correctly copy file: path=%s", file.Path())
	}

	if err := replaceFile(tmpPath, dst); err != nil {
		if rmErr := os.Remove(tmpPath); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after replace failure",
				"path", tmpPath,
				"error", rmErr,
			)
		}
		return nil, err
	}

	return newContextFileWithHash(dst, tmpHash)
}
