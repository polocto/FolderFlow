//go:build windows

package classify

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.Write([]byte(content)); err != nil {
		f.Close()
		t.Fatal(err)
	}

	if err := f.Sync(); err != nil {
		f.Close()
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	return string(b)
}

func TestExecuteRegroup_Windows_Copy(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "regroup", "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "copy"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "hello" {
		t.Fatalf("copy failed, got %q", got)
	}
}

func TestExecuteRegroup_Windows_CreatesTargetDir(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "a", "b", "c", "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "copy"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("target file not created")
	}
}

func TestExecuteRegroup_Windows_InvalidMode(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	err := executeRegroup(src, dst, "invalid-mode")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid regroup mode") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestExecuteRegroup_Windows_HardlinkOrCopy(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "hardlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "hello" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestExecuteRegroup_Windows_SymlinkFallbackChain(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "symlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "hello" {
		t.Fatalf("unexpected content after fallback chain")
	}
}

func TestExecuteRegroup_Windows_TargetAlreadyExists(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")
	writeTestFile(t, dst, "existing")

	if err := executeRegroup(src, dst, "copy"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "existing" {
		t.Fatalf("destination should not be overwritten")
	}
}
