package app

import (
	"fmt"

	"github.com/longwavee/effective-mobile-test/internal/config"
)

type (
	App struct {
		cfg *config.Config
		di  *DIContainer
	}
)

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}

	di := NewDIContainer(cfg)

	return &App{
		cfg: cfg,
		di:  di,
	}, nil
}

func (a *App) Start() error {
	return nil
}

func (a *App) GracefullStop() error {
	return nil
}
