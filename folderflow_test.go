package folderflow_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/polocto/FolderFlow/internal/fsutil"
)

func mustCopyDir(t *testing.T, src, dst string) {
	t.Helper()
	if err := fsutil.CopyDir(src, dst); err != nil {
		t.Fatalf("copyDir: %v", err)
	}
}

func assertDirEquals(t *testing.T, expected, actual string) {
	t.Helper()
	// Walk both dirs
	// Compare:
	// - paths
	// - file names
	// - file contents (if needed)
}

func TestAllConfigs(t *testing.T) {
	root := "testdata"

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() != "config.yaml" {
			return nil
		}

		caseDir := filepath.Dir(path)
		expected := filepath.Join(caseDir, "expected")

		t.Run(strings.TrimPrefix(caseDir, root+string(os.PathSeparator)), func(t *testing.T) {
			t.Parallel()

			work := t.TempDir()

			// Arrange
			src := filepath.Join(work, "source")
			mustCopyDir(t, filepath.Join(root, "input/source"), src)

			// Act
			if err := RunFlow(path, src, work); err != nil {
				t.Fatal(err)
			}

			// Assert
			assertDirEquals(t, expected, work)
		})

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
