package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Regroup struct {
	Path string `yaml:"path"`
	Mode string `yaml:"mode,omitempty"` // "symlink", "hardlink" or "copy"
	// Strategy *strategy.Strategy `yaml:"strategy,omitempty"` // "date", "dirchain", etc.
}

type Config struct {
	SourceDirs []string `yaml:"source_dirs"`
	// GlobalStategy *strategy.Strategy `yaml:"global_strategy,omitempty"` // Default strategy if not specified in DestDir
	DestDirs   map[string]DestDir `yaml:"dest_dirs"` // Map destination name to DestDir
	Regroup    `yaml:"regroup,omitempty"`
	MaxWorkers int `yaml:"max_workers"` // Maximum number of concurrent workers
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

func (c *Config) SetDefault() {

	if c.Regroup.Mode == "" && c.Regroup.Path != "" {
		c.Regroup.Mode = "symlink"
	}
}
