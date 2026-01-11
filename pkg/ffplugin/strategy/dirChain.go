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
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

// plugin/date_strategy.go
type DirChainStrategy struct {
}

func (s *DirChainStrategy) Selector() string {
	return "dirchain"
}

func (s *DirChainStrategy) LoadConfig(config map[string]interface{}) error {
	return nil
}

func (s *DirChainStrategy) FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error) {
	if info.IsDir() {
		return "", fmt.Errorf("filePath %s is a directory, expected a file", filePath)
	}
	// Nettoyer les chemins pour éviter les problèmes avec les slashes finaux
	srcDir = filepath.Clean(srcDir)
	destDir = filepath.Clean(destDir)
	filePath = filepath.Clean(filePath)
	fileDir := filepath.Dir(filePath)

	relFromSrc, err := filepath.Rel(srcDir, fileDir)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path from srcDir: %w", err)
	}

	if relFromSrc == ".." || strings.HasPrefix(relFromSrc, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf(
			"filePath %s is not a subdirectory of srcDir %s",
			filePath, srcDir,
		)
	}

	// Construire le chemin de destination
	if relFromSrc == "." {
		// Fichier à la racine du source
		return destDir, nil
	}

	finalDest := filepath.Join(destDir, relFromSrc)

	// Vérifier que la destination reste dans destDir (défense en profondeur)
	relFromDest, err := filepath.Rel(destDir, finalDest)
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
