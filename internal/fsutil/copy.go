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

package fsutil

import (
	"crypto/sha256"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/stats"
)

func CopyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return CopyFileAtomic(path, target)
	})
}

func CopyFile(src, dst string, s *stats.Stats) ([]byte, error) {
	if s != nil {
		defer s.Time(&s.Timing.Hash)()
	}
	in, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := in.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := out.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	hasher := sha256.New()

	// Write to file AND hash simultaneously
	writer := io.MultiWriter(out, hasher)

	if size, err := io.Copy(writer, in); err != nil {
		return nil, err
	} else if s != nil {
		s.FileCopied(size)
		s.HashComputed()
	}

	if err := out.Sync(); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func CopyFileAtomic(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := in.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	tmp := dst + ".tmp"

	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		if cerr := out.Close(); cerr != nil {
			slog.Warn(
				"failed to close file after copy failure",
				"path", tmp,
				"error", cerr,
			)
		}
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after copy failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	if err := out.Sync(); err != nil {
		if cerr := out.Close(); cerr != nil {
			slog.Warn(
				"failed to close file after sync failure",
				"path", tmp,
				"error", cerr,
			)
		}
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after sync failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	if err := out.Close(); err != nil {
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after close failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	if err := ReplaceFile(tmp, dst); err != nil {
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after replace failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	return nil
}
