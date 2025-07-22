package handlers

import (
	"net/http"
	"strings"

	"social-network/pkg/models"
	"social-network/pkg/utils"
)

/* ------------------------ fetch all posts for user ------------------------ */
func (handler *Handler) AllPosts(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all posts
	posts, errPosts := handler.Repos.PostRepo.GetAll(userId)
	if errPosts != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get post author info attached
	if err := AttachAuthors(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get comment info for each post
	if err := AttachComments(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithPosts(w, posts, 200)
}

func (handler *Handler) UserPosts(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	currentUserId := r.Context().Value(utils.UserKey).(string)

	// get user id from request
	query := r.URL.Query()
	userId := query.Get("id")
	if userId == "" {
		utils.RespondWithError(w, "Error user id", 200)
		return
	}
	// request user posts
	posts, errPosts := handler.Repos.PostRepo.GetUserPosts(userId, currentUserId)
	if errPosts != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get post author info attached
	if err := AttachAuthors(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get comment info for each post
	if err := AttachComments(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithPosts(w, posts, 200)
}

/* ----------------------------- create new post ---------------------------- */
func (handler *Handler) NewPost(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	err := r.ParseMultipartForm(3145728) // 3MB
	if err != nil {
		utils.RespondWithError(w, "Error in form validation", 200)
		return
	}
	// configure data
	userId := r.Context().Value(utils.UserKey).(string)
	visibility := strings.ToUpper(r.PostFormValue("privacy"))
	visibility = strings.Replace(visibility, "-", "_", -1)
	// create new post instance
	newPost := models.Post{
		ID:         utils.UniqueId(),
		Content:    r.PostFormValue("body"),
		GroupID:    r.PostFormValue("groupId"),
		Visibility: visibility,
		AuthorID:   userId,
	}
	// save image in filesystem
	newPost.ImagePath = utils.SaveImage(r)
	// save post in database
	errDB := handler.Repos.PostRepo.New(newPost)
	if errDB != nil {
		utils.RespondWithError(w, "Error in form validation", 200)
		return
	}
	// in case of "almost private post", save users with access
	if newPost.Visibility == "ALMOST_PRIVATE" {
		accessListRaw := r.PostFormValue("checkedfollowers")
		accessList := strings.Split(accessListRaw, ",")
		for i := 0; i < len(accessList); i++ {
			// save each follower in db
			err = handler.Repos.PostRepo.SaveAccess(newPost.ID, accessList[i])
			if errDB != nil {
				utils.RespondWithError(w, "Internal server error", 200)
				return
			}
		}

	}
	utils.RespondWithSuccess(w, "New post created", 200)
}

/* -------------------------------------------------------------------------- */
/*                                   helpers                                  */
/* -------------------------------------------------------------------------- */

func AttachAuthors(handler *Handler, posts *[]models.Post) error {
	for i := 0; i < len(*posts); i++ {
		userId := (*posts)[i].AuthorID
		author, err := handler.Repos.UserRepo.GetDataMin(userId)
		if err != nil {
			return err
		}
		(*posts)[i].Author = author
	}
	return nil
}

func AttachComments(handler *Handler, posts *[]models.Post) error {
	for i := 0; i < len(*posts); i++ {
		postId := (*posts)[i].ID
		comments, err := handler.Repos.CommentRepo.Get(postId)
		if err != nil {
			return err
		}
		// add author
		for i := 0; i < len(comments); i++ {
			author, err := handler.Repos.UserRepo.GetDataMin(comments[i].AuthorID)
			if err != nil {
				return err
			}
			comments[i].Author = author
		}
		(*posts)[i].Comments = comments
	}
	return nil
}
