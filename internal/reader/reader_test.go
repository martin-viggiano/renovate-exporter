package reader

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	tmpDir := t.TempDir()

	lines := [][]byte{}

	lineFn := func(ctx context.Context, data []byte) error {
		lines = append(lines, data)

		return nil
	}

	reader := NewReader(lineFn, Options{
		IdleTimeout: 3 * time.Second,
	})
	require.NotNil(t, reader)

	file, err := os.Create(filepath.Join(tmpDir, "file1.log"))
	require.NoError(t, err)

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		err := reader.Tail(t.Context(), filepath.Join(tmpDir, "file1.log"))
		errCh <- err
	}()

	time.Sleep(time.Second)

	// Write a complete line
	_, err = fmt.Fprintf(file, "line 1\n")
	require.NoError(t, err)

	select {
	case err := <-errCh:
		require.NoError(t, err, "Tail should not have returned")
	default:
	}

	assert.Eventually(t, func() bool {
		return slices.ContainsFunc(lines, func(line []byte) bool {
			return bytes.Equal(line, []byte("line 1"))
		})
	}, 1*time.Second, 100*time.Millisecond)

	// Write a partial line
	_, err = fmt.Fprintf(file, "line 2 not complete")
	require.NoError(t, err)

	select {
	case err := <-errCh:
		require.NoError(t, err, "Tail should not have returned")
	default:
	}

	assert.Never(t, func() bool {
		return slices.ContainsFunc(lines, func(line []byte) bool {
			return bytes.Equal(line, []byte("line 2 not complete"))
		})
	}, 1*time.Second, 100*time.Millisecond)

	_, err = fmt.Fprintf(file, " completed\n")
	require.NoError(t, err)

	select {
	case err := <-errCh:
		require.NoError(t, err, "Tail should not have returned")
	default:
	}

	assert.Eventually(t, func() bool {
		return slices.ContainsFunc(lines, func(line []byte) bool {
			return bytes.Equal(line, []byte("line 2 not complete completed"))
		})
	}, 1*time.Second, 100*time.Millisecond)

	// Close file
	require.NoError(t, file.Close())

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		assert.Fail(t, "IDLE timer should have returned")
	}
}
