package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigValidation(t *testing.T) {
	data := `
source_dirs: []
dest_dirs: []
`

	var cfg Config
	err := yaml.Unmarshal([]byte(data), &cfg)

	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}

func TestConfigValidation_OK(t *testing.T) {
	data := `
source_dirs:
  - /tmp
dest_dirs:
  - name: out
    path: /dest
    strategy:
      name: dirchain
`

	var cfg Config
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
