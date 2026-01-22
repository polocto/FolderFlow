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
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

func Equal(file1, file2 Context) (bool, error) {
	if file1 == nil || file2 == nil {
		return false, ErrContextIsNil
	}
	if file1.Size() != file2.Size() {
		return false, nil
	}
	hash1, err := file1.GetHash()
	if err != nil {
		return false, fmt.Errorf("failed to get hash of the first file: %w", err)
	}
	hash2, err := file2.GetHash()
	if err != nil {
		return false, fmt.Errorf("failed to get hash of the second file: %w", err)
	}

	if hash1 != hash2 {
		return false, nil
	}
	return true, nil
}

func ListDuplicates(files []Context) ([][]Context, error) {
	bySize := make(map[int64][]Context)

	// Group by file size
	for _, f := range files {
		if f == nil || f.Kind() != KindRegular {
			continue
		}
		bySize[f.Size()] = append(bySize[f.Size()], f)
	}

	duplicateGroups := make([][]Context, 0, len(files)/2)

	// Within same-size groups, group by hash
	for _, group := range bySize {
		if len(group) < 2 {
			continue
		}

		byHash := make(map[[sha256.Size]byte][]Context)

		for _, f := range group {
			hash, err := f.GetHash()
			if err != nil {
				return nil, err
			}
			byHash[hash] = append(byHash[hash], f)
		}

		// Keep only real duplicates
		for _, dupGroup := range byHash {
			if len(dupGroup) > 1 {
				duplicateGroups = append(duplicateGroups, dupGroup)
			}
		}
	}

	return duplicateGroups, nil
}

// ListSymlinks returns all symlink contexts from the list of files.
// If checkTargets is true, it will also check the target existence and separate
// functional, broken, and unknown symlinks.
func ListSymlinks(files []Context, checkTargets bool) (functional, broken, unknown []Context) {
	functional = make([]Context, 0, len(files))
	broken = make([]Context, 0, len(files))
	unknown = make([]Context, 0, len(files))

	for _, f := range files {
		if f == nil || f.Kind() != KindSymlink {
			continue
		}

		if !checkTargets {
			functional = append(functional, f)
			continue
		}

		target, err := os.Readlink(f.Path())
		if err != nil {
			// Invalid symlink
			broken = append(broken, f)
			continue
		}

		// Resolve relative symlink paths
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(f.Path()), target)
		}

		if _, err := os.Stat(target); errors.Is(err, fs.ErrNotExist) {
			broken = append(broken, f)
		} else if err != nil {
			unknown = append(unknown, f)
			slog.Warn("cannot access symlink target", "symlink", f.Path(), "target", target, "error", err)
		} else {
			functional = append(functional, f)
		}
	}

	return
}

func SplitFiles(files []Context) (regular, funcSymlinks, brokenSymlinks, unknownSymlinks []Context) {
	regular = make([]Context, 0, len(files))
	funcSymlinks = make([]Context, 0, len(files))
	brokenSymlinks = make([]Context, 0, len(files))
	unknownSymlinks = make([]Context, 0, len(files))

	for _, f := range files {
		if f == nil {
			continue
		}

		switch f.Kind() {
		case KindRegular:
			regular = append(regular, f)
		case KindSymlink:
			// Check if symlink target exists
			target, err := os.Readlink(f.Path())
			if err != nil {
				// Invalid symlink
				brokenSymlinks = append(brokenSymlinks, f)
				continue
			}
			// Resolve relative symlink paths
			if !filepath.IsAbs(target) {
				target = filepath.Join(filepath.Dir(f.Path()), target)
			}

			if _, err := os.Stat(target); errors.Is(err, fs.ErrNotExist) {
				// Target does not exist
				brokenSymlinks = append(brokenSymlinks, f)
			} else if err != nil {
				// Target exists, but we could not access it (permissions, I/O errors, etc.)
				unknownSymlinks = append(unknownSymlinks, f)
				slog.Warn("cannot access symlink target", "symlink", f.Path(), "target", target, "error", err)
			} else {
				// Target exists
				funcSymlinks = append(funcSymlinks, f)
			}
		}
	}

	return
}
