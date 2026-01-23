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

package strategy

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

// plugin/date_strategy.go
type DirChainStrategy struct{}

func (s *DirChainStrategy) Selector() string {
	return "dirchain"
}

func (s *DirChainStrategy) LoadConfig(config map[string]interface{}) error {
	return nil
}

func (s *DirChainStrategy) FinalDirPath(ctx Context) (string, error) {
	if ctx.Info().IsDir() {
		return "", fmt.Errorf("filePath %s is a directory, expected a file", ctx.PathFromSource())
	}
	// Nettoyer les chemins pour éviter les problèmes avec les slashes finaux

	finalDest := filepath.Join(ctx.DstDir(), ctx.PathFromSource())

	// Vérifier que la destination reste dans destDir (défense en profondeur)
	relFromDest, err := filepath.Rel(ctx.DstDir(), finalDest)
	if err != nil {
		return "", fmt.Errorf("failed to validate destination path: %w", err)
	}

	if relFromDest == ".." || strings.HasPrefix(relFromDest, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf(
			"computed destination path is outside of destination directory: %s",
			finalDest,
		)
	}

	return finalDest, nil
}

func init() {
	RegisterStrategy("dirchain", func() Strategy {
		slog.Debug("Create a strategy", "name", "dirchain")
		return &DirChainStrategy{}
	})
}
