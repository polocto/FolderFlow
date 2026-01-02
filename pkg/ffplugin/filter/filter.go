// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


package filter

import (
	"fmt"
	"io/fs"
)

// Filter defines the interface for custom file filters.
type Filter interface {
	// Match performs the filter's logic, returning true if a correspondance has been found
	Match(path string, info fs.FileInfo) (bool, error)
	// Selector returns a unique identifier for the strategy (e.g., "extension", "date")
	Selector() string
	// LoadConfig allows the filter to be configured from the YAML config
	LoadConfig(config map[string]interface{}) error
}

var filterRegistry = make(map[string]func() Filter)

func NewFilter(name string) (Filter, error) {
	factory, ok := filterRegistry[name]
	if !ok {
		return nil, fmt.Errorf("unknown filter: %s", name)
	}
	return factory(), nil
}

func RegisterFilter(name string, factory func() Filter) {
	if name == "" {
		panic("filter name cannot be empty")
	}
	if _, exists := filterRegistry[name]; exists {
		panic(fmt.Sprintf("filter '%s' is already registered", name))
	}
	filterRegistry[name] = factory
}
