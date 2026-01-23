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

package filehandler

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetUniquePath_FileDoesNotExist(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "file.txt")
	unique := GetUniquePath(path)

	if unique != path {
		t.Fatalf("expected %q, got %q", path, unique)
	}
}

func TestGetUniquePath_FileExists(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	unique := GetUniquePath(path)
	expected := filepath.Join(dir, "file_1.txt")

	if unique != expected {
		t.Fatalf("expected %q, got %q", expected, unique)
	}
}

func TestGetUniquePath_MultipleConflicts(t *testing.T) {
	dir := t.TempDir()

	base := filepath.Join(dir, "file.txt")

	// Create file.txt, file_1.txt, file_2.txt
	for i := 0; i <= 2; i++ {
		var p string
		if i == 0 {
			p = base
		} else {
			p = filepath.Join(dir, fmt.Sprintf("file_%d.txt", i))
		}
		if err := os.WriteFile(p, []byte("data"), 0o644); err != nil {
			t.Fatalf("failed to create file %s: %v", p, err)
		}
	}

	unique := GetUniquePath(base)
	expected := filepath.Join(dir, "file_3.txt")

	if unique != expected {
		t.Fatalf("expected %q, got %q", expected, unique)
	}
}

func TestGetUniquePath_NoExtension(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "file")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	unique := GetUniquePath(path)
	expected := filepath.Join(dir, "file_1")

	if unique != expected {
		t.Fatalf("expected %q, got %q", expected, unique)
	}
}

func TestGetUniquePath_NestedExtension(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "archive.tar.gz")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	unique := GetUniquePath(path)
	expected := filepath.Join(dir, "archive.tar_1.gz")

	if unique != expected {
		t.Fatalf("expected %q, got %q", expected, unique)
	}
}
