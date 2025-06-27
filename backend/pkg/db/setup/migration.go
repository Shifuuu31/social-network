package setup

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // sqlite3 driver for migration
	_ "github.com/golang-migrate/migrate/v4/source/file"      // migration source driver
	_ "github.com/mattn/go-sqlite3"                           // sqlite3 driver for database/sql
)

func ConnectAndMigrate(dbPath string, migrationsPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("failed to enable pragma: %w", err)
	}

	// enable concurrent sqlite writes
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to enable pragma: %w", err)
	}
	// Create migrate instance
	m, err := migrate.New(
		"file://"+migrationsPath,
		"sqlite3://"+dbPath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}
