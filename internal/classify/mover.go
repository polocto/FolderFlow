package classify

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/fsutil"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

type MoveAction int

const (
	MoveSkipped MoveAction = iota
	MoveMoved
	MoveOverwritten
	MoveRenamed
	MoveCopy
	MoveSkippedIdentical
	MoveFailed
)

func destPath(sourceDir, destDir, filePath string, info fs.FileInfo, strat strategy.Strategy) (string, error) {
	finalDir, err := strat.FinalDirPath(sourceDir, destDir, filePath, info)
	if err != nil {
		slog.Error("Strategy failed to compute destination path", "strategy", strat.Selector(), "err", err)
		return "", err
	}

	if !fsutil.IsSubDirectory(destDir, finalDir) {
		slog.Error("Computed destination path is outside of destination directory", "computedPath", finalDir, "destDir", destDir)
		return "", fmt.Errorf("computed destination path is outside of destination directory : computedPath=%s destDir=%s", finalDir, destDir)
	}

	destFile := filepath.Join(finalDir, filepath.Base(filePath))

	return destFile, nil
}

func resolveConflict(srcPath, destPath, onConflict string) (string, MoveAction, error) {
	var action MoveAction
	switch onConflict {
	case "skip":
		action = MoveSkipped
	case "overwrite":
		action = MoveOverwritten
	case "rename": // rename
		if ok, err := fsutil.FilesEqual(srcPath, destPath); err != nil {
			slog.Error("Failed to compare files for equality", "source", srcPath, "dest", destPath, "err", err)
			return "", MoveFailed, err
		} else if ok {
			slog.Warn("Source and destination files are identical, skipping move", "source", srcPath, "dest", destPath)
			action = MoveSkippedIdentical
		} else {
			temp := destPath
			destPath = fsutil.GetUniquePath(destPath)
			slog.Warn("Renaming destination file to avoid conflict", "originalDest", temp, "newDest", destPath)
			action = MoveRenamed
		}
	default:
		slog.Error("Unknown conflict resolution strategy", "strategy", onConflict)
		return "", MoveFailed, fmt.Errorf("unknown conflict resolution strategy: %s", onConflict)
	}

	return destPath, action, nil
}

func moveFile(srcPath, destPath, onConflict string, dryRun bool) (MoveAction, error) {
	action := MoveMoved

	if _, err := os.Stat(destPath); err == nil {
		slog.Debug("Conflic found resolving it")
		if destPath, action, err = resolveConflict(srcPath, destPath, onConflict); err != nil {
			return action, fmt.Errorf("failed to resolve conflict at %s", destPath)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return MoveFailed, err
	}
	if dryRun {
		slog.Debug("Dry run enabled, not moving file", "source", srcPath, "dest", destPath)
		return action, nil
	}

	if action == MoveSkipped {
		return action, nil
	}

	if err := executeMove(srcPath, destPath); err != nil {
		return MoveSkipped, err
	}
	slog.Debug("File moved with success", "src", srcPath, "dst", destPath)
	return action, nil
}

func executeMove(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	// Tentative rapide et atomique
	if err := fsutil.ReplaceFile(src, dst); err == nil {
		return nil
	} else if !fsutil.IsCrossDeviceError(err) {
		// Erreur autre que EXDEV → vraie erreur
		slog.Error("Failed to rename file", "source", src, "dest", dst, "error", err)
		return err
	}
	// Fallback : copy + fsync + remove
	if err := fsutil.CopyFileAtomic(src, dst); err != nil {
		return err
	}

	return os.Remove(src)
}
