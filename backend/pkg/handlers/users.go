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

type Handler struct {
	Repos *models.Repositories
}

/* -------------------------------------------------------------------------- */
/*                                    users                                   */
/* -------------------------------------------------------------------------- */
// Find all users and they relation with current user
func (handler *Handler) AllUsers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all users exccept current + relations
	users, errUsers := handler.Repos.UserRepo.GetAllAndFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, users, 200)
}

// Returns user nickname, id and path to avatar
func (handler *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	user, err := handler.Repos.UserRepo.GetDataMin(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, []models.User{user}, 200)
}

// Returns user data based on public / private profile and user_id from request
// waits for GET request with query "userId" ->user client is looking for
//
//	can be used both on own profile and other users
func (handler *Handler) UserData(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get user id from request
	query := r.URL.Query()
	userId := query.Get("userId")
	// get if profile public or private
	status, err := handler.Repos.UserRepo.ProfileStatus(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// check if client looking for own profile
	currentUser := (currentUserId == userId)
	var following bool
	if !currentUser {
		// check if current user following user he is looking for
		following, err = handler.Repos.UserRepo.IsFollowing(userId, currentUserId)
		if err != nil {
			fmt.Println("Throwing an err!", err)
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
	}
	// request user info based on following/ and profile status
	// if public or current user or if following  get large data set
	// if private and not following => get small data set
	var user models.User
	if currentUser || following || status == "PUBLIC" { // get full data set
		user, err = handler.Repos.UserRepo.GetProfileMax(userId)
	} else {
		user, _ = handler.Repos.UserRepo.GetProfileMin(userId)
		/* -------------------- check if follow status is pending ------------------- */
		notif := models.Notification{Type: "FOLLOW", Content: currentUserId, TargetID: userId}
		user.FollowRequestPending, err = handler.Repos.NotifRepo.CheckIfExists(notif)
	}
	if err != nil {
		fmt.Println("Throwing an err!12", err)
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// tie together stats to user object
	user.Following = following
	user.CurrentUser = currentUser
	user.Status = status

	utils.RespondWithUsers(w, []models.User{user}, 200)
}

// changes user status in db return status
// in case of turning to PUBLIC -> also accept follow requests
func (handler *Handler) UserStatus(w http.ResponseWriter, r *http.Request) {
	statusList := []string{"PUBLIC", "PRIVATE"} // possible status
	var client models.User

	w = utils.ConfigHeader(w)
	// access user id
	client.ID = r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqStatus := strings.ToUpper(query.Get("status"))

	// check if valid value and asign to user
	if reqStatus == statusList[0] {
		client.Status = statusList[0]
	} else if reqStatus == statusList[1] {
		client.Status = statusList[1]
	} else {
		utils.RespondWithError(w, "Requested status not valid", 200)
		return
	}
	// request current status from db
	currentStatus, err := handler.Repos.UserRepo.GetStatus(client.ID)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// check if requested status is not the same as current
	if currentStatus == client.Status {
		utils.RespondWithError(w, "Status change not valid", 200)
		return
	}
	// Set new status
	err = handler.Repos.UserRepo.SetStatus(client)
	if err != nil {
		utils.RespondWithError(w, "Error on saving status", 200)
		return
	}
	// if new status is public -> also accept pending follow requests
	// responds with success and newly created status
	utils.RespondWithSuccess(w, client.Status, 200)
}

/* -------------------------------------------------------------------------- */
/*                                  followers                                 */
/* -------------------------------------------------------------------------- */
// Find all followers
func (handler *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get userId from request
	query := r.URL.Query()
	userId := query.Get("userId")
	// request all  following users
	followers, errUsers := handler.Repos.UserRepo.GetFollowers(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, followers, 200)
}

// Find all who client are following
func (handler *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get userId from request
	query := r.URL.Query()
	userId := query.Get("userId")
	// request all  following users
	followers, errUsers := handler.Repos.UserRepo.GetFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, followers, 200)
}

func (handler *Handler) Follow(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqUserId := query.Get("userId")
	/* ----------------- safety check -> if request already made ---------------- */
	alreadyFollowing, _ := handler.Repos.UserRepo.IsFollowing(reqUserId, currentUserId)
	if alreadyFollowing {
		utils.RespondWithError(w, "User already is following", 200)
		return
	}
	// get target user profile status -> public or private
	reqUserStatus, err := handler.Repos.UserRepo.GetStatus(reqUserId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if reqUserStatus == "PUBLIC" {
		// SAVE AS FOLLOWER
		err := handler.Repos.UserRepo.SaveFollower(reqUserId, currentUserId)
		if err != nil {
			utils.RespondWithError(w, "Error on saving follower", 200)
			return
		}
	} else if reqUserStatus == "PRIVATE" {
		// SAVE IN NOTIFICATIONS as pending folllow request
		notification := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: reqUserId,
			Type:     "FOLLOW",
			Content:  currentUserId,
			Sender:   currentUserId,
		}
		err := handler.Repos.NotifRepo.Save(notification)
		if err != nil {
			utils.RespondWithError(w, "Error on save", 200)
			return
		}
		// if user online send notification about follow request
		for client := range wsServer.Clients {
			if client.ID == reqUserId {
				client.SendNotification(notification)
			}
		}

	}
	utils.RespondWithSuccess(w, "Following successful", 200)
}

func (handler *Handler) CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqUserId := query.Get("userId")
	// delete notification corresponding to follow request
	notif := models.Notification{
		Type:     "FOLLOW",
		TargetID: reqUserId,
		Content:  currentUserId,
	}
	if err := handler.Repos.NotifRepo.DeleteByType(notif); err != nil {
		utils.RespondWithError(w, "Error on canceling request", 200)
		return
	}
	utils.RespondWithSuccess(w, "Follow request canceled successfuly", 200)
}

func (handler *Handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get status from request
	query := r.URL.Query()
	reqUserId := query.Get("userId")

	if err := handler.Repos.UserRepo.DeleteFollower(reqUserId, currentUserId); err != nil {
		utils.RespondWithError(w, "Error on deleting follower", 200)
		return
	}
	utils.RespondWithSuccess(w, "Unfollowing successful", 200)
}

// not tested
// wait for POST request with notification Id and response -"ACCEPT" or "DECLINE"
func (handler *Handler) ResponseFollowRequest(w http.ResponseWriter, r *http.Request) {
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
	// get other user id from notification
	followerId, err := handler.Repos.NotifRepo.GetUserFromRequest(resp.RequestID)
	userId := r.Context().Value(utils.UserKey).(string)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	if strings.ToUpper(resp.Response) == "ACCEPT" {
		err = handler.Repos.UserRepo.SaveFollower(userId, followerId)
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

/* -------------------------------------------------------------------------- */
/*                                  chat List                                 */
/* -------------------------------------------------------------------------- */
func (handler *Handler) ChatList(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get userId from request
	query := r.URL.Query()
	userId := query.Get("userId")
	// request all  following users
	followers, errUsers := handler.Repos.UserRepo.GetFollowing(userId)
	if errUsers != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// get users_ids that have a chat history
	ids, errIds := handler.Repos.MsgRepo.GetChatHistoryIds(userId)
	if errIds != nil {
		utils.RespondWithError(w, "Error on getting chat history", 200)
		return
	}
	// loop over chat history ids
	// compare with followers
	// if not found in folllowers, fetch user data and add to the list
	for currentId := range ids {
		isPresent := ContainsUser(followers, currentId)
		if !isPresent {
			user, err := handler.Repos.UserRepo.GetDataMin(currentId)
			if err != nil {
				utils.RespondWithError(w, "Error on getting chat history data", 200)
				return
			}
			followers = append(followers, user)
		}
	}

	utils.RespondWithUsers(w, followers, 200)
}

/* --------------------------------- helper --------------------------------- */
func ContainsUser(list []models.User, id string) bool {
	for _, value := range list {
		if value.ID == id {
			return true
		}
	}
	return false
}
