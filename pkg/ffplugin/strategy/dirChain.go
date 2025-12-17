package strategy

import (
	"io/fs"
	"log/slog"
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

func (s *DirChainStrategy) FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error) {

	finalDest := filepath.Join(destDir, filepath.Base(srcDir))
	return finalDest, nil
}

func init() {
	RegisterStrategy("dirchain", func() Strategy {
		slog.Debug("Create a strategy", "name", "date")
		return &DirChainStrategy{}
	})
}
