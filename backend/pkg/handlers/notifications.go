package handlers

import (
	"net/http"

	"social-network/pkg/utils"
)

func (handler *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	notifs, err := handler.Repos.NotifRepo.GetAll(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* ------------------ populate additional notification data ----------------- */
	for i := 0; i < len(notifs); i++ {
		// get user || group || invite for notif
		switch notifs[i].Type {
		case "GROUP_INVITE":
			notifs[i].Group, _ = handler.Repos.GroupRepo.GetGroupData(notifs[i].Content)
			notifs[i].User, _ = handler.Repos.UserRepo.GetDataMin(notifs[i].Sender)
		case "FOLLOW":
			notifs[i].User, _ = handler.Repos.UserRepo.GetDataMin(notifs[i].Content)
		case "EVENT":
			notifs[i].Event, _ = handler.Repos.EventRepo.GetData(notifs[i].Content)
			notifs[i].User, _ = handler.Repos.UserRepo.GetDataMin(notifs[i].Sender)
			notifs[i].Group, _ = handler.Repos.GroupRepo.GetGroupData(notifs[i].Event.GroupID)
		case "GROUP_REQUEST":
			notifs[i].User, _ = handler.Repos.UserRepo.GetDataMin(notifs[i].Content)
			notifs[i].Group, _ = handler.Repos.GroupRepo.GetGroupData(notifs[i].TargetID)
		}
		// change msg
		utils.DefineNotificationMsg(&notifs[i])
	}
	utils.RespondWithNotifications(w, notifs, 200)
}
