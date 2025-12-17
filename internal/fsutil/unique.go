package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetUniquePath returns a unique path based on destPath.
// If destPath exists, it appends a numeric suffix like "_1", "_2", etc.
func GetUniquePath(destPath string) string {
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return destPath
	}

	ext := filepath.Ext(destPath)
	base := strings.TrimSuffix(destPath, ext)
	counter := 1

	for {
		newPath := fmt.Sprintf("%s_%d%s", base, counter, ext)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}
