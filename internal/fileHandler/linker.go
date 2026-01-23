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

package filehandler

import (
	"fmt"
	"os"
)

func Hardlink(src Context, dstPath string) (Context, error) {
	if src == nil {
		return nil, fmt.Errorf("failed to create a hard link: %w", ErrContextIsNil)
	}

	if err := os.Link(src.Path(), dstPath); err != nil {
		return nil, fmt.Errorf(
			"failed to create hard link from %q to %q: %w",
			src.Path(),
			dstPath,
			err,
		)
	}

	return NewContextFile(dstPath)
}

func Symlink(src Context, dstPath string) (Context, error) {
	if src == nil {
		return nil, fmt.Errorf("failed to create a symbolic link: %w", ErrContextIsNil)
	}

	if err := os.Symlink(src.Path(), dstPath); err != nil {
		return nil, fmt.Errorf(
			"failed to create symbolic link from %q to %q: %w",
			src.Path(),
			dstPath,
			err,
		)
	}

	return NewContextFile(dstPath)
}
