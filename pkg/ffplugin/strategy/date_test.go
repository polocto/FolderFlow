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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDateStrategy(t *testing.T) {
	s := &DateStrategy{}
	require.NoError(t, s.LoadConfig(map[string]interface{}{}))

	file := filepath.Join(t.TempDir(), "a.txt")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o644))
	info, _ := os.Stat(file)

	fileCtx := &ContextStrategy{
		relPath: filepath.Join("Important", "Famille", "fichier.txt"),
		dstDir:  filepath.Join("srv", "backup"),
		info:    info,
	}
	path, err := s.FinalDirPath(fileCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if path == "" {
		t.Fatalf("expected non-empty path")
	}
}
