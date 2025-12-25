//go:build !windows

package classify

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
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

func TestExecuteRegroup_Symlink(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "regroup", "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "symlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Lstat(dst)
	if err != nil {
		t.Fatalf("lstat failed: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected symlink, got %v", info.Mode())
	}
}

func TestExecuteRegroup_SymlinkAlreadyExists(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")
	writeTestFile(t, dst, "existing")

	if err := executeRegroup(src, dst, "symlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "existing" {
		t.Fatalf("destination file should not be modified")
	}
}

func TestExecuteRegroup_Hardlink(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "hardlink"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	if !info.Mode().IsRegular() {
		t.Fatalf("expected regular file")
	}
}

func TestExecuteRegroup_Copy(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "nested", "dst.txt")

	writeTestFile(t, src, "hello")

	if err := executeRegroup(src, dst, "copy"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, dst); got != "hello" {
		t.Fatalf("copy failed")
	}
}

func TestExecuteRegroup_CreatesTargetDir(t *testing.T) {
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

func TestExecuteRegroup_InvalidMode(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	writeTestFile(t, src, "hello")

	err := executeRegroup(src, dst, "invalid-mode")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

}
