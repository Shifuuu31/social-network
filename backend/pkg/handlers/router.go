package handlers

import (
	// "fmt"
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

	mainMux := http.NewServeMux()
	// mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("ROOT REQUEST: %s %s", r.Method, r.URL.Path)
	// })
	rt.SetupPostRoutes(mainMux)

	// Mount sub-muxes under prefixes
	mainMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mainMux.Handle("/users/", http.StripPrefix("/users", userHandler))

	return mainMux
}
