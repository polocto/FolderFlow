// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


/*
Package stats provides thread-safe runtime statistics collection
for FolderFlow executions.

Design principles:

- Stats are observational only; they must not affect control flow.
- Each file produces exactly one terminal event:
  - FileRenamed
  - FileCopied
  - FileSkipped
  - FileFailed
- FileSeen must be called exactly once per file.
- Decision methods explain why a terminal action occurred.
- All counters are safe for concurrent use.

Typical usage:

	stats.StartRun()
	defer stats.EndRun()

	stats.FileSeen(size)

	if sameFS {
		stats.DecisionSameFS()
		stats.FileRenamed(size)
	} else {
		stats.DecisionCrossFS()
		stats.FileCopied(size)
		stats.HashComputed()
	}

Timing:

	defer stats.Time(&stats.Timing.Move)()

Stats must never be mutated directly.
All updates must go through methods on *Stats.
*/

package stats

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type FileAction string

const (
	ActionSeen   FileAction = "seen"
	ActionRename FileAction = "rename"
	ActionCopy   FileAction = "copy"
	ActionSkip   FileAction = "skip"
	ActionFail   FileAction = "fail"
)

type RunStats struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration

	FilesSeen    int64
	FilesMatched int64
	FilesMoved   int64
	FilesCopied  int64
	FilesRenamed int64
	FilesSkipped int64
	FilesFailed  int64

	BytesRead    int64
	BytesWritten int64

	Errors int64
}

type OperationStats struct {
	Rename int64
	Copy   int64
	Link   int64
	Delete int64
}

type DecisionStats struct {
	SameFS       int64
	CrossFS      int64
	HasConflict  int64
	Duplicate    int64
	HashComputed int64
	HashVerified int64
}

type FileEvent struct {
	Path     string
	Action   string // move, copy, skip, fail
	Size     int64
	Duration time.Duration
	Error    string `json:",omitempty"`
}

type HashStats struct {
	Computed int64
	Verified int64
	Skipped  int64
	Algo     string
}

type TimingStats struct {
	Walk     atomic.Int64
	Classify atomic.Int64
	Move     atomic.Int64
	Hash     atomic.Int64
}

type ErrorStats struct {
	Total int64

	ByKind map[string]int64
	mu     sync.Mutex
}

type Stats struct {
	Run        RunStats
	Operations OperationStats
	Decisions  DecisionStats
	Hash       HashStats
	Timing     TimingStats
	Errors     ErrorStats
}

func (s *Stats) FileMoved(size int64) {
	atomic.AddInt64(&s.Run.FilesMoved, 1)
	atomic.AddInt64(&s.Run.BytesWritten, size)
}

func (s *Stats) FileMatched() {
	atomic.AddInt64(&s.Run.FilesMatched, 1)
}

func (s *Stats) StartRun() {
	s.Run.StartTime = time.Now()
}

func (s *Stats) EndRun() {
	s.Run.EndTime = time.Now()
	s.Run.Duration = s.Run.EndTime.Sub(s.Run.StartTime)
}

func (s *Stats) FileSeen(size int64) {
	atomic.AddInt64(&s.Run.FilesSeen, 1)
	atomic.AddInt64(&s.Run.BytesRead, size)
}

func (s *Stats) FileRenamed(size int64) {
	atomic.AddInt64(&s.Run.FilesMoved, 1)
	atomic.AddInt64(&s.Run.FilesRenamed, 1)
	atomic.AddInt64(&s.Operations.Rename, 1)
}

func (s *Stats) FileOverwrtitten(size int64) {
	atomic.AddInt64(&s.Run.FilesMoved, 1)
	atomic.AddInt64(&s.Run.BytesWritten, size)
	atomic.AddInt64(&s.Operations.Delete, 1)
}

func (s *Stats) FileCopied(size int64) {
	atomic.AddInt64(&s.Run.FilesMoved, 1)
	atomic.AddInt64(&s.Run.FilesCopied, 1)
	atomic.AddInt64(&s.Run.BytesWritten, size)
	atomic.AddInt64(&s.Operations.Copy, 1)
}

func (s *Stats) FileSkipped() {
	atomic.AddInt64(&s.Run.FilesSkipped, 1)
}

func (s *Stats) DecisionSameFS() {
	atomic.AddInt64(&s.Decisions.SameFS, 1)
}

func (s *Stats) DecisionCrossFS() {
	atomic.AddInt64(&s.Decisions.CrossFS, 1)
}

func (s *Stats) DecisionConflict() {
	atomic.AddInt64(&s.Decisions.HasConflict, 1)
}

func (s *Stats) DecisionDuplicate() {
	atomic.AddInt64(&s.Decisions.Duplicate, 1)
}

func (s *Stats) HashComputed() {
	atomic.AddInt64(&s.Hash.Computed, 1)
}

func (s *Stats) HashVerified() {
	atomic.AddInt64(&s.Hash.Verified, 1)
}

func (s *Stats) HashSkipped() {
	atomic.AddInt64(&s.Hash.Skipped, 1)
}

func (s *Stats) Time(section *atomic.Int64) func() {
	start := time.Now()
	return func() {
		section.Add(time.Since(start).Nanoseconds())
	}
}

func (s *Stats) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Run duration: %s\n", s.Run.Duration.Truncate(time.Microsecond))

	fmt.Fprintf(&b, "Files:\n")
	fmt.Fprintf(&b, "  Seen:    %d\n", s.Run.FilesSeen)
	fmt.Fprintf(&b, "  Moved:   %d\n", s.Run.FilesMoved)
	fmt.Fprintf(&b, "    Renamed: %d\n", s.Run.FilesRenamed)
	fmt.Fprintf(&b, "    Copied:  %d\n", s.Run.FilesCopied)
	fmt.Fprintf(&b, "  Skipped: %d\n", s.Run.FilesSkipped)
	fmt.Fprintf(&b, "  Failed:  %d\n", s.Run.FilesFailed)
	fmt.Fprintf(&b, "Hashing: %d computed (%s)\n", s.Hash.Computed, s.Hash.Algo)

	if s.Run.BytesWritten > 0 {
		fmt.Fprintf(&b, "Data moved: %s\n", humanBytes(s.Run.BytesWritten))
	}

	if s.Hash.Computed > 0 {
		fmt.Fprintf(
			&b,
			"Hashing: %d computed (%s)\n",
			s.Hash.Computed,
			s.Hash.Algo,
		)
	}

	if s.Run.Errors > 0 {
		fmt.Fprintf(&b, "Errors: %d\n", s.Run.Errors)
	}

	return b.String()
}

func (s *Stats) Error(err error) {
	if err == nil {
		return
	}

	atomic.AddInt64(&s.Run.Errors, 1)
	atomic.AddInt64(&s.Run.FilesFailed, 1)

	atomic.AddInt64(&s.Errors.Total, 1)

	kind := classifyError(err)

	s.Errors.mu.Lock()
	if s.Errors.ByKind == nil {
		s.Errors.ByKind = make(map[string]int64)
	}
	s.Errors.ByKind[kind]++
	s.Errors.mu.Unlock()
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for n >= div*unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(n)/float64(div),
		"KMGTPE"[exp],
	)
}

func classifyError(err error) string {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "not_found"
	case errors.Is(err, fs.ErrPermission):
		return "permission"
	case errors.Is(err, context.Canceled):
		return "canceled"
	default:
		return "other"
	}
}
