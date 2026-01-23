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
	"errors"
	"fmt"
	"log/slog"
	"syscall"
)

func copyAndRemove(src Context, dstPath string) (Context, error) {

	dst, err := CopyFileAtomic(src, dstPath)
	if err != nil {
		return nil, err
	}
	err = Remove(src)
	if err != nil {
		slog.Warn("failed to remove file after copying it", "path", src.Path())
	}
	return dst, err
}

func Replace(src Context, dstPath string) (Context, error) {
	if src == nil {
		return nil, fmt.Errorf("failed to replace file: dst=%s err=%w", dstPath, ErrContextIsNil)
	}

	if err := replaceFile(src.Path(), dstPath); err != nil {
		if errors.Is(err, syscall.EXDEV) {
			slog.Warn("cannot move file, different filesystems trying copy and remove", "file", src.Path(), "destination", dstPath)
			return copyAndRemove(src, dstPath)
		}
		return nil, fmt.Errorf("failed to replace file: src=%s dst=%s err=%w", src.Path(), dstPath, err)
	}
	src.setPath(dstPath)
	return src, nil
}
