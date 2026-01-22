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

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/polocto/FolderFlow/internal/fsutil"
	"github.com/polocto/FolderFlow/pkg/ffplugin/strategy"
)

func destPath(file filehandler.Context, sourceDir, destDir string, strat strategy.Strategy) (string, error) {

	ctx, err := strategy.NewContextStrategy(file, sourceDir, destDir)

	finalDst, err := strat.FinalDirPath(ctx)
	if err != nil {
		return "", fmt.Errorf("strategy failed to compute destination path : strategy=%s err=%w", strat.Selector(), err)
	}

	if !fsutil.IsSubDirectory(destDir, finalDst) {
		return "", fmt.Errorf("computed destination path is outside of destination directory : computedPath=%s destDir=%s", finalDst, destDir)
	}

	return finalDst, nil
}

func (c *Classifier) runStartegy(file filehandler.Context, sourceDir, destDir string, strat strategy.Strategy) (finalDst string, err error) {
	err = c.safeRun("strategy", func() (err error) {
		finalDst, err = destPath(file, sourceDir, destDir, strat)
		return err
	})
	return finalDst, err
}
