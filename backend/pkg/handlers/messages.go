package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"social-network/pkg/models"
	"social-network/pkg/utils"
	ws "social-network/pkg/wsServer"
)

// get all previous messages for chat
// waits for POST request with RECEIVER as target and TYPE
// respondes with all messages through simple http response
func (handler *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	/* ------------------- // get incoming data in msg format ------------------- */
	var msgIn models.ChatMessage
	err := json.NewDecoder(r.Body).Decode(&msgIn)
	if err != nil {
		utils.RespondWithError(w, "Error on reading the incomming message", 200)
		return
	}
	msgIn.SenderId = r.Context().Value(utils.UserKey).(string)

	/* ----------------------- get massages form database ----------------------- */
	var messages []models.ChatMessage
	if msgIn.Type == "PERSON" {
		messages, err = handler.Repos.MsgRepo.GetAll(msgIn)
		if err != nil {
			utils.RespondWithError(w, "Error on getting the messages", 200)
			return
		}
		// mark as read
		for i := 0; i < len(messages); i++ {
			// if current user is also sender of message, then skip
			if messages[i].SenderId == msgIn.SenderId {
				continue
			}
			err = handler.Repos.MsgRepo.MarkAsRead(messages[i])
			if err != nil {
				utils.RespondWithError(w, "Error on marking message as read", 200)
				return
			}
		}
		// if no messages so far, check if request made and add message
		if len(messages) == 0 {
			requetExists, err := handler.Repos.NotifRepo.CheckIfChatRequestExists(msgIn.SenderId, msgIn.ReceiverId)
			if err != nil {
				utils.RespondWithError(w, "Error on checking chat history", 200)
				return
			}
			if requetExists {
				msgContent, _ := handler.Repos.NotifRepo.GetContentFromChatRequest(msgIn.SenderId, msgIn.ReceiverId)
				newMessage := models.ChatMessage{ID: "0", SenderId: msgIn.SenderId, ReceiverId: msgIn.ReceiverId, Content: msgContent, Type: "PERSON"}
				messages = append(messages, newMessage)
			}
		}
	} else if msgIn.Type == "GROUP" {
		messages, err = handler.Repos.MsgRepo.GetAllGroup(msgIn.SenderId, msgIn.ReceiverId)
		if err != nil {
			utils.RespondWithError(w, "Error on getting the messages", 200)
			return
		}
		// mark as read
		for i := 0; i < len(messages); i++ {
			// if current user is also sender of message, then skip
			if messages[i].SenderId == msgIn.SenderId {
				continue
			}
			err = handler.Repos.MsgRepo.MarkAsReadGroup(models.ChatMessage{ID: messages[i].ID, ReceiverId: msgIn.SenderId})
			if err != nil {
				utils.RespondWithError(w, "Error on marking message as read", 200)
				return
			}
		}
	}
	/* --------------------------- attach sender data --------------------------- */
	for i := 0; i < len(messages); i++ {
		messages[i].Sender, _ = handler.Repos.UserRepo.GetDataMin(messages[i].SenderId)
	}

	utils.RespondWithMessages(w, messages, 200)
}

