package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gilab.com/estate-agency-api/internal/app"
	"gilab.com/estate-agency-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	logger := get_logger(cfg.Env)

	app := app.New(cfg, logger)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Start server on " + cfg.HTTPServerConfig.Address)

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	logger.Info("server started")

	<-done
	logger.Info("stopping server")

	if err := app.Shutdown(); err != nil {
		logger.Error("failed to stop server " + err.Error())

		return
	}

	logger.Info("server stopped")
}

func get_logger(env string) *slog.Logger {

	var level slog.Level

	switch env {
	case "debug":
		level = slog.LevelDebug
	case "prod":
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)

	return logger
}
