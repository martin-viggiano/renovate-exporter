package tailer

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

	// Unblock scanner on cancel
	go func() {
		<-tailCtx.Done()
		f.Close()
	}()

	timer := time.AfterFunc(idleTimeout, cancel)
	defer timer.Stop()

	line := make(chan result)
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
		case result, ok := <-line:
			if !ok {
				return nil // Channel closed
			}

			if result.error != nil {
				return result.error
			}

			if result.entry == nil {
				// EOF
				return nil
			}

			entry, err := logentry.Parse(result.entry)
			if err != nil {
				slog.Warn("skipping malformed line", slog.String("path", path), slog.String("err", err.Error()))
				continue
			}

			timer.Reset(idleTimeout)
			t.matcher.Process(entry)
		case <-tailCtx.Done():
			// Parent context cancelled
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// Tail context cancelled, return nil
			return nil
		}
	}
}
