package main

import (
	"fmt"
	"log"
	"net/http"

	"social-network/pkg/db/sqlite"
	"social-network/pkg/handlers"
	"social-network/pkg/utils"
	ws "social-network/pkg/wsServer"
)

func main() {
	db,repos, err := sqlite.ConnectAndMigrate()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	wsServer := ws.StartServer(repos)

	server := &http.Server{
		Addr:    ":8081",
		Handler: Router(&handlers.Handler{Repos: repos}, wsServer),
	}

	fmt.Printf("Server started at http://localhost" + server.Addr + "\n")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error", err)
	}
}

// Set up all routes
func Router(handler *handlers.Handler, wsServer *ws.Server) http.Handler {
	mux := http.NewServeMux()
	/* ------------------------------ image server ------------------------------ */
	fs := http.FileServer(http.Dir("./imageUpload"))
	mux.Handle("/imageUpload/", http.StripPrefix("/imageUpload/", utils.ConfigFSHeader(fs)))
	/* ------------------------------- auth route ------------------------------- */
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/signin", handler.Signin)
	mux.HandleFunc("/logout", handler.Auth(handler.Logout))
	mux.HandleFunc("/sessionActive", handler.SessionActive)

	/* ---------------------------------- users --------------------------------- */
	mux.HandleFunc("/allUsers", handler.Auth(handler.AllUsers))       // all users + info except current
	mux.HandleFunc("/followers", handler.Auth(handler.GetFollowers))  // follower list
	mux.HandleFunc("/following", handler.Auth(handler.GetFollowing))  // following list
	mux.HandleFunc("/currentUser", handler.Auth(handler.CurrentUser)) // current user data
	mux.HandleFunc("/userData", handler.Auth(handler.UserData))       // userd data based on following status
	mux.HandleFunc("/changeStatus", handler.Auth(handler.UserStatus)) // change status

	mux.HandleFunc("/follow", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.Follow(wsServer, w, r)
	})) // follow user
	mux.HandleFunc("/cancelFollowRequest", handler.Auth(handler.CancelFollowRequest))
	mux.HandleFunc("/unfollow", handler.Auth(handler.Unfollow))
	mux.HandleFunc("/responseFollowRequest", handler.Auth(handler.ResponseFollowRequest))

	/* ---------------------------------- posts --------------------------------- */
	mux.HandleFunc("/allPosts", handler.Auth(handler.AllPosts))   // all posts- main page
	mux.HandleFunc("/userPosts", handler.Auth(handler.UserPosts)) // all user posts - user page
	mux.HandleFunc("/newPost", handler.Auth(handler.NewPost))     // create route

	/* -------------------------------- comments -------------------------------- */
	mux.HandleFunc("/newComment", handler.Auth(handler.NewComment)) // create route

	/* --------------------------------- groups --------------------------------- */
	mux.HandleFunc("/allGroups", handler.Auth(handler.AllGroups))             // group list
	mux.HandleFunc("/userGroups", handler.Auth(handler.UserGroups))           // group list of user groups
	mux.HandleFunc("/otherUserGroups", handler.Auth(handler.OtherUserGroups)) // group list for specific user

	mux.HandleFunc("/groupInfo", handler.Auth(handler.GroupInfo))                     // get group info
	mux.HandleFunc("/groupMembers", handler.Auth(handler.GroupMembers))               // get group members
	mux.HandleFunc("/groupEvents", handler.Auth(handler.GroupEvents))                 // get group events
	mux.HandleFunc("/groupPosts", handler.Auth(handler.GroupPosts))                   // get group posts
	mux.HandleFunc("/groupRequests", handler.Auth(handler.GroupRequests))             // get group member requests
	mux.HandleFunc("/cancelGroupRequests", handler.Auth(handler.CancelGroupRequests)) // cancel request or joing group

	mux.HandleFunc("/newGroup", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewGroup(wsServer, w, r)
	})) // create new group
	mux.HandleFunc("/newGroupPost", handler.Auth(handler.NewGroupPost))                           // create new group post
	mux.HandleFunc("/newGroupInvite", handler.Auth(func(w http.ResponseWriter, r *http.Request) { // invite new users to group
		handler.NewGroupInvite(wsServer, w, r)
	}))
	mux.HandleFunc("/newGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) { // invite new users to group
		handler.NewGroupRequest(wsServer, w, r)
	}))
	mux.HandleFunc("/responseGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.ResponseGroupRequest(wsServer, w, r)
	})) // response to join request
	mux.HandleFunc("/responseInviteRequest", handler.Auth(handler.ResponseInviteRequest)) // response to invite request

	/* --------------------------------- events --------------------------------- */
	mux.HandleFunc("/newEvent", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewEvent(wsServer, w, r)
	})) // create new
	mux.HandleFunc("/participate", handler.Auth(handler.Participate)) // react to participation in event

	/* ------------------------------ notifications ----------------------------- */
	mux.HandleFunc("/notifications", handler.Auth(handler.Notifications)) // get all notifs from db on login

	/* ------------------------------ chat messages ----------------------------- */
	mux.HandleFunc("/messages", handler.Auth(handler.Messages))             // get all chat messages for specific chat
	mux.HandleFunc("/unreadMessages", handler.Auth(handler.UnreadMessages)) // get list of messages that isn't read
	mux.HandleFunc("/messageRead", handler.Auth(handler.MessageRead))       // mark message as read
	mux.HandleFunc("/newMessage", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewMessage(wsServer, w, r)
	})) // new chat message
	mux.HandleFunc("/chatList", handler.Auth(handler.ChatList))                       // get list of users to display in chatbox
	mux.HandleFunc("/responseChatRequest", handler.Auth(handler.ResponseChatRequest)) // response to chat request

	/* ---------------------------- websocket server ---------------------------- */
	mux.HandleFunc("/ws", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.SocketHandler(wsServer, w, r)
	}))

	return mux
}
