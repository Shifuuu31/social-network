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
	usersHandler := rt.NewUsersHandler()
	groupsHandler := rt.NewGroupsHandler()

	mainMux := http.NewServeMux()

	// Mount sub-muxes under prefixes
	mainMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mainMux.Handle("/users/", http.StripPrefix("/users", usersHandler))
	mainMux.Handle("/groups/", http.StripPrefix("/groups", groupsHandler))

	return mainMux
}
