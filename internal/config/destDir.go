// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


package config

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
	"gopkg.in/yaml.v3"
)

type DestDir struct {
	Name       string            `yaml:"name"`
	Path       string            `yaml:"path"`
	Filters    []filter.Filter   // File extensions to include
	Strategy   strategy.Strategy // "date", "dirchain", etc.
	OnConflict string            `yaml:"on_conflict,omitempty"` // "skip", "overwrite", "rename"
}

// func (d *DestDir) LoadPlugins() ([]filter.Filter, strategy.Strategy, error) {

// 	return filters, *strat, nil
// }

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (d *DestDir) UnmarshalYAML(node *yaml.Node) error {
	// Define a temporary struct to unmarshal the YAML into
	type tempDestDir struct {
		Name       string         `yaml:"name"`
		Path       string         `yaml:"path"`
		Filters    []filterConfig `yaml:"filters,omitempty"`
		Strategy   strategyConfig `yaml:"strategy,omitempty"`
		OnConflict string         `yaml:"on_conflict,omitempty"`
	}
	var temp tempDestDir
	if err := node.Decode(&temp); err != nil {
		return err
	}
	if temp.Path == "" {
		return fmt.Errorf("dest_dir path cannot be empty")
	}

	switch temp.OnConflict {
	case "":
		temp.OnConflict = "rename" // default value
	case "skip", "overwrite", "rename":
		// valid options
	default:
		return fmt.Errorf("invalid on_conflict option '%s', must be 'skip', 'overwrite' or 'rename'", temp.OnConflict)
	}

	// Copy simple fields
	d.Name = temp.Name
	if path, err := filepath.Abs(temp.Path); err != nil {
		return err
	} else {
		d.Path = path
	}
	d.OnConflict = temp.OnConflict
	// Load filters
	for _, fc := range temp.Filters {
		f, err := filter.NewFilter(fc.Name)
		if err != nil {
			return fmt.Errorf("failed to create filter '%s': %v", fc.Name, err)
		}
		if err := f.LoadConfig(fc.Config); err != nil {
			return fmt.Errorf("failed to load config for filter '%s': %v", fc.Name, err)
		}
		d.Filters = append(d.Filters, f)
	}
	slog.Debug("Loaded filters successfully", "dst", temp.Name, "Number of loaded filters", len(d.Filters))
	// Load strategy
	strat, err := strategy.NewStrategy(temp.Strategy.Name)
	if err != nil {
		return fmt.Errorf("failed to create strategy '%s': %v", temp.Strategy.Name, err)
	}
	if err := strat.LoadConfig(temp.Strategy.Config); err != nil {
		return fmt.Errorf("failed to load config for strategy '%s': %v", temp.Strategy.Name, err)
	}

	d.Strategy = strat
	slog.Debug("Loaded strategy successfully", "dst", temp.Name, "strategy", d.Strategy.Selector())
	return nil
}

// filterConfig is a helper struct to unmarshal filter configurations
type filterConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

// strategyConfig is a helper struct to unmarshal strategy configurations
type strategyConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}
