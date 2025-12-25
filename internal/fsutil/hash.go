package fsutil

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"

	"github.com/polocto/FolderFlow/internal/stats"
)

func FileHash(path string, s *stats.Stats) ([]byte, error) {
	if s != nil {
		defer s.Time(&s.Timing.Hash)()
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	} else if s != nil {
		s.HashComputed()
	}
	return h.Sum(nil), nil
}

func HashEqual(f1, f2 []byte, s *stats.Stats) bool {
	if !bytes.Equal(f1, f2) {
		return false
	}

	if s != nil {
		s.HashVerified()
	}
	return true
}
