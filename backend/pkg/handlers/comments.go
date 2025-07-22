package handlers

import (
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/utils"
)

func (handler *Handler) NewComment(w http.ResponseWriter, r *http.Request) {
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
	// create new comment instance
	newComment := models.Comment{
		ID:       utils.UniqueId(),
		PostID:   r.PostFormValue("postid"),
		Content:  r.PostFormValue("body"),
		AuthorID: userId,
	}
	// save image in filesystem
	newComment.ImagePath = utils.SaveImage(r)
	// save comment in database
	errDB := handler.Repos.CommentRepo.New(newComment)
	if errDB != nil {
		utils.RespondWithError(w, "Error on saving data", 200)
		return
	}
	utils.RespondWithSuccess(w, "New comment created", 200)
}
