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

package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRegroupUnmarshal_Defaults(t *testing.T) {
	data := `path: /tmp/regroup`

	var r Regroup
	if err := yaml.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("failed to unmarshal regroup: %v", err)
	}

	if r.Mode != "symlink" {
		t.Fatalf("expected default mode 'symlink'")
	}
}
