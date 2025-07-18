package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"os"
	"path/filepath"
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
	fmt.Println("aalo")
	postMux := http.NewServeMux()
	postMux.HandleFunc("POST /feed", rt.GetFeedPosts)
	postMux.HandleFunc("POST /new", rt.NewPost)
	postMux.HandleFunc("GET /followers", rt.GetFollowers) // it already in profile&follow.go

	// Add comment routes
	postMux.HandleFunc("GET /{post_id}/comments", rt.GetFeedComments)
	postMux.HandleFunc("POST /{post_id}/comments/new", rt.NewComment)

	// log.Println("Mounting post multiplexer at /post/")
	mux.Handle("/post/", http.StripPrefix("/post", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		postMux.ServeHTTP(w, r)
	})))
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

	currentUserId := app.DL.GetRequesterID(w, r)

	if filter.Type != "user" {
		filter.Id = currentUserId
	}

	// Set the filter ID to the current user's ID
 
	// filter.Type = "followers"
	// Validate filter using tools function
	fmt.Println("filter:", filter)
	if err := filter.Validate(); err != nil {
		fmt.Println("GetFeedPosts filter validation error:", err, "filter:", filter)
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid filter parameters",
		})
		return
	}

	posts, err := app.DL.Posts.GetPosts(filter)
	// fmt.Println("waaaaaaaaaaaaaaaaaaaaaaaazi", posts)
	// fmt.Println("GetFeedPosts filter:", filter, "Posts:", posts)

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
	contentType := r.Header.Get("Content-Type")

	fmt.Println("asasas", post.Content, "@@@@")
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
			log.Printf("ParsePostFromForm failed with status: %d", status)
			log.Printf("Form values: owner_id=%s, content=%s, privacy=%s",
				r.FormValue("owner_id"), r.FormValue("content"), r.FormValue("privacy"))
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid form data",
			})
			return
		}

		// Debug: Log the parsed post data
		log.Printf("Parsed post: OwnerId=%d, Content='%s', Privacy='%s'",
			post.OwnerId, post.Content, post.Privacy)
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
	if err := post.Validate(); err != nil {
		log.Printf("Post validation failed: %v", err)
		log.Printf("Post data: OwnerId=%d, Content='%s', Privacy='%s', GroupId=%d",
			post.OwnerId, post.Content, post.Privacy, post.GroupId)
		tools.RespondError(w, "Invalid post data", http.StatusBadRequest)
		return
	}
	// fmt.Println(post.Content, post.Image_url, "xccxc")
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
	// imagePath="asd/sad/sd"
	if post.Privacy == "almost_private" {
		post.Privacy = "followers"
	}
	fmt.Println("gggggggggggggggggggergregergergerggg", post.OwnerId, post.GroupId, post.Content, imagePath, post.Privacy)
	// Insert post

	result, err := tx.Exec(`
		INSERT INTO posts (user_id, group_id, content, image_url,privacy,created_at)
		VALUES (?,?,?,?,?, ?)`,
		post.OwnerId, post.GroupId, post.Content, imagePath, post.Privacy, time.Now())
	if err != nil {
		log.Printf("Error inserting post cause: %v", err)
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

func (app *Root) GetFeedComments(w http.ResponseWriter, r *http.Request) {
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

	post_id, err := strconv.Atoi(post)
	if err != nil {
		log.Printf("Error converting post_id to integer: %v", err)
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid post_id format",
		})
		return
	}

	comments, err := app.DL.Comments.GetComments(post_id)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch comments",
		})
		return
	}

	// Return empty array instead of 404 for no comments
	if comments == nil {
		log.Println("No comments found for post_id:", post_id)
		tools.EncodeJSON(w, http.StatusOK, []models.Comment{})
		return
	}

	tools.EncodeJSON(w, http.StatusOK, comments)
}

// ServeImage serves uploaded images
func (app *Root) ServeImage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ServeImage called with path: %s\n", r.URL.Path)

	filename := r.PathValue("filename")
	if filename == "" {
		fmt.Printf("No filename provided\n")
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Filename is required",
		})
		return
	}

	fmt.Printf("Serving image: %s\n", filename)

	// Construct the full path to the image
	imagePath := filepath.Join("uploads", filename)
	fmt.Printf("Full image path: %s\n", imagePath)

	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Printf("Image file not found: %s\n", imagePath)
		tools.EncodeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Image not found",
		})
		return
	}

	fmt.Printf("Image file found, serving: %s\n", imagePath)

	// Set proper headers for image serving
	w.Header().Set("Content-Type", "image/jpeg")                // Default to JPEG, could be made dynamic
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	// Serve the file
	http.ServeFile(w, r, imagePath)
}

