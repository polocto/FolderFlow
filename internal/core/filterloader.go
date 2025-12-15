package core

import (
	"fmt"
	"plugin" // Go standard library

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

// LoadGoPlugin loads a Go plugin and configures it.
func LoadGoPlugin(pluginPath string, config map[string]interface{}) (filter.Filter, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %v", err)
	}
	sym, err := p.Lookup("NewFilter")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup NewFilter: %v", err)
	}
	newFilter, ok := sym.(func() filter.Filter)
	if !ok {
		return nil, fmt.Errorf("NewFilter is not a func() filter.Filter")
	}
	pluginFilter := newFilter()
	if err := pluginFilter.LoadConfig(config); err != nil {
		return nil, fmt.Errorf("failed to set plugin config: %v", err)
	}
	return pluginFilter, nil
}
