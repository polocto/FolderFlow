package strategy

import (
	"io/fs"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockFileInfo creates a mock fs.FileInfo for testing
type mockFileInfo struct {
	isDir bool
	name  string
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 1024 }
func (m mockFileInfo) Mode() fs.FileMode  { return 0644 }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestDirChainStrategy_FinalDirPath(t *testing.T) {
	s := &DirChainStrategy{}

	testCases := []struct {
		name      string
		srcDir    string
		destDir   string
		filePath  string
		info      fs.FileInfo
		expected  string
		shouldErr bool
	}{
		{
			name:      "Basic relative path",
			srcDir:    filepath.Join("home", "polocto", "Document"),
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("home", "polocto", "Document", "Important", "Famille", "fichier.txt"),
			info:      mockFileInfo{isDir: false, name: "fichier.txt"},
			expected:  filepath.Join("srv", "backup", "Important", "Famille"),
			shouldErr: false,
		},
		{
			name:      "File in root of srcDir",
			srcDir:    filepath.Join("home", "polocto", "Document"),
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("home", "polocto", "Document", "fichier.txt"),
			info:      mockFileInfo{isDir: false, name: "fichier.txt"},
			expected:  filepath.Join("srv", "backup"),
			shouldErr: false,
		},
		{
			name:      "Path with spaces",
			srcDir:    filepath.Join("home", "polocto", "My Documents"),
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("home", "polocto", "My Documents", "Important Project", "File with spaces.txt"),
			info:      mockFileInfo{isDir: false, name: "File with spaces.txt"},
			expected:  filepath.Join("srv", "backup", "Important Project"),
			shouldErr: false,
		},
		{
			name:      "Invalid path (not a subdirectory)",
			srcDir:    filepath.Join("home", "polocto", "Document"),
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("other", "path", "file.txt"),
			info:      mockFileInfo{isDir: false, name: "file.txt"},
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Directory instead of file",
			srcDir:    filepath.Join("home", "polocto", "Document"),
			destDir:   filepath.Join("srv", "backup"),
			filePath:  filepath.Join("home", "polocto", "Document", "Folder"),
			info:      mockFileInfo{isDir: true, name: "Folder"},
			expected:  "",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dest, err := s.FinalDirPath(tc.srcDir, tc.destDir, tc.filePath, tc.info)
			if tc.shouldErr {
				assert.Error(t, err, "Expected an error for %s", tc.name)
			} else {
				assert.NoError(t, err, "Unexpected error for %s", tc.name)
				assert.Equal(t, tc.expected, dest, "Unexpected destination path for %s", tc.name)
			}
		})
	}
}

func TestDirChainStrategy_Selector(t *testing.T) {
	s := &DirChainStrategy{}
	assert.Equal(t, "dirchain", s.Selector(), "Selector should return 'dirchain'")
}

func TestDirChainStrategy_LoadConfig(t *testing.T) {
	s := &DirChainStrategy{}
	err := s.LoadConfig(map[string]interface{}{"some": "config"})
	assert.NoError(t, err, "LoadConfig should not return an error")
}

func TestDirChainStrategy_Registration(t *testing.T) {
	// Save the original registry
	originalRegistry := strategyRegistry
	defer func() { strategyRegistry = originalRegistry }()

	// Reset the registry
	strategyRegistry = make(map[string]func() Strategy)

	// Re-register the strategy (as in init())
	RegisterStrategy("dirchain", func() Strategy {
		return &DirChainStrategy{}
	})

	// Verify that the strategy is registered correctly
	strat, err := NewStrategy("dirchain")
	assert.NoError(t, err, "NewStrategy should not return an error")
	assert.Equal(t, "dirchain", strat.Selector(), "Strategy selector should be 'dirchain'")
}
