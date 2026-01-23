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
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/polocto/FolderFlow/internal/classify"
	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/internal/stats"

	_ "github.com/polocto/FolderFlow/internal/filter"
	_ "github.com/polocto/FolderFlow/internal/strategy"
)

// FileInfo contient les informations nécessaires pour comparer deux fichiers
type FileInfo struct {
	Path      string
	Name      string
	Hash      [sha256.Size]byte
	HashError error
}

func TestMain(m *testing.M) {
	// Sauvegarder le logger actuel
	oldLogger := slog.Default()

	// Rediriger les logs vers io.Discard
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Exécuter les tests
	code := m.Run()

	// Rétablir le logger d'origine
	slog.SetDefault(oldLogger)

	os.Exit(code)
}

func mustCopyDir(t *testing.T, src, dst string) {
	t.Helper()
	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("copyDir: %v", err)
	}
}

func mockDirs(t *testing.T, cfg *config.Config, workDir string) string {
	t.Helper()
	if cfg == nil {
		t.Fatal("config is nil")
	}
	source := t.TempDir()
	destination := t.TempDir()

	for i, src := range cfg.SourceDirs {
		src = filepath.Join(workDir, src)
		tmp := filepath.Join(source, filepath.Base(src))
		mustCopyDir(t, src, tmp)
		cfg.SourceDirs[i] = tmp
	}
	dir, _ := os.Getwd()
	for i, destDir := range cfg.DestDirs {
		desDirPath, _ := filepath.Rel(dir, destDir.Path)
		cfg.DestDirs[i].Path = filepath.Join(destination, desDirPath)
	}

	cfg.Regroup.Path = filepath.Join(
		destination,
		filepath.Base(cfg.Regroup.Path),
	)

	return destination
}

func assertDirEquals(t *testing.T, expectedPath, resultPath string) {
	// Get files info for each result and expected directories
	expectedFiles, err := walkDir(expectedPath)
	if err != nil {
		t.Error(err)
		return
	}
	resultFiles, err := walkDir(resultPath)
	if err != nil {
		t.Error(err)
		return
	}

	// Sort list to facilitate comparaison
	sort.Slice(
		expectedFiles,
		func(i, j int) bool { return expectedFiles[i].Path < expectedFiles[j].Path },
	)
	sort.Slice(
		resultFiles,
		func(i, j int) bool { return resultFiles[i].Path < resultFiles[j].Path },
	)

	// Files name, path and files comparaison
	if len(expectedFiles) != len(resultFiles) {
		t.Errorf(
			"Not the same number of paths. Expected: %d\tResult: %d",
			len(expectedFiles),
			len(resultFiles),
		)
		return
	}

	for i := 0; i < len(expectedFiles); i++ {
		if expectedFiles[i].Path != resultFiles[i].Path ||
			expectedFiles[i].Hash != resultFiles[i].Hash {
			t.Errorf(
				"Files are different | Path1 : %s  & Path2 : %s\n",
				expectedFiles[i].Path,
				resultFiles[i].Path,
			)
			return
		}
	}
}

// walkDir walk through a directory and returns a list of FileInfo
func walkDir(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(root, path)
			hash, err := FileHash(path, nil)
			files = append(files, FileInfo{
				Path:      relPath,
				Name:      info.Name(),
				Hash:      hash,
				HashError: err,
			})
		}
		return nil
	})

	return files, err
}

func TestAllClassifyConfigs(t *testing.T) {
	// Root directory for test datas
	root := "../testdata"
	maxDepth := 2 // defin max depth level to 2

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Identify actual depth
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		depth := 0
		if relPath != "." { // Do not count actual dir as a depth level
			depth = len(filepath.SplitList(relPath))
		}

		// Stop recursion when max depth level is reached
		if depth >= maxDepth && d.IsDir() {
			return fs.SkipDir
		}

		if d.IsDir() && d.Name() == "input" {
			return fs.SkipDir
		}

		if d.Name() != "config.yaml" {
			return nil
		}

		caseDir := filepath.Dir(path)

		testName := filepath.Base(caseDir)

		// Exécuter un sous-test pour chaque config.yaml
		t.Run(testName, func(t *testing.T) {
			cfg, err := config.LoadConfig(filepath.Join(caseDir, "config.yaml"))
			if err != nil {
				t.Fatal(err)
			}

			// Arrange
			result := mockDirs(t, cfg, root)

			// Run
			var s stats.Stats
			class, err := classify.NewClassifier(*cfg, &s, false)
			if err != nil {
				t.Fatal(err)
			}

			if err := class.Classify(); err != nil {
				t.Error(err)
			}

			// Assert
			expected := filepath.Join(caseDir, "expected")

			assertDirEquals(t, expected, result)

			if t.Failed() {
				t.Logf("Actual Test Dir : %s", caseDir)
				child, _ := os.ReadDir(expected)
				t.Logf("Expected sub-dirs : %s", child)
				child, _ = os.ReadDir(result)
				t.Logf("Result sub-dirs : %s", child)
				t.Log(s.String())
			}
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
