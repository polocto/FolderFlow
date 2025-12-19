package classify

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/fsutil"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

// matchFile checks if a file matches all the rules in DestDir.
func matchFile(path string, info fs.FileInfo, filters []filter.Filter) (bool, error) {
	// If no filters are provided, match all files
	if len(filters) == 0 {
		return true, nil
	}

	// Run all filters
	for _, f := range filters {
		matched, err := f.Match(path, info)
		if err != nil {
			slog.Warn("Filter error", "filter", f.Selector(), "path", path, "err", err)
			return false, err
		}
		if !matched {
			slog.Debug("File rejected by filter", "filter", f.Selector(), "path", path)
			return false, nil
		}
	}
	slog.Debug("File matched all filters", "path", path)
	return true, nil
}

func destPath(sourceDir, destDir, filePath string, info fs.FileInfo, strat strategy.Strategy) (string, error) {
	var destPath string
	destPath, err := strat.FinalDirPath(sourceDir, destDir, filePath, info)
	if err != nil {
		slog.Error("Strategy failed to compute destination path", "strategy", strat.Selector(), "err", err)
		return "", err
	}

	if !IsSubDirectory(destDir, destPath) {
		slog.Error("Computed destination path is outside of destination directory", "computedPath", destPath, "destDir", destDir)
		return "", fmt.Errorf("computed destination path is outside of destination directory : computedPath=%s destDir=%s", destPath, destDir)
	}

	slog.Debug("Creating directory", "path", destPath)

	destFile := filepath.Join(destPath, info.Name())
	slog.Debug("Would move", "source", filepath.Dir(filePath), "dest", destFile)

	return destFile, nil
}

func moveFile(srcPath, destPath, onConflict string, dryRun bool) error {
	slog.Debug("Moving file", "source", srcPath, "dest", destPath)

	if dryRun {
		slog.Debug("Dry run enabled, not moving file", "source", srcPath, "dest", destPath)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		slog.Error("Failed to create destination directory", "dir", filepath.Dir(destPath), "err", err)
		return err
	}

	if _, err := os.Stat(destPath); err == nil {
		switch onConflict {
		case "skip":
			slog.Debug("Destination file already exists, skipping", "dest", destPath)
			return nil
		case "overwrite":
			if err := os.Remove(destPath); err != nil {
				return err
			}
		case "rename": // rename
			if ok, err := fsutil.FilesEqual(srcPath, destPath); err != nil {
				slog.Error("Failed to compare files for equality", "source", srcPath, "dest", destPath, "err", err)
				return err
			} else if ok {
				slog.Debug("Source and destination files are identical, skipping move", "source", srcPath, "dest", destPath)
				return nil
			}
			destPath = fsutil.GetUniquePath(destPath)
		default:
			slog.Error("Unknown conflict resolution strategy", "strategy", onConflict)
			return fmt.Errorf("unknown conflict resolution strategy: %s", onConflict)
		}
	}

	if err := os.Rename(srcPath, destPath); err != nil {
		slog.Error("Failed to move file", "source", srcPath, "dest", destPath, "err", err)
		return err
	}

	return nil
}
