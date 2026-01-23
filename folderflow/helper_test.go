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

package testdata_test

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
)

// CopyDir recursively copies src directory to dst
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		return os.Chmod(destPath, info.Mode())
	})
}

// FileHash returns the SHA256 hash of the file
func FileHash(path string, _ interface{}) ([sha256.Size]byte, error) {
	var h [sha256.Size]byte
	f, err := os.Open(path)
	if err != nil {
		return h, err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return h, err
	}

	copy(h[:], hash.Sum(nil))
	return h, nil
}
