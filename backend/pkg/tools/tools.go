package tools

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func SliceHasPrefix(prifixes []string, path string) bool {
	for _, prefix := range prifixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func IsAllowedFile(filename string, file multipart.File) bool {
	// Check extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return false
	}

	// Read first 512 bytes for content sniffing
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false
	}
	// Reset file pointer ofr future reads
	file.Seek(0, 0)

	// Detect MIME type from sample (buffer)
	mimeType := http.DetectContentType(buffer)
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	return slices.Contains(allowedTypes, mimeType)
}

func UploadHandler(file multipart.File, handler *multipart.FileHeader) (string, int) {

	// // Parse up to 10MB of form data
	// err := r.ParseMultipartForm(10 << 20) // 10MB
	// if err != nil {
	// 	http.Error(w, "Failed to parse form", http.StatusBadRequest)
	// 	return
	// }

	// Check if it's an image
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		fmt.Println("Invalid file type")
		return "", 400
	}

	// Make sure the uploads directory exists
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return "", 500
	}

	// Save the file
	dst, err := os.Create(filepath.Join("uploads", handler.Filename)) // TODO might wann add user spicific folder assignment
	if err != nil {
		return "", 500
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", 500
	}
	// = filepath.Join("images", handler.Filename)

	return filepath.Join("uploads", handler.Filename), 200
}

// func ValidateComment(db *sql.DB, comment *models.Comment) int {
// 	if comment.OwnerId <= 0 || comment.Post_id <= 0 {
// 		return 400
// 	}

// 	comment.Content = strings.TrimSpace(comment.Content)
// 	if len(comment.Content) > 300 || comment.Content == "" {
// 		return 400
// 	}
// 	return 200
// }

// ValidatePostFilter validates the post filter and returns HTTP status code
// func ValidatePostFilter(filter *models.PostFilter) int {
// 	if filter == nil {
// 		return 400
// 	}

// 	// Validate Type if provided
// 	if filter.Type != "" {
// 		validTypes := map[string]bool{
// 			"feed":   true,
// 			"user":   true,
// 			"group":  true,
// 			"public": true,
// 		}
// 		if !validTypes[filter.Type] {
// 			return 400
// 		}
// 	}

// 	// Validate Start (pagination offset)
// 	if filter.Start < 0 {
// 		return 400
// 	}

// 	// Validate NPost (number of posts to fetch)
// 	if filter.NPost <= 0 {
// 		filter.NPost = 10 // Default value
// 	}
// 	if filter.NPost > 100 { // Limit maximum posts per request
// 		return 400
// 	}

// 	// Validate Id if provided (for specific post or user/group filtering)
// 	if filter.Id < 0 {
// 		return 400
// 	}

// 	return 200 // OK
// }

// ParsePostFromForm parses post data from multipart form and returns HTTP status code

// func ValidatePost(post *models.Post) int {
// 	fmt.Println("mok")
// 	if post.OwnerId <= 0 { //TODO i need to check for a not exesting id
// 		return 400
// 	}
// 	// 0 means no group
// 	if post.GroupId < 0 { //TODO i need to check for a not exesting id
// 		return 400
// 	}

// 	post.Content = strings.TrimSpace(post.Content)
// 	if post.Content == "" {
// 		return 400
// 	}
// 	if len(post.Content) > 1000 {
// 		return 400
// 	}

// 	if post.Privacy != "public" && post.Privacy != "almost_private" && post.Privacy != "private" {
// 		return 400
// 	}

// 	if post.Privacy == "private" {
// 		if len(post.ChosenUsersIds) == 0 {
// 			return 400
// 		}
// 		userIdMap := make(map[int]bool)
// 		for _, id := range post.ChosenUsersIds {
// 			if id <= 0 || userIdMap[id] {
// 				return 400
// 			}
// 			userIdMap[id] = true
// 		}
// 	} else if len(post.ChosenUsersIds) > 0 {
// 		return 400
// 	}

// 	return 200
// }
