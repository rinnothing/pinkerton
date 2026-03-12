package controller

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/rinnothing/pinkerton/config"
	"github.com/rinnothing/pinkerton/internal/controller/service"
	"github.com/rinnothing/pinkerton/internal/interactor/pinger"
	"github.com/rinnothing/pinkerton/internal/repository/storage"
	"github.com/rinnothing/pinkerton/internal/usecase/healthcheck"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pngr := pinger.New(ctx, cfg.HealthCheck.Timeout)
	strg := storage.New()

	uc := healthcheck.New(ctx, cfg.InitModels, cfg.Threads, pngr, strg)

	apiMux := service.New(uc)

	addr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: apiMux,
	}

	go func() {
		slog.Info("starting to listen for incoming requests", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed listening ")
		}
	}()

	<-ctx.Done()
	cancel()
	slog.Info("shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("got error when shutting down server gracefully", "error", err)

		err = server.Close()
		if err != nil {
			slog.Error("got error when shutting down server forcefully", "error", err)
		}
	}

	slog.Info("server shutdown gracefully")
}
