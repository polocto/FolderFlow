//go:build windows

package classify

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/fsutil"
)

func executeRegroup(source, target, mode string) error {

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	switch mode {
	case "symlink":
		// Try symlink
		rel := source
		if r, err := filepath.Rel(filepath.Dir(target), source); err == nil {
			rel = r
		}
		if err := os.Symlink(rel, target); err == nil {
			break
		} else if errors.Is(err, fs.ErrExist) {
			slog.Debug("Regroup target already exists, skipping", "path", target)
			return nil
		}
		slog.Warn("Symlink failed on Windows, falling back to hardlink",
			"source", source,
			"dest", target,
		)

		// Try hardlink
		if err := os.Link(source, target); err == nil {
			break
		} else if errors.Is(err, fs.ErrExist) {
			slog.Debug("Regroup target already exists, skipping", "path", target)
			return nil
		}
		slog.Warn("Hardlink failed on Windows, falling back to copy",
			"source", source,
			"dest", target,
		)
		if err := fsutil.CopyFileAtomic(source, target); err != nil {
			slog.Error("Failed to copy file for regrouping", "from", source, "to", target, "err", err)
			return err
		}

	case "hardlink":
		if err := os.Link(source, target); err == nil {
			break
		} else if errors.Is(err, fs.ErrExist) {
			slog.Debug("Regroup target already exists, skipping", "path", target)
			return nil
		}
		slog.Warn("Hardlink failed on Windows, falling back to copy",
			"source", source,
			"dest", target,
		)
		if err := fsutil.CopyFileAtomic(source, target); err != nil {
			return err
		}

	case "copy":
		if err := fsutil.CopyFileAtomic(source, target); err != nil {
			return err
		}

	default:
		return ErrInvalidRegroupMode(mode)
	}
	slog.Debug(
		"File regrouped successfully",
		"source", source,
		"target", target,
		"mode", mode,
	)

	return nil
}
