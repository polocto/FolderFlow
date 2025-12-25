package classify

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/internal/stats"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

func newClassifier(cfg config.Config, dryRun bool) *Classifier {
	s := &stats.Stats{}
	c, _ := NewClassifier(cfg, s, dryRun)
	return c
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestNewClassifier(t *testing.T) {
	c, err := NewClassifier(config.Config{}, &stats.Stats{}, false)
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("classifier is nil")
	}
}

func TestClassify_NoSources(t *testing.T) {
	c := newClassifier(config.Config{
		SourceDirs: nil,
		DestDirs:   []config.DestDir{{Path: "/tmp"}},
	}, false)

	if err := c.Classify(); err == nil {
		t.Fatal("expected error")
	}
}

func TestClassify_NoDestinations(t *testing.T) {
	c := newClassifier(config.Config{
		SourceDirs: []string{"/tmp"},
		DestDirs:   nil,
	}, false)

	if err := c.Classify(); err == nil {
		t.Fatal("expected error")
	}
}

func TestClassify_SkipInvalidSources(t *testing.T) {
	tmp := t.TempDir()

	c := newClassifier(config.Config{
		SourceDirs: []string{"", filepath.Join(tmp, "missing")},
		DestDirs:   []config.DestDir{{Path: tmp}},
	}, false)

	if err := c.Classify(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

////////////////////////
/// processSourceDir ///
////////////////////////

func TestProcessSourceDir_WalkError(t *testing.T) {
	c := newClassifier(config.Config{
		MaxWorkers: 1,
	}, false)

	err := c.processSourceDir(string([]byte{0}))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProcessSourceDir_SkipDirs(t *testing.T) {
	tmp := t.TempDir()
	os.Mkdir(filepath.Join(tmp, ".git"), 0755)
	os.Mkdir(filepath.Join(tmp, "node_modules"), 0755)

	c := newClassifier(config.Config{
		MaxWorkers: 1,
	}, false)

	if err := c.processSourceDir(tmp); err != nil {
		t.Fatal(err)
	}
}

func TestProcessFile_NoMatch(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "a.txt")
	writeFile(t, src)

	c := newClassifier(config.Config{
		DestDirs: []config.DestDir{
			{
				Path:    filepath.Join(tmp, "dest"),
				Filters: []filter.Filter{&mockFilter{match: false}},
				Strategy: &mockStrategy{
					dest: filepath.Join(tmp, "dest"),
				},
			},
		},
	}, false)

	err := c.processFile(tmp, src, mockFileInfo{name: "a.txt"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessFile_MoveError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("Skipping on %s: POSIX permissions are not enforced", runtime.GOOS)

	}

	tmp := t.TempDir()
	src := filepath.Join(tmp, "a.txt")
	writeFile(t, src)

	destDir := filepath.Join(tmp, "dest")
	if err := os.Mkdir(destDir, 0555); err != nil { // read-only on POSIX
		t.Fatal(err)
	}

	c := newClassifier(config.Config{
		DestDirs: []config.DestDir{
			{
				Path: destDir,
				Filters: []filter.Filter{
					&mockFilter{match: true},
				},
				Strategy: &mockStrategy{
					dest: destDir,
				},
				OnConflict: "overwrite",
			},
		},
	}, false)

	err := c.processFile(tmp, src, mockFileInfo{name: "a.txt"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProcessFile_RegroupEnabled(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("File move + regroup chain is not reliable on Windows due to file locking")
	}
	tmp := t.TempDir()
	src := filepath.Join(tmp, "a.txt")
	writeFile(t, src)

	dest := filepath.Join(tmp, "dest")

	c := newClassifier(config.Config{
		DestDirs: []config.DestDir{
			{
				Path:    dest,
				Filters: []filter.Filter{&mockFilter{match: true}},
				Strategy: &mockStrategy{
					dest: dest,
				},
				OnConflict: "skip",
			},
		},
		Regroup: &config.Regroup{
			Path:     filepath.Join(tmp, "regroup"),
			Mode:     "copy",
			Strategy: &mockStrategy{dest: filepath.Join(tmp, "regroup")},
		},
	}, false)

	if err := c.processFile(tmp, src, mockFileInfo{name: "a.txt"}); err != nil {
		t.Fatal(err)
	}
}
