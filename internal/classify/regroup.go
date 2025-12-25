package classify

import (
	"io/fs"
	"log/slog"
)

func (c *Classifier) regroupFile(srcDir, originalPath, finalFile string, info fs.FileInfo) error {

	if c.cfg.Regroup == nil {
		return nil
	}

	// Note: regroup path is computed from original source path to preserve structure
	regroupPath, err := destPath(srcDir, c.cfg.Regroup.Path, originalPath, info, c.cfg.Regroup.Strategy)

	if err != nil {
		slog.Error("Failed to compute final directory path for regrouping", "file", finalFile, "err", err)
		return err
	}
	// Implementation of regrouping logic goes here
	if c.dryRun {
		slog.Debug("Dry run: would regroup file", "originalPath", finalFile, "to", regroupPath, "mode", c.cfg.Regroup.Mode)
		return nil
	}

	return executeRegroup(finalFile, regroupPath, c.cfg.Regroup.Mode)
}
