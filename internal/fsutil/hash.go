package fsutil

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
)

// Optional: use hash for verification or caching
func FilesEqualHash(path1, path2 string) (bool, error) {
	hash1, err := fileHash(path1)
	if err != nil {
		return false, err
	}
	hash2, err := fileHash(path2)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hash1, hash2), nil
}

func fileHash(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
