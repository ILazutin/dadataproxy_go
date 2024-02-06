package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "github.com/ilazutin/dadataproxy_go/internal/api/rest_server"
	"github.com/ilazutin/dadataproxy_go/internal/config"
	"github.com/ilazutin/dadataproxy_go/internal/service"
	"github.com/ilazutin/dadataproxy_go/internal/service/dadata"
	storage "github.com/ilazutin/dadataproxy_go/internal/storage/redis"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Starting dadata proxy", slog.String("env", cfg.Env))

	redis := storage.New(context.Background(), cfg.Redis.Url, cfg.Redis.Password, log, cfg.Redis.Expire)

	dadata := dadata.New(cfg.DaData.Token, cfg.DaData.SecretKey)

	proxy := service.New(dadata, redis, log)

	server := server.New(cfg.HTTPServer.Address, cfg.HTTPServer.Timeout, cfg.HTTPServer.IdleTimeout, proxy, log)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(done)

	err := server.Run(context.Background())
	if err != nil {
		log.Error("Couldn't run server", slog.String("error", err.Error()))
	}

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(ctx.Err())
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
