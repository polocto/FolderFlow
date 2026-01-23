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
	"io/fs"
	"time"

	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

type mockStrategy struct {
	selector     string
	finalDirPath string
	err          error
}

// FinalDirPath computes the final directory path for a file based on the strategy.
func (m *mockStrategy) FinalDirPath(ctx strategy.Context) (string, error) {
	return m.finalDirPath, m.err
}

// Selector returns a unique identifier for the strategy (e.g., "date", "dirchain")
func (m *mockStrategy) Selector() string {
	return m.selector
}

// LoadConfig allows the strategy to be configured from the YAML config
func (m *mockStrategy) LoadConfig(config map[string]interface{}) error {
	return nil
}

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

type mockContext struct {
	pathFromSource       string
	destinationDirectory string
	info                 fs.FileInfo
}

func (mc *mockContext) PathFromSource() string {
	return mc.pathFromSource
}

func (mc *mockContext) DstDir() string {
	return mc.destinationDirectory
}

func (mc *mockContext) Info() fs.FileInfo {
	return mc.info
}
