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

package filter_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/polocto/FolderFlow/internal/filter"
	"github.com/stretchr/testify/assert"
)

// Helper to create a temporary ContextFile
func newTempContextFile(t *testing.T, name string, content []byte) filehandler.Context {
	t.Helper()
	tmpFile := t.TempDir() + "/" + name
	err := os.WriteFile(tmpFile, content, 0o644)
	assert.NoError(t, err)

	ctx, err := filehandler.NewContextFile(tmpFile)
	assert.NoError(t, err)
	return ctx
}

func TestNewContext(t *testing.T) {
	ctxFile := newTempContextFile(t, "file.txt", []byte(""))

	ctx, err := filter.NewContextFilter(ctxFile)
	assert.NoError(t, err)
	assert.Equal(t, "file.txt", ctx.BaseName())
	assert.False(t, ctx.IsDir())
	assert.Equal(t, int64(0), ctx.Size())
}

func TestWithInput(t *testing.T) {
	content := []byte("Hello World")
	ctxFile := newTempContextFile(t, "file.txt", content)

	ctx, err := filter.NewContextFilter(ctxFile)
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = ctx.WithInput(func(r io.Reader) error {
		_, err := io.Copy(&buf, r)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, content, buf.Bytes())
}

func TestWithInputLimited(t *testing.T) {
	content := []byte("Hello World")
	ctxFile := newTempContextFile(t, "file.txt", content)

	ctx, err := filter.NewContextFilter(ctxFile)
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = ctx.WithInputLimited(5, func(r io.Reader) error {
		_, err := io.Copy(&buf, r)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello"), buf.Bytes())
}

func TestReadChunks(t *testing.T) {
	content := []byte("0123456789")
	ctxFile := newTempContextFile(t, "file.txt", content)

	ctx, err := filter.NewContextFilter(ctxFile)
	assert.NoError(t, err)

	var chunks [][]byte
	err = ctx.ReadChunks(4, func(b []byte) error {
		cpy := make([]byte, len(b))
		copy(cpy, b)
		chunks = append(chunks, cpy)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{
		[]byte("0123"),
		[]byte("4567"),
		[]byte("89"),
	}, chunks)
}

func TestReadChunks_ErrorPropagation(t *testing.T) {
	content := []byte("012345")
	ctxFile := newTempContextFile(t, "file.txt", content)

	ctx, err := filter.NewContextFilter(ctxFile)
	assert.NoError(t, err)

	expectedErr := errors.New("stop early")
	err = ctx.ReadChunks(3, func(b []byte) error {
		return expectedErr
	})
	assert.ErrorIs(t, err, expectedErr)
}

func TestWithInput_OnDir(t *testing.T) {
	tmpDir := t.TempDir()
	ctxFile, err := filehandler.NewContextFile(tmpDir)
	assert.NoError(t, err)

	ctx, err := filter.NewContextFilter(ctxFile)
	assert.NoError(t, err)

	err = ctx.WithInput(func(r io.Reader) error { return nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot open directory")
}
