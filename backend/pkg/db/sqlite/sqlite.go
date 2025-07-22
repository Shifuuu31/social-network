package sqlite

import (
	"database/sql"
	"fmt"

	"social-network/pkg/models"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // sqlite3 driver for migration
	_ "github.com/golang-migrate/migrate/v4/source/file"      // migration source driver
	_ "github.com/mattn/go-sqlite3"                           // sqlite3 driver for database/sql
)

func ConnectAndMigrate() (*sql.DB, *models.Repositories, error) {
	dbPath, migrationsPath := "./pkg/db/data.db", "./pkg/db/migration/sqlite"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return db, nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return db, nil, fmt.Errorf("failed to enable pragma: %w", err)
	}

	// enable concurrent sqlite writes
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return db, nil, fmt.Errorf("failed to enable pragma: %w", err)
	}
	// Create migrate instance
	m, err := migrate.New(
		"file://"+migrationsPath,
		"sqlite3://"+dbPath,
	)
	if err != nil {
		return db, nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return db, nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, &models.Repositories{
		UserRepo:    &UserRepository{DB: db},
		SessionRepo: &SessionRepository{DB: db},
		GroupRepo:   &GroupRepository{DB: db},
		PostRepo:    &PostRepository{DB: db},
		CommentRepo: &CommentRepository{DB: db},
		NotifRepo:   &NotifRepository{DB: db},
		EventRepo:   &EventRepository{DB: db},
		MsgRepo:     &MsgRepository{DB: db},
	}, nil
}
