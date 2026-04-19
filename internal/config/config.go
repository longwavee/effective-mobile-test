package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type (
	Config struct {
		Env string `env:"ENV" env-default:"prod"`
	}
)

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env failed: %w", err)
	}

	return &cfg, nil
}
