package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

// ImageHandler handles image upload and serving
type ImageHandler struct {
	*Root
}

// NewImageHandler creates a new image handler
func (rt *Root) NewImageHandler() *ImageHandler {
	return &ImageHandler{Root: rt}
}

// SetupImageRoutes sets up image-related routes
func (ih *ImageHandler) SetupImageRoutes() *http.ServeMux {
	imageMux := http.NewServeMux()

	// Image upload routes
	imageMux.HandleFunc("/post", ih.UploadPostImage)
	imageMux.HandleFunc("/profile", ih.UploadProfileImage)
	imageMux.HandleFunc("/comment", ih.UploadCommentImage)

	// Image serving routes
	imageMux.HandleFunc("/serve/", ih.ServeImage)
	imageMux.HandleFunc("/delete/", ih.DeleteImage)

	return imageMux
}

// UploadPostImage handles post image uploads
func (ih *ImageHandler) UploadPostImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tools.EncodeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	requesterID := ih.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.EncodeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Failed to parse form",
		})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "No image file provided",
		})
		return
	}
	defer file.Close()

	// Upload image using unified handler
	imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
	image, err := imageHandler.UploadImage(file, header, models.ImageTypePost, requesterID)
	if err != nil {
		log.Printf("Error uploading post image: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to upload image",
		})
		return
	}

	tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Image uploaded successfully",
		"image":   image,
	})
}

// UploadProfileImage handles profile image uploads
func (ih *ImageHandler) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tools.EncodeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	requesterID := ih.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.EncodeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Failed to parse form",
		})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "No image file provided",
		})
		return
	}
	defer file.Close()

	// Upload image using unified handler
	imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
	image, err := imageHandler.UploadImage(file, header, models.ImageTypeProfile, requesterID)
	if err != nil {
		log.Printf("Error uploading profile image: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to upload image",
		})
		return
	}

	// Update user's profile with new image
	if err := ih.updateUserProfileImage(requesterID, image.ID); err != nil {
		log.Printf("Error updating user profile image: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update profile",
		})
		return
	}

	tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Profile image updated successfully",
		"image":   image,
	})
}

// UploadCommentImage handles comment image uploads
func (ih *ImageHandler) UploadCommentImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tools.EncodeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	requesterID := ih.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.EncodeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Failed to parse form",
		})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "No image file provided",
		})
		return
	}
	defer file.Close()

	// Upload image using unified handler
	imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
	image, err := imageHandler.UploadImage(file, header, models.ImageTypeComment, requesterID)
	if err != nil {
		log.Printf("Error uploading comment image: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to upload image",
		})
		return
	}

	tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Image uploaded successfully",
		"image":   image,
	})
}

// ServeImage serves an image file
func (ih *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	imageIDStr := r.PathValue("image_id")
	if imageIDStr == "" {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Image ID required",
		})
		return
	}

	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid image ID",
		})
		return
	}

	requesterID := ih.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.EncodeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Serve image using unified handler
	imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
	if err := imageHandler.ServeImage(w, imageID, requesterID); err != nil {
		log.Printf("Error serving image: %v", err)
		tools.EncodeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Image not found or unauthorized",
		})
		return
	}
}

// DeleteImage deletes an image
func (ih *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	imageIDStr := r.PathValue("image_id")
	if imageIDStr == "" {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Image ID required",
		})
		return
	}

	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid image ID",
		})
		return
	}

	requesterID := ih.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.EncodeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Delete image using unified handler
	imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
	if err := imageHandler.DeleteImage(imageID, requesterID); err != nil {
		log.Printf("Error deleting image: %v", err)
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete image",
		})
		return
	}

	tools.EncodeJSON(w, http.StatusOK, map[string]string{
		"message": "Image deleted successfully",
	})
}

// updateUserProfileImage updates a user's profile image
func (ih *ImageHandler) updateUserProfileImage(userID, imageID int) error {
	// Get the image to construct the URL
	imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
	image, err := imageHandler.GetImageByID(imageID)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	// Construct the avatar URL
	avatarURL := fmt.Sprintf("/images/%s", image.Filename)

	query := `UPDATE users SET avatar_url = ? WHERE id = ?`
	_, err = ih.DL.Posts.DB.Exec(query, avatarURL, userID)
	return err
}
