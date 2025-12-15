package example

import (
	"io/fs"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

// CustomFilter is an example custom filter.
type CustomFilter struct {
	Custom string `yaml:"custom"`
}

func (f *CustomFilter) Match(path string, info fs.FileInfo) (bool, error) {
	// Example: Only match files larger than 1MB
	return info.Size() > 1048576, nil
}

func (f *CustomFilter) Selector() string {
	return "CustomFilter"
}

// LoadConfig allows setting configuration for the filter.
func (f *CustomFilter) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for CustomFilter
	return nil
}

func init() {
	filter.RegisterFilter("CustomFilter", func() filter.Filter {
		return &CustomFilter{}
	})
}
