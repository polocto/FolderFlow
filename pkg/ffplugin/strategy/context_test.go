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

// strategy/context_test.go
// strategy/context_test.go
package strategy

import (
	"crypto/sha256"
	"io/fs"
	"path/filepath"
	"testing"
	"time"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/stretchr/testify/assert"
)

// --------------------------
// Mock FileInfo for testing
// --------------------------
type mockFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() fs.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return m.sys }

// --------------------------
// Mock filehandler.Context
// --------------------------
type mockFileHandlerContext struct {
	path      string
	info      fs.FileInfo
	kind      filehandler.FileKind
	isRegular bool
	hash      [sha256.Size]byte
	err       error
}

func (m *mockFileHandlerContext) Path() string {
	return m.path
}

func (m *mockFileHandlerContext) Info() fs.FileInfo {
	return m.info
}

func (m *mockFileHandlerContext) GetHash() ([sha256.Size]byte, error) {
	return m.hash, m.err
}
func (m *mockFileHandlerContext) IsRegular() bool {
	return m.isRegular
}
func (m *mockFileHandlerContext) Kind() filehandler.FileKind {
	return m.kind
}

// --------------------------
// Tests for NewContext
// --------------------------
func TestNewContext_NilFile(t *testing.T) {
	ctx, err := NewContextStrategy(nil, "/src", "/dst")
	assert.Nil(t, ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create strategy context")
}

func TestNewContext_InvalidRelPath(t *testing.T) {
	file := &mockFileHandlerContext{
		path: "/file.txt",
		info: mockFileInfo{name: "file.txt"},
	}
	// Pass srcDir as an invalid path to trigger filepath.Rel error (using invalid UTF-8 char)
	_, err := NewContextStrategy(file, string([]byte{0xff}), "/dst")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create strategy context")
}

func TestNewContext_ValidFile(t *testing.T) {
	file := &mockFileHandlerContext{
		path: "/home/user/src/folder/file.txt",
		info: mockFileInfo{name: "file.txt"},
	}

	srcDir := "/home/user/src"
	dstDir := "/home/user/dst"

	ctx, err := NewContextStrategy(file, srcDir, dstDir)
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	// Check relative path
	assert.Equal(t, "folder/file.txt", ctx.PathFromSource())
	assert.Equal(t, dstDir, ctx.DstDir())
	assert.Equal(t, file.info, ctx.Info())
}

// --------------------------
// Tests for getters
// --------------------------
func TestContext_Getters(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}
	ctx := &ContextStrategy{
		path:    "/some/path/file.txt",
		relPath: "file.txt",
		dstDir:  "/dst",
		info:    info,
	}

	assert.Equal(t, "file.txt", ctx.PathFromSource())
	assert.Equal(t, "/dst", ctx.DstDir())
	assert.Equal(t, info, ctx.Info())
}

// --------------------------
// Test complex relative paths
// --------------------------
func TestNewContext_RelativePaths(t *testing.T) {
	file := &mockFileHandlerContext{
		path: "/home/user/src/folder/sub/file.txt",
		info: mockFileInfo{name: "file.txt"},
	}

	srcDir := "/home/user/src"
	dstDir := "/backup"

	ctx, err := NewContextStrategy(file, srcDir, dstDir)
	assert.NoError(t, err)
	assert.Equal(t, "folder/sub/file.txt", ctx.PathFromSource())

	// Destination path computation for finalDir
	final := filepath.Join(ctx.DstDir(), filepath.Dir(ctx.PathFromSource()))
	assert.Equal(t, "/backup/folder/sub", final)
}

// --------------------------
// Test edge case: file at root of srcDir
// --------------------------
func TestNewContext_FileAtRoot(t *testing.T) {
	file := &mockFileHandlerContext{
		path: "/home/user/src/file.txt",
		info: mockFileInfo{name: "file.txt"},
	}

	srcDir := "/home/user/src"
	dstDir := "/backup"

	ctx, err := NewContextStrategy(file, srcDir, dstDir)
	assert.NoError(t, err)
	assert.Equal(t, "file.txt", ctx.PathFromSource())

	// Final destination path
	final := filepath.Join(ctx.DstDir(), filepath.Dir(ctx.PathFromSource()))
	assert.Equal(t, "/backup", final)
}

// --------------------------
// Test source outside srcDir
// --------------------------
func TestNewContext_FileOutsideSrcDir(t *testing.T) {
	file := &mockFileHandlerContext{
		path: "/home/user/other/file.txt",
		info: mockFileInfo{name: "file.txt"},
	}

	srcDir := "/home/user/src"
	dstDir := "/backup"

	ctx, err := NewContextStrategy(file, srcDir, dstDir)
	assert.NoError(t, err)
	assert.Equal(t, "../other/file.txt", ctx.PathFromSource())
}
