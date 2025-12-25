package fsutil

import (
	"crypto/sha256"
	"io"
	"log/slog"
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
	defer func() {
		_ = in.Close()
	}()

	tmp := dst + ".tmp"

	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after copy failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	if err := out.Sync(); err != nil {
		_ = out.Close()
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after sync failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	if err := out.Close(); err != nil {
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after close failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	if err := ReplaceFile(tmp, dst); err != nil {
		if rmErr := os.Remove(tmp); rmErr != nil {
			slog.Warn(
				"failed to delete temp file after replace failure",
				"path", tmp,
				"error", rmErr,
			)
		}
		return err
	}

	return nil
}
