package handlers

import (
	"net/http"

	"social-network/pkg/middleware"
)

type Root struct {
	DL  *middleware.DataLayer
	Hub *WSHub
}

func (rt *Root) Router() http.Handler {
	var (
		wsHanler             = rt.NewWSHandler()
		authMux              = rt.NewAuthHandler()
		usersHandler         = rt.NewUsersHandler()
		groupsHandler        = rt.NewGroupsHandler()
		filesHandler         = rt.NewServeFilesHandler()
		postsHandler         = rt.SetupPostRoutes()
		notificationsHandler = rt.NewNotificationsHandler()
	)

	mainMux := http.NewServeMux()

	mainMux.Handle("/connect", wsHanler)
	mainMux.Handle("/posts/", http.StripPrefix("/posts", postsHandler))
	mainMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mainMux.Handle("/users/", http.StripPrefix("/users", usersHandler))
	mainMux.Handle("/groups/", http.StripPrefix("/groups", groupsHandler))
	mainMux.Handle("/get/", http.StripPrefix("/get", filesHandler))
	mainMux.Handle("/notifications/", http.StripPrefix("/notifications", notificationsHandler))

	return rt.DL.GlobalMiddleware(mainMux)
}
