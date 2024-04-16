package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/JamshedJ/position-tracking-service/internal/app"
	"github.com/JamshedJ/position-tracking-service/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("running application")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath)
	go application.GRPCSrv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	
	<-stop
	application.GRPCSrv.Stop()
	log.Info("application stopped")
	
	// TODO: run GRPC server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
