package handlers

import (
	"net/http"

	"social-network/pkg/middleware"
)

type Root struct {
	DL  *middleware.DataLayer
	Hub *WSHub
}

func (rt *Root) Router() (uh *http.ServeMux) {
	var (
		wsHanler      = rt.NewWSHandler()
		authMux       = rt.NewAuthHandler()
		usersHandler  = rt.NewUsersHandler()
		groupsHandler = rt.NewGroupsHandler()
		filesHandler  = rt.NewServeFilesHandler()
		postsHandler  = rt.SetupPostRoutes()
	)

	mainMux := http.NewServeMux()

	// Mount sub-muxes under prefixes
	mainMux.Handle("/posts/", http.StripPrefix("/posts", postsHandler))
	mainMux.Handle("/", http.StripPrefix("/ ", wsHanler))
	mainMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mainMux.Handle("/users/", http.StripPrefix("/users", usersHandler))
	mainMux.Handle("/groups/", http.StripPrefix("/groups", groupsHandler))
	mainMux.Handle("/get/", http.StripPrefix("/get", filesHandler))

	return mainMux
}
