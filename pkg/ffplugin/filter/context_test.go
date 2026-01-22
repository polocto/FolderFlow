package filter_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"

	filehandler "github.com/polocto/FolderFlow/internal/fileHandler"
	"github.com/polocto/FolderFlow/pkg/ffplugin/filter"
	"github.com/stretchr/testify/assert"
)

// --------------------------
// Mock fs.FileInfo
// --------------------------
type mockFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() fs.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return m.sys }

// --------------------------
// Mock filehandler.Context
// --------------------------
type mockFileHandlerContext struct {
	path      string
	info      fs.FileInfo
	kind      filehandler.FileKind
	isRegular bool
	hash      [32]byte
	err       error
}

func (m *mockFileHandlerContext) Path() string               { return m.path }
func (m *mockFileHandlerContext) Info() fs.FileInfo          { return m.info }
func (m *mockFileHandlerContext) GetHash() ([32]byte, error) { return m.hash, m.err }
func (m *mockFileHandlerContext) IsRegular() bool            { return m.isRegular }
func (m *mockFileHandlerContext) Kind() filehandler.FileKind { return m.kind }

// --------------------------
// Tests for filter.Context
// --------------------------
func TestNewContext(t *testing.T) {
	mockFH := &mockFileHandlerContext{
		path:      "/fake/file.txt",
		info:      mockFileInfo{name: "file.txt"},
		isRegular: true,
		kind:      filehandler.KindRegular,
	}

	ctx, err := filter.NewContextFilter(mockFH)
	assert.NoError(t, err)
	assert.Equal(t, "file.txt", ctx.BaseName())
	assert.False(t, ctx.IsDir())
	assert.Equal(t, int64(0), ctx.Size()) // mock size defaults to 0
}

func TestWithInput(t *testing.T) {
	content := []byte("Hello World")

	// Use a temporary real file since WithInput opens os.Open
	tmpFile := t.TempDir() + "/file.txt"
	err := os.WriteFile(tmpFile, content, 0644)
	assert.NoError(t, err)

	info, err := os.Stat(tmpFile)
	assert.NoError(t, err)

	mockFH := &mockFileHandlerContext{
		path:      tmpFile,
		info:      info,
		isRegular: true,
		kind:      filehandler.KindRegular,
	}

	ctx, err := filter.NewContextFilter(mockFH)
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = ctx.WithInput(func(r io.Reader) error {
		_, err := io.Copy(&buf, r)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, content, buf.Bytes())
}

func TestWithInput_OnDir(t *testing.T) {
	mockFH := &mockFileHandlerContext{
		path:      "/fake/dir",
		info:      mockFileInfo{name: "dir", isDir: true},
		isRegular: false,
		kind:      filehandler.KindDir,
	}

	ctx, err := filter.NewContextFilter(mockFH)
	assert.NoError(t, err)

	err = ctx.WithInput(func(r io.Reader) error { return nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot open directory")
}

func TestWithInputLimited(t *testing.T) {
	content := []byte("Hello World")
	tmpFile := t.TempDir() + "/file.txt"
	err := os.WriteFile(tmpFile, content, 0644)
	assert.NoError(t, err)

	info, err := os.Stat(tmpFile)
	assert.NoError(t, err)

	mockFH := &mockFileHandlerContext{
		path:      tmpFile,
		info:      info,
		isRegular: true,
		kind:      filehandler.KindRegular,
	}

	ctx, err := filter.NewContextFilter(mockFH)
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = ctx.WithInputLimited(5, func(r io.Reader) error {
		_, err := io.Copy(&buf, r)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello"), buf.Bytes())
}

func TestReadChunks(t *testing.T) {
	content := []byte("0123456789")
	tmpFile := t.TempDir() + "/file.txt"
	err := os.WriteFile(tmpFile, content, 0644)
	assert.NoError(t, err)

	info, err := os.Stat(tmpFile)
	assert.NoError(t, err)

	mockFH := &mockFileHandlerContext{
		path:      tmpFile,
		info:      info,
		isRegular: true,
		kind:      filehandler.KindRegular,
	}

	ctx, err := filter.NewContextFilter(mockFH)
	assert.NoError(t, err)

	var chunks [][]byte
	err = ctx.ReadChunks(4, func(b []byte) error {
		cpy := make([]byte, len(b))
		copy(cpy, b)
		chunks = append(chunks, cpy)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{
		[]byte("0123"),
		[]byte("4567"),
		[]byte("89"),
	}, chunks)
}

func TestReadChunks_ErrorPropagation(t *testing.T) {
	content := []byte("012345")
	tmpFile := t.TempDir() + "/file.txt"
	err := os.WriteFile(tmpFile, content, 0644)
	assert.NoError(t, err)

	info, err := os.Stat(tmpFile)
	assert.NoError(t, err)

	mockFH := &mockFileHandlerContext{
		path:      tmpFile,
		info:      info,
		isRegular: true,
		kind:      filehandler.KindRegular,
	}

	ctx, err := filter.NewContextFilter(mockFH)
	assert.NoError(t, err)

	expectedErr := errors.New("stop early")
	err = ctx.ReadChunks(3, func(b []byte) error {
		return expectedErr
	})
	assert.ErrorIs(t, err, expectedErr)
}
