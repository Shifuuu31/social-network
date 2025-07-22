package ws

import (
	"log"
	"social-network/pkg/models"
	"social-network/pkg/utils"

	"github.com/gorilla/websocket"
)

// represents single websocket client
type Client struct {
	ID    string
	conn  *websocket.Conn      //ws connection
	send  chan []byte          //sedn channel for outgoing messages
	repos *models.Repositories //connection to db actions
}

func NewClient(conn *websocket.Conn, repos *models.Repositories, ID string) *Client {
	return &Client{
		ID:    ID,
		conn:  conn,
		send:  make(chan []byte, 256),
		repos: repos,
	}
}

/* -------------------------------------------------------------------------- */
/*                           client action functions                          */
/* -------------------------------------------------------------------------- */

// Configure the notification with additional data about sender || group
// Change the content to reusable sentence
// send notification to client
func (client *Client) SendNotification(notif models.Notification) {
	switch notif.Type {
	case "GROUP_INVITE":
		notif.Group, _ = client.repos.GroupRepo.GetGroupData(notif.Content)
		notif.User, _ = client.repos.UserRepo.GetDataMin(notif.Sender)
	case "FOLLOW":
		notif.User, _ = client.repos.UserRepo.GetDataMin(notif.Content)
	case "EVENT":
		notif.Event, _ = client.repos.EventRepo.GetData(notif.Content)
		notif.User, _ = client.repos.UserRepo.GetDataMin(notif.Sender)
		notif.Group,_ = client.repos.GroupRepo.GetGroupData(notif.Event.GroupID)
	case "GROUP_REQUEST":
		notif.User, _ = client.repos.UserRepo.GetDataMin(notif.Content)
		notif.Group, _ = client.repos.GroupRepo.GetGroupData(notif.TargetID)
	case "CHAT_REQUEST":
		notif.User, _ = client.repos.UserRepo.GetDataMin(notif.Sender)
	}
	/* ---------------------------- add message text ---------------------------- */
	utils.DefineNotificationMsg(&notif)

	/* ---------------------------- construct message --------------------------- */
	message := WsMessage{
		Action:       NotificationAction,
		Notification: notif,
	}
	/* ---------------------------------- send ---------------------------------- */
	client.send <- message.encode()
}

func (client *Client) SendChatMessage(msg models.ChatMessage, flag string) {
	message := WsMessage{
		Action:      ChatAction,
		ChatMessage: msg,
		Message: flag,
	}

	client.send <- message.encode()
}

func (client *Client) SendGroupRequestAccept(groupId string) {
	message := WsMessage{
		Action:  GroupAcceptAction,
		Message: groupId,
	}

	client.send <- message.encode()
}

/* -------------------------------------------------------------------------- */
/*                    basic reader and writer for websocket conn              */
/* -------------------------------------------------------------------------- */
// define a writer which will send
// new messages to our WebSocket endpoint
func (client *Client) Writer() {
	for {

		message, ok := <-client.send
		if !ok {
			log.Println("err on writing message")
			return
		}
		w, err := client.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		_, err = w.Write(message)
		if err != nil {
			log.Println("Line 91", err)
			return
		}

		if err := w.Close(); err != nil {
			log.Println("Line 95", err)
			return
		}
	}
}

// define a reader which will listen for
// new messages being sent to our WebSocketendpoint
// Unregister client when client disconnect
func (client *Client) Reader(wsServer *Server) {
	defer wsServer.UnregisterClient(client)
	for {
		// read in a message
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// log.Println(msg)
	}
}
