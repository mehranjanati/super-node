package migrations

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Config holds the configuration for the migrations.
type Config struct {
	PostgresURL string
}

func Run(cfg *Config) {
	fmt.Println("Running migrations...")

	m, err := migrate.New(
		"file://./migrations",
		cfg.PostgresURL,
	)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database: %v", err)
	}

	log.Println("Migrations completed successfully.")
}
