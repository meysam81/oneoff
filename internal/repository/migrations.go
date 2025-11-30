package repository

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/meysam81/oneoff/internal/logging"
)

// RunMigrations runs database migrations
func RunMigrations(dbPath string, migrationsPath string, direction string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	logging.Info().
		Uint("current_version", currentVersion).
		Bool("dirty", dirty).
		Msg("Current migration status")

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migration up failed: %w", err)
		}
		logging.Info().Msg("Migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migration down failed: %w", err)
		}
		logging.Info().Msg("Migrations rolled back successfully")

	default:
		return fmt.Errorf("invalid migration direction: %s (must be 'up' or 'down')", direction)
	}

	newVersion, _, _ := m.Version()
	logging.Info().Uint("new_version", newVersion).Msg("Migration completed")

	return nil
}
