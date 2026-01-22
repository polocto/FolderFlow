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

// Go standard library

// ScriptFilter runs an external script to filter files.
type ScriptFilter struct {
	ScriptPath string `yaml:"scriptPath"`
}

func (sf *ScriptFilter) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for ScriptFilter
	return nil
}

func (sf *ScriptFilter) Match(ctx Context) (bool, error) {
	// ... (script execution logic)
	return true, nil
}

func (sf *ScriptFilter) Selector() string {
	return "script"
}

func init() {
	RegisterFilter("script", func() Filter {
		return &ScriptFilter{}
	})
}
