package strategy

import (
	"fmt"
	"io/fs"
	"log/slog"
)

type Strategy interface {
	// Apply performs the strategy's logic (e.g., move, link, organize)
	Apply(srcPath, destPath string, info fs.FileInfo, dryrun bool) error
	// Selector returns a unique identifier for the strategy (e.g., "date", "dirchain")
	Selector() string
	// LoadConfig allows the strategy to be configured from the YAML config
	LoadConfig(config map[string]interface{}) error
}

var strategyRegistry = make(map[string]func() Strategy)

func RegisterStrategy(name string, factory func() Strategy) {
	strategyRegistry[name] = factory
	slog.Debug("Create a new strategy", "name", name)
}

func NewStrategy(name string) (Strategy, error) {
	factory, ok := strategyRegistry[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
	return factory(), nil
}
