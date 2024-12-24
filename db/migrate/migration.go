package db

import (
	"database/sql"
	"github/Shubhpreet-Rana/projects/internal/logging"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDatabase handles database migrations (up or down)
func MigrateDatabase(db *sql.DB, cmd string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logging.ErrorLogger.Fatalf("Failed to set up migration driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://db/migrate/migrations", "postgres", driver)
	if err != nil {
		logging.ErrorLogger.Fatalf("Failed to create migration instance: %v", err)
	}

	switch cmd {
	case "up":
		logging.InfoLogger.Println("Applying migrations...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logging.ErrorLogger.Fatalf("Migration failed: %v", err)
		}
	case "down":
		logging.InfoLogger.Println("Reverting migrations...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			logging.ErrorLogger.Fatalf("Migration revert failed: %v", err)
		}
	default:
		logging.InfoLogger.Println("No migration command specified")
	}
}