// ServeImageHandler handles image serving with traditional URL pattern
func (app *Root) ServeImageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	// URL will be like /images/filename.jpg
	path := r.URL.Path
	filename := path[len("/images/"):] // Remove "/images/" prefix

	if filename == "" {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Filename is required",
		})
		return
	}

	fmt.Printf("ServeImageHandler called with filename: %s\n", filename)

	// Try to find the image in different subdirectories
	possiblePaths := []string{
		filepath.Join("uploads", "profiles", filename),
		filepath.Join("uploads", "posts", filename),
		filepath.Join("uploads", "comments", filename),
		filepath.Join("uploads", filename), // Fallback to old format
	}

	var imagePath string
	var found bool

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			imagePath = path
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Image file not found in any directory: %s\n", filename)
		tools.EncodeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Image not found",
		})
		return
	}

	fmt.Printf("Image file found, serving: %s\n", imagePath)

	// Determine MIME type based on file extension
	ext := strings.ToLower(filepath.Ext(filename))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	default:
		mimeType = "image/jpeg" // Default fallback
	}

	// Set proper headers for image serving
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	// Serve the file
	http.ServeFile(w, r, imagePath)
}

func (app *Root) NewComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	var hasFile bool

	fmt.Println("NewComment handler accessed")
	if r.Method != http.MethodPost {
		tools.EncodeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	// CRITICAL FIX: Extract post_id from URL path
	post := r.PathValue("post_id")
	fmt.Println("post id is:", post)
	if post == "" {
		log.Println("Error: post_id is required")
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "post_id is required",
		})
		return
	}

	post_id, err := strconv.Atoi(post)
	if err != nil {
		log.Printf("Error converting post_id to integer: %v", err)
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid post_id format",
		})
		return
	}

	// Set the post_id in the comment struct
	comment.Post_id = post_id

	contentType := r.Header.Get("Content-Type")
	fmt.Println(contentType, ":contenttype issssssssssssssssss")

	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
			log.Printf("Error parsing multipart form: %v", err)
			tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid form data",
			})
			return
		}
		hasFile = true

		// Parse comment data from form values
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

	// Ensure post_id is set (override any value from form/JSON)
	comment.Post_id = post_id

	// Validate comment
	if ok, status := models.ValidateComment(app.DL.Comments.DB, &comment); !ok {
		tools.EncodeJSON(w, status, map[string]string{
			"error": "Invalid comment data",
		})
		return
	}

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

			// Validate file type
			if !tools.IsAllowedFile(handler.Filename, file) {
				tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "File format is not supported",
				})
				return
			}

			// Upload file
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
				"error": "Error processing image upload",
			})
			return
		}
	}

	// Begin transaction
	tx, err := app.DL.Comments.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
		return
	}
	defer tx.Rollback()

	// Insert comment
	result, err := tx.Exec(`
		INSERT INTO comments (user_id, post_id, content, image_url, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		comment.OwnerId, comment.Post_id, comment.Content, imagePath, time.Now())
	if err != nil {
		log.Printf("Error inserting comment: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create comment",
		})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create comment",
		})
		return
	}

	// Get the created comment ID and prepare response
	commentID, _ := result.LastInsertId()
	comment.Id = int(commentID)
	comment.Image = imagePath
	comment.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	// Fetch the owner's name for the response
	var ownerName string
	err = app.DL.Comments.DB.QueryRow(`
		SELECT COALESCE(NULLIF(nickname, ''), first_name || ' ' || last_name) as owner_name
		FROM users WHERE id = ?
	`, comment.OwnerId).Scan(&ownerName)
	if err != nil {
		log.Printf("Error fetching owner name: %v", err)
		ownerName = "Unknown User"
	}
	comment.Owner = ownerName

	// Return the created comment
	tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Comment created successfully",
		"comment": comment,
	})
}
