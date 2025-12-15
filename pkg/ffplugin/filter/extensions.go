package filter

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

// CustomFilter is an example custom filter.
type ExtensionFilter struct {
	Extensions []string `yaml:"extensions"`
}

func (f *ExtensionFilter) Match(path string, info fs.FileInfo) (bool, error) {
	ext := strings.ToLower(filepath.Ext(path))
	for _, allowedExt := range f.Extensions {
		if ext == strings.ToLower(allowedExt) {
			slog.Debug("Match found", "file's path", path, "matched extension", allowedExt)
			return true, nil
		}
	}
	slog.Debug("No extensions matched the file's one", "filter", f.Extensions, "path", path, "file's extension", ext)
	return false, nil
}

func (f *ExtensionFilter) Selector() string {
	return "extensions"
}

func (f *ExtensionFilter) LoadConfig(config map[string]interface{}) error {
	if exts, ok := config["extensions"].([]string); ok {
		f.Extensions = exts
	} else {
		slog.Error("Failed to load extensions", "config", config)
		return fmt.Errorf("invalid or missing 'extensions' config")
	}
	slog.Debug("Loading extensions was successful", "extensions", f.Extensions)
	return nil
}

func init() {
	RegisterFilter("extensions", func() Filter {
		return &ExtensionFilter{}
	})
}