// new chat message wits for POST requet with SENDER, RECEIVER AND TYPE
// function saves new message and responds
// to sender with regular http response
// to recievers through websocket connection
func (handler *Handler) NewMessage(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	/* --------------------------- read incoming data --------------------------- */
	var msg models.ChatMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		utils.RespondWithError(w, "Error on reading the incomming message", 200)
		return
	}

	/* -------------------- attach sender id ------------------------------------ */
	msg.SenderId = r.Context().Value(utils.UserKey).(string)

	newChatFlag := "" // flag is raised with valu "NEW" if tha chat does not exist for user yet

	// check if receiver is following current user
	isFollowingBack, err := handler.Repos.UserRepo.IsFollowing(msg.SenderId, msg.ReceiverId)
	if err != nil {
		utils.RespondWithError(w, "Error on saving checking status", 200)
		return
	}

	isGroupMember, err := handler.Repos.GroupRepo.IsGroupMember(msg.ReceiverId, msg.SenderId)
	if err != nil {
		utils.RespondWithError(w, "Error on checking if user is member", 200)
		return
	}

	isGroupAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(msg.ReceiverId, msg.SenderId)
	if err != nil {
		utils.RespondWithError(w, "Error on checking if user is admin", 200)
		return
	}

	// if he is private and have no chat history and not group member, create notification insted of saving msg
	if !isFollowingBack && !isGroupMember && !isGroupAdmin {
		status, err := handler.Repos.UserRepo.GetStatus(msg.ReceiverId)
		if err != nil {
			utils.RespondWithError(w, "Error on saving checking status", 200)
			return
		}
		hasHistory, err := handler.Repos.MsgRepo.HasHistory(msg.SenderId, msg.ReceiverId)
		if err != nil {
			utils.RespondWithError(w, "Error on checking chat history", 200)
			return
		}
		if status == "PRIVATE" && !hasHistory {
			// check if request is already made
			requestExists, err := handler.Repos.NotifRepo.CheckIfChatRequestExists(msg.SenderId, msg.ReceiverId)
			if err != nil {
				utils.RespondWithError(w, "Internal server error", 200)
				return
			}
			if requestExists {
				utils.RespondWithError(w, "Chat request already saved.\n Wait for user to respond to your request.", 200)
				return
			}
			// save msg in notification table
			newNotif := models.Notification{
				ID:       utils.UniqueId(),
				TargetID: msg.ReceiverId,
				Type:     "CHAT_REQUEST",
				Content:  msg.Content,
				Sender:   msg.SenderId,
			}
			err = handler.Repos.NotifRepo.Save(newNotif)
			if err != nil {
				utils.RespondWithError(w, "Internal server error", 200)
				return
			}
			// NOTIFY  RECEIVER ABOUT THE NEW CHAT REQUEST IF ONLINE
			for client := range wsServer.Clients {
				if client.ID == newNotif.TargetID {
					client.SendNotification(newNotif)
				}
			}
			utils.RespondWithSuccess(w, "New request saved", 200)
			return
		} else if status == "PUBLIC" && !hasHistory {
			newChatFlag = "NEW"
		}
	}
	/* --------------------------- generate message id -------------------------- */
	msg.ID = utils.UniqueId()
	/* ---------------------------- save in database ---------------------------- */
	err = handler.Repos.MsgRepo.Save(msg)
	if err != nil {
		fmt.Println("MSG", msg)
		fmt.Println("ERR", err)

		utils.RespondWithError(w, "Error on saving message", 200)
		return
	}
	/* --------------------------- attach sender  info -------------------------- */
	msg.Sender, _ = handler.Repos.UserRepo.GetDataMin(msg.SenderId)
	// send message respond with new message to sender
	utils.RespondWithMessages(w, []models.ChatMessage{msg}, 200)

	/* ------------------ respond through websocket to receiver ----------------- */
	if msg.Type == "PERSON" {
		for client := range wsServer.Clients {
			if client.ID == msg.ReceiverId {
				client.SendChatMessage(msg, newChatFlag)
			}
		}
	} else if msg.Type == "GROUP" { // incase of group find and respond to all recievers
		// find all group members + admin except the sender
		allMembers, err := handler.Repos.GroupRepo.GetGroupMembers(msg.ReceiverId)
		if err != nil {
			utils.RespondWithError(w, "Error on geting group members", 200)
			return
		}
		for i := 0; i < len(allMembers); i++ {
			if allMembers[i].ID != msg.SenderId {
				/* -------- save also in group_messages table -------- */
				err = handler.Repos.MsgRepo.SaveGroupMsg(models.ChatMessage{ID: msg.ID, ReceiverId: allMembers[i].ID})
			}
			for client := range wsServer.Clients {
				if client.ID == msg.SenderId {
					continue
				}
				if client.ID == allMembers[i].ID {
					client.SendChatMessage(msg, "")
				}
			}

		}
	}
}

// respond with list of messages, that user has missed
// in response include message_id, senderId (group id or user id) type -> group or person
func (handler *Handler) UnreadMessages(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	userId := r.Context().Value(utils.UserKey).(string)
	/* ---------- collect chat stats from db for prvate and group chats --------- */
	messages, err := handler.Repos.MsgRepo.GetUnread(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting the unread messages", 200)
		return
	}
	groupMessages, err := handler.Repos.MsgRepo.GetUnreadGroup(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting the unread messages", 200)
		return
	}

	/* ------------------------- combine all chat stats ------------------------- */
	allUnreadMessages := append(messages, groupMessages...)

	/* ---------------------------- respond to client --------------------------- */
	utils.RespondWithChatStats(w, allUnreadMessages, 200)
}

// handler needs data about msg ->  id and type
// and it marks it as read in database
func (handler *Handler) MessageRead(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	/* --------------------------- read incoming data --------------------------- */
	var msg models.ChatMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		utils.RespondWithError(w, "Error on reading the incomming message", 200)
		return
	}
	// attach current user id
	msg.ReceiverId = r.Context().Value(utils.UserKey).(string)
	if msg.Type == "GROUP" {
		err = handler.Repos.MsgRepo.MarkAsReadGroup(msg)
		if err != nil {
			utils.RespondWithError(w, "Error on marking message as read", 200)
			return
		}
	} else if msg.Type == "PERSON" {
		err = handler.Repos.MsgRepo.MarkAsRead(msg)
		if err != nil {
			utils.RespondWithError(w, "Error on marking message as read", 200)
			return
		}
	} else {
		utils.RespondWithError(w, "Error. Message type not provided or not recognized", 200)
		return
	}
	utils.RespondWithSuccess(w, "Message marked as read successfuly", 200)
}

func (handler *Handler) ResponseChatRequest(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new response
	type Response struct {
		RequestID string `json:"requestId"`
		Response  string `json:"response"` // ACCEPT or DECLINE
	}
	var resp Response
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	if strings.ToUpper(resp.Response) == "ACCEPT" {
		notifData, err := handler.Repos.NotifRepo.GetCahtNotifById(resp.RequestID)
		if err != nil {
			utils.RespondWithError(w, "Error on getting notification", 200)
			return
		}
		// save new message
		newMsg := models.ChatMessage{
			ID:         utils.UniqueId(),
			SenderId:   notifData.Sender,
			ReceiverId: notifData.TargetID,
			Type:       "PERSON",
			Content:    notifData.Content,
		}
		err = handler.Repos.MsgRepo.Save(newMsg)
		if err != nil {
			utils.RespondWithError(w, "Error on saving message", 200)
			return
		}
	}
	/* ----------------------- delete pending notification ---------------------- */
	err = handler.Repos.NotifRepo.Delete(resp.RequestID)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}

	utils.RespondWithSuccess(w, "Response successful", 200)
}
