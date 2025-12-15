package strategy

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// plugin/date_strategy.go
type DateStrategy struct {
	Format string `yaml:"format"`
}

func (s *DateStrategy) Selector() string {
	return "date"
}

func (s *DateStrategy) LoadConfig(config map[string]interface{}) error {
	if format, ok := config["format"].(string); ok {
		s.Format = format
	} else {
		s.Format = "2006/01" // default
	}
	return nil
}

func (s *DateStrategy) Apply(srcPath, destPath string, info fs.FileInfo, dryRun bool) error {
	yearMonth := info.ModTime().Format(s.Format)
	finalDest := filepath.Join(destPath, yearMonth, filepath.Base(srcPath))
	return os.Rename(srcPath, finalDest)
}

func init() {
	RegisterStrategy("date", func() Strategy {
		slog.Debug("Create a strategy", "name", "date")
		return &DateStrategy{}
	})
}
