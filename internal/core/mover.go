package core

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/utils"
)

func Classify(sourceDir string, destDirs map[string][]string, verbose bool, dryRun bool) error {
	slog.Info("Starting classification")
	return filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			slog.Warn("Skipping file due to error", "path", path, "err", err)
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)

			for dest, extensions := range destDirs {

				for _, e := range extensions {
					if ext == e {
						relPath, err := filepath.Rel(sourceDir, filepath.Dir(path))
						if err != nil {
							return err
						}
						destPath := filepath.Join(dest, relPath)

						slog.Debug("Creating directory", "path", destPath)
						if !dryRun {
							if err := os.MkdirAll(destPath, 0755); err != nil {
								return err
							}
						}
						destFile := filepath.Join(destPath, info.Name())
						slog.Debug("Would move", "source", path, "dest", destFile)
						if !dryRun {
							if err := os.Rename(path, destFile); err != nil {
								if err := utils.CopyFile(path, destFile); err != nil {
									return err
								}
								os.Remove(path)
							}
						}
					}
				}
			}
		}
		return nil
	})
}
