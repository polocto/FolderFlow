// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


//go:build !windows

package fsutil

import "os"

func ReplaceFile(src, dst string) error {
	return os.Rename(src, dst)
}
