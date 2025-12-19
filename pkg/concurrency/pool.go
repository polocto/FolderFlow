package concurrency

import (
	"fmt"
	"runtime"
	"sync"
)

type WorkerPool struct {
	sem    chan struct{}
	wg     sync.WaitGroup
	errMu  sync.Mutex
	errors []error
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	if maxWorkers <= 0 {
		// Calculer automatiquement en fonction du CPU et du type de tâche
		ioScalingFactor := 4
		maxWorkers = runtime.NumCPU() * ioScalingFactor
		if maxWorkers > 32 { // Limite supérieure
			maxWorkers = 32
		}
	}
	return &WorkerPool{
		sem: make(chan struct{}, maxWorkers),
	}
}

func (wp *WorkerPool) Add() {
	wp.wg.Add(1)
	wp.sem <- struct{}{} // Bloque si le nombre max de workers est atteint
}

func (wp *WorkerPool) Done() {
	<-wp.sem
	wp.wg.Done()
}

func (wp *WorkerPool) ReportError(err error) {
	wp.errMu.Lock()
	wp.errors = append(wp.errors, err)
	wp.errMu.Unlock()
}

func (wp *WorkerPool) Wait() error {
	wp.wg.Wait()
	if len(wp.errors) > 0 {
		return fmt.Errorf("%d errors occurred: %v", len(wp.errors), wp.errors)
	}
	return nil
}
