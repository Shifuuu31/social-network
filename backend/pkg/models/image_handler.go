package models

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageHandler handles all image operations for the social network
type ImageHandler struct {
	DB *sql.DB
}

// ImageType represents different types of images in the system
type ImageType string

const (
	ImageTypePost    ImageType = "posts"
	ImageTypeProfile ImageType = "profiles"
	ImageTypeComment ImageType = "comments"
	ImageTypeGroup   ImageType = "groups"
)

// Supported image types
var allowedMIMETypes = map[string]string{
	"image/jpeg": "jpg",
	"image/jpg":  "jpg",
	"image/png":  "png",
	"image/gif":  "gif",
}

// Image represents an image in the database
type Image struct {
	ID           int       `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	FileSize     int64     `json:"file_size"`
	UploadPath   string    `json:"upload_path"`
	UploadedBy   int       `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// UploadImage handles image upload for any type (posts, profiles, comments)
func (ih *ImageHandler) UploadImage(file multipart.File, header *multipart.FileHeader, imageType ImageType, uploadedBy int) (*Image, error) {
	// 1. Validate file
	if err := ih.validateImage(header); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Generate secure filename
	filename, err := ih.generateSecureFilename(header)
	if err != nil {
		return nil, fmt.Errorf("filename generation failed: %w", err)
	}

	// 3. Create upload directory
	uploadDir := filepath.Join("uploads", string(imageType))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 4. Save file to disk
	filePath := filepath.Join(uploadDir, filename)
	if err := ih.saveFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 5. Save metadata to database
	image := &Image{
		Filename:     filename,
		OriginalName: header.Filename,
		MimeType:     header.Header.Get("Content-Type"),
		FileSize:     header.Size,
		UploadPath:   filepath.Join(string(imageType), filename),
		UploadedBy:   uploadedBy,
		CreatedAt:    time.Now(),
	}

	if err := ih.saveImageMetadata(image); err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image metadata: %w", err)
	}

	return image, nil
}

// validateImage checks if the uploaded file is valid
func (ih *ImageHandler) validateImage(header *multipart.FileHeader) error {
	// Check file size (5MB limit)
	const maxSize = 5 << 20 // 5MB
	if header.Size > maxSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", header.Size, maxSize)
	}

	// Check MIME type
	mimeType := header.Header.Get("Content-Type")
	if _, allowed := allowedMIMETypes[mimeType]; !allowed {
		return fmt.Errorf("unsupported MIME type: %s", mimeType)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		return fmt.Errorf("no file extension")
	}

	expectedExt := "." + allowedMIMETypes[mimeType]
	if ext != expectedExt {
		return fmt.Errorf("file extension %s doesn't match MIME type %s", ext, mimeType)
	}

	return nil
}

// generateSecureFilename creates a secure, unique filename
func (ih *ImageHandler) generateSecureFilename(header *multipart.FileHeader) (string, error) {
	// Generate random bytes for uniqueness
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Get extension from MIME type
	mimeType := header.Header.Get("Content-Type")
	ext := allowedMIMETypes[mimeType]

	// Create filename: timestamp_randombytes.extension
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%x.%s", timestamp, randomBytes, ext)

	return filename, nil
}

// saveFile saves the uploaded file to disk
func (ih *ImageHandler) saveFile(file multipart.File, filePath string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	return err
}

// saveImageMetadata saves image metadata to database
func (ih *ImageHandler) saveImageMetadata(image *Image) error {
	query := `
		INSERT INTO images (filename, original_name, mime_type, file_size, upload_path, uploaded_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := ih.DB.Exec(query,
		image.Filename,
		image.OriginalName,
		image.MimeType,
		image.FileSize,
		image.UploadPath,
		image.UploadedBy,
		image.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	image.ID = int(id)
	return nil
}

// ServeImage serves an image file with proper headers
func (ih *ImageHandler) ServeImage(w http.ResponseWriter, imageID int, requesterID int) error {
	// Get image metadata
	image, err := ih.getImageByID(imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	// Check authorization based on image type
	if err := ih.checkImageAuthorization(image, requesterID); err != nil {
		return fmt.Errorf("unauthorized: %w", err)
	}

	// Serve the file
	filePath := filepath.Join("uploads", image.UploadPath)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Set proper headers
	w.Header().Set("Content-Type", image.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", image.FileSize))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	// Copy file to response
	_, err = io.Copy(w, file)
	return err
}

// getImageByID retrieves image metadata from database
func (ih *ImageHandler) getImageByID(imageID int) (*Image, error) {
	query := `SELECT id, filename, original_name, mime_type, file_size, upload_path, uploaded_by, created_at 
			  FROM images WHERE id = ?`

	image := &Image{}
	err := ih.DB.QueryRow(query, imageID).Scan(
		&image.ID, &image.Filename, &image.OriginalName, &image.MimeType,
		&image.FileSize, &image.UploadPath, &image.UploadedBy, &image.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return image, nil
}

// checkImageAuthorization checks if requester can view the image
func (ih *ImageHandler) checkImageAuthorization(image *Image, requesterID int) error {
	// For now, allow access if user uploaded the image
	// You can extend this with more complex authorization logic
	if image.UploadedBy == requesterID {
		return nil
	}

	// TODO: Add more authorization logic based on your privacy rules
	// - Profile images: check if profile is public or requester is follower
	// - Post images: check post privacy settings
	// - Comment images: check parent post privacy settings

	return fmt.Errorf("unauthorized access")
}

// DeleteImage deletes an image file and its metadata
func (ih *ImageHandler) DeleteImage(imageID int, requesterID int) error {
	image, err := ih.getImageByID(imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	// Check if user can delete this image
	if image.UploadedBy != requesterID {
		return fmt.Errorf("unauthorized to delete this image")
	}

	// Delete file from disk
	filePath := filepath.Join("uploads", image.UploadPath)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete metadata from database
	query := `DELETE FROM images WHERE id = ?`
	_, err = ih.DB.Exec(query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image metadata: %w", err)
	}

	return nil
}

// GetImageByID returns image metadata
func (ih *ImageHandler) GetImageByID(imageID int) (*Image, error) {
	return ih.getImageByID(imageID)
}
