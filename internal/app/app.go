package app

import (
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
	a.log.Info("app started successfully")
	return nil
}

func (a *App) GracefullStop() error {
	a.log.Info("app start stopping gracefully...")
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
