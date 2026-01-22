package classify

import (
	"errors"
	"os"
	"testing"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/stretchr/testify/require"
)

// --------------------
// mockFilter for testing
// --------------------
type mockFilter struct {
	match  bool
	err    error
	called *int
}

func (m *mockFilter) Match(ctx filter.Context) (bool, error) {
	if m.called != nil {
		*m.called++
	}
	return m.match, m.err
}

func (m *mockFilter) Selector() string                        { return "mock" }
func (m *mockFilter) LoadConfig(map[string]interface{}) error { return nil }

// --------------------
// Test helpers
// --------------------
func createContextFile(t *testing.T, content []byte) filehandler.Context {
	tmp := t.TempDir()
	file := tmp + "/file.txt"
	err := os.WriteFile(file, content, 0644)
	require.NoError(t, err)

	ctx, err := filehandler.NewContextFile(file)
	require.NoError(t, err)
	return ctx
}

// --------------------
// Tests
// --------------------
func TestMatchFile_NoFilters(t *testing.T) {
	ctx := createContextFile(t, []byte("Hello"))

	ok, err := matchFile(ctx, nil)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestMatchFile_FilterRejects(t *testing.T) {
	ctx := createContextFile(t, []byte("Hello"))

	mf := &mockFilter{match: false}
	ok, err := matchFile(ctx, []filter.Filter{mf})
	require.NoError(t, err)
	require.False(t, ok)
}

func TestMatchFile_SingleFilterMatch(t *testing.T) {
	ctx := createContextFile(t, []byte("Hello"))

	f := &mockFilter{match: true}
	ok, err := matchFile(ctx, []filter.Filter{f})
	require.NoError(t, err)
	require.True(t, ok)
}

func TestMatchFile_MultipleFilters_AllMatch(t *testing.T) {
	ctx := createContextFile(t, []byte("Hello"))

	filters := []filter.Filter{
		&mockFilter{match: true},
		&mockFilter{match: true},
	}
	ok, err := matchFile(ctx, filters)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestMatchFile_StopsOnFirstNoMatch(t *testing.T) {
	ctx := createContextFile(t, []byte("Hello"))

	var called1, called2 int
	filters := []filter.Filter{
		&mockFilter{match: false, called: &called1},
		&mockFilter{match: true, called: &called2},
	}

	ok, err := matchFile(ctx, filters)
	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, 1, called1)
	require.Equal(t, 0, called2)
}

func TestMatchFile_FilterError(t *testing.T) {
	ctx := createContextFile(t, []byte("Hello"))

	expectedErr := errors.New("filter error")
	f := &mockFilter{err: expectedErr}

	ok, err := matchFile(ctx, []filter.Filter{f})
	require.ErrorIs(t, err, expectedErr)
	require.False(t, ok)
}

func TestMatchFile_StopsOnError(t *testing.T) {
	ctx := createContextFile(t, []byte("Hello"))

	expectedErr := errors.New("boom")
	var called int
	filters := []filter.Filter{
		&mockFilter{err: expectedErr},
		&mockFilter{match: true, called: &called},
	}

	ok, err := matchFile(ctx, filters)
	require.ErrorIs(t, err, expectedErr)
	require.False(t, ok)
	require.Equal(t, 0, called)
}
