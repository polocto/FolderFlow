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
	"fmt"
	"io"
	"os"

	"github.com/polocto/FolderFlow/internal/stats"
)

func FileHash(path string, s *stats.Stats) ([sha256.Size]byte, error) {
	var result [sha256.Size]byte

	if s != nil {
		defer s.Time(&s.Timing.Hash)()
	}

	f, err := os.Open(path)
	if err != nil {
		return result, fmt.Errorf("cannot open file %s: %v", path, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("error while closing file %s: %v", path, cerr)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return result, fmt.Errorf("error while reading file%s: %v", path, err)
	}

	// Copie le hachage calculé dans un tableau de 32 octets
	hash := h.Sum(nil)
	copy(result[:], hash)

	if s != nil {
		s.HashComputed()
	}

	return result, nil
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
