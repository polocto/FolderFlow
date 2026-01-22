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

package filehandler_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()

	path := filepath.Join(dir, name)

	require.NoError(t, os.WriteFile(path, content, 0644))

	return path
}

func tempSubDir(t *testing.T, parent, name string) string {
	t.Helper()

	path := filepath.Join(parent, name)
	require.NoError(t, os.MkdirAll(path, 0755))
	return path
}

func helloWorld() (str string) {
	str = "Hello world!"
	return str
}

func TestNewContextExistingFile(t *testing.T) {
	filePath := tempFile(t, t.TempDir(), "file.txt", []byte(helloWorld()))

	file, err := filehandler.NewContextFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file == nil {
		t.Fatalf("unexpected error: %v", filehandler.ErrContextIsNil)
	}

	// Validate Info fields
	if file.Name() != "file.txt" {
		t.Fatalf("expected filename %q, got %q", "file.txt", file.Name())
	}

	expectedSize := int64(len(helloWorld()))
	if file.Size() != expectedSize {
		t.Fatalf("expected size %d, got %d", expectedSize, file.Size())
	}

	if !file.Mode().IsRegular() {
		t.Fatalf("expected regular file, got mode %v", file.Mode())
	}

}

func TestNewContextNoneExistingFile(t *testing.T) {
	filePath := filepath.Join("tmp", "testing", "folderflow", "new")

	file, err := filehandler.NewContextFile(filePath)
	if err == nil {
		t.Fatalf("expected error but nil")
	}
	if file != nil {
		t.Fatalf("unexpected context: %s", file.Name())
	}

	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("unexpected error: %v", err)
	}

}

func TestNewContextExistingFileWithWrongPermission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission test is not reliable on Windows")
	}

	dir := t.TempDir()

	filePath := tempFile(t, dir, "file.txt", []byte(helloWorld()))

	// Remove execute permission from directory
	require.NoError(t, os.Chmod(dir, 0600))
	t.Cleanup(func() {
		_ = os.Chmod(dir, 0700) // restore for cleanup
	})

	file, err := filehandler.NewContextFile(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if file != nil {
		t.Fatal("expected nil context")
	}

	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected fs.ErrPermission, got %v", err)
	}

}
