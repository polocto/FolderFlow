package core

import (
	"testing"
)

func TestStats(t *testing.T) {
	var s stats

	// Simuler le traitement de fichiers
	s.RecordFile(true, true, false)  // Fichier matched et moved
	s.RecordFile(false, false, true) // Erreur
	s.RecordFile(true, false, false) // Fichier matched mais pas moved

	if s.totalFiles != 3 {
		t.Errorf("Expected 3 total files, got %d", s.totalFiles)
	}
	if s.matchedFiles != 2 {
		t.Errorf("Expected 2 matched files, got %d", s.matchedFiles)
	}
	if s.movedFiles != 1 {
		t.Errorf("Expected 1 moved file, got %d", s.movedFiles)
	}
	if s.errors != 1 {
		t.Errorf("Expected 1 error, got %d", s.errors)
	}
}
