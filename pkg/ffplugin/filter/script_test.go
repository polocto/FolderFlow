package filter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScriptFilter_Match(t *testing.T) {
	f := &ScriptFilter{ScriptPath: "/bin/true"}

	file := filepath.Join(t.TempDir(), "a.txt")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0644))
	info, _ := os.Stat(file)

	ok, err := f.Match(file, info)
	if err != nil || !ok {
		t.Fatalf("script filter should return true")
	}
}
