package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

// Posts
// After a user is logged in he/she can create posts and comments on already created posts. While creating a post or a comment, the user can include an image or GIF.

// The user must be able to specify the privacy of the post:

// public (all users in the social network will be able to see the post) [no condition to fetch]
// almost private (only followers of the creator of the post will be able to see the post)
// private (only the followers chosen by the creator of the post will be able to see it)

func (rt *Root) SetupPostRoutes(mux *http.ServeMux) {
	postMux := http.NewServeMux()

	postMux.HandleFunc("POST /new", rt.NewPost)
	postMux.HandleFunc("POST /feed", rt.GetFeedPosts)
	postMux.HandleFunc("GET /followers", rt.GetFollowers)

	log.Println("Mounting post multiplexer at /post/")

	mux.Handle("/post/", http.StripPrefix("/post", postMux))
}

func (app *Root) GetFollowers(w http.ResponseWriter, r *http.Request) {
	// userId := r.PathValue("user_id")
	userId := app.DL.GetRequesterID(w, r)
	// id, err := strconv.Atoi(userId)
	if userId <= 0 {
		log.Println("Error: user_id not provided or invalid")
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "user_id is required",
		})
		return
	}

	followers, err := app.DL.Follows.GetFollows(userId, "followers")
	fmt.Println("GetFollowers userId:", userId, "Followers:", followers)
	if err != nil {
		log.Printf("Error ff followers for user %d: %v", userId, err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal Error!",
		})
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, followers); err != nil {
		log.Printf("Error encoding response: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal Error!",
		})
	}
}

