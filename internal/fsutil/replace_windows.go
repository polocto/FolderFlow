//go:build windows

package fsutil

import (
	"golang.org/x/sys/windows"
)

func ReplaceFile(src, dst string) error {
	return windows.MoveFileEx(
		windows.StringToUTF16Ptr(src),
		windows.StringToUTF16Ptr(dst),
		windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH,
	)
}
