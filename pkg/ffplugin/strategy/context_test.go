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

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/stretchr/testify/assert"
)

// Helper to create temporary ContextFile
func newTempContextFile(t *testing.T, name string, content []byte) filehandler.Context {
	t.Helper()
	tmpFile := t.TempDir() + "/" + name
	err := os.WriteFile(tmpFile, content, 0o644)
	assert.NoError(t, err)

	ctx, err := filehandler.NewContextFile(tmpFile)
	assert.NoError(t, err)
	return ctx
}

func TestNewContext_NilFile(t *testing.T) {
	ctx, err := NewContextStrategy(nil, "/src", "/dst")
	assert.Nil(t, ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create strategy context")
}

func TestNewContext_ValidFile(t *testing.T) {
	ctxFile := newTempContextFile(t, "file.txt", []byte(""))

	srcDir := filepath.Dir(ctxFile.Path()) // parent directory of file
	dstDir := t.TempDir()

	ctx, err := NewContextStrategy(ctxFile, srcDir, dstDir)
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	// Check relative path
	expectedRel, _ := filepath.Rel(srcDir, ctxFile.Path())
	assert.Equal(t, expectedRel, ctx.PathFromSource())
	assert.Equal(t, dstDir, ctx.DstDir())
	assert.Equal(t, ctxFile, ctx.Info())
}

func TestContext_Getters(t *testing.T) {
	ctxFile := newTempContextFile(t, "file.txt", []byte(""))

	ctx, err := NewContextStrategy(ctxFile, filepath.Dir(ctxFile.Path()), t.TempDir())
	assert.NoError(t, err)

	assert.Equal(t, filepath.Base(ctxFile.Path()), filepath.Base(ctx.PathFromSource()))
	assert.Equal(t, ctx.Info(), ctxFile)
}

func TestNewContext_RelativePaths(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "folder", "sub")
	assert.NoError(t, os.MkdirAll(subDir, 0o755))

	filePath := filepath.Join(subDir, "file.txt")
	assert.NoError(t, os.WriteFile(filePath, []byte("content"), 0o644))

	ctxFile, err := filehandler.NewContextFile(filePath)
	assert.NoError(t, err)

	ctx, err := NewContextStrategy(ctxFile, tmpDir, t.TempDir())
	assert.NoError(t, err)

	expectedRel, _ := filepath.Rel(tmpDir, filePath)
	assert.Equal(t, expectedRel, ctx.PathFromSource())
}

func TestNewContext_FileAtRoot(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	assert.NoError(t, os.WriteFile(filePath, []byte("data"), 0o644))

	ctxFile, err := filehandler.NewContextFile(filePath)
	assert.NoError(t, err)

	ctx, err := NewContextStrategy(ctxFile, tmpDir, t.TempDir())
	assert.NoError(t, err)
	assert.Equal(t, "file.txt", ctx.PathFromSource())
}

func TestNewContext_FileOutsideSrcDir(t *testing.T) {
	tmpDir := t.TempDir()
	otherDir := t.TempDir() // outside srcDir

	filePath := filepath.Join(otherDir, "file.txt")
	assert.NoError(t, os.WriteFile(filePath, []byte("data"), 0o644))

	ctxFile, err := filehandler.NewContextFile(filePath)
	assert.NoError(t, err)

	ctx, err := NewContextStrategy(ctxFile, tmpDir, t.TempDir())
	assert.NoError(t, err)

	rel, _ := filepath.Rel(tmpDir, filePath)
	assert.Equal(t, rel, ctx.PathFromSource())
}
