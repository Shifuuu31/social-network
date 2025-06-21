package main

import (
	"log"

	"social-network/pkg/db/sqlite"
)

func main() {
	dbPath := "./pkg/db/data.db"
	migrationsPath := "./pkg/db/migrations/sqlite"

	db, err := sqlite.ConnectAndMigrate(dbPath, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Now you can use `db` for your queries
	// Start your server, handlers, etc.
}
