package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDestDirUnmarshal_Minimal(t *testing.T) {
	data := `
path: /tmp
strategy:
  name: dirchain
`

	var d DestDir
	if err := yaml.Unmarshal([]byte(data), &d); err != nil {
		t.Fatalf("failed to unmarshal DestDir: %v", err)
	}
}
