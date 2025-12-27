package reader

import (
	"context"
	"log/slog"
	"time"

	"github.com/hpcloud/tail"
)

type Reader struct {
	processLineFunc func(ctx context.Context, data []byte) error
	opts            Options
}

type Options struct {
	IdleTimeout time.Duration
}

func NewReader(processLineFunc func(ctx context.Context, data []byte) error, opts Options) *Reader {
	return &Reader{
		processLineFunc: processLineFunc,
		opts:            opts,
	}
}

func (t *Reader) Tail(ctx context.Context, path string) error {
	tf, err := tail.TailFile(path, tail.Config{
		MustExist: true,
		Follow:    true,
		ReOpen:    false,
	})
	if err != nil {
		slog.Error("failed to open file", slog.String("path", path), slog.String("error", err.Error()))
		return err
	}
	defer tf.Cleanup()

	idleTimer := time.NewTimer(t.opts.IdleTimeout)
	defer idleTimer.Stop()

	for {
		select {
		case line := <-tf.Lines:
			if line.Err != nil {
				return line.Err
			}

			idleTimer.Reset(t.opts.IdleTimeout)

			if err := t.processLineFunc(ctx, []byte(line.Text)); err != nil {
				slog.Warn("error wile processing line", slog.String("path", path), slog.String("err", err.Error()))
			}

		case <-tf.Dead():
			return nil

		case <-idleTimer.C:
			slog.Debug("idle timeout reached", slog.String("path", path))
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
