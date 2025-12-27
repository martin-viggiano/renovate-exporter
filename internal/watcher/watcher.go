package watcher

import (
	"context"
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/martin-viggiano/renovate-exporter/internal/matcher"
)

type Watcher struct {
	matcher *matcher.Engine
}

func NewWatcher(matcher *matcher.Engine) *Watcher {
	return &Watcher{matcher: matcher}
}

func (w *Watcher) Watch(ctx context.Context, path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("failed to start watcher", slog.String("path", path), slog.String("err", err.Error()))
		return err
	}

	if err := watcher.AddWith(path); err != nil { // TODO: add any subdirectories recurively.
		slog.Error("failed to add base path to watcher", slog.String("path", path), slog.String("err", err.Error()))
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if !event.Has(fsnotify.Create) {
				continue
			}

			slog.Debug("file created", slog.String("path", event.Name))

			info, err := os.Stat(event.Name)
			if err != nil {
				continue
			}

			if info.IsDir() {
				if err := watcher.Add(event.Name); err != nil {
					return err
				}
				continue
			}

			tailer := NewTailer(w.matcher)

			go func() {
				if err := tailer.Tail(ctx, event.Name); err != nil {
					slog.Error("error while tailing file", slog.String("path", event.Name), slog.String("err", err.Error()))
				}
			}()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
