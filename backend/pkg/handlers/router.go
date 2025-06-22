package handlers

import (
	"net/http"

	"social-network/pkg/middleware"
)

// TODO add routes...

type Root struct {
	DL *middleware.DataLayer
}

func (rt *Root)Router() (uh *http.ServeMux) {
	userMux := rt.NewUserHandler()
	// userHandler := handlers.NewUserHandler()

	mainMux := http.NewServeMux()

	// Mount sub-muxes under prefixes
	mainMux.Handle("/user/", http.StripPrefix("/user", userMux))
	// mainMux.Handle("/user/", http.StripPrefix("/user", userHandler.Mux))

	return mainMux
}
