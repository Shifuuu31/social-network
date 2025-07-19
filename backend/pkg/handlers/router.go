package handlers

import (
	"fmt"
	"net/http"

	"social-network/pkg/middleware"
	"social-network/pkg/websocket"
)

// TODO add routes...

type Root struct {
	DL  *middleware.DataLayer
	Hub *websocket.Hub
}

func (rt *Root) Router() (uh *http.ServeMux) {
	authMux := rt.NewAuthHandler()
	userHandler := rt.NewUsersHandler()
	imageHandler := rt.NewImageHandler()
	chatHandler := NewChatHandler(rt.DL, rt.Hub)

	mainMux := http.NewServeMux()
	// mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("ROOT REQUEST: %s %s", r.Method, r.URL.Path)
	// })
	rt.SetupPostRoutes(mainMux)

	// Setup image routes
	imageMux := imageHandler.SetupImageRoutes()

	// Setup chat routes
	chatMux := chatHandler.SetupChatRoutes()

	// Add a test route
	mainMux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test route working!")
	})

	// Add image serving route to main router using traditional pattern
	mainMux.HandleFunc("/images/", rt.ServeImageHandler)

	// Mount sub-muxes under prefixes
	mainMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mainMux.Handle("/users/", http.StripPrefix("/users", userHandler))
	mainMux.Handle("/upload/", http.StripPrefix("/upload", imageMux))
	mainMux.Handle("/chat/", http.StripPrefix("/chat", chatMux))

	return mainMux
}
