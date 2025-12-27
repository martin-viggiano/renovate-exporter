package fswatch

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher   *fsnotify.Watcher
	path      string
	files     map[string]struct{}
	newFileFn func(ctx context.Context, path string)
}

func New(path string, onNewFile func(ctx context.Context, path string)) (*Watcher, error) {
	w := &Watcher{
		files:     make(map[string]struct{}),
		path:      path,
		newFileFn: onNewFile,
	}

	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Close closes the underlying [fsnotify.Watcher].
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

// addRecursive walks the given path and starts monitoring any subdirectories.
func (w *Watcher) addRecursive(path string) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return w.watcher.Add(path)
		}

		return nil
	})
}

// newFilesFromDir walks the given path and starts monitoring any subdirectories while
// also calling the new file function on any file.
func (w *Watcher) newFilesFromDir(ctx context.Context, path string) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return w.watcher.Add(path)
		}

		go w.newFileFn(ctx, path)

		return nil
	})
}

// Watch starts watching the given path and processing events.
//
// Returns on critical errors or if the context is cancelled.
func (w *Watcher) Watch(ctx context.Context) error {
	if err := w.watcher.Add(w.path); err != nil {
		slog.Error("failed to add base path to watcher", slog.String("path", w.path), slog.String("err", err.Error()))
		return err
	}

	// [fsnotify.Watcher.Add] is not recursive, so we add any subdirectory to the watcher.
	if err := w.addRecursive(w.path); err != nil {
		slog.Error("failed to add subdirectories to watcher", slog.String("path", w.path), slog.String("err", err.Error()))
		return err
	}

	for {
		select {
		case event := <-w.watcher.Events:
			if !event.Has(fsnotify.Create) {
				continue
			}

			slog.Debug("event received", slog.String("path", event.Name))

			info, err := os.Stat(event.Name)
			if err != nil {
				continue
			}

			if info.IsDir() {
				if err := w.newFilesFromDir(ctx, event.Name); err != nil {
					slog.Warn("failed to watch new subdirectory", slog.String("path", event.Name), slog.String("err", err.Error()))
				}
				continue
			}

			go w.newFileFn(ctx, event.Name)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
