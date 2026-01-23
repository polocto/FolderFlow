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
)

func (c *Classifier) safeRun(name string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil && c.stats != nil {
			err = fmt.Errorf("%s panicked: %v", name, r)
			c.stats.Error(err)
		}
	}()
	err = fn()
	if err != nil && c.stats != nil {
		c.stats.Error(err)
	}
	return err
}
