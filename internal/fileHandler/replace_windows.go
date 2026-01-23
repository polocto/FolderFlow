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

//go:build windows

package filehandler

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func replaceFile(src, dst string) error {
	ptrSrc := windows.StringToUTF16Ptr(src)
	ptrDst := windows.StringToUTF16Ptr(dst)

	err := windows.MoveFileEx(ptrSrc, ptrDst,
		windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH)
	if err != nil {
		return fmt.Errorf("failed to replace %s with %s: %w", dst, src, err)
	}
	return nil
}
