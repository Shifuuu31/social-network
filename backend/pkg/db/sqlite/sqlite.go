package sqlite

import (
    "database/sql"
    "fmt"

    _ "github.com/golang-migrate/migrate/v4/source/file"          // migration source driver
    _ "github.com/golang-migrate/migrate/v4/database/sqlite3"    // sqlite3 driver for migration
    _ "github.com/mattn/go-sqlite3"                              // sqlite3 driver for database/sql
    "github.com/golang-migrate/migrate/v4"
)

func ConnectAndMigrate(dbPath string, migrationsPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
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
