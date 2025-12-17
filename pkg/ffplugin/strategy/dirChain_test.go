package strategy

import (
	"io/fs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockFileInfo crée un fs.FileInfo mock pour les tests
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
		info      fs.FileInfo // Utiliser un fs.FileInfo mock
		expected  string
		shouldErr bool
	}{
		{
			name:      "Basic relative path",
			srcDir:    "/home/polocto/Document/",
			destDir:   "/srv/backup/",
			filePath:  "/home/polocto/Document/Important/Famille/fichier.txt",
			info:      mockFileInfo{isDir: false, name: "fichier.txt"},
			expected:  "/srv/backup/Important/Famille",
			shouldErr: false,
		},
		{
			name:      "File in root of srcDir",
			srcDir:    "/home/polocto/Document/",
			destDir:   "/srv/backup/",
			filePath:  "/home/polocto/Document/fichier.txt",
			info:      mockFileInfo{isDir: false, name: "fichier.txt"},
			expected:  "/srv/backup",
			shouldErr: false,
		},
		{
			name:      "Path with spaces",
			srcDir:    "/home/polocto/My Documents/",
			destDir:   "/srv/backup/",
			filePath:  "/home/polocto/My Documents/Important Project/File with spaces.txt",
			info:      mockFileInfo{isDir: false, name: "File with spaces.txt"},
			expected:  "/srv/backup/Important Project",
			shouldErr: false,
		},
		{
			name:      "Invalid path (not a subdirectory)",
			srcDir:    "/home/polocto/Document/",
			destDir:   "/srv/backup/",
			filePath:  "/other/path/file.txt",
			info:      mockFileInfo{isDir: false, name: "file.txt"},
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Directory instead of file",
			srcDir:    "/home/polocto/Document/",
			destDir:   "/srv/backup/",
			filePath:  "/home/polocto/Document/Folder/",
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
	// Sauvegarder le registre original
	originalRegistry := strategyRegistry
	defer func() { strategyRegistry = originalRegistry }()

	// Réinitialiser le registre
	strategyRegistry = make(map[string]func() Strategy)

	// Réenregistrer la stratégie (comme dans init())
	RegisterStrategy("dirchain", func() Strategy {
		return &DirChainStrategy{}
	})

	// Vérifier que la stratégie est bien enregistrée
	strat, err := NewStrategy("dirchain")
	assert.NoError(t, err, "NewStrategy should not return an error")
	assert.Equal(t, "dirchain", strat.Selector(), "Strategy selector should be 'dirchain'")
}
