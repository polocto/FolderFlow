// Copyright (c) 2026 Paul Sade.
//
// This file is part of the FolderFlow project.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3,
// as published by the Free Software Foundation (see the LICENSE file).
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU General Public License for more details.

package stats

import (
	"context"
	"errors"
	"io/fs"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunTiming(t *testing.T) {
	var s Stats

	s.StartRun()
	time.Sleep(10 * time.Millisecond)
	s.EndRun()

	if s.Run.StartTime.IsZero() {
		t.Fatal("StartTime not set")
	}
	if s.Run.EndTime.IsZero() {
		t.Fatal("EndTime not set")
	}
	if s.Run.Duration <= 0 {
		t.Fatalf("expected positive duration, got %s", s.Run.Duration)
	}
}

func TestFileLifecycleCounters(t *testing.T) {
	var s Stats

	const size = int64(1234)

	s.FileSeen(size)
	s.FileMatched()
	s.DecisionSameFS()
	s.FileRenamed(size)

	if s.Run.FilesSeen != 1 {
		t.Errorf("FilesSeen = %d, want 1", s.Run.FilesSeen)
	}
	if s.Run.FilesMatched != 1 {
		t.Errorf("FilesMatched = %d, want 1", s.Run.FilesMatched)
	}
	if s.Run.FilesMoved != 1 {
		t.Errorf("FilesMoved = %d, want 1", s.Run.FilesMoved)
	}
	if s.Run.FilesRenamed != 1 {
		t.Errorf("FilesRenamed = %d, want 1", s.Run.FilesRenamed)
	}
	if s.Decisions.SameFS != 1 {
		t.Errorf("SameFS = %d, want 1", s.Decisions.SameFS)
	}
}

func TestCopyAcrossFilesystems(t *testing.T) {
	var s Stats
	const size = int64(4096)

	s.FileSeen(size)
	s.DecisionCrossFS()
	s.FileCopied(size)
	s.HashComputed()

	if s.Decisions.CrossFS != 1 {
		t.Errorf("CrossFS = %d, want 1", s.Decisions.CrossFS)
	}
	if s.Run.FilesCopied != 1 {
		t.Errorf("FilesCopied = %d, want 1", s.Run.FilesCopied)
	}
	if s.Run.BytesWritten != size {
		t.Errorf("BytesWritten = %d, want %d", s.Run.BytesWritten, size)
	}
	if s.Hash.Computed != 1 {
		t.Errorf("HashComputed = %d, want 1", s.Hash.Computed)
	}
}

func TestFileSkipped(t *testing.T) {
	var s Stats

	s.FileSeen(100)
	s.DecisionDuplicate()
	s.FileSkipped()

	if s.Run.FilesSkipped != 1 {
		t.Errorf("FilesSkipped = %d, want 1", s.Run.FilesSkipped)
	}
	if s.Decisions.Duplicate != 1 {
		t.Errorf("Duplicate = %d, want 1", s.Decisions.Duplicate)
	}
}

func TestErrorClassification(t *testing.T) {
	var s Stats

	s.Error(fs.ErrNotExist)
	s.Error(fs.ErrPermission)
	s.Error(context.Canceled)
	s.Error(errors.New("boom"))

	if s.Run.Errors != 4 {
		t.Fatalf("Errors = %d, want 4", s.Run.Errors)
	}
	if s.Run.FilesFailed != 4 {
		t.Fatalf("FilesFailed = %d, want 4", s.Run.FilesFailed)
	}

	want := map[string]int64{
		"not_found":  1,
		"permission": 1,
		"canceled":   1,
		"other":      1,
	}

	for k, v := range want {
		if got := s.Errors.ByKind[k]; got != v {
			t.Errorf("Errors.ByKind[%q] = %d, want %d", k, got, v)
		}
	}
}

func TestTimingAccumulation(t *testing.T) {
	var s Stats

	func() {
		defer s.Time(&s.Timing.Move)()
		time.Sleep(5 * time.Millisecond)
	}()

	if s.Timing.Move.Load() <= 0 {
		t.Fatalf("Timing.Move = %d, want > 0", s.Timing.Move.Load())
	}
}

func TestConcurrentUpdates(t *testing.T) {
	var s Stats
	const goroutines = 100
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.FileSeen(1)
				s.FileCopied(1)
				s.HashComputed()
			}
		}()
	}

	wg.Wait()

	wantFiles := int64(goroutines * iterations)

	if atomic.LoadInt64(&s.Run.FilesSeen) != wantFiles {
		t.Errorf("FilesSeen = %d, want %d", s.Run.FilesSeen, wantFiles)
	}
	if atomic.LoadInt64(&s.Run.FilesCopied) != wantFiles {
		t.Errorf("FilesCopied = %d, want %d", s.Run.FilesCopied, wantFiles)
	}
	if atomic.LoadInt64(&s.Hash.Computed) != wantFiles {
		t.Errorf("HashComputed = %d, want %d", s.Hash.Computed, wantFiles)
	}
}

func TestStringDoesNotPanic(t *testing.T) {
	var s Stats

	s.StartRun()
	s.FileSeen(10)
	s.FileCopied(10)
	s.HashComputed()
	s.EndRun()

	_ = s.String() // must not panic
}
