package strategy

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// plugin/date_strategy.go
type DirChainStrategy struct {
}

func (s *DirChainStrategy) Selector() string {
	return "dirchain"
}

func (s *DirChainStrategy) LoadConfig(config map[string]interface{}) error {
	return nil
}

func (s *DirChainStrategy) Apply(srcPath, destPath string, info fs.FileInfo, dryRun bool) error {

	finalDest := filepath.Join("hello", filepath.Base(srcPath))
	slog.Debug("Moving file", "source", srcPath, "destination", finalDest)
	if dryRun {
		return nil
	}
	return os.Rename(srcPath, finalDest)
}

func init() {
	RegisterStrategy("dirchain", func() Strategy {
		slog.Debug("Create a strategy", "name", "date")
		return &DirChainStrategy{}
	})
}
