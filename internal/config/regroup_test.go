package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRegroupUnmarshal_Defaults(t *testing.T) {
	data := `path: /tmp/regroup`

	var r Regroup
	if err := yaml.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("failed to unmarshal regroup: %v", err)
	}

	if r.Mode != "symlink" {
		t.Fatalf("expected default mode 'symlink'")
	}
}
