package filter

import (
	"io/fs"
	// Go standard library
)

// ScriptFilter runs an external script to filter files.
type ScriptFilter struct {
	ScriptPath string `yaml:"scriptPath"`
}

func (sf *ScriptFilter) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for ScriptFilter
	return nil
}

func (sf *ScriptFilter) Match(path string, info fs.FileInfo) (bool, error) {
	// ... (script execution logic)
	return true, nil
}

func (sf *ScriptFilter) Selector() string {
	return "script"
}

func init() {
	RegisterFilter("script", func() Filter {
		return &ScriptFilter{}
	})
}
