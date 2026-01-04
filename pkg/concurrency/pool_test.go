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

package concurrency

import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewWorkerPool_AutoWorkers(t *testing.T) {
	wp := NewWorkerPool(-1) // Doit calculer automatiquement
	expected := runtime.NumCPU() * 4
	if expected > 32 {
		expected = 32
	}
	if cap(wp.sem) != expected {
		t.Errorf("Expected %d workers, got %d", expected, cap(wp.sem))
	}
}

func TestNewWorkerPool_CustomWorkers(t *testing.T) {
	wp := NewWorkerPool(10) // Valeur personnalisée
	if cap(wp.sem) != 10 {
		t.Errorf("Expected 10 workers, got %d", cap(wp.sem))
	}
}

func TestWorkerPool(t *testing.T) {
	t.Run("should limit the number of workers", func(t *testing.T) {
		wp := NewWorkerPool(2) // 2 workers max
		var wg sync.WaitGroup
		var count int
		var mu sync.Mutex

		for i := 0; i < 10; i++ {
			wp.Add()
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				count++
				mu.Unlock()
				time.Sleep(100 * time.Millisecond) // Simuler un travail
				wp.Done()
			}()
		}

		wg.Wait()
		if count != 10 {
			t.Errorf("Expected 10 tasks to complete, got %d", count)
		}
	})

	t.Run("should report errors", func(t *testing.T) {
		wp := NewWorkerPool(1)
		wp.Add()
		go func() {
			defer wp.Done()
			wp.ReportError(errors.New("test error"))
		}()

		err := wp.Wait()
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "test error") {
			t.Errorf("Expected error to contain 'test error', got %v", err)
		}
	})
}
