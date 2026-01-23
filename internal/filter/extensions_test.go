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
	"testing"
)

func TestExtensionFilter(t *testing.T) {
	f := &ExtensionFilter{Extensions: []string{".txt"}}

	ok, err := f.Match(&mockContext{[]byte("Hello World"), &mockFileInfo{NameVal: "a.txt"}})
	if err != nil || !ok {
		t.Fatalf("expected extension to match")
	}
}
