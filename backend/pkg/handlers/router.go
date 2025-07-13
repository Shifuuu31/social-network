package handlers

import (
	"fmt"
	"net/http"

	"social-network/pkg/middleware"
)

// TODO add routes...

type Root struct {
	DL *middleware.DataLayer
}

func (rt *Root) Router() (uh *http.ServeMux) {
	authMux := rt.NewAuthHandler()
	userHandler := rt.NewUsersHandler()
	imageHandler := rt.NewImageHandler()

	mainMux := http.NewServeMux()
	// mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("ROOT REQUEST: %s %s", r.Method, r.URL.Path)
	// })
	rt.SetupPostRoutes(mainMux)

	// Setup image routes
	imageMux := imageHandler.SetupImageRoutes()

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

	return mainMux
}
