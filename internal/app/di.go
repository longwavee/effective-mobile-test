package app

import (
	"log/slog"

	"github.com/longwavee/effective-mobile-test/internal/config"
)

type (
	DIContainer struct {
		cfg *config.Config
		log *slog.Logger
	}
)

func NewDIContainer(cfg *config.Config, log *slog.Logger) *DIContainer {
	return &DIContainer{
		cfg: cfg,
		log: log,
	}
}
