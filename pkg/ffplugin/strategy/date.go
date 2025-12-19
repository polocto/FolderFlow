package strategy

import (
	"io/fs"
	"log/slog"
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

func (s *DateStrategy) FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error) {
	yearMonth := info.ModTime().Format(s.Format)
	finalDest := filepath.Join(destDir, yearMonth, filepath.Base(srcDir))
	return finalDest, nil
}

func init() {
	RegisterStrategy("date", func() Strategy {
		slog.Debug("Create a strategy", "name", "date")
		return &DateStrategy{}
	})
}
