package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to create a temp file with content
func writeTempFile(t *testing.T, dir, name string, data []byte) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write temp file %s: %v", path, err)
	}
	return path
}

func TestFilesEqual_SameContent(t *testing.T) {
	dir := t.TempDir()

	content := []byte("hello world\nthis is a test")
	f1 := writeTempFile(t, dir, "file1.txt", content)
	f2 := writeTempFile(t, dir, "file2.txt", content)

	equal, err := FilesEqual(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !equal {
		t.Fatalf("expected files to be equal")
	}
}

func TestFilesEqual_DifferentContent(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte("hello world"))
	f2 := writeTempFile(t, dir, "file2.txt", []byte("hello WORLD"))

	equal, err := FilesEqual(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if equal {
		t.Fatalf("expected files to be different")
	}
}

func TestFilesEqual_DifferentSizes(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte("short"))
	f2 := writeTempFile(t, dir, "file2.txt", []byte("this is longer"))

	equal, err := FilesEqual(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if equal {
		t.Fatalf("expected files with different sizes to be different")
	}
}

func TestFilesEqual_EmptyFiles(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte{})
	f2 := writeTempFile(t, dir, "file2.txt", []byte{})

	equal, err := FilesEqual(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !equal {
		t.Fatalf("expected empty files to be equal")
	}
}

func TestFilesEqual_NonExistentFile(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte("data"))
	f2 := filepath.Join(dir, "missing.txt")

	_, err := FilesEqual(f1, f2)
	if err == nil {
		t.Fatalf("expected error for missing file, got nil")
	}
}

func TestFilesEqual_LargeFile(t *testing.T) {
	dir := t.TempDir()

	// Create content larger than one chunk (64KB)
	large := make([]byte, 200*1024)
	for i := range large {
		large[i] = byte(i % 256)
	}

	f1 := writeTempFile(t, dir, "file1.bin", large)
	f2 := writeTempFile(t, dir, "file2.bin", large)

	equal, err := FilesEqual(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !equal {
		t.Fatalf("expected large files to be equal")
	}
}
