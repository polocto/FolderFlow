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
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SourceDirs []string `yaml:"source_dirs"`
	// GlobalStategy *strategy.Strategy `yaml:"global_strategy,omitempty"` // Default strategy if not specified in DestDir
	DestDirs   []DestDir `yaml:"dest_dirs"`         // Map destination name to DestDir
	Regroup    *Regroup  `yaml:"regroup,omitempty"` // Optional regroup configuration
	MaxWorkers int       `yaml:"max_workers"`       // Maximum number of concurrent workers 0, negatif or omitted means default
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config // Avoid recursion
	raw := rawConfig{}
	slog.Debug("Loading Config")
	if err := unmarshal(&raw); err != nil {
		return err
	}

	if len(raw.SourceDirs) == 0 {
		return fmt.Errorf("at least one source_dir must be specified")
	}

	if len(raw.DestDirs) == 0 {
		return fmt.Errorf("at least one dest_dir must be specified")
	}

	*cfg = Config(raw)
	slog.Debug("Config unmarshaling successful", "SourceDirs", cfg.SourceDirs, "Number of DestDirs", len(cfg.DestDirs))
	return nil
}
