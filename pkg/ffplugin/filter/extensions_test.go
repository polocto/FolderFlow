package filter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtensionFilter(t *testing.T) {
	f := &ExtensionFilter{Extensions: []string{".txt"}}

	file := filepath.Join(t.TempDir(), "a.txt")
	os.WriteFile(file, []byte("x"), 0644)
	info, _ := os.Stat(file)

	ok, err := f.Match(file, info)
	if err != nil || !ok {
		t.Fatalf("expected extension to match")
	}
}
