package watcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatcher(t *testing.T) {
	tmpDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "subdir"), os.ModePerm), "must be able to create a subdirectory")

	files := map[string]struct{}{}

	newFileFn := func(ctx context.Context, path string) {
		files[path] = struct{}{}
	}

	watcher, err := NewWatcher(tmpDir, newFileFn)
	require.NoError(t, err)

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		err := watcher.Watch(t.Context())
		errCh <- err
	}()

	time.Sleep(1 * time.Second) // Wait until file watch starts.

	// Add a file in the root directory
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file1.log"), []byte("file1.log"), os.ModePerm))

	select {
	case err := <-errCh:
		require.NoError(t, err, "Watch should not have returned")
	default:
	}

	assert.Eventually(t, func() bool {
		_, ok := files[filepath.Join(tmpDir, "file1.log")]
		return ok
	}, 5*time.Second, 100*time.Millisecond)

	// Add a file in a pre-existing subdirectory
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "subdir", "file2.log"), []byte("file2.log"), os.ModePerm))

	select {
	case err := <-errCh:
		require.NoError(t, err, "Watch should not have returned")
	default:
	}

	assert.Eventually(t, func() bool {
		_, ok := files[filepath.Join(tmpDir, "subdir", "file2.log")]
		return ok
	}, 5*time.Second, 100*time.Millisecond)

	// Add a new subdirectory and a new file
	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "subdir2"), os.ModePerm))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "subdir2", "file3.log"), []byte("file3.log"), os.ModePerm))

	select {
	case err := <-errCh:
		require.NoError(t, err, "Watch should not have returned")
	default:
	}

	assert.Eventually(t, func() bool {
		_, ok := files[filepath.Join(tmpDir, "subdir2", "file3.log")]
		return ok
	}, 5*time.Second, 100*time.Millisecond)

	// Add a new file in a the new subdirectory
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "subdir2", "file4.log"), []byte("file4.log"), os.ModePerm))

	select {
	case err := <-errCh:
		require.NoError(t, err, "Watch should not have returned")
	default:
	}

	assert.Eventually(t, func() bool {
		_, ok := files[filepath.Join(tmpDir, "subdir2", "file4.log")]
		return ok
	}, 5*time.Second, 100*time.Millisecond)
}
