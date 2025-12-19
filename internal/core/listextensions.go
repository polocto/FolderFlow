package core

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
)

func ListAllFilesExtensions(dir string) ([]string, error) {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		slog.Error("Directory does not exist", "dir", dir)
		return nil, err
	}

	extMap := make(map[string]bool) // Use a map to track unique extensions

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Get the file extension (e.g., ".jpg")
		ext := filepath.Ext(path)
		if ext != "" { // Ignore files without extensions
			extMap[ext] = true
		} else {
			slog.Warn("File has no extension", "path", path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert the map keys to a sorted slice
	extensions := make([]string, 0, len(extMap))
	for ext := range extMap {
		extensions = append(extensions, ext)
	}
	sort.Strings(extensions) // Sort alphabetically

	return extensions, nil
}
