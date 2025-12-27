package reader

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"time"
)

const (
	idleTimeout = 5 * time.Minute
)

type result struct {
	entry []byte
	error error
}

type Reader struct {
	processLineFunc func(ctx context.Context, data []byte) error
}

func NewReader(processLineFunc func(ctx context.Context, data []byte) error) *Reader {
	return &Reader{
		processLineFunc: processLineFunc,
	}
}

func (t *Reader) Tail(ctx context.Context, path string) error {
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

			timer.Reset(idleTimeout)

			err := t.processLineFunc(ctx, result.entry)
			if err != nil {
				slog.Warn("error wile processing line", slog.String("path", path), slog.String("err", err.Error()))
				continue
			}
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
