package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var connectionStr string
	var migrationsPath string
	var migrationsTable string

	flag.StringVar(&connectionStr, "connection-str", "", "(e.g. postgres://user:pass@localhost:5432/db?sslmode=disable)")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations folder")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if connectionStr == "" {
		log.Error("connection-str is required")
		os.Exit(1)
	}
	if migrationsPath == "" {
		log.Error("migrations-path is required")
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("%s&x-migrations-table=%s", connectionStr, migrationsTable),
	)
	if err != nil {
		log.Error("failed to initialize migrator", "err", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply")
			return
		}

		log.Error("failed to apply migrations", "err", err)
		os.Exit(1)
	}

	log.Info("migrations applied successfully")
}
