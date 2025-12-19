package classify

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile_Success(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "dest.txt")

	content := []byte("hello world")

	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read dest file: %v", err)
	}

	if string(got) != string(content) {
		t.Fatalf("copied content mismatch: got %q want %q", got, content)
	}
}

func TestCopyFile_SourceDoesNotExist(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "missing.txt")
	dst := filepath.Join(dir, "dest.txt")

	if err := CopyFile(src, dst); err == nil {
		t.Fatal("expected error when source file does not exist")
	}
}

func TestCopyFile_OverwriteDestination(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "dest.txt")

	if err := os.WriteFile(src, []byte("new content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != "new content" {
		t.Fatalf("destination not overwritten, got %q", got)
	}
}
