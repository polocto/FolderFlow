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
	"fmt"
	"io/fs"
	"path/filepath"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
)

type Context interface {
	PathFromSource() string
	DstDir() string
	Info() fs.FileInfo
}

type ContextStrategy struct {
	path    string      // File's path
	relPath string      // Relative Path from the Source Directory
	dstDir  string      // Destination Directory
	info    fs.FileInfo // File's informations
}

// NewContext creates a new Context for the given path.
// It returns an error if the path does not exist or cannot be stat-ed.
func NewContextStrategy(file filehandler.Context, srcDir, dstDir string) (Context, error) {
	if file == nil {
		return nil, fmt.Errorf("cannot create strategy context: %w", filehandler.ErrContextIsNil)
	}

	subPath, err := filepath.Rel(srcDir, file.Path())
	if err != nil {
		return nil, fmt.Errorf("cannot create strategy context: %w", err)
	}

	return &ContextStrategy{
		path:    file.Path(),
		relPath: subPath,
		dstDir:  dstDir,
		info:    file,
	}, nil
}

func (ctx *ContextStrategy) PathFromSource() string {
	return ctx.relPath
}

func (ctx *ContextStrategy) DstDir() string {
	return ctx.dstDir
}

func (ctx *ContextStrategy) Info() fs.FileInfo {
	return ctx.info
}
