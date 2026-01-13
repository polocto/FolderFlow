// Copyright (c) 2026 Paul Sade.
//
// This file is part of the FolderFlow project.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3,
// as published by the Free Software Foundation (see the LICENSE file).
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU General Public License for more details.

package classify

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/fsutil"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

func destPath(sourceDir, destDir, filePath string, info fs.FileInfo, strat strategy.Strategy) (string, error) {
	finalDir, err := strat.FinalDirPath(sourceDir, destDir, filePath, info)
	if err != nil {
		return "", fmt.Errorf("strategy failed to compute destination path : strategy=%s err=%w", strat.Selector(), err)
	}

	if !fsutil.IsSubDirectory(destDir, finalDir) {
		return "", fmt.Errorf("computed destination path is outside of destination directory : computedPath=%s destDir=%s", finalDir, destDir)
	}

	destFile := filepath.Join(finalDir, filepath.Base(filePath))

	return destFile, nil
}

func (c *Classifier) runStartegy(sourceDir, destDir, filePath string, info fs.FileInfo, strat strategy.Strategy) (finalDst string, err error) {
	err = c.safeRun("strategy", func() (err error) {
		finalDst, err = destPath(sourceDir, destDir, filePath, info, strat)
		return err
	})
	return finalDst, err
}
