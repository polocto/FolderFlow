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

package filter

import (
	"bytes"
	"io"
	"io/fs"
	"time"
)

// mockFileInfo is a mock implementation of fs.FileInfo for testing.
type mockFileInfo struct {
	NameVal    string
	SizeVal    int64
	ModeVal    fs.FileMode
	ModTimeVal time.Time
	IsDirVal   bool
}

func (m *mockFileInfo) Name() string       { return m.NameVal }
func (m *mockFileInfo) Size() int64        { return m.SizeVal }
func (m *mockFileInfo) Mode() fs.FileMode  { return m.ModeVal }
func (m *mockFileInfo) ModTime() time.Time { return m.ModTimeVal }
func (m *mockFileInfo) IsDir() bool        { return m.IsDirVal }
func (m *mockFileInfo) Sys() interface{}   { return nil }

type mockContext struct {
	content []byte
	info    fs.FileInfo
}

// helper method for clarity
func (mc *mockContext) IsDir() bool {
	return mc.info.IsDir()
}

func (mc *mockContext) BaseName() string {
	return mc.info.Name()
}

func (mc *mockContext) Size() int64 {
	return mc.info.Size()
}

func (mc *mockContext) ModTime() time.Time {
	return mc.info.ModTime()
}

func (mc *mockContext) Info() fs.FileInfo {
	return mc.info
}

func (mc *mockContext) WithInput(fn func(r io.Reader) error) error {
	if mc.content == nil {
		return fs.ErrInvalid
	}
	return fn(io.NopCloser(bytes.NewReader(mc.content)))
}

func (mc *mockContext) WithInputLimited(maxBytes int64, fn func(r io.Reader) error) error {
	if mc.content == nil {
		return fs.ErrInvalid
	}

	var reader io.Reader = bytes.NewReader(mc.content)
	if int64(len(mc.content)) > maxBytes {
		reader = io.LimitReader(reader, maxBytes)
	}
	return fn(reader)
}

func (mc *mockContext) ReadChunks(chunkSize int, fn func(chunk []byte) error) error {
	if mc.content == nil {
		return fs.ErrInvalid
	}

	for start := 0; start < len(mc.content); start += chunkSize {
		end := start + chunkSize
		if end > len(mc.content) {
			end = len(mc.content)
		}
		if err := fn(mc.content[start:end]); err != nil {
			return err
		}
	}
	return nil
}
