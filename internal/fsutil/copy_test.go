package fsutil

import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"

	"github.com/polocto/FolderFlow/internal/stats"
)

//
// --------------------
// Test helpers
// --------------------
//

func writeTestFile(t *testing.T, path string, content []byte) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) []byte {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed for %s: %v", path, err)
	}
	return b
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected %s to NOT exist", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("unexpected stat error for %s: %v", path, err)
	}
}

//
// --------------------
// CopyFile tests
// --------------------
//

func TestCopyFile_NoStats(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	content := []byte("hello world")
	writeTestFile(t, src, content)

	hash, err := CopyFile(src, dst, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := sha256.Sum256(content)
	if !bytes.Equal(hash, expected[:]) {
		t.Fatalf("hash mismatch")
	}

	if got := readTestFile(t, dst); !bytes.Equal(got, content) {
		t.Fatalf("copied content mismatch")
	}
}

func TestCopyFile_WithStats(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	content := []byte("hello stats")
	writeTestFile(t, src, content)

	s := &stats.Stats{}

	hash, err := CopyFile(src, dst, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := sha256.Sum256(content)
	if !bytes.Equal(hash, expected[:]) {
		t.Fatalf("hash mismatch")
	}

	if s.Operations.Copy == 0 {
		t.Fatalf("FilesCopied counter not incremented")
	}
	if s.Hash.Computed == 0 {
		t.Fatalf("HashComputed counter not incremented")
	}
}

//
// --------------------
// CopyFileAtomic tests
// --------------------
//

func TestCopyFileAtomic_Success_Invariants(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")
	tmpFile := dst + ".tmp"

	content := []byte("atomic success")
	writeTestFile(t, src, content)

	if err := CopyFileAtomic(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Invariants
	assertExists(t, src)        // source must remain
	assertExists(t, dst)        // destination must exist
	assertNotExists(t, tmpFile) // tmp file must be cleaned

	if got := readTestFile(t, dst); !bytes.Equal(got, content) {
		t.Fatalf("destination content mismatch")
	}
}

func TestCopyFileAtomic_Overwrite(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")
	tmpFile := dst + ".tmp"

	writeTestFile(t, src, []byte("new content"))
	writeTestFile(t, dst, []byte("old content"))

	if err := CopyFileAtomic(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertExists(t, src)
	assertExists(t, dst)
	assertNotExists(t, tmpFile)

	if got := readTestFile(t, dst); string(got) != "new content" {
		t.Fatalf("destination not overwritten")
	}
}

func TestCopyFileAtomic_Error_CleansTmpAndKeepsSource(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dstDir := filepath.Join(tmp, "dst") // directory => invalid target
	tmpFile := dstDir + ".tmp"

	writeTestFile(t, src, []byte("fail case"))

	if err := os.Mkdir(dstDir, 0755); err != nil {
		t.Fatal(err)
	}

	err := CopyFileAtomic(src, dstDir)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Invariants after failure
	assertExists(t, src)        // source must remain
	assertNotExists(t, tmpFile) // tmp must be cleaned
	assertExists(t, dstDir)
}

func TestCopyFileAtomic_SourceUnchanged(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	content := []byte("immutable source")
	writeTestFile(t, src, content)

	if err := CopyFileAtomic(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readTestFile(t, src); !bytes.Equal(got, content) {
		t.Fatalf("source file was modified")
	}
}
