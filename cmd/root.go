package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/martin-viggiano/renovate-exporter/internal/analyzer"
	"github.com/martin-viggiano/renovate-exporter/internal/analyzer/matchers"
	"github.com/martin-viggiano/renovate-exporter/internal/fswatch"
	"github.com/martin-viggiano/renovate-exporter/internal/reader"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("RENOVATE_EXPORTER")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	rootCmd.Flags().StringVarP(&watchDir, "path", "p", "", "Directory to watch")
	rootCmd.MarkFlagRequired("path")

	rootCmd.Flags().StringVar(&metricsAddress, "address", ":9090", "Metrics server address")

	viper.BindPFlags(rootCmd.Flags())
}

var (
	watchDir       string
	metricsAddress string

	rootCmd = &cobra.Command{
		Use:   "renovate-exporter",
		Short: "renovate-exporter extracts metrics from your Renovate logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			handlerOptions := &slog.HandlerOptions{
				Level: slog.LevelDebug, // TODO: set configurable log level
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

			slog.SetDefault(logger)

			reg := prometheus.NewRegistry()

			mux := http.NewServeMux()
			mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

			server := &http.Server{
				Addr:    metricsAddress,
				Handler: mux,
			}

			httpErrCh := make(chan error)
			go func() {
				slog.Info("starting metrics server", slog.String("address", metricsAddress))

				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					httpErrCh <- fmt.Errorf("metrics server: %w", err)
				}
			}()

			matchers := matchers.DefaultMatchers()

			matcher, err := analyzer.NewEngine(reg, matchers)
			if err != nil {
				return err
			}

			watcher, err := fswatch.New(watchDir, func(ctx context.Context, path string) {
				lineFn := func(ctx context.Context, data []byte) error {
					return matcher.Process(data)
				}

				t := reader.NewReader(lineFn, reader.Options{
					IdleTimeout: time.Minute,
				})

				t.Tail(ctx, path)
			})
			if err != nil {
				return err
			}

			watchErr := make(chan error)
			go func() {
				slog.Info("watching directory", slog.String("path", watchDir))
				if err := watcher.Watch(ctx); err != nil && ctx.Err() == nil {
					slog.Error("watcher failed", slog.String("error", err.Error()))
					watchErr <- fmt.Errorf("watcher: %w", err)
				}
			}()

			select {
			case err := <-httpErrCh:
				slog.Error("metrics server failed", slog.String("error", err.Error()))
				stop()
				return err
			case err := <-watchErr:
				slog.Error("watcher failed", slog.String("error", err.Error()))
				stop()
				return err

			case <-ctx.Done():
				slog.Info("shutdown started")

				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := server.Shutdown(shutdownCtx); err != nil {
					slog.Error("failed to shut down metrics server", slog.String("err", err.Error()))
				}

				slog.Info("shutdown complete")
				return nil
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
