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
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetUniquePath returns a unique path based on destPath.
// If destPath exists, it appends a numeric suffix like "_1", "_2", etc.
func GetUniquePath(destPath string) string {
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return destPath
	}

	ext := filepath.Ext(destPath)
	base := strings.TrimSuffix(destPath, ext)
	counter := 1

	for {
		newPath := fmt.Sprintf("%s_%d%s", base, counter, ext)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}
