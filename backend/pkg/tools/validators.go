package tools

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func ImageUpload(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	fmt.Println(r.FormFile("image"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return "", errors.New("file extension not valid")
	}
	if handler.Size > 5<<20 { // 5 MB limit
		return "", errors.New("size too big")
	}

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", errors.New("couldn't read the file")
	}
	// Reset file pointer ofr future reads
	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer)
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif"} // TODO add all types
	if !slices.Contains(allowedTypes, mimeType) {
		return "", errors.New("file format is not valid")
	}
	err = os.MkdirAll("pkg/data", os.ModePerm)
	if err != nil {
		return "", errors.New("failed to create uploads dir")
	}

	// Save the file
	dst, err := os.Create(filepath.Join("pkg/data", handler.Filename))
	if err != nil {
		return "", errors.New("failed to create file")
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", errors.New("failed to save file")
	}
	// = filepath.Join("images", handler.Filename)

	return filepath.Join("uploads", handler.Filename), nil
}

// func UploadHandler(file multipart.File, handler *multipart.FileHeader) (string, int) {
// 	// // Parse up to 10MB of form data
// 	// err := r.ParseMultipartForm(10 << 20) // 10MB
// 	// if err != nil {
// 	// 	http.Error(w, "Failed to parse form", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	// // Check if it's an image
// 	// ext := strings.ToLower(filepath.Ext(handler.Filename))
// 	// if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
// 	// 	fmt.Println("Invalid file type")
// 	// 	return "", 400
// 	// }

// 	// Make sure the uploads directory exists
// 	err := os.MkdirAll("uploads", os.ModePerm)
// 	if err != nil {
// 		return "", 500
// 	}

// 	// Save the file
// 	dst, err := os.Create(filepath.Join("uploads", handler.Filename))
// 	if err != nil {
// 		return "", 500
// 	}
// 	defer dst.Close()

// 	_, err = io.Copy(dst, file)
// 	if err != nil {
// 		return "", 500
// 	}
// 	// = filepath.Join("images", handler.Filename)

// 	return filepath.Join("uploads", handler.Filename), 200
// }

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
