package filter

import (
	"fmt"
	"io/fs"
	"log/slog"
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
	filterRegistry[name] = factory
	slog.Debug("Register a new filter", "name", name)
}
