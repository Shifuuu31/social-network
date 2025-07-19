package main

import (
	"log"
	"net/http"
	"os"

	"social-network/pkg/db/setup"
	"social-network/pkg/handlers"
	"social-network/pkg/middleware"
	"social-network/pkg/models"
	"social-network/pkg/websocket"
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

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	AppRoot := handlers.Root{
		DL: &middleware.DataLayer{
			Users: &models.UserModel{
				DB: db,
			},
			Sessions: &models.SessionModel{
				DB: db,
			},
			Posts: &models.PostModel{
				DB: db,
			},
			Comments: &models.CommentModel{
				DB: db,
			},
			Follows: &models.FollowRequestModel{
				DB: db,
			},
			Messages: &models.MessageModel{
				DB: db,
			},
			Logger: &models.LoggerModel{
				DB: db,
			},
		},
		Hub: hub,
	}

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += "8080"
	}
	server := http.Server{
		Addr:    port,
		Handler: AppRoot.DL.CORSMiddleware(AppRoot.DL.AccessMiddleware(AppRoot.Router())),
	}

	log.Println("server listening on http://localhost" + port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

	// Now you can use `db` for your queries
	// Start your server, handlers, etc.
}
