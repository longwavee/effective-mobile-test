package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/longwavee/effective-mobile-test/internal/config"
)

type (
	App struct {
		cfg *config.Config
		log *slog.Logger

		di *DIContainer
	}
)

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}
	log := setupLogger(cfg.Env, cfg.Logger.Level)

	di := NewDIContainer(cfg, log)

	return &App{
		cfg: cfg,
		log: log,

		di: di,
	}, nil
}

func (a *App) Start() error {
	err := a.di.SubscriptionsRepo().Check(context.Background())
	if err != nil {
		return fmt.Errorf("check subscription repo failed: %w", err)
	}

	err = a.di.HTTPServer().Start()
	if err != nil {
		return fmt.Errorf("start http server failed: %w", err)
	}
	a.log.Info("http server started", "port", a.cfg.HTTPServer.Port)

	a.log.Info("app started successfully")
	return nil
}

func (a *App) GracefullStop() error {
	a.log.Info("app start stopping gracefully...")

	err := a.di.HTTPServer().Stop()
	if err != nil {
		return fmt.Errorf("stop http server failed: %w", err)
	}
	a.log.Info("http server stopped", "port", a.cfg.HTTPServer.Port)

	a.di.SubscriptionsRepo().Close()
	a.log.Info("subscriptions repo closed")

	a.log.Info("app gracefully stopped")
	return nil
}

func setupLogger(env string, level int) *slog.Logger {
	log := slog.Default()

	switch env {
	case config.EnvDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(level)}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(level)}),
		)
	}
	return log
}
