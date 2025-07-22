package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"social-network/pkg/models"
	"social-network/pkg/utils"
	ws "social-network/pkg/wsServer"
)

func (handler *Handler) NewEvent(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to Event
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	event.ID = utils.UniqueId()
	event.AuthorID = r.Context().Value(utils.UserKey).(string)
	/* -------------------- check if user is a meber of group ------------------- */
	isMember := false
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(event.GroupID, event.AuthorID)
	if err != nil {
		utils.RespondWithError(w, "Error on reading role", 200)
		return
	}
	if !isAdmin {
		isMember, err = handler.Repos.GroupRepo.IsGroupMember(event.GroupID, event.AuthorID)
		if err != nil {
			utils.RespondWithError(w, "Error on checking if is group member", 200)
			return
		}
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	/* ------------------------- save event in database ------------------------- */
	if err = handler.Repos.EventRepo.Save(event); err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	/* ----------------- if user going also save as participant ----------------- */
	if strings.ToUpper(event.Going) == "YES" {
		if err = handler.Repos.EventRepo.AddParticipant(event.ID, event.AuthorID); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}
	/* -------------------- save new notification about event ------------------- */
	// get all group members
	members, err := handler.Repos.GroupRepo.GetGroupMembers(event.GroupID)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}

	// for each member create notification
	for i := 0; i < len(members); i++ {
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: members[i].ID,
			Type:     "EVENT",
			Content:  event.ID,
			Sender:   event.AuthorID,
		}
		// save notification in database
		if members[i].ID != event.AuthorID {
			err = handler.Repos.NotifRepo.Save(newNotif)
			if err != nil {
				utils.RespondWithError(w, "Internal server error", 200)
				return
			}
		}

		// NOTIFY  GROUP MEMBER ABOUT THE NEW EVENT IF ONLINE
		for client := range wsServer.Clients {
			if client.ID == members[i].ID && client.ID != event.AuthorID {
				client.SendNotification(newNotif)
			}
		}

	}
	utils.RespondWithEvents(w, []models.Event{event}, 200)
}

// requestId: notification.id,
// eventId: notification.event.id,
// response: reqResponse,
// Handles clients reaction to participation in event
// waits for POST req with eventID as "id" and user status "going" with response YES or NO
func (handler *Handler) Participate(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* --------------------------- read incoming data --------------------------- */
	type Response struct {
		EventID   string `json:"eventId"`
		RequestID string `json:"requestId"` // notif id
		Response  string `json:"response"`  // YES || NO
	}
	var response Response
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	// get current user
	userId := r.Context().Value(utils.UserKey).(string)

	/* ---------------- check that event id and response provided --------------- */
	if len(response.EventID) == 0 || len(response.Response) == 0 {
		utils.RespondWithError(w, "Provided incomplete data", 200)
		return
	}

	/* ------------------- check if response alredy registerd ------------------- */
	isParticipating, err := handler.Repos.EventRepo.IsParticipating(response.EventID, userId)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	/* ----------------------------- handle response ---------------------------- */
	if response.Response == "YES" && !isParticipating {
		if err = handler.Repos.EventRepo.AddParticipant(response.EventID, userId); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	} else if strings.ToUpper(response.Response) == "NO" && isParticipating {
		if err = handler.Repos.EventRepo.RemoveParticipant(response.EventID, userId); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}
	/* --------------------------- remove notificaton -------------------------- */
	if len(response.RequestID) != 0 { // participation activated form notification
		if err = handler.Repos.NotifRepo.Delete(response.RequestID); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	} else { // participation activated without noification
		// delete notification if exists
		notif := models.Notification{Type: "EVENT", TargetID: userId, Content: response.EventID}
		if err = handler.Repos.NotifRepo.DeleteByType(notif); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

	}
	utils.RespondWithSuccess(w, "Data saved successfully", 200)
}
