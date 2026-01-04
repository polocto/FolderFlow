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
	"path/filepath"
	"testing"
)

func TestIsSubDirectory(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tests := []struct {
		name   string
		parent string
		child  string
		want   bool
	}{
		{
			name:   "same directory absolute",
			parent: filepath.Join(wd, "test_data/destination/images"),
			child:  filepath.Join(wd, "test_data/destination/images"),
			want:   true,
		},
		{
			name:   "same directory relative vs dot relative",
			parent: "./test_data/destination/images",
			child:  "test_data/destination/images",
			want:   true,
		},
		{
			name:   "child is subdirectory",
			parent: "test_data/destination",
			child:  "test_data/destination/images",
			want:   true,
		},
		{
			name:   "child escapes with parent traversal",
			parent: "test_data/destination/images",
			child:  "test_data/destination/images/../..",
			want:   false,
		},
		{
			name:   "child is sibling directory",
			parent: "test_data/destination/images",
			child:  "test_data/destination/documents",
			want:   false,
		},
		{
			name:   "completely outside directory",
			parent: "test_data/destination/images",
			child:  "/tmp",
			want:   false,
		},
		{
			name:   "dot child",
			parent: "test_data/destination/images",
			child:  "test_data/destination/images/.",
			want:   true,
		},
		{
			name:   "cleaned traversal inside",
			parent: "test_data/destination",
			child:  "test_data/destination/a/../b",
			want:   true,
		},
		{
			name:   "pre-dot child",
			parent: "test_data/destination/images",
			child:  "./test_data/destination/images/.",
			want:   true,
		},
		{
			name:   "pre-dot parent",
			parent: "./test_data/destination/images",
			child:  "test_data/destination/images",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSubDirectory(tt.parent, tt.child)
			if got != tt.want {
				t.Errorf(
					"IsSubDirectory(%q, %q) = %v, want %v",
					tt.parent, tt.child, got, tt.want,
				)
			}
		})
	}
}
