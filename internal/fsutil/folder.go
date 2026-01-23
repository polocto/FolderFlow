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
	"log/slog"
	"path/filepath"
	"strings"
)

func IsSubDirectory(parent, child string) bool {
	// Normaliser
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)

	// Convertir en chemins absolus
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		slog.Error("Failed to get parent absolute path", "path", parent, "err", err)
		return false
	}

	childAbs, err := filepath.Abs(child)
	if err != nil {
		slog.Error("Failed to get child absolute path", "path", child, "err", err)
		return false
	}

	// Calculer le relatif
	rel, err := filepath.Rel(parentAbs, childAbs)
	if err != nil {
		slog.Error(
			"Failed to compute relative path",
			"parent",
			parentAbs,
			"child",
			childAbs,
			"err",
			err,
		)
		return false
	}

	// Vérifier la sortie
	if rel == "." {
		return true
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		slog.Debug(
			"Child is not a subdirectory of parent",
			"parent",
			parentAbs,
			"child",
			childAbs,
			"relative",
			rel,
		)
		return false
	}

	return true
}
