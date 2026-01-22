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

package classify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()

	path := filepath.Join(dir, name)

	require.NoError(t, os.WriteFile(path, content, 0644))

	return path
}

func tempSubDir(t *testing.T, parent, name string) string {
	t.Helper()

	path := filepath.Join(parent, name)
	require.NoError(t, os.MkdirAll(path, 0755))
	return path
}
