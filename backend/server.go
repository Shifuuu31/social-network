package main

import (
	"log"
	"net/http"
	"os"

	"social-network/pkg/db/setup"
	"social-network/pkg/handlers"
	"social-network/pkg/middleware"
	"social-network/pkg/models"
)

func main() {
	dbPath := "./pkg/db/data.db"
	migrationsPath := "./pkg/db/migrations/sqlite"

	db, err := setup.ConnectAndMigrate(dbPath, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_ = db

	AppRoot := handlers.Root{
		DL: &middleware.DataLayer{
			Users: &models.UserModel{
				DB: db,
			},
			Sessions: &models.SessionModel{
				DB: db,
			},
			Logger: &models.LoggerModel{
				DB: db,
			},
		},
	}

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += "8080"
	}
	server := http.Server{
		Addr:    port,
		Handler: AppRoot.Router(),
	}

	log.Println("server listening on http://localhost" + port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

	// Now you can use `db` for your queries
	// Start your server, handlers, etc.
}
