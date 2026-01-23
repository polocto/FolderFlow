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

func Remove(file Context) error {
	if file == nil {
		return fmt.Errorf("failed to remove file: %w", ErrContextIsNil)
	}

	if err := os.Remove(file.Path()); err != nil {
		return fmt.Errorf("failed to delete file: path=%q err=%w", file.Path(), err)
	}
	file.delete()
	return nil
}
