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

package folderflow_test

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
	"github.com/polocto/FolderFlow/internal/fsutil"
	"github.com/polocto/FolderFlow/internal/stats"
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
	if err := fsutil.CopyDir(src, dst); err != nil {
		t.Fatalf("copyDir: %v", err)
	}
}

func mockDirs(t *testing.T, cfg *config.Config) string {
	t.Helper()
	if cfg == nil {
		t.Fatal("config is nil")
	}
	source := t.TempDir()
	destination := t.TempDir()

	for i, src := range cfg.SourceDirs {
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

func assertDirEquals(t *testing.T, expected, result string) {
	files1, err := walkDir(expected)
	if err != nil {
		t.Error(err)
		return
	}
	files2, err := walkDir(result)
	if err != nil {
		t.Error(err)
		return
	}
	// Tri des listes pour faciliter la comparaison
	sort.Slice(files1, func(i, j int) bool { return files1[i].Path < files1[j].Path })
	sort.Slice(files2, func(i, j int) bool { return files2[i].Path < files2[j].Path })

	// Comparaison des chemins et noms de fichiers
	i := 0

	if len(files1) != len(files2) {
		t.Errorf("Not the same number of paths. Expected: %d\tResult: %d", len(files1), len(files2))
		return
	}

	for i < len(files1) {
		if files1[i].Path != files2[i].Path || files1[i].Hash != files2[i].Hash {
			t.Errorf("Files are different | Path1 : %s  & Path2 : %s\n", files1[i].Path, files2[i].Path)
		}
		i++
	}
}

// walkDir parcourt un répertoire et retourne une liste de FileInfo
func walkDir(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(root, path)
			hash, err := fsutil.FileHash(path, nil)
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
	root := "testdata"
	maxDepth := 2

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculer la profondeur actuelle
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		depth := 0
		if relPath != "." { // Éviter de compter le dossier racine comme un niveau
			depth = len(filepath.SplitList(relPath))
		}

		// Arrêter la récursion si la profondeur maximale est atteinte
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
		expected := filepath.Join(caseDir, "expected")

		cfg, err := config.LoadConfig(filepath.Join(caseDir, "config.yaml"))
		if err != nil {
			t.Fatal(err)
		}

		// Arrange
		result := mockDirs(t, cfg)
		var s stats.Stats
		class, err := classify.NewClassifier(*cfg, &s, false)
		if err != nil {
			t.Fatal(err)
		}

		if err := class.Classify(); err != nil {
			t.Error(err)
		}

		// Assert
		// assertDirEquals(t, filepath.Join(expected, "destination"), filepath.Join(result, "destination"))
		// assertDirEquals(t, filepath.Join(expected, "regrouped"), filepath.Join(result, "regrouped"))
		assertDirEquals(t, expected, result)
		if t.Failed() {
			t.Logf("Actual Test Dir : %s", caseDir)
			child, _ := os.ReadDir(expected)
			t.Logf("Expected sub-dirs : %s", child)
			child, _ = os.ReadDir(result)
			t.Logf("Result sub-dirs : %s", child)
			t.Log(s.String())
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
