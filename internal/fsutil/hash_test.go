package fsutil

import (
	"path/filepath"
	"testing"
)

func TestFilesEqualHash_SameContent(t *testing.T) {
	dir := t.TempDir()

	content := []byte("hello world\nthis is a test")
	f1 := writeTempFile(t, dir, "file1.txt", content)
	f2 := writeTempFile(t, dir, "file2.txt", content)

	equal, err := FilesEqualHash(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !equal {
		t.Fatalf("expected files to be equal")
	}
}

func TestFilesEqualHash_DifferentContent(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte("hello world"))
	f2 := writeTempFile(t, dir, "file2.txt", []byte("hello WORLD"))

	equal, err := FilesEqualHash(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if equal {
		t.Fatalf("expected files to be different")
	}
}

func TestFilesEqualHash_EmptyFiles(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte{})
	f2 := writeTempFile(t, dir, "file2.txt", []byte{})

	equal, err := FilesEqualHash(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !equal {
		t.Fatalf("expected empty files to be equal")
	}
}

func TestFilesEqualHash_NonExistentFile(t *testing.T) {
	dir := t.TempDir()

	f1 := writeTempFile(t, dir, "file1.txt", []byte("data"))
	f2 := filepath.Join(dir, "missing.txt")

	_, err := FilesEqualHash(f1, f2)
	if err == nil {
		t.Fatalf("expected error for missing file, got nil")
	}
}

func TestFilesEqualHash_LargeFile(t *testing.T) {
	dir := t.TempDir()

	// Create content larger than typical buffer sizes
	large := make([]byte, 512*1024) // 512 KB
	for i := range large {
		large[i] = byte(i % 251)
	}

	f1 := writeTempFile(t, dir, "file1.bin", large)
	f2 := writeTempFile(t, dir, "file2.bin", large)

	equal, err := FilesEqualHash(f1, f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !equal {
		t.Fatalf("expected large files to be equal")
	}
}
