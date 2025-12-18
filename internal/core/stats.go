package core

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

// FileStat represents metadata for a file.
type FileStat struct {
	Path      string        `json:"path"`
	Size      int64         `json:"size"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
}

// ErrorStat represents an error encountered during file processing.
type ErrorStat struct {
	Path    string    `json:"path"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// Stats tracks file processing statistics.
type Stats struct {
	sync.Mutex
	TotalFiles   map[string]FileStat
	MatchedFiles map[string]FileStat
	MovedFiles   map[string]FileStat
	Errors       []ErrorStat
}

// NewStats creates a new Stats instance.
func NewStats() *Stats {
	return &Stats{
		TotalFiles:   make(map[string]FileStat),
		MatchedFiles: make(map[string]FileStat),
		MovedFiles:   make(map[string]FileStat),
		Errors:       []ErrorStat{},
	}
}

// Merge combines the statistics from another Stats object into this one.
func (s *Stats) Merge(other *Stats) {
	s.Lock()
	defer s.Unlock()

	for k, v := range other.TotalFiles {
		s.TotalFiles[k] = v
	}
	for k, v := range other.MatchedFiles {
		s.MatchedFiles[k] = v
	}
	for k, v := range other.MovedFiles {
		s.MovedFiles[k] = v
	}
	s.Errors = append(s.Errors, other.Errors...)
}

// Reset clears all statistics.
func (s *Stats) Reset() {
	s.Lock()
	defer s.Unlock()

	s.TotalFiles = make(map[string]FileStat)
	s.MatchedFiles = make(map[string]FileStat)
	s.MovedFiles = make(map[string]FileStat)
	s.Errors = []ErrorStat{}
}

// String returns a human-readable summary of the statistics.
func (s *Stats) String() string {
	s.Lock()
	defer s.Unlock()

	totalSize := int64(0)
	for _, stat := range s.MovedFiles {
		totalSize += stat.Size
	}

	avgDuration := time.Duration(0)
	if len(s.MovedFiles) > 0 {
		for _, stat := range s.MovedFiles {
			avgDuration += stat.Duration
		}
		avgDuration /= time.Duration(len(s.MovedFiles))
	}

	return fmt.Sprintf(
		"Total: %d, Matched: %d, Moved: %d, Errors: %d, Total size moved: %s, Avg move time: %s",
		len(s.TotalFiles),
		len(s.MatchedFiles),
		len(s.MovedFiles),
		len(s.Errors),
		humanize.Bytes(uint64(totalSize)),
		avgDuration,
	)
}

// TotalFile records a file as processed.
func (s *Stats) TotalFile(path string, size int64) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.TotalFiles[path]; !exists {
		s.TotalFiles[path] = FileStat{
			Path:      path,
			Size:      size,
			StartTime: time.Now(),
		}
	}
}

// MatchedFile records a file as matched.
func (s *Stats) MatchedFile(path string, size int64) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.MatchedFiles[path]; !exists {
		s.MatchedFiles[path] = FileStat{
			Path:      path,
			Size:      size,
			StartTime: time.Now(),
		}
	}
}

// MovedFile records a file as moved.
func (s *Stats) MovedFile(path string, size int64, startTime time.Time) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.MovedFiles[path]; !exists {
		s.MovedFiles[path] = FileStat{
			Path:      path,
			Size:      size,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime),
		}
	}
}

// RecordError records an error encountered during file processing.
func (s *Stats) RecordError(path, message string) {
	s.Lock()
	defer s.Unlock()
	s.Errors = append(s.Errors, ErrorStat{
		Path:    path,
		Message: message,
		Time:    time.Now(),
	})
}

// ToJSON exports the statistics as JSON.
func (s *Stats) ToJSON(pretty bool) ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	type statsJSON struct {
		TotalFiles   map[string]FileStat `json:"total_files"`
		MatchedFiles map[string]FileStat `json:"matched_files"`
		MovedFiles   map[string]FileStat `json:"moved_files"`
		Errors       []ErrorStat         `json:"errors"`
	}

	data := statsJSON{
		TotalFiles:   s.TotalFiles,
		MatchedFiles: s.MatchedFiles,
		MovedFiles:   s.MovedFiles,
		Errors:       s.Errors,
	}

	if pretty {
		return json.MarshalIndent(data, "", "  ")
	}
	return json.Marshal(data)
}

// WriteToFile writes the statistics to a file in JSON format.
func (s *Stats) WriteToFile(filename string, pretty bool) error {
	jsonData, err := s.ToJSON(pretty)
	if err != nil {
		return fmt.Errorf("failed to marshal stats to JSON: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write stats to file: %w", err)
	}

	return nil
}

// TotalSize returns the total size of all moved files.
func (s *Stats) TotalSize() int64 {
	s.Lock()
	defer s.Unlock()

	var total int64
	for _, stat := range s.MovedFiles {
		total += stat.Size
	}
	return total
}

// AvgDuration returns the average duration of all moved files.
func (s *Stats) AvgDuration() time.Duration {
	s.Lock()
	defer s.Unlock()

	if len(s.MovedFiles) == 0 {
		return 0
	}

	var total time.Duration
	for _, stat := range s.MovedFiles {
		total += stat.Duration
	}

	return total / time.Duration(len(s.MovedFiles))
}
