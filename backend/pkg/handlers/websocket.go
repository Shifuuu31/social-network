package handlers

import (
	"log"
	"net/http"
	"social-network/pkg/utils"
	ws "social-network/pkg/wsServer"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func (handler *Handler) SocketHandler(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //for CORS err

	// access user id
	userId := r.Context().Value(utils.UserKey).(string)

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	// crete new client
	client := ws.NewClient(conn, wsServer.Repos, userId)
	// register the clinet in wsServer
	wsServer.RegisterNewClient(client)

	// put in action infinit read and write functions
	go client.Writer()
	go client.Reader(wsServer)
}
