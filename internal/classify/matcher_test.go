package classify

import (
	"errors"
	"testing"

	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

// --------------------
// matchFile tests
// --------------------
//

func TestMatchFile_NoFilters(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}

	ok, err := matchFile("file.txt", info, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected file to match when no filters are provided")
	}
}

func TestMatchFile_FilterRejects(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}

	filters := []filter.Filter{
		&mockFilter{match: false},
	}

	ok, err := matchFile("file.txt", info, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected file to be rejected")
	}
}

func TestMatchFile_SingleFilterMatch(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}
	f := &mockFilter{match: true, err: nil}

	ok, err := matchFile("file.txt", info, []filter.Filter{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected file to match")
	}
}

func TestMatchFile_MultipleFilters_AllMatch(t *testing.T) {
	info := mockFileInfo{name: "file.txt"}
	filters := []filter.Filter{
		&mockFilter{match: true},
		&mockFilter{match: true},
	}

	ok, err := matchFile("file.txt", info, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected all filters to match")
	}
}

func TestMatchFile_StopsOnFirstNoMatch(t *testing.T) {
	called := 0

	info := mockFileInfo{name: "file.txt"}
	filters := []filter.Filter{
		&mockFilter{match: false},
		&mockFilter{match: true, called: &called},
	}

	ok, err := matchFile("file.txt", info, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected file to not match")
	}
	if called != 0 {
		t.Fatal("expected second filter not to be called")
	}
}

func TestMatchFile_FilterError(t *testing.T) {
	expectedErr := errors.New("filter error")
	info := mockFileInfo{name: "file.txt"}

	f := &mockFilter{name: "err", err: expectedErr}

	ok, err := matchFile("file.txt", info, []filter.Filter{f})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected match to be false on error")
	}
}

func TestMatchFile_StopsOnError(t *testing.T) {
	called := 0
	expectedErr := errors.New("boom")
	info := mockFileInfo{name: "file.txt"}

	filters := []filter.Filter{
		&mockFilter{name: "f1", err: expectedErr},
		&mockFilter{name: "f2", match: true, called: &called},
	}

	ok, err := matchFile("file.txt", info, filters)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected match to be false")
	}
	if called != 0 {
		t.Fatal("expected second filter not to be called")
	}
}
