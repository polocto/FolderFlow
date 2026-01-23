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
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileKind int

const (
	KindUnknown FileKind = iota
	KindRegular
	KindSymlink
	KindDir
)

// Context represents a filesystem object.
// Once deleted, any method call except IsDeleted() will panic.
type Context interface {
	Path() string
	setPath(newPath string)
	GetHash() ([sha256.Size]byte, error)
	IsRegular() bool
	Kind() FileKind
	delete()
	IsDeleted() bool
	fs.FileInfo
}

type ContextFile struct {
	absPath     string // File's path
	fs.FileInfo        // File's informations
	kind        FileKind
	hash        [sha256.Size]byte
	hasHash     bool
}

type deletedFileInfo struct{}

func (deletedFileInfo) Name() string       { panic("use of deleted ContextFile") }
func (deletedFileInfo) Size() int64        { panic("use of deleted ContextFile") }
func (deletedFileInfo) Mode() fs.FileMode  { panic("use of deleted ContextFile") }
func (deletedFileInfo) ModTime() time.Time { panic("use of deleted ContextFile") }
func (deletedFileInfo) IsDir() bool        { panic("use of deleted ContextFile") }
func (deletedFileInfo) Sys() any           { panic("use of deleted ContextFile") }

// NewContext creates a new Context for the given path.
// It returns an error if the path does not exist or cannot be stat-ed.
func NewContextFile(filePath string) (Context, error) {
	file, err := newContextFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func newContextFile(filePath string) (*ContextFile, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat the file %s: err=%w", filePath, err)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of %q: %w", filePath, err)
	}

	kind := KindUnknown
	switch {
	case info.Mode().IsRegular():
		kind = KindRegular
	case info.Mode()&os.ModeSymlink != 0:
		kind = KindSymlink
	case info.IsDir():
		kind = KindDir
	}

	return &ContextFile{
		absPath:  absPath,
		FileInfo: info,
		kind:     kind,
	}, nil
}

func newContextFileWithHash(filePath string, hash [sha256.Size]byte) (Context, error) {
	file, err := newContextFile(filePath)
	if err != nil {
		return nil, err
	}
	file.hash = hash
	file.hasHash = true

	return file, err
}

func (c *ContextFile) Path() string {
	if c.IsDeleted() {
		panic("use of deleted ContextFile")
	}
	return c.absPath
}

func (c *ContextFile) setPath(newPath string) {
	if c.IsDeleted() {
		panic("use of deleted ContextFile")
	}
	c.absPath = newPath
}

func (c *ContextFile) GetHash() ([sha256.Size]byte, error) {
	if c.IsDeleted() {
		panic("use of deleted ContextFile")
	}
	if c.hasHash {
		return c.hash, nil
	}

	if c.kind != KindRegular {
		return c.hash, fmt.Errorf("cannot get hash: path %q is not a regular file", c.absPath)
	}

	f, err := os.Open(c.absPath)
	if err != nil {
		return c.hash, fmt.Errorf("cannot open file %q: %w", c.absPath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Warn("failed to close file : ", "path", c.absPath, "err", err)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return c.hash, fmt.Errorf("error while reading file %q: %w", c.absPath, err)
	}

	copy(c.hash[:], h.Sum(nil))
	c.hasHash = true

	return c.hash, nil
}

func (c *ContextFile) IsRegular() bool {
	if c.IsDeleted() {
		panic("use of deleted ContextFile")
	}
	return c.Mode().IsRegular()
}

func (c *ContextFile) IsSymLink() bool {
	if c.IsDeleted() {
		panic("use of deleted ContextFile")
	}
	return c.Mode()&os.ModeSymlink != 0
}

func (c *ContextFile) IsDeleted() bool {
	return c.absPath == ""
}

func (c *ContextFile) IsSubDirectory(parent string) bool {
	if c.IsDeleted() {
		panic("use of deleted ContextFile")
	}
	// Convert to absolute path
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		slog.Error("Failed to get parent absolute path", "path", parent, "err", err)
		return false
	}

	childAbs, err := filepath.Abs(c.absPath)
	if err != nil {
		slog.Error("Failed to get child absolute path", "path", childAbs, "err", err)
		return false
	}
	// Get relative path from parent to child
	rel, err := filepath.Rel(parentAbs, childAbs)
	if err != nil {
		slog.Error("Failed to compute relative path", "parent", parentAbs, "child", childAbs, "err", err)
		return false
	}

	// If the relative path starts with "..", the child is outside the parent
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}

	return true
}
func (c *ContextFile) Kind() FileKind {
	if c.IsDeleted() {
		panic("use of deleted ContextFile")
	}
	return c.kind
}

func (c *ContextFile) delete() {
	if c.IsDeleted() {
		return
	}
	// Invalidate path
	c.absPath = ""

	// Clear cached data
	c.FileInfo = deletedFileInfo{}
	c.kind = KindUnknown
	c.hasHash = false
	c.hash = [sha256.Size]byte{}
}
