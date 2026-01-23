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

package duplicates

import (
	"io/fs"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/config"
	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/polocto/FolderFlow/internal/stats"
)

type Duplicates struct {
	cfg    config.Config
	stats  *stats.Stats
	dryRun bool
}

func New(cfg config.Config, s *stats.Stats, dryRun bool) (*Duplicates, error) {
	if err := config.ValidateConfiguration(cfg); err != nil {
		return nil, err
	}

	return &Duplicates{
		cfg:    cfg,
		stats:  s,
		dryRun: dryRun,
	}, nil
}

func (d *Duplicates) ListDuplicates() (list [][]string, err error) {

	var listFile []filehandler.Context
	filepath.Walk("", func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := filehandler.NewContextFile(path)

		if err != nil {
			return nil
		}

		_, err = file.GetHash()

		if err != nil {
			return err
		}

		listFile = append(listFile, file)

		return err
	})

	regular, _, _, _ := filehandler.SplitFiles(listFile)

	// Get duplicates as [][]filehandler.Context
	duplicatesCtx, err := filehandler.ListDuplicates(regular)
	if err != nil {
		return nil, err
	}

	// Map [][]filehandler.Context to [][]string (file paths)
	duplicates := make([][]string, len(duplicatesCtx))
	for i, group := range duplicatesCtx {
		duplicates[i] = make([]string, len(group))
		for j, file := range group {
			duplicates[i][j] = file.Path()
		}
	}
	return duplicates, nil
}
