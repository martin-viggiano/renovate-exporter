package watcher

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/martin-viggiano/renovate-exporter/internal/logentry"
	"github.com/martin-viggiano/renovate-exporter/internal/matcher"
)

const (
	idleTimeout = 5 * time.Minute
)

type result struct {
	entry []byte
	error error
}

type Tailer struct {
	matcher *matcher.Engine
}

func NewTailer(matcher *matcher.Engine) *Tailer {
	return &Tailer{
		matcher: matcher,
	}
}

func (t *Tailer) Tail(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		slog.Error("failed to open file", slog.String("path", path), slog.String("error", err.Error()))
		return err
	}
	defer f.Close()

	tailCtx, cancel := context.WithCancel(ctx)

	line := make(chan result)
	timer := time.AfterFunc(idleTimeout, cancel)

	go func() {
		defer close(line)

		r := bufio.NewScanner(f)

		for r.Scan() {
			l := append([]byte(nil), r.Bytes()...)

			select {
			case line <- result{
				entry: l,
				error: nil,
			}:
			case <-tailCtx.Done():
				return
			}
		}

		select {
		case line <- result{
			entry: nil,
			error: r.Err(),
		}:
		case <-tailCtx.Done():
			return
		}
	}()

	for {
		select {
		case result := <-line:
			if result.error != nil {
				return result.error
			}

			if result.entry == nil {
				// EOF
				return nil
			}

			entry, err := logentry.Parse(result.entry)
			if err != nil {
				slog.Error("failed to parse line", slog.String("path", path), slog.String("err", err.Error()))
				return err
			}

			timer.Reset(idleTimeout) // Update IDLE deadline
			t.matcher.Process(entry)
		case <-tailCtx.Done():
			return tailCtx.Err()
		}
	}
}
