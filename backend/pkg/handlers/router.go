package handlers

import (
	"fmt"
	"net/http"

	"social-network/pkg/middleware"
)

type Root struct {
	DL *middleware.DataLayer
}

func (rt *Root) Router() (uh *http.ServeMux) {
	fmt.Println("router")
	authMux := rt.NewAuthHandler()
	usersHandler := rt.NewUsersHandler()
	groupsHandler := rt.NewGroupsHandler()
	postsHandler := rt.NewPostHandler()

	mainMux := http.NewServeMux()

	// Mount sub-muxes under prefixes
	mainMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mainMux.Handle("/users/", http.StripPrefix("/users", usersHandler))
	mainMux.Handle("/groups/", http.StripPrefix("/groups", groupsHandler))
	mainMux.Handle("/posts/", http.StripPrefix("/posts", postsHandler))

	return mainMux
}
