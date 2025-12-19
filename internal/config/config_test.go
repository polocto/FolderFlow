package config

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary YAML config for testing
	configContent := `
source_dirs:
  - /path/to/source1
  - /path/to/source2
dest_dirs:
  photos:
    path: /path/to/photos
    filters:
      - type: extension
        extensions: [".jpg", ".png"]
    strategy:
      type: dirchain
    on_conflict: rename
  videos:
    path: /path/to/videos
    filters:
      - type: extension
        extensions: [".mp4", ".mov"]
regroup:
  path: /path/to/regroup
max_workers: 4
`

	// Write the config to a temporary file
	tmpfile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Load the config
	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Assert the loaded config
	assert.Equal(t, []string{"/path/to/source1", "/path/to/source2"}, cfg.SourceDirs)
	assert.Equal(t, "/path/to/regroup", cfg.Regroup.Path)
	assert.Equal(t, "", cfg.Regroup.Mode) // Not set yet, should be set by SetDefault
	assert.Equal(t, 4, cfg.MaxWorkers)

	// Test DestDir
	photosDir, ok := cfg.DestDirs["photos"]
	assert.True(t, ok)
	assert.Equal(t, "/path/to/photos", photosDir.Path)
	assert.Equal(t, "rename", photosDir.OnConflict)
	assert.Len(t, photosDir.Filters, 1)
	assert.NotNil(t, photosDir.Strategy)

	videosDir, ok := cfg.DestDirs["videos"]
	assert.True(t, ok)
	assert.Equal(t, "/path/to/videos", videosDir.Path)
	assert.Equal(t, "", videosDir.OnConflict) // Not set, should default to "rename" after SetDefaults
}

func TestSetDefault(t *testing.T) {
	cfg := &Config{
		Regroup: Regroup{
			Path: "/path/to/regroup",
		},
	}
	cfg.SetDefault()
	assert.Equal(t, "symlink", cfg.Regroup.Mode)
}

func TestDestDirSetDefaults(t *testing.T) {
	dir := &DestDir{
		Path: "/path/to/dir",
	}
	dir.SetDefaults()
	assert.Equal(t, "rename", dir.OnConflict)
}

func TestDestDirLoadPlugins(t *testing.T) {
	// Create a DestDir with FilterYAML and StrategyYAML
	dir := &DestDir{
		Path: "/path/to/dir",
		Filters: []*filter.FilterYAML{
			{
				Name: "extensions",
				Config: map[string]interface{}{
					"extensions": []string{".jpg", ".png"},
				},
			},
		},
		Strategy: &strategy.StrategyYAML{
			Name:   "dirchain",
			Config: map[string]interface{}{
				// Add any necessary config for the strategy
			},
		},
	}

	// Call LoadPlugins
	filters, strat, err := dir.LoadPlugins()
	require.NoError(t, err)

	// Assertions
	assert.NotNil(t, filters)
	assert.Len(t, filters, 1)
	assert.NotNil(t, strat)
}

// Mock types for testing
type MockFilter struct{}

func (m *MockFilter) Match(path string, info fs.FileInfo) (bool, error) {
	return true, nil
}

func (m *MockFilter) Selector() string {
	return "mock"
}

func (m *MockFilter) LoadConfig(config map[string]interface{}) error {
	return nil
}

type MockStrategy struct{}

func (m *MockStrategy) FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error) {
	return destDir, nil
}

func (m *MockStrategy) Selector() string {
	return "mock"
}

func (m *MockStrategy) LoadConfig(config map[string]interface{}) error {
	return nil
}

// Helper variables for mocking
var (
	filterToFilter     = filterToFilterImpl
	strategyToStrategy = strategyToStrategyImpl
)

func filterToFilterImpl(fy *filter.FilterYAML) (filter.Filter, error) {
	return fy.ToFilter()
}

func strategyToStrategyImpl(sy *strategy.StrategyYAML) (strategy.Strategy, error) {
	return sy.ToStrategy()
}
