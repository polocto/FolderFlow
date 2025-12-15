package example

import (
	"io/fs"

	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

type CustomStrategy struct {
	Custom string `yaml:"custom"`
}

func (s *CustomStrategy) Apply(srcPath, destPath string, info fs.FileInfo, dryrun bool) error {
	// Example: Custom file operation logic
	return nil
}

func (s *CustomStrategy) Selector() string {
	return "CustomStrategy"
}

// LoadConfig allows setting configuration for the strategy.
func (s *CustomStrategy) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for CustomStrategy
	return nil
}

func init() {
	strategy.RegisterStrategy("CustomStrategy", func() strategy.Strategy {
		return &CustomStrategy{}
	})
}