func (app *Root) GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	var filter *models.PostFilter

	if err := tools.DecodeJSON(r, &filter); err != nil {
		log.Printf("Error decoding filter: %v", err)
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// // Validate filter using tools function
	// if status := tools.ValidatePostFilter(filter); status != 200 {
	// 	tools.EncodeJson(w, http.StatusBadRequest, map[string]string{
	// 		"error": "Invalid filter parameters",
	// 	})
	// 	return
	// }

	posts, err := app.DL.Posts.GetPosts(filter)
	fmt.Println("GetFeedPosts filter:", filter, "Posts:", posts)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch posts",
		})
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, posts); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (app *Root) NewPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	var hasFile bool
	if r.Method != http.MethodPost {
		tools.EncodeJSON(w, 403, nil)
		return
	}
	contentType := r.Header.Get("Content-Type")

	// Handle different content types
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse multipart form data
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
			log.Printf("Error parsing multipart form: %v", err)
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid form data",
			})
			return
		}
		hasFile = true

		// Parse post data from form values using tools function
		if status := models.ParsePostFromForm(r, &post); status != 200 {
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid form data",
			})
			return
		}
	} else if strings.Contains(contentType, "application/json") {
		// Handle JSON request
		if err := tools.DecodeJSON(r, &post); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid JSON format",
			})
			return
		}
	} else {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Content-Type must be application/json or multipart/form-data",
		})
		return
	}

	// Validate post using existing tools function
	if status := models.ValidatePost(&post); status != 200 {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid post data",
		})
		return
	}

	// Handle file upload only if it's multipart form data
	var imagePath string
	if hasFile {
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			// Validate file size
			if handler.Size > 5<<20 { // 5 MB limit
				tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "Image file too large (max 5MB)",
				})
				return
			}

			// Validate file type using existing tools function
			if !tools.IsAllowedFile(handler.Filename, file) {
				tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "file format is not supported",
				})
				return
			}

			// Upload file using existing tools function
			uploadedPath, status := tools.UploadHandler(file, handler)
			if status != 200 {
				tools.EncodeJSON(w, status, map[string]string{
					"error": "Failed to upload image",
				})
				return
			}
			imagePath = uploadedPath
		} else if err != http.ErrMissingFile {
			log.Printf("Error handling file upload: %v", err)
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
				// "error": "Error processing image upload",
			})
			return
		}
	}

	// Begin transaction
	tx, err := app.DL.Posts.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
		return
	}
	defer tx.Rollback()
	fmt.Println(post.OwnerId, post.GroupId, len(post.Content), imagePath, post.Privacy)
	// Insert post
	result, err := tx.Exec(`
		INSERT INTO posts (user_id, group_id, content, image, privacy, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		post.OwnerId, post.GroupId, post.Content, imagePath, post.Privacy, time.Now())

	if err != nil {
		log.Printf("Error inserting post: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create post",
		})
		return
	}

	// Get the inserted post ID
	postId, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting post ID: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create post",
		})
		return
	}

	// Handle private post permissions
	if post.Privacy == "private" && len(post.ChosenUsersIds) > 0 {
		stmt, err := tx.Prepare("INSERT INTO post_privacy (chosen_id, post_id) VALUES (?, ?)")
		if err != nil {
			log.Printf("Error preparing privacy statement: %v", err)
			tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to set post privacy",
			})
			return
		}
		defer stmt.Close()

		for _, userId := range post.ChosenUsersIds {
			if _, err := stmt.Exec(userId, postId); err != nil {
				log.Printf("Error inserting post privacy: %v", err)
				tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Failed to set post privacy",
				})
				return
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create post",
		})
		return
	}

	tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Post created successfully",
		"post_id": postId,
	})
}

func (app *Root) SetupCommentRoutes(mux *http.ServeMux) {
	commentMux := http.NewServeMux()

	commentMux.HandleFunc("POST /{post_id}/new", app.NewComment)
	commentMux.HandleFunc("GET /{post_id}/comment", app.GetFeedComments)
	// commentMux.HandleFunc("POST /new/upload", UploadHandler) // TODO need to handel image

	log.Println("Mounting post multiplexer at /comments/")

	mux.Handle("/comments/", http.StripPrefix("/comments", commentMux))
}

func (app *Root) GetFeedComments(w http.ResponseWriter, r *http.Request) {
	// TODO need to specify the methode
	fmt.Println("GetFeedComments path accessed:", r.URL.Path, "Method:", r.Method)
	if r.Method != http.MethodGet {
		tools.EncodeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}
	post := r.PathValue("post_id")
	if post == "" {
		log.Println("Error: post_id is required")
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "post_id is required",
		})
		return
	}
	// var err error
	post_id, er := strconv.Atoi(post)
	if er != nil {
		log.Printf("Error converting post_id to integer: %v", er)
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid post_id format",
		})
		return
	}

	// comments := []models.Comment{}
	comments, err := app.DL.Comments.GetComments(post_id)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch comments",
		})
		return
	}
	if comments == nil {
		log.Println("No comments found for post_id:", post_id)
		tools.EncodeJSON(w, 404, nil)

		return
	}
	tools.EncodeJSON(w, http.StatusOK, comments)
	// fmt.Fprintln(w, "Listing all posts")
}

func (app *Root) NewComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	var hasFile bool
	if r.Method != http.MethodPost {
		tools.EncodeJSON(w, 403, nil)
		return
	}
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
			log.Printf("Error parsing multipart form: %v", err)
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid form data",
			})
			return
		}
		hasFile = true

		// Parse post data from form values using tools function
		if status := models.ParseCommentFromForm(r, &comment); status != 200 {
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid form data",
			})
			return
		}
	} else if strings.Contains(contentType, "application/json") {
		if err := tools.DecodeJSON(r, &comment); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid JSON format",
			})
			return
		}
	} else {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Content-Type must be application/json or multipart/form-data",
		})
		return
	}
	fmt.Println("mok", comment)

	if ok, status := models.ValidateComment(app.DL.Comments.DB, &comment); !ok {
		tools.EncodeJSON(w, status, map[string]string{
			"error": "Invalid post data",
		})
		return
	}
	fmt.Println("mok", comment)

	var imagePath string
	if hasFile {
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			// Validate file size
			if handler.Size > 5<<20 { // 5 MB limit
				tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "Image file too large (max 5MB)",
				})
				return
			}

			// Validate file type using existing tools function
			if !tools.IsAllowedFile(handler.Filename, file) {
				tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "file format is not supported",
				})
				return
			}

			// Upload file using existing tools function
			uploadedPath, status := tools.UploadHandler(file, handler)
			if status != 200 {
				tools.EncodeJSON(w, status, map[string]string{
					"error": "Failed to upload image",
				})
				return
			}
			imagePath = uploadedPath
		} else if err != http.ErrMissingFile {
			log.Printf("Error handling file upload: %v", err)
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
				// "error": "Error processing image upload",
			})
			return
		}
	}

	// Begin transaction
	tx, err := app.DL.Posts.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
		return
	}
	defer tx.Rollback()
	// fmt.Println(post.OwnerId, post.GroupId, len(post.Content), imagePath, post.Privacy)
	// Insert post
	_, err = tx.Exec(`
		INSERT INTO comments (user_id, post_id, content, image, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		comment.OwnerId, comment.Post_id, comment.Content, imagePath, time.Now())

	if err != nil {
		log.Printf("Error inserting post: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create post",
		})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create post",
		})
		return
	}

	tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "comment created successfully",
	})
}

// func (app *SocialApp) NewComment(w http.ResponseWriter, r *http.Request) {
// 	log.Printf("New comment path accessed: %s", r.URL.Path)
// 	// if r.URL.Path != "/new" {
// 	// 	tools.EncodeJSON(w, 500, nil)
// 	// 	return
// 	// }
// 	var comment models.Comment
// 	if err := tools.DecodeJSON(r, &comment); err != nil {
// 		tools.EncodeJSON(w, 500, nil)
// 		return
// 	}
// 	stmt, err := app.Posts.DB.Prepare(`
//     INSERT INTO comments (user_id, post_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)`)
// 	if err != nil {
// 		tools.EncodeJson(w, 500, nil)
// 		return
// 	}
// 	defer stmt.Close()

// 	file, handler, err := r.FormFile("image")
// 	if err != nil {
// 		tools.EncodeJson(w, 500, nil)
// 		return
// 	}

// 	defer file.Close()

// 	temp, status := "", 0
// 	if handler.Filename != "" {
// 		temp, status = tools.UploadHandler(file, handler)
// 		if status != 200 {
// 			tools.EncodeJson(w, status, nil)
// 			return
// 		}
// 	}

// 	if _, err = stmt.Exec(comment.OwnerId, comment.Post_id, comment.Content, temp, time.Now()); err != nil {
// 		tools.EncodeJson(w, 200, nil)
// 		return
// 	}

// 	tools.EncodeJson(w, 200, nil)
// }
