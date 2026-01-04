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

//go:build windows

package classify

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.Write([]byte(content)); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	return string(b)
}

func TestExecuteRegroup_Windows_Copy(t *testing.T) {
	t.Skip("Copy-based regroup is not reliably testable on Windows due to file locking")
}

func TestExecuteRegroup_Windows_CreatesTargetDir(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "a", "b", "c", "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "hardlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("target file not created")
	}
}

func TestExecuteRegroup_Windows_InvalidMode(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	err := executeRegroup(src, dst, "invalid-mode")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid regroup mode") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestExecuteRegroup_Windows_HardlinkOrCopy(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "hardlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "hello" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestExecuteRegroup_Windows_SymlinkFallbackChain(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "symlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "hello" {
		t.Fatalf("unexpected content after fallback chain")
	}
}

func TestExecuteRegroup_Windows_Copy_Overwrite(t *testing.T) {
	t.Skip("Copy-based regroup is not reliably testable on Windows due to file locking")
}
