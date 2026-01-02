// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


import "io/fs"

// ExampleFilter is an example custom filter.
type ExampleFilter struct {
	Custom string `yaml:"examplefilter"`
}

func (f *ExampleFilter) Match(path string, info fs.FileInfo) (bool, error) {
	// Example: Only match files larger than 1MB
	return info.Size() > 1048576, nil
}

func (f *ExampleFilter) Selector() string {
	return "ExampleFilter"
}

// LoadConfig allows setting configuration for the filter.
func (f *ExampleFilter) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for CustomFilter
	return nil
}

func init() {
	RegisterFilter("CustomFilter", func() Filter {
		return &ExampleFilter{}
	})
}
