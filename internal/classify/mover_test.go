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
	"errors"
	"path/filepath"
	"testing"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
	"github.com/stretchr/testify/require"
)

// --------------------
// Mock Strategy
// --------------------

type mockStrategy struct {
	dest string
	err  error
}

func (m *mockStrategy) FinalDirPath(ctx strategy.Context) (string, error) {
	return m.dest, m.err
}

func (m *mockStrategy) Selector() string                        { return "mock" }
func (m *mockStrategy) LoadConfig(map[string]interface{}) error { return nil }

// --------------------
// destPath tests
// --------------------

func TestDestPath_Success(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	file := tempFile(t, srcDir, "file.txt", []byte("Hello world!"))

	fhCtx, err := filehandler.NewContextFile(file)
	require.NoError(t, err)

	mockStrat := &mockStrategy{
		dest: filepath.Join(destDir, "subdir", "file.txt"),
	}

	out, err := destPath(fhCtx, srcDir, destDir, mockStrat)
	require.NoError(t, err)

	expected := filepath.Join(destDir, "subdir", "file.txt")
	require.Equal(t, expected, out)
}

func TestDestPath_StrategyError(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	file := tempFile(t, srcDir, "file.txt", []byte("Hello world!"))

	fhCtx, err := filehandler.NewContextFile(file)
	require.NoError(t, err)

	mockStrat := &mockStrategy{err: errors.New("boom")}

	_, err = destPath(fhCtx, srcDir, destDir, mockStrat)
	require.Error(t, err)
}

func TestDestPath_OutsideDestination(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	file := tempFile(t, srcDir, "file.txt", []byte("Hello world!"))

	fhCtx, err := filehandler.NewContextFile(file)
	require.NoError(t, err)

	mockStrat := &mockStrategy{dest: filepath.Join("/evil", "file.txt")}

	_, err = destPath(fhCtx, srcDir, destDir, mockStrat)
	require.Error(t, err)
}

// --------------------
// Conflict resolution tests
// --------------------

func TestResolveConflict_Skip(t *testing.T) {
	tmp := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("data"))
	dst := tempFile(t, tmp, "dst.txt", []byte("data"))

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	fhDst, err := filehandler.NewContextFile(dst)
	require.NoError(t, err)

	dst, action, err := resolveConflict(fhSrc, fhDst, "skip")
	require.NoError(t, err)
	require.Equal(t, MoveSkipped, action)
	require.Equal(t, "src.txt", filepath.Base(src))
}

func TestResolveConflict_Overwrite(t *testing.T) {
	tmp := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("data"))
	dst := tempFile(t, tmp, "dst.txt", []byte("data"))

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	fhDst, err := filehandler.NewContextFile(dst)
	require.NoError(t, err)

	dst, action, err := resolveConflict(fhSrc, fhDst, "overwrite")
	require.NoError(t, err)
	require.Equal(t, MoveOverwritten, action)
	require.Equal(t, dst, fhDst.Path())
}

func TestResolveConflict_Rename_Identical(t *testing.T) {
	tmp := t.TempDir()
	tmp2 := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("data"))
	dst := tempFile(t, tmp2, "src.txt", []byte("data"))

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	fhDst, err := filehandler.NewContextFile(dst)
	require.NoError(t, err)

	newDst, action, err := resolveConflict(fhSrc, fhDst, "rename")
	require.NoError(t, err)
	require.Equal(t, MoveSkippedIdentical, action)
	expectedPath := filepath.Join(filepath.Dir(dst), "src.txt")
	require.Equal(t, expectedPath, newDst)
}

func TestResolveConflict_Rename_Different(t *testing.T) {
	tmp := t.TempDir()
	tmp2 := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("data"))
	dst := tempFile(t, tmp2, "src.txt", []byte("data2"))

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	fhDst, err := filehandler.NewContextFile(dst)
	require.NoError(t, err)

	newDst, action, err := resolveConflict(fhSrc, fhDst, "rename")
	require.NoError(t, err)
	require.Equal(t, MoveRenamed, action)
	require.NotEqual(t, dst, newDst)
}

func TestResolveConflict_UnknownMode(t *testing.T) {
	tmp := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("data"))
	dst := tempFile(t, tmp, "src.txt", []byte("data"))

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	fhDst, err := filehandler.NewContextFile(dst)
	require.NoError(t, err)

	_, action, err := resolveConflict(fhSrc, fhDst, "???")
	require.Error(t, err)
	require.Equal(t, MoveFailed, action)
}

// --------------------
// moveFile tests
// --------------------

func TestMoveFile_NoConflict(t *testing.T) {
	tmp := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("data"))
	dst := filepath.Join(tmp, "dst.txt")

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	action, copy, err := moveFile(fhSrc, dst, "rename", false)
	require.NoError(t, err)
	require.Equal(t, MoveMoved, action)
	require.FileExists(t, dst)
	require.Nil(t, copy)
}

func TestMoveFile_ConflictSkip(t *testing.T) {
	tmp := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("A"))
	dst := tempFile(t, tmp, "dst.txt", []byte("B"))

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	action, copy, err := moveFile(fhSrc, dst, "skip", false)
	require.NoError(t, err)
	require.Equal(t, MoveSkipped, action)
	require.Nil(t, copy)
}

func TestMoveFile_Overwrite(t *testing.T) {
	tmp := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("src"))
	dst := tempFile(t, tmp, "dst.txt", []byte("dst"))

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	action, copy, err := moveFile(fhSrc, dst, "overwrite", false)
	require.NoError(t, err)
	require.NoFileExists(t, src)
	require.FileExists(t, dst)
	require.Equal(t, MoveOverwritten, action)
	require.Nil(t, copy)
}

func TestMoveFile_DryRun(t *testing.T) {
	tmp := t.TempDir()
	src := tempFile(t, tmp, "src.txt", []byte("data"))
	dst := filepath.Join(tmp, "dst.txt")

	fhSrc, err := filehandler.NewContextFile(src)
	require.NoError(t, err)
	action, copy, err := moveFile(fhSrc, dst, "overwrite", true)
	require.NoError(t, err)
	require.Equal(t, MoveMoved, action)
	require.FileExists(t, src)
	require.NoFileExists(t, dst)
	require.Nil(t, copy)
}
