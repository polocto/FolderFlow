package config

import (
	"log/slog"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

type DestDir struct {
	Path       string                 `yaml:"path"`
	Filters    []*filter.FilterYAML   `yaml:"filters,omitempty"`     // File extensions to include
	Strategy   *strategy.StrategyYAML `yaml:"strategy,omitempty"`    // "date", "dirchain", etc.
	OnConflict string                 `yaml:"on_conflict,omitempty"` // "skip", "overwrite", "rename"
}

func (d *DestDir) LoadPlugins() ([]filter.Filter, strategy.Strategy, error) {
	slog.Debug("Loading plugins")
	var filters []filter.Filter
	for _, fy := range d.Filters {
		filter, err := fy.ToFilter()
		if err != nil {
			return nil, nil, err
		}
		filters = append(filters, filter)
	}

	var strategy strategy.Strategy
	if d.Strategy != nil {
		strat, err := d.Strategy.ToStrategy()
		if err != nil {
			return nil, nil, err
		}
		strategy = strat
	}

	return filters, strategy, nil
}

func (d *DestDir) SetDefaults() {
	if d.OnConflict == "" {
		d.OnConflict = "rename"
	}
}
