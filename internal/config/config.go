package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type (
	Config struct {
		Env        string `env:"ENV" env-default:"prod"`
		Logger     Logger
		HTTPServer HTTPServer
		Postgres   Postgres
	}

	Logger struct {
		Level int `env:"LOGGER_LEVEL" env-default:"0"`
	}

	HTTPServer struct {
		Port              string        `env:"HTTP_SERVER_PORT" env-default:"8080"`
		ReadTimeout       time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" env-default:"1s"`
		ReadHeaderTimeout time.Duration `env:"HTTP_SERVER_READ_HEADER_TIMEOUT" env-default:"1s"`
		WriteTimeout      time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" env-default:"1s"`
		IdleTimeout       time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
		MaxHeaderBytes    int           `env:"HTTP_SERVER_MAX_HEADER_BYTES" env-default:"500"`
		ShutdownTimeout   time.Duration `env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" env-default:"1s"`

		ReadinessTimeout time.Duration `env:"HTTP_SERVER_READINESS_TIMEOUT" env-default:"30s"`
	}

	Postgres struct {
		Host     string `env:"POSTGRES_HOST" env-default:"postgres"`
		Port     string `env:"POSTGRES_PORT" env-default:"5432"`
		User     string `env:"POSTGRES_USER" env-required:"true"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		DBName   string `env:"POSTGRES_DB" env-required:"true"`
		SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
	}
)

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env failed: %w", err)
	}

	return &cfg, nil
}

func (p *Postgres) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.DBName,
		p.SSLMode,
	)
}
