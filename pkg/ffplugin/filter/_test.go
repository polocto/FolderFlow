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

// ExampleFilter is an example custom filter.
type ExampleFilter struct {
	Custom string `yaml:"examplefilter"`
}

func (f *ExampleFilter) Match(path string, info fs.FileInfo) (bool, error) {
	// Example: Only match files larger than 1MB
	return info.Size() > 1048576, nil
}

func (f *ExampleFilter) Selector() string {
	return "ExampleFilter"
}

// LoadConfig allows setting configuration for the filter.
func (f *ExampleFilter) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for CustomFilter
	return nil
}

func init() {
	RegisterFilter("CustomFilter", func() Filter {
		return &ExampleFilter{}
	})
}
