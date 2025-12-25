package fsutil

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/polocto/FolderFlow/internal/stats"
)

func TestCopyFile_Success(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "dest.txt")

	content := []byte("hello world")

	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	if _, err := CopyFile(src, dst, &stats.Stats{}); err != nil {
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

	if _, err := CopyFile(src, dst, &stats.Stats{}); err == nil {
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

	if _, err := CopyFile(src, dst, &stats.Stats{}); err != nil {
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

func TestCopyFile_HashEquality(t *testing.T) {
	dir := t.TempDir()

	src1 := filepath.Join(dir, "a.txt")
	src2 := filepath.Join(dir, "b.txt")
	dst1 := filepath.Join(dir, "out1.txt")
	dst2 := filepath.Join(dir, "out2.txt")

	if err := os.WriteFile(src1, []byte("same"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(src2, []byte("same"), 0644); err != nil {
		t.Fatal(err)
	}

	h1, err := CopyFile(src1, dst1, &stats.Stats{})
	if err != nil {
		t.Fatal(err)
	}
	h2, err := CopyFile(src2, dst2, &stats.Stats{})
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(h1, h2) {
		t.Fatalf("expected hashes to match, got %s and %s", h1, h2)
	}
}

func TestCopyFile_LargeFile(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "big.bin")
	dst := filepath.Join(dir, "big.out")

	f, err := os.Create(src)
	if err != nil {
		t.Fatal(err)
	}

	// 100MB file
	buf := make([]byte, 1024*1024)
	for i := 0; i < 100; i++ {
		if _, err := f.Write(buf); err != nil {
			t.Fatal(err)
		}
	}
	f.Close()

	if _, err := CopyFile(src, dst, &stats.Stats{}); err != nil {
		t.Fatal(err)
	}

	info1, _ := os.Stat(src)
	info2, _ := os.Stat(dst)

	if info1.Size() != info2.Size() {
		t.Fatalf("size mismatch: %d vs %d", info1.Size(), info2.Size())
	}
}
