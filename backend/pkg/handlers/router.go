package handlers

import (
	"net/http"

	"social-network/pkg/middleware"
)



type Root struct {
	DL *middleware.DataLayer
}

func (rt *Root) Router() (uh *http.ServeMux) {
	authMux := rt.NewAuthHandler()
	usersHandler := rt.NewUsersHandler()
	groupsHandler := rt.NewGroupsHandler()

	mainMux := http.NewServeMux()

	// Mount sub-muxes under prefixes
	mainMux.Handle("/api/auth/", http.StripPrefix("/api/auth", authMux))
	mainMux.Handle("/api/users/", http.StripPrefix("/api/users", usersHandler))
	mainMux.Handle("/api/groups/", http.StripPrefix("/api/groups", groupsHandler))

	return mainMux
}
