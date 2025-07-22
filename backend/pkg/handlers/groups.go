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

/* -------------------------------------------------------------------------- */
/*                                  get data                                  */
/* -------------------------------------------------------------------------- */

func (handler *Handler) AllGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all groups + relations
	groups, errGroups := handler.Repos.GroupRepo.GetAllAndRelations(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns all groups that current user is a member of or admin
func (handler *Handler) UserGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request user Groups
	groups, errGroups := handler.Repos.GroupRepo.GetUserGroups(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns all groups that specified user is a member of or admin
func (handler *Handler) OtherUserGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	query := r.URL.Query()
	userId := query.Get("userId")
	if userId == "" { // check if user id provided in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// request user Groups
	groups, errGroups := handler.Repos.GroupRepo.GetUserGroups(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns info about group - > name, description, id and administrator id
// also includes group status for current user -> admin / member or pending member request
func (handler *Handler) GroupInfo(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { // check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	group, err := handler.Repos.GroupRepo.GetGroupData(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* --------------- // add additional data about member status --------------- */
	// check if admin, or member or member request pending
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	if userId == group.AdminID {
		group.Administrator = true
	} else {
		group.Member, err = handler.Repos.GroupRepo.IsGroupMember(group.ID, userId)
		if err != nil {
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
		if !group.Member {
			notification := models.Notification{
				TargetID: group.ID,
				Type:     "GROUP_REQUEST",
				Content:  userId,
			}
			group.RequestPending, err = handler.Repos.NotifRepo.CheckIfExists(notification)
			if err != nil {
				utils.RespondWithError(w, "Error on getting data", 200)
				return
			}
		}

	}
	utils.RespondWithGroups(w, []models.Group{group}, 200)
}

// returns list of all group members and administrator
func (handler *Handler) GroupMembers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { // check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	members, err := handler.Repos.GroupRepo.GetGroupMembers(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, members, 200)
}

// returns list of events for group
// GET request with group_id included in query
// in case if user is not a member of the group, returns error
func (handler *Handler) GroupEvents(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { // check if group id exists in request
		utils.RespondWithError(w, "Error on reading group ID", 200)
		return
	}
	// check if current user is a member or admin  of the group
	isMember := false
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on reading role", 200)
		return
	}
	if !isAdmin {
		isMember, err = handler.Repos.GroupRepo.IsGroupMember(groupId, userId)
		if err != nil {
			utils.RespondWithError(w, "Error on checking if is group member", 200)
			return
		}
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	// current user is a member or admin -> get events
	events, err := handler.Repos.EventRepo.GetAll(groupId)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, "Error on getting event data", 200)
		return
	}
	/* ----------------------- attach author to each event ---------------------- */
	for i := 0; i < len(events); i++ {
		events[i].Author, _ = handler.Repos.UserRepo.GetDataMin(events[i].AuthorID)
	}
	/* -------------------- attach participation to each event ------------------- */
	for i := 0; i < len(events); i++ {
		going, _ := handler.Repos.EventRepo.IsParticipating(events[i].ID, userId)
		if going {
			events[i].Going = "YES"
		} else {
			events[i].Going = "NO"
		}
	}
	utils.RespondWithEvents(w, events, 200)
}

// returns all posts that belongs to group
func (handler *Handler) GroupPosts(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { // check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(groupId, userId)
	isMember, err := handler.Repos.GroupRepo.IsGroupMember(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	/* ------------- current user is a member or admin -> get posts ------------- */
	posts, err := handler.Repos.PostRepo.GetGroupPosts(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get post author info attached
	if err = AttachAuthors(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get comment info for each post
	if err = AttachComments(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithPosts(w, posts, 200)
}

// returns pending requests to join to group, only for admin
// for others respond with error
func (handler *Handler) GroupRequests(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { // check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* ------------------------- check if user is admin ------------------------- */
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isAdmin {
		utils.RespondWithError(w, "Unauthorized", 200)
		return
	}
	/* ---------------------- get pending requests form db ---------------------- */
	notifications, err := handler.Repos.NotifRepo.GetGroupRequests(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	for i := 0; i < len(notifications); i++ {
		notifications[i].User, err = handler.Repos.UserRepo.GetDataMin(notifications[i].Content)
		if err != nil {
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
	}
	utils.RespondWithNotifications(w, notifications, 200)
}

func (handler *Handler) CancelGroupRequests(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { // check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// delete notification
	notif := models.Notification{
		Type:     "GROUP_REQUEST",
		TargetID: groupId,
		Content:  currentUserId,
	}
	if err := handler.Repos.NotifRepo.DeleteByType(notif); err != nil {
		utils.RespondWithError(w, "Error on canceling request", 200)
		return
	}
	utils.RespondWithSuccess(w, "gROUP request canceled successfuly", 200)
}

/* -------------------------------------------------------------------------- */
/*                                save new data                               */
/* -------------------------------------------------------------------------- */

func (handler *Handler) NewGroup(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new Group
	var newGroup models.Group
	err := json.NewDecoder(r.Body).Decode(&newGroup)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ------------------------- attach addidtional data ------------------------ */
	// access user id
	newGroup.AdminID = r.Context().Value(utils.UserKey).(string)
	newGroup.ID = utils.UniqueId()

	/* ------------------------------- save in db ------------------------------- */
	err = handler.Repos.GroupRepo.NewGroup(newGroup)
	if err != nil {
		utils.RespondWithError(w, "Error on saving group", 200)
		return
	}
	/* ------------------------- save invitations in db ------------------------- */
	for i := 0; i < len(newGroup.Invitations); i++ {
		// save each follower in db
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: newGroup.Invitations[i],
			Type:     "GROUP_INVITE",
			Content:  newGroup.ID,
			Sender:   newGroup.AdminID,
		}
		err = handler.Repos.NotifRepo.Save(newNotif)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

		for client := range wsServer.Clients {
			if client.ID == newNotif.TargetID {
				client.SendNotification(newNotif)
			}
		}

		// save as a new member of group
		// if err = handler.Repos.GroupRepo.SaveGroupMember(newGroup.Invitations[i], newGroup.ID); err != nil {
		// 	utils.RespondWithError(w, "Internal server error", 200)
		// 	return
		// }
	}
	newGroup.Administrator = true
	// NOTIFY WEBSOCKET ABOUT NEW NOTIFICATION
	utils.RespondWithGroups(w, []models.Group{newGroup}, 200)
}

// NOT TESTED
func (handler *Handler) NewGroupPost(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	err := r.ParseMultipartForm(3145728) // 3MB
	if err != nil {
		utils.RespondWithError(w, "Error in form validation", 200)
		return
	}
	userId := r.Context().Value(utils.UserKey).(string)
	/* ------------------------ create new post instance ------------------------ */
	newPost := models.Post{
		ID:       utils.UniqueId(),
		Content:  r.PostFormValue("body"),
		GroupID:  r.PostFormValue("groupId"),
		AuthorID: userId,
	}
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(newPost.GroupID, userId)
	isMember, err := handler.Repos.GroupRepo.IsGroupMember(newPost.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	/* ------------------------ save image in filesystem ------------------------ */
	newPost.ImagePath = utils.SaveImage(r)
	/* -------------------------- save post in database ------------------------- */
	errDB := handler.Repos.PostRepo.New(newPost)
	if errDB != nil {
		utils.RespondWithError(w, "Error on saving post", 200)
		return
	}
	utils.RespondWithSuccess(w, "New post created", 200)
}

// handle when new user wants to join the group
func (handler *Handler) NewGroupRequest(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { // check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* --------------------------- check if not admin --------------------------- */
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if isAdmin {
		utils.RespondWithError(w, "Invalid request", 200)
		return
	}
	/* -------------------- create new notification instance -------------------- */
	notification := models.Notification{
		ID:       utils.UniqueId(),
		TargetID: groupId,
		Type:     "GROUP_REQUEST",
		Content:  userId,
		Sender:   userId,
	}
	/* --------------------- check if request already exists -------------------- */
	exists, err := handler.Repos.NotifRepo.CheckIfExists(notification)
	if err != nil {
		utils.RespondWithError(w, "Error on saving request", 200)
		return
	}
	if exists {
		utils.RespondWithSuccess(w, "Request already saved", 200)
		return
	}
	/* ------------------------- save notification in db ------------------------ */
	err = handler.Repos.NotifRepo.Save(notification)
	if err != nil {
		utils.RespondWithError(w, "Error on saving request", 200)
	}

	// SEND MESSAGE TO TO GROUP ADMIN IF IS ONLINE
	admin, err := handler.Repos.GroupRepo.GetGroupAdmin(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on finding group admin", 200)
		return
	}
	for client := range wsServer.Clients {
		if client.ID == admin {
			client.SendNotification(notification)
		}
	}
	utils.RespondWithSuccess(w, "Request saved successfuly", 200)
}

// NOT TESTED
// handle response from group administrator for requests to join group
// waits for requestId and response -accept/decline
func (handler *Handler) ResponseGroupRequest(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on data submittion", 200)
		return
	}
	/* --------------------------- read incoming data --------------------------- */
	type Response struct {
		GroupID   string `json:"groupId"`
		RequestID string `json:"requestId"`
		Response  string `json:"response"`
	}
	var response Response
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------- check if all fields provided ---------------------- */
	if response.Response != "accept" && response.Response != "decline" {
		utils.RespondWithError(w, "Response unvalid", 200)
		return
	}
	if response.GroupID == "" || response.RequestID == "" {
		utils.RespondWithError(w, "Response incomplete", 200)
		return
	}
	/* ------------------- chack if curren user is group admin ------------------ */
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(response.GroupID, userId)
	if !isAdmin {
		utils.RespondWithError(w, "Unauthorized", 200)
		return
	}
	/* ----------------------------- handle response ---------------------------- */
	// if accepted -> save as new member
	if response.Response == "accept" {
		// get id of member that requests to join
		joinerId, err := handler.Repos.NotifRepo.GetUserFromRequest(response.RequestID)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
		// save as a new member of group
		if err = handler.Repos.GroupRepo.SaveGroupMember(joinerId, response.GroupID); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
		// if joiner online, send updated group status
		for client := range wsServer.Clients {
			if client.ID == joinerId {
				client.SendGroupRequestAccept(response.GroupID)
			}
		}
	}
	// delete from pending notification table
	if err = handler.Repos.NotifRepo.Delete(response.RequestID); err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	utils.RespondWithSuccess(w, "Response was successful", 200)
}

// NOT TESTED
// waits for POST request with groupId and Invitation list included
func (handler *Handler) NewGroupInvite(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new Group
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	userId := r.Context().Value(utils.UserKey).(string)
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.Repos.GroupRepo.IsGroupAdmin(group.ID, userId)
	isMember, err := handler.Repos.GroupRepo.IsGroupMember(group.ID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	for i := 0; i < len(group.Invitations); i++ {
		// save each invitation in db
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: group.Invitations[i],
			Type:     "GROUP_INVITE",
			Content:  group.ID,
			Sender:   userId,
		}
		err = handler.Repos.NotifRepo.Save(newNotif)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

		// search if user has open ws connection and send notification
		for client := range wsServer.Clients {
			if client.ID == group.Invitations[i] {
				client.SendNotification(newNotif)
			}
		}
	}
	utils.RespondWithSuccess(w, "Invitations saved", 200)
}

// NOT TESTED
func (handler *Handler) ResponseInviteRequest(w http.ResponseWriter, r *http.Request) {
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
	// get group id from notification
	groupId, err := handler.Repos.NotifRepo.GetGroupId(resp.RequestID)
	userId := r.Context().Value(utils.UserKey).(string)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	if strings.ToUpper(resp.Response) == "ACCEPT" {
		err = handler.Repos.GroupRepo.SaveGroupMember(userId, groupId)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}
	/* ----------------------- delete pending notification ---------------------- */
	err = handler.Repos.NotifRepo.Delete(resp.RequestID)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	// notify websocket about notification changes
	utils.RespondWithSuccess(w, "Response successful", 200)
}
