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
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirChainStrategy_FinalDirPath(t *testing.T) {
	s := &DirChainStrategy{}

	testCases := []struct {
		name      string
		relPath   string
		destDir   string
		filePath  string
		info      fs.FileInfo
		expected  string
		shouldErr bool
	}{
		{
			name:    "Basic relative path",
			relPath: filepath.Join("Important", "Famille", "fichier.txt"),
			destDir: filepath.Join("srv", "backup"),
			filePath: filepath.Join(
				"home",
				"polocto",
				"Document",
				"Important",
				"Famille",
				"fichier.txt",
			),
			info:      mockFileInfo{isDir: false, name: "fichier.txt"},
			expected:  filepath.Join("srv", "backup", "Important", "Famille", "fichier.txt"),
			shouldErr: false,
		},
		{
			name:      "File in root of srcDir",
			relPath:   "fichier.txt",
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("home", "polocto", "Document", "fichier.txt"),
			info:      mockFileInfo{isDir: false, name: "fichier.txt"},
			expected:  filepath.Join("srv", "backup", "fichier.txt"),
			shouldErr: false,
		},
		{
			name:    "Path with spaces",
			relPath: filepath.Join("Important Project", "File with spaces.txt"),
			destDir: filepath.Join("srv", "backup"),
			filePath: filepath.Join(
				"home",
				"polocto",
				"My Documents",
				"Important Project",
				"File with spaces.txt",
			),
			info:      mockFileInfo{isDir: false, name: "File with spaces.txt"},
			expected:  filepath.Join("srv", "backup", "Important Project", "File with spaces.txt"),
			shouldErr: false,
		},
		{
			name:      "Invalid path (not a subdirectory)",
			relPath:   filepath.Join("..", "..", "..", "other", "path", "file.txt"),
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("other", "path", "file.txt"),
			info:      mockFileInfo{isDir: false, name: "file.txt"},
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Directory instead of file",
			relPath:   "Folder",
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("home", "polocto", "Document", "Folder"),
			info:      mockFileInfo{isDir: true, name: "Folder"},
			expected:  "",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileCtx := &ContextStrategy{
				relPath: tc.relPath,
				dstDir:  tc.destDir,
				info:    tc.info,
			}
			dest, err := s.FinalDirPath(fileCtx)
			if tc.shouldErr {
				assert.Error(t, err, "Expected an error for %s", tc.name)
			} else {
				assert.NoError(t, err, "Unexpected error for %s", tc.name)
				assert.Equal(t, tc.expected, dest, "Unexpected destination path for %s", tc.name)
			}
		})
	}
}

func TestDirChainStrategy_Selector(t *testing.T) {
	s := &DirChainStrategy{}
	assert.Equal(t, "dirchain", s.Selector(), "Selector should return 'dirchain'")
}

func TestDirChainStrategy_LoadConfig(t *testing.T) {
	s := &DirChainStrategy{}
	err := s.LoadConfig(map[string]interface{}{"some": "config"})
	assert.NoError(t, err, "LoadConfig should not return an error")
}

func TestDirChainStrategy_Registration(t *testing.T) {
	// Save the original registry
	originalRegistry := strategyRegistry
	defer func() { strategyRegistry = originalRegistry }()

	// Reset the registry
	strategyRegistry = make(map[string]func() Strategy)

	// Re-register the strategy (as in init())
	RegisterStrategy("dirchain", func() Strategy {
		return &DirChainStrategy{}
	})

	// Verify that the strategy is registered correctly
	strat, err := NewStrategy("dirchain")
	assert.NoError(t, err, "NewStrategy should not return an error")
	assert.Equal(t, "dirchain", strat.Selector(), "Strategy selector should be 'dirchain'")
}
