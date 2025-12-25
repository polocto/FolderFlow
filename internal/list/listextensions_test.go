package list

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListAllFilesExtensions(t *testing.T) {
	// Create a temporary directory for testing
	dir := t.TempDir()

	// Create test files
	testFiles := []string{"file1.txt", "file2.jpg", "file3.PNG", "file4"}
	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(dir, file))
		if err != nil {
			t.Fatal(err)
		}
		require.NoError(t, f.Close())
	}

	// Call the function
	extensions, err := ListAllFilesExtensions(dir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expected extensions (sorted)
	expected := []string{".PNG", ".jpg", ".txt"}
	if len(extensions) != len(expected) {
		t.Fatalf("Expected %d extensions, got %d", len(expected), len(extensions))
	}

	for i, ext := range expected {
		if extensions[i] != ext {
			t.Errorf("Expected %s, got %s", ext, extensions[i])
		}
	}
}
func TestListAllFilesExtensions_NonExistentDir(t *testing.T) {
	_, err := ListAllFilesExtensions("non_existent_dir")
	if err == nil {
		t.Fatal("Expected error for non-existent directory, got nil")
	}
}
func TestListAllFilesExtensions_NoExtensions(t *testing.T) {
	// Create a temporary directory for testing
	dir := t.TempDir()

	// Create test files without extensions
	testFiles := []string{"file1", "file2", "file3"}
	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(dir, file))
		if err != nil {
			t.Fatal(err)
		}
		require.NoError(t, f.Close())
	}

	// Call the function
	extensions, err := ListAllFilesExtensions(dir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expect no extensions
	if len(extensions) != 0 {
		t.Fatalf("Expected 0 extensions, got %d", len(extensions))
	}
}
