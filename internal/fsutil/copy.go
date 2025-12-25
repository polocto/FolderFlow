package fsutil

import (
	"crypto/sha256"
	"io"
	"os"

	"github.com/polocto/FolderFlow/internal/stats"
)

func CopyFile(src, dst string, s *stats.Stats) ([]byte, error) {
	if s != nil {
		defer s.Time(&s.Timing.Hash)()
	}
	in, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	hasher := sha256.New()

	// Write to file AND hash simultaneously
	writer := io.MultiWriter(out, hasher)

	if size, err := io.Copy(writer, in); err != nil {
		return nil, err
	} else if s != nil {
		s.FileCopied(size)
		s.HashComputed()
	}

	if err := out.Sync(); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func CopyFileAtomic(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp := dst + ".tmp"

	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}

	if err := out.Sync(); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}

	if err := out.Close(); err != nil {
		os.Remove(tmp)
		return err
	}

	return ReplaceFile(src, dst)
}
