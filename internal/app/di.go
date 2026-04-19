package app

import (
	"github.com/longwavee/effective-mobile-test/internal/config"
)

type (
	DIContainer struct {
		cfg *config.Config
	}
)

func NewDIContainer(cfg *config.Config) *DIContainer {
	return &DIContainer{
		cfg: cfg,
	}
}
