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

package fsutil

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/polocto/FolderFlow/internal/stats"
)

func TestFilesEqualHash_SameContent(t *testing.T) {
	dir := t.TempDir()

	content := []byte("hello world\nthis is a test")
	f1 := writeTempFile(t, dir, "file1.txt", content)
	f2 := writeTempFile(t, dir, "file2.txt", content)

	hash1, err := FileHash(f1, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hash2, err := FileHash(f2, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(hash1, hash2) {
		t.Fatalf("expected files to be equal")
	}
}

func TestFilesEqualHash_DifferentContent(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte("hello world"))
	f2 := writeTempFile(t, dir, "file2.txt", []byte("hello WORLD"))

	hash1, err := FileHash(f1, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hash2, err := FileHash(f2, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bytes.Equal(hash1, hash2) {
		t.Fatalf("expected files to be different")
	}
}

func TestFilesEqualHash_EmptyFiles(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte{})
	f2 := writeTempFile(t, dir, "file2.txt", []byte{})

	hash1, err := FileHash(f1, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hash2, err := FileHash(f2, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(hash1, hash2) {
		t.Fatalf("expected empty files to be equal")
	}
}

func TestFilesEqualHash_NonExistentFile(t *testing.T) {
	dir := t.TempDir()

	f1 := filepath.Join(dir, "missing.txt")

	_, err := FileHash(f1, &stats.Stats{})
	if err == nil {
		t.Fatalf("expected error for missing file, got nil")
	}
}

func TestFilesEqualHash_LargeFile(t *testing.T) {
	dir := t.TempDir()

	// Create content larger than typical buffer sizes
	large := make([]byte, 512*1024) // 512 KB
	for i := range large {
		large[i] = byte(i % 251)
	}

	f1 := writeTempFile(t, dir, "file1.bin", large)
	f2 := writeTempFile(t, dir, "file2.bin", large)

	hash1, err := FileHash(f1, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hash2, err := FileHash(f2, &stats.Stats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(hash1, hash2) {
		t.Fatalf("expected large files to be equal")
	}
}
