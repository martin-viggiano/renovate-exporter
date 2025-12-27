package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/martin-viggiano/renovate-exporter/internal/matcher"
	"github.com/martin-viggiano/renovate-exporter/internal/registry"
	"github.com/martin-viggiano/renovate-exporter/internal/watcher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().StringVarP(&watchDir, "path", "p", "", "Directory to watch")
	rootCmd.MarkFlagRequired("path")
}

var (
	watchDir string

	rootCmd = &cobra.Command{
		Use:   "renovate-exporter",
		Short: "renovate-exporter extracts metrics from your Renovate logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			handlerOptions := &slog.HandlerOptions{
				Level: slog.LevelDebug, // TODO: set configurable log level
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

			slog.SetDefault(logger)

			reg := prometheus.NewRegistry()

			registry, err := registry.New(reg)
			if err != nil {
				return err
			}

			matcher := matcher.NewEngine(registry)

			watcher, err := watcher.NewWatcher(watchDir, func(path string) {
				t := watcher.NewTailer(matcher)
				t.Tail(ctx, path)
			})
			if err != nil {
				return err
			}

			slog.Info("watching directory", slog.String("path", watchDir))
			return watcher.Watch(ctx)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
