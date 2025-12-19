package classify

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStats(t *testing.T) {
	stats := NewStats()
	assert.NotNil(t, stats)
	assert.Empty(t, stats.TotalFiles)
	assert.Empty(t, stats.MatchedFiles)
	assert.Empty(t, stats.MovedFiles)
	assert.Empty(t, stats.Errors)
}

func TestTotalFile(t *testing.T) {
	stats := NewStats()
	path := "/path/to/file1"
	size := int64(1024)

	stats.TotalFile(path, size)

	assert.Len(t, stats.TotalFiles, 1)
	assert.Contains(t, stats.TotalFiles, path)
	assert.Equal(t, size, stats.TotalFiles[path].Size)
	assert.NotZero(t, stats.TotalFiles[path].StartTime)
}

func TestMatchedFile(t *testing.T) {
	stats := NewStats()
	path := "/path/to/file1"
	size := int64(1024)

	stats.MatchedFile(path, size)

	assert.Len(t, stats.MatchedFiles, 1)
	assert.Contains(t, stats.MatchedFiles, path)
	assert.Equal(t, size, stats.MatchedFiles[path].Size)
	assert.NotZero(t, stats.MatchedFiles[path].StartTime)
}

func TestMovedFile(t *testing.T) {
	stats := NewStats()
	path := "/path/to/file1"
	size := int64(1024)
	startTime := time.Now().Add(-time.Second)

	stats.MovedFile(path, size, startTime)

	assert.Len(t, stats.MovedFiles, 1)
	assert.Contains(t, stats.MovedFiles, path)
	assert.Equal(t, size, stats.MovedFiles[path].Size)
	assert.NotZero(t, stats.MovedFiles[path].StartTime)
	assert.NotZero(t, stats.MovedFiles[path].EndTime)
	assert.NotZero(t, stats.MovedFiles[path].Duration)
}

func TestRecordError(t *testing.T) {
	stats := NewStats()
	path := "/path/to/file1"
	message := "Failed to move file"

	stats.RecordError(path, message)

	assert.Len(t, stats.Errors, 1)
	assert.Equal(t, path, stats.Errors[0].Path)
	assert.Equal(t, message, stats.Errors[0].Message)
	assert.NotZero(t, stats.Errors[0].Time)
}

func TestMerge(t *testing.T) {
	stats1 := NewStats()
	stats2 := NewStats()

	stats1.TotalFile("/path/to/file1", 1024)
	stats1.MatchedFile("/path/to/file1", 1024)
	stats1.MovedFile("/path/to/file1", 1024, time.Now().Add(-time.Second))
	stats1.RecordError("/path/to/file1", "Error 1")

	stats2.TotalFile("/path/to/file2", 2048)
	stats2.MatchedFile("/path/to/file2", 2048)
	stats2.MovedFile("/path/to/file2", 2048, time.Now().Add(-time.Second))
	stats2.RecordError("/path/to/file2", "Error 2")

	stats1.Merge(stats2)

	assert.Len(t, stats1.TotalFiles, 2)
	assert.Len(t, stats1.MatchedFiles, 2)
	assert.Len(t, stats1.MovedFiles, 2)
	assert.Len(t, stats1.Errors, 2)
}

func TestReset(t *testing.T) {
	stats := NewStats()

	stats.TotalFile("/path/to/file1", 1024)
	stats.MatchedFile("/path/to/file1", 1024)
	stats.MovedFile("/path/to/file1", 1024, time.Now().Add(-time.Second))
	stats.RecordError("/path/to/file1", "Error 1")

	assert.Len(t, stats.TotalFiles, 1)
	assert.Len(t, stats.MatchedFiles, 1)
	assert.Len(t, stats.MovedFiles, 1)
	assert.Len(t, stats.Errors, 1)

	stats.Reset()

	assert.Empty(t, stats.TotalFiles)
	assert.Empty(t, stats.MatchedFiles)
	assert.Empty(t, stats.MovedFiles)
	assert.Empty(t, stats.Errors)
}

func TestString(t *testing.T) {
	stats := NewStats()

	stats.TotalFile("/path/to/file1", 1024)
	stats.MatchedFile("/path/to/file1", 1024)
	stats.MovedFile("/path/to/file1", 1024, time.Now().Add(-time.Second))

	str := stats.String()
	assert.Contains(t, str, "Total: 1")
	assert.Contains(t, str, "Matched: 1")
	assert.Contains(t, str, "Moved: 1")
	assert.Contains(t, str, "Errors: 0")
	assert.Contains(t, str, "1.0 kB")
}

func TestToJSON(t *testing.T) {
	stats := NewStats()

	stats.TotalFile("/path/to/file1", 1024)
	stats.MatchedFile("/path/to/file1", 1024)
	stats.MovedFile("/path/to/file1", 1024, time.Now().Add(-time.Second))
	stats.RecordError("/path/to/file1", "Error 1")

	jsonData, err := stats.ToJSON(true)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "total_files")
	assert.Contains(t, result, "matched_files")
	assert.Contains(t, result, "moved_files")
	assert.Contains(t, result, "errors")
}

func TestWriteToFile(t *testing.T) {
	stats := NewStats()
	filename := "test_stats.json"

	stats.TotalFile("/path/to/file1", 1024)
	stats.MatchedFile("/path/to/file1", 1024)
	stats.MovedFile("/path/to/file1", 1024, time.Now().Add(-time.Second))
	stats.RecordError("/path/to/file1", "Error 1")

	err := stats.WriteToFile(filename, true)
	require.NoError(t, err)

	fileContent, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.NotEmpty(t, fileContent)

	var result map[string]interface{}
	err = json.Unmarshal(fileContent, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "total_files")
	assert.Contains(t, result, "matched_files")
	assert.Contains(t, result, "moved_files")
	assert.Contains(t, result, "errors")

	// Clean up
	os.Remove(filename)
}

func TestTotalSize(t *testing.T) {
	stats := NewStats()

	stats.MovedFile("/path/to/file1", 1024, time.Now().Add(-time.Second))
	stats.MovedFile("/path/to/file2", 2048, time.Now().Add(-time.Second))

	totalSize := stats.TotalSize()
	assert.Equal(t, int64(3072), totalSize)
}

func TestAvgDuration(t *testing.T) {
	stats := NewStats()
	now := time.Now()

	stats.MovedFile("/path/to/file1", 1024, now.Add(-2*time.Second))
	stats.MovedFile("/path/to/file2", 2048, now.Add(-4*time.Second))

	avgDuration := stats.AvgDuration()
	assert.Greater(t, avgDuration, time.Second)
	assert.Less(t, avgDuration, 4*time.Second)
}

func TestConcurrentAccess(t *testing.T) {
	stats := NewStats()
	numGoroutines := 100
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			path := fmt.Sprintf("/path/to/file%d", i)
			stats.TotalFile(path, int64(i))
			stats.MatchedFile(path, int64(i))
			stats.MovedFile(path, int64(i), time.Now().Add(-time.Second))
			done <- true
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	assert.Len(t, stats.TotalFiles, numGoroutines)
	assert.Len(t, stats.MatchedFiles, numGoroutines)
	assert.Len(t, stats.MovedFiles, numGoroutines)
}
