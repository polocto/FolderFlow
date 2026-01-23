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

package classify

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/internal/stats"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

func newClassifier(cfg config.Config, dryRun bool) *Classifier {
	s := &stats.Stats{}
	c, _ := NewClassifier(cfg, s, dryRun)
	return c
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestNewClassifier_EmptyConfiguration(t *testing.T) {
	c, err := NewClassifier(config.Config{}, &stats.Stats{}, false)
	if err == nil {
		t.Fatal("no errors is return on an empty configuration")
	}
	if c != nil {
		t.Fatal("classifier is not nil while config has no parameters")
	}
}

////////////////////////
/// processSourceDir ///
////////////////////////

func TestProcessSourceDir_WalkError(t *testing.T) {
	t.SkipNow()
	c := newClassifier(config.Config{
		MaxWorkers: 1,
	}, false)

	err := c.processSourceDir(string([]byte{0}))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProcessFile_NoMatch(t *testing.T) {
	t.SkipNow()
	tmp := t.TempDir()
	src := filepath.Join(tmp, "a.txt")
	writeFile(t, src)

	c := newClassifier(config.Config{
		DestDirs: []config.DestDir{
			{
				Path:    filepath.Join(tmp, "dest"),
				Filters: []filter.Filter{&mockFilter{match: false}},
				Strategy: &mockStrategy{
					dest: filepath.Join(tmp, "dest"),
				},
			},
		},
	}, false)

	err := c.processFile(tmp, src)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessFile_MoveError(t *testing.T) {
	t.SkipNow()
	if runtime.GOOS == "windows" {
		t.Skipf("Skipping on %s: POSIX permissions are not enforced", runtime.GOOS)
	}

	tmp := t.TempDir()
	src := filepath.Join(tmp, "a.txt")
	writeFile(t, src)

	destDir := filepath.Join(tmp, "dest")
	if err := os.Mkdir(destDir, 0o555); err != nil { // read-only on POSIX
		t.Fatal(err)
	}
	c := newClassifier(config.Config{
		DestDirs: []config.DestDir{
			{
				Path: destDir,
				Filters: []filter.Filter{
					&mockFilter{match: true},
				},
				Strategy: &mockStrategy{
					dest: destDir,
				},
				OnConflict: "overwrite",
			},
		},
	}, false)

	err := c.processFile(tmp, src)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProcessFile_RegroupEnabled(t *testing.T) {
	t.SkipNow()
	if runtime.GOOS == "windows" {
		t.Skip("File move + regroup chain is not reliable on Windows due to file locking")
	}
	tmp := t.TempDir()
	src := filepath.Join(tmp, "a.txt")
	writeFile(t, src)

	dest := filepath.Join(tmp, "dest")

	c := newClassifier(config.Config{
		DestDirs: []config.DestDir{
			{
				Path:    dest,
				Filters: []filter.Filter{&mockFilter{match: true}},
				Strategy: &mockStrategy{
					dest: dest,
				},
				OnConflict: "skip",
			},
		},
		Regroup: &config.Regroup{
			Path:     filepath.Join(tmp, "regroup"),
			Mode:     "copy",
			Strategy: &mockStrategy{dest: filepath.Join(tmp, "regroup")},
		},
	}, false)

	if err := c.processFile(tmp, src); err != nil {
		t.Fatal(err)
	}
}
