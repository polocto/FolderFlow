package classify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/polocto/FolderFlow/internal/config"
)

func TestRegroupFile_NoRegroupConfig(t *testing.T) {
	c := &Classifier{
		cfg: config.Config{},
	}

	err := c.regroupFile(
		"/src",
		"/src/file.txt",
		"/dst/file.txt",
		mockFileInfo{name: "file.txt"},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertError() error {
	return os.ErrInvalid
}

func TestRegroupFile_StrategyError(t *testing.T) {
	c := &Classifier{
		cfg: config.Config{
			Regroup: &config.Regroup{
				Path: "/regroup",
				Mode: "copy",
				Strategy: &mockStrategy{
					err: assertError(),
				},
			},
		},
	}

	err := c.regroupFile(
		"/src",
		"/src/file.txt",
		"/dst/file.txt",
		mockFileInfo{name: "file.txt"},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegroupFile_DryRun(t *testing.T) {
	tmp := t.TempDir()

	c := &Classifier{
		cfg: config.Config{
			Regroup: &config.Regroup{
				Path:     filepath.Join(tmp, "regroup"),
				Mode:     "copy",
				Strategy: &mockStrategy{dest: filepath.Join(tmp, "regroup")},
			},
		},
		dryRun: true,
	}

	err := c.regroupFile(
		tmp,
		filepath.Join(tmp, "src.txt"),
		filepath.Join(tmp, "final.txt"),
		mockFileInfo{name: "final.txt"},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegroupFile_Success(t *testing.T) {
	tmp := t.TempDir()

	final := filepath.Join(tmp, "final.txt")
	if err := os.WriteFile(final, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := &Classifier{
		cfg: config.Config{
			Regroup: &config.Regroup{
				Path:     filepath.Join(tmp, "regroup"),
				Mode:     "hardlink",
				Strategy: &mockStrategy{dest: filepath.Join(tmp, "regroup")},
			},
		},
	}

	err := c.regroupFile(
		tmp,
		"/src/original.txt",
		final,
		mockFileInfo{name: "final.txt"},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegroupFile_ExecuteRegroupError(t *testing.T) {
	tmp := t.TempDir()

	final := filepath.Join(tmp, "final.txt")
	if err := os.WriteFile(final, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := &Classifier{
		cfg: config.Config{
			Regroup: &config.Regroup{
				Path:     filepath.Join(tmp, "regroup"),
				Mode:     "invalid-mode",
				Strategy: &mockStrategy{dest: filepath.Join(tmp, "regroup")},
			},
		},
	}

	err := c.regroupFile(
		tmp,
		"/src/original.txt",
		final,
		mockFileInfo{name: "final.txt"},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
