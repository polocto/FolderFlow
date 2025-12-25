package fsutil

import (
	"bytes"
	"io"
	"os"
)

// filesEqual compares two files for exact equality.
func FilesEqual(path1, path2 string) (bool, error) {
	// Get file info
	fi1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}
	fi2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}

	// Quick check: file size
	if fi1.Size() != fi2.Size() {
		return false, nil
	}

	// Open files
	f1, err := os.Open(path1)
	if err != nil {
		return false, err
	}
	defer func() {
		if cerr := f1.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	f2, err := os.Open(path2)
	if err != nil {
		return false, err
	}
	defer func() {
		if cerr := f2.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Compare chunks
	const chunkSize = 64 * 1024 // 64 KB
	buf1 := make([]byte, chunkSize)
	buf2 := make([]byte, chunkSize)

	for {
		n1, err1 := f1.Read(buf1)
		n2, err2 := f2.Read(buf2)

		if n1 != n2 || !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}

		if err1 == io.EOF && err2 == io.EOF {
			break
		}
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}
	}

	return true, nil
}
