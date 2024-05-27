package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/oustrix/ozon_journal/config"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// migrateUp applies all migrations to the database
func migrateUp(cfg *config.Postgres) error {
	var (
		attempts = cfg.ConnAttempts
		err      error
		m        *migrate.Migrate
	)

	// Try to connect to postgres.
	for attempts > 0 {
		m, err = migrate.New("file://migrations", cfg.DSN)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(cfg.ConnTimeout) * time.Second)
		attempts--
	}

	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if m == nil {
		return fmt.Errorf("migration is nil. check attempts to connect to postgres")
	}

	// Apply migrations.
	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
