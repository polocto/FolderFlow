package core

import (
	"fmt"
	"sync"
)

type stats struct {
	sync.Mutex
	totalFiles   int
	matchedFiles int
	movedFiles   int
	errors       int
}

func (s *stats) String() string {
	return fmt.Sprintf("Total files: %d, Matched files: %d, Moved files: %d, Errors: %d",
		s.totalFiles, s.matchedFiles, s.movedFiles, s.errors)
}

func (s *stats) RecordFile(matched bool, moved bool, errOccurred bool) {
	s.Lock()
	defer s.Unlock()
	s.totalFiles++
	if matched {
		s.matchedFiles++
	}
	if moved {
		s.movedFiles++
	}
	if errOccurred {
		s.errors++
	}
}
