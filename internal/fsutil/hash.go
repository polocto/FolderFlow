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

package fsutil

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"

	"github.com/polocto/FolderFlow/internal/stats"
)

func FileHash(path string, s *stats.Stats) ([]byte, error) {
	if s != nil {
		defer s.Time(&s.Timing.Hash)()
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	} else if s != nil {
		s.HashComputed()
	}
	return h.Sum(nil), nil
}

func HashEqual(f1, f2 []byte, s *stats.Stats) bool {
	if !bytes.Equal(f1, f2) {
		return false
	}

	if s != nil {
		s.HashVerified()
	}
	return true
}
