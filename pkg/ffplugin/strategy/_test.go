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

import "io/fs"

type ExampleStartegy struct {
	Custom string `yaml:"examplestartegy"`
}

func (s *ExampleStartegy) FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error) {
	// Example: Custom file operation logic
	return "", nil
}

func (s *ExampleStartegy) Selector() string {
	return "ExampleStartegy"
}

// LoadConfig allows setting configuration for the strategy.
func (s *ExampleStartegy) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for ExampleStartegy
	return nil
}

func init() {
	RegisterStrategy("ExampleStartegy", func() Strategy {
		return &ExampleStartegy{}
	})
}
