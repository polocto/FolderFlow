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
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
)

type Context interface {
	// helper method for clarity
	IsDir() bool
	BaseName() string
	Size() int64
	ModTime() time.Time
	Info() fs.FileInfo
	WithInput(fn func(r io.Reader) error) error
	WithInputLimited(maxBytes int64, fn func(r io.Reader) error) error
	ReadChunks(chunkSize int, fn func([]byte) error) error
}

// ContextFilter represents a file or directory to be filtered.
type ContextFilter struct {
	path string
	info fs.FileInfo
}

// --- sync.Pool for reusable buffers ---
var chunkPool = sync.Pool{
	New: func() interface{} {
		// Default buffer size, can be overridden per ReadChunks call
		return make([]byte, 4096)
	},
}

// NewContext creates a new Context for the given path.
// It returns an error if the path does not exist or cannot be stat-ed.
func NewContextFilter(file filehandler.Context) (Context, error) {
	if file == nil {
		return nil, fmt.Errorf("cannot create filter context because file context is nil: %w", filehandler.ErrContextIsNil)
	}

	return &ContextFilter{
		path: file.Path(),
		info: file,
	}, nil
}

// helper method for clarity
func (c *ContextFilter) IsDir() bool        { return c.info.IsDir() }
func (c *ContextFilter) BaseName() string   { return c.info.Name() }
func (c *ContextFilter) Size() int64        { return c.info.Size() }
func (c *ContextFilter) ModTime() time.Time { return c.info.ModTime() }
func (c *ContextFilter) Info() fs.FileInfo  { return c.info }

// WithInput opens the file for reading and passes it to the callback.
func (c *ContextFilter) WithInput(fn func(r io.Reader) error) error {
	if c.IsDir() {
		return fmt.Errorf("cannot open directory %q for reading", c.path)
	}
	f, err := os.Open(c.path)
	if err != nil {
		return fmt.Errorf("cannot open file %q: %w", c.path, err)
	}
	defer f.Close()

	return fn(f)
}

// WithInputLimited reads only the first maxBytes bytes of the file.
func (c *ContextFilter) WithInputLimited(maxBytes int64, fn func(r io.Reader) error) error {
	return c.WithInput(func(f io.Reader) error {
		err := fn(io.LimitReader(f, maxBytes))
		if err != nil {
			return fmt.Errorf("error reading limited input from %q: %w", c.BaseName(), err)
		}
		return nil
	})
}

// ReadChunks reads the file in chunks of the given size and passes each chunk to fn.
//
// For best performance, it is recommended to use a chunkSize of 4096 bytes,
// which matches the internal buffer pool size. Using this size allows buffer reuse
// and reduces memory allocations. Other sizes work normally but will allocate new buffers.

func (c *ContextFilter) ReadChunks(chunkSize int, fn func([]byte) error) error {
	return c.WithInput(func(r io.Reader) error {
		// Get a buffer from the pool
		var buf []byte
		if chunkSize == 4096 {
			buf = chunkPool.Get().([]byte)
			defer chunkPool.Put(buf)
		} else {
			// Custom-sized buffer, allocate normally
			buf = make([]byte, chunkSize)
		}
		for {
			n, err := r.Read(buf)
			if n > 0 {
				if err := fn(buf[:n]); err != nil {
					return err
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
		return nil
	})
}
