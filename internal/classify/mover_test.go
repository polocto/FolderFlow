package classify

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	name   string
	match  bool
	err    error
	called *int
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

// --------------------
// destPath tests
// --------------------
func TestDestPath_Success(t *testing.T) {
	strat := &mockStrategy{
		dest: "/dest/sub",
	}

	out, err := destPath(
		"/src",
		"/dest",
		"/src/file.txt",
		mockFileInfo{name: "file.txt"},
		strat,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join("/dest/sub", "file.txt")
	if out != expected {
		t.Fatalf("expected %s, got %s", expected, out)
	}
}

func TestDestPath_StrategyError(t *testing.T) {
	strat := &mockStrategy{
		err: errors.New("boom"),
	}

	_, err := destPath(
		"/src",
		"/dest",
		"/src/file.txt",
		mockFileInfo{name: "file.txt"},
		strat,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDestPath_OutsideDestination(t *testing.T) {
	strat := &mockStrategy{
		dest: "/evil",
	}

	_, err := destPath(
		"/src",
		"/dest",
		"/src/file.txt",
		mockFileInfo{name: "file.txt"},
		strat,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveConflict_Skip(t *testing.T) {
	dst, action, err := resolveConflict("src", "dst", "skip")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != MoveSkipped {
		t.Fatalf("expected MoveSkipped, got %v", action)
	}
	if dst != "dst" {
		t.Fatal("destination path changed unexpectedly")
	}
}

func TestResolveConflict_Overwrite(t *testing.T) {
	dst, action, err := resolveConflict("src", "dst", "overwrite")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != MoveOverwritten {
		t.Fatalf("expected MoveOverwritten, got %v", action)
	}
	if dst != "dst" {
		t.Fatal("destination path changed unexpectedly")
	}
}

func TestResolveConflict_Rename_Identical(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "a.txt")
	dst := filepath.Join(tmp, "b.txt")

	require.NoError(t, os.WriteFile(src, []byte("same"), 0644))
	require.NoError(t, os.WriteFile(dst, []byte("same"), 0644))

	newDst, action, err := resolveConflict(src, dst, "rename")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != MoveSkippedIdentical {
		t.Fatalf("expected MoveSkippedIdentical, got %v", action)
	}
	if newDst != dst {
		t.Fatal("destination path should not change")
	}
}

func TestResolveConflict_Rename_Different(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "a.txt")
	dst := filepath.Join(tmp, "b.txt")

	os.WriteFile(src, []byte("A"), 0644)
	os.WriteFile(dst, []byte("B"), 0644)

	newDst, action, err := resolveConflict(src, dst, "rename")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != MoveRenamed {
		t.Fatalf("expected MoveRenamed, got %v", action)
	}
	if newDst == dst {
		t.Fatal("expected destination to be renamed")
	}
}

func TestResolveConflict_UnknownMode(t *testing.T) {
	_, action, err := resolveConflict("src", "dst", "???")

	if err == nil {
		t.Fatal("expected error")
	}
	if action != MoveFailed {
		t.Fatalf("expected MoveFailed, got %v", action)
	}
}

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

// --------------------
// moveFile tests
// --------------------
func TestMoveFile_NoConflict(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	os.WriteFile(src, []byte("data"), 0644)

	action, err := moveFile(src, dst, "skip", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != MoveMoved {
		t.Fatalf("expected MoveMoved, got %v", action)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Fatal("destination file missing")
	}
}

func TestMoveFile_ConflictSkip(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	os.WriteFile(src, []byte("A"), 0644)
	os.WriteFile(dst, []byte("B"), 0644)

	action, err := moveFile(src, dst, "skip", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != MoveSkipped {
		t.Fatalf("expected MoveSkipped, got %v", action)
	}
}

func TestMoveFile_DryRun(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	os.WriteFile(src, []byte("data"), 0644)

	action, err := moveFile(src, dst, "skip", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != MoveMoved {
		t.Fatalf("expected MoveMoved, got %v", action)
	}
	if _, err := os.Stat(dst); err == nil {
		t.Fatal("file should not have been moved")
	}
}

func TestMoveFile_Overwrite(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	os.WriteFile(src, []byte("src"), 0644)
	os.WriteFile(dst, []byte("dst"), 0644)

	if _, err := moveFile(src, dst, "overwrite", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("source file should be moved")
	}
}

func TestMoveFile_StatError(t *testing.T) {
	action, err := moveFile(
		"/does/not/matter",
		string([]byte{0}),
		"skip",
		false,
	)

	if err == nil {
		t.Fatal("expected error")
	}
	if action != MoveFailed {
		t.Fatalf("expected MoveFailed, got %v", action)
	}
}

func TestExecuteMove_Success(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "sub/dst.txt")

	os.WriteFile(src, []byte("hello"), 0644)

	if err := executeMove(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Fatal("destination file missing")
	}
}

func TestMoveFile_Skip(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	os.WriteFile(src, []byte("src"), 0644)
	os.WriteFile(dst, []byte("dst"), 0644)

	if _, err := moveFile(src, dst, "overwrite", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(src); err != nil {
		t.Fatal("source file should still exist when skipping")
	}
}
