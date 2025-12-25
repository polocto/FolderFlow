package strategy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDateStrategy(t *testing.T) {
	s := &DateStrategy{}
	s.LoadConfig(map[string]interface{}{})

	file := filepath.Join(t.TempDir(), "a.txt")
	os.WriteFile(file, []byte("x"), 0644)
	info, _ := os.Stat(file)

	path, err := s.FinalDirPath("src", "dest", file, info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if path == "" {
		t.Fatalf("expected non-empty path")
	}
}
