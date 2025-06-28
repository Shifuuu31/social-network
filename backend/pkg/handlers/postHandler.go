package handlers

import (
	"fmt"
	"net/http"

	"social-network/pkg/models"

	"social-network/pkg/tools"
)

func (rt *Root) NewPostHandler() (postsMux *http.ServeMux) {
	postsMux = http.NewServeMux()

	postsMux.HandleFunc("POST /addpost", rt.NewPost)
	// postsMux.HandleFunc("GET /posts", rt.InviteToJoinGroup)

	return postsMux
}

func (rt *Root) NewPost(w http.ResponseWriter, r *http.Request) {
	userID := rt.DL.GetRequesterID(w, r)

	var post models.Post

	if err := tools.DecodeJSON(r, &post); err != nil {
	}

	post.UserID = userID
	fmt.Println(&post)

	if err := models.ValidatePost(&post); err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, nil)
	}

	if err := rt.DL.Posts.InsertPost(post); err != nil {
		tools.EncodeJSON(w, http.StatusInternalServerError, nil)
	}

	tools.EncodeJSON(w, http.StatusOK, nil)
}
