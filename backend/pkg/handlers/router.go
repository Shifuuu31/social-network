package handlers

import (
	"net/http"

	"social-network/pkg/middleware"
)

// TODO add routes...

type Root struct {
	DL *middleware.DataLayer
}

func (rt *Root) Router() (uh *http.ServeMux) {
	authMux := rt.NewAuthHandler()
	userHandler := rt.NewUserHandler()

	mainMux := http.NewServeMux()

	// Mount sub-muxes under prefixes
	mainMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mainMux.Handle("/user/", http.StripPrefix("/user", userHandler))

	return mainMux
}
