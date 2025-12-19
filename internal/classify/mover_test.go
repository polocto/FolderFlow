package classify

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

//
// --------------------
// Test helpers
// --------------------
//

// mockFileInfo implements fs.FileInfo
type mockFileInfo struct {
	name string
	mode fs.FileMode
	size int64
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() fs.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return m.mode.IsDir() }
func (m mockFileInfo) Sys() any           { return nil }

// --------------------
// mockFilter
// --------------------

type mockFilter struct {
	match bool
	err   error
}

func (m *mockFilter) Match(string, fs.FileInfo) (bool, error) {
	return m.match, m.err
}

func (m *mockFilter) Selector() string {
	return "mock"
}

func (m *mockFilter) LoadConfig(map[string]interface{}) error {
	return nil
}

// --------------------
// mockStrategy
// --------------------

type mockStrategy struct {
	dest string
	err  error
}

func (m *mockStrategy) FinalDirPath(_, _, _ string, _ fs.FileInfo) (string, error) {
	return m.dest, m.err
}

func (m *mockStrategy) Selector() string {
	return "mock"
}

func (m *mockStrategy) LoadConfig(map[string]interface{}) error {
	return nil
}

//
// --------------------
// matchFile tests
// --------------------
//

func TestMatchFile_NoFilters(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}

	ok, err := matchFile("file.txt", info, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected file to match when no filters are provided")
	}
}

func TestMatchFile_FilterRejects(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}

	filters := []filter.Filter{
		&mockFilter{match: false},
	}

	ok, err := matchFile("file.txt", info, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected file to be rejected")
	}
}

func TestMatchFile_FilterError(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}

	filters := []filter.Filter{
		&mockFilter{match: false, err: fs.ErrInvalid},
	}

	_, err := matchFile("file.txt", info, filters)
	if err == nil {
		t.Fatal("expected error from filter")
	}
}

//
// --------------------
// destPath tests
// --------------------
//

func TestDestPath_ValidSubdirectory(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	info := mockFileInfo{name: "file.txt"}

	strat := &mockStrategy{
		dest: filepath.Join(destDir, "subdir"),
	}

	result, err := destPath(
		srcDir,
		destDir,
		filepath.Join(srcDir, "file.txt"),
		info,
		strat,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(destDir, "subdir", "file.txt")
	if result != expected {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestDestPath_OutsideDestinationDir(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	info := mockFileInfo{name: "file.txt"}

	strat := &mockStrategy{
		dest: filepath.Dir(destDir), // outside destDir
	}

	_, err := destPath(
		srcDir,
		destDir,
		filepath.Join(srcDir, "file.txt"),
		info,
		strat,
	)
	if err == nil {
		t.Fatal("expected error when destination is outside destDir")
	}
}

//
// --------------------
// moveFile tests
// --------------------
//

func TestMoveFile_DryRun(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := moveFile(src, dst, "overwrite", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(src); err != nil {
		t.Fatal("source file should still exist in dry run")
	}
}

func TestMoveFile_Overwrite(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	os.WriteFile(src, []byte("src"), 0644)
	os.WriteFile(dst, []byte("dst"), 0644)

	if err := moveFile(src, dst, "overwrite", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("source file should be moved")
	}
}

func TestMoveFile_Skip(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	os.WriteFile(src, []byte("src"), 0644)
	os.WriteFile(dst, []byte("dst"), 0644)

	if err := moveFile(src, dst, "skip", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(src); err != nil {
		t.Fatal("source file should still exist when skipping")
	}
}
