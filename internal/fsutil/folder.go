// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


package fsutil

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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
		slog.Error("Failed to compute relative path", "parent", parentAbs, "child", childAbs, "err", err)
		return false
	}

	// Vérifier la sortie
	if rel == "." {
		return true
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		slog.Debug("Child is not a subdirectory of parent", "parent", parentAbs, "child", childAbs, "relative", rel)
		return false
	}

	return true
}

func IsCrossDeviceError(err error) bool {
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		return errors.Is(linkErr.Err, syscall.EXDEV)
	}
	return false
}
