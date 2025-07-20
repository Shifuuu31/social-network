package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"social-network/pkg/db/setup"
	"social-network/pkg/handlers"
	"social-network/pkg/middleware"
	"social-network/pkg/models"
)

func NewApp(db *sql.DB) *handlers.Root {
	return &handlers.Root{
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
			Follows: &models.FollowRequestModel{
				DB: db,
			},
			Groups: &models.GroupModel{
				DB: db,
			},
			Members: &models.GroupMemberModel{
				DB: db,
			},
			Events: &models.EventModel{
				DB: db,
			},
			Votes: &models.EventVoteModel{
				DB: db,
			},
			Messages: &models.MessageModel{
				DB: db,
			},
			Notifications: &models.NotificationModel{
				DB: db,
			},
			Images: &models.ImageModel{
				DB: db,
			},

			Logger: &models.LoggerModel{
				DB: db,
			},
		},
		Hub: handlers.NewHub(),
	}
}

func main() {
	dbPath := "./pkg/db/data.db"
	migrationsPath := "./pkg/db/migrations/sqlite"

	db, err := setup.ConnectAndMigrate(dbPath, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	App := NewApp(db)

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += "8080"
	}

	// Apply middleware to the router
	router := App.Router()
	handlerWithMiddleware := App.DL.CORSMiddleware(
		App.DL.RecoverMiddleware(router),
	)

	server := http.Server{
		Addr:    port,
		Handler: handlerWithMiddleware,
	}

	log.Println("server listening on http://localhost" + port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

	// Now you can use `db` for your queries
	// Start your server, handlers, etc.
}
