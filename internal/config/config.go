package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type DestDir struct {
	Path       string   `yaml:"path"`
	Extensions []string `yaml:"extensions"`
}

type Config struct {
	SourceDir string             `yaml:"source_dir"`
	DestDirs  map[string]DestDir `yaml:"dest_dirs"` // Map destination name to DestDir
	Regroup   struct {
		Enable   bool   `yaml:"enable"`
		Path     string `yaml:"path"`
		LinkType string `yaml:"link_type"` // "symlink" or "hardlink"
	} `yaml:"regroup"`
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
