package config

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
	"gopkg.in/yaml.v3"
)

type Regroup struct {
	Path     string            `yaml:"path"`
	Mode     string            `yaml:"mode,omitempty"` // "symlink", "hardlink" or "copy"
	Strategy strategy.Strategy // "date", "dirchain", etc.
}

func (r *Regroup) UnmarshalYAML(node *yaml.Node) error {
	var raw regroupYAML

	if err := node.Decode(&raw); err != nil {
		return err
	}

	if raw.Path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	switch raw.Mode {
	case "":
		slog.Debug("No mode specified for Regroup, using default 'symlink'")
		raw.Mode = "symlink"
	case "symlink", "hardlink", "copy":
		// ok
	default:
		return fmt.Errorf(
			"invalid mode '%s', must be 'symlink', 'hardlink' or 'copy'",
			raw.Mode,
		)
	}

	var strat strategy.Strategy
	if raw.Strategy.Name == "" {
		slog.Debug("No strategy specified for Regroup, using default 'dirchain' strategy")
		s, err := strategy.NewStrategy("dirchain")
		if err != nil {
			return fmt.Errorf(
				"failed to create default strategy 'dirchain': %w",
				err,
			)
		}
		strat = s
	} else {

		s, err := strategy.NewStrategy(raw.Strategy.Name)
		if err != nil {
			return fmt.Errorf(
				"failed to create strategy '%s': %w",
				raw.Strategy.Name,
				err,
			)
		}

		if err := s.LoadConfig(raw.Strategy.Config); err != nil {
			return fmt.Errorf(
				"failed to load config for strategy '%s': %w",
				raw.Strategy.Name,
				err,
			)
		}

		strat = s
	}

	if path, err := filepath.Abs(raw.Path); err != nil {
		return err
	} else {
		r.Path = path
	}

	r.Mode = raw.Mode
	r.Strategy = strat

	slog.Debug("Loaded regroup successfully",
		"path", r.Path,
		"mode", r.Mode,
		"strategy", r.Strategy.Selector(),
	)

	return nil
}

type regroupYAML struct {
	Path     string         `yaml:"path"`
	Mode     string         `yaml:"mode,omitempty"`
	Strategy strategyConfig `yaml:"strategy,omitempty"`
}
