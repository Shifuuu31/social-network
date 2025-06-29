package tools

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const UploadDir = "./pkg/db/data/uploads"

// Allowed MIME types and their corresponding extensions
var allowedImageMIMEs = map[string]string{
	"image/png":  "png",
	"image/jpeg": "jpg",
	"image/jpg":  "jpg",
	"image/gif":  "gif",
	"image/webp": "webp",
}

func EnsureUploadDir() error {
	// hadi 3abita usi nichan makdirall 7it recursive
	return os.MkdirAll(UploadDir, os.ModePerm)
}

// GetImageExtension returns the file extension for a given MIME type
func GetImageExtension(mimeType string) (string, error) {
	ext, ok := allowedImageMIMEs[mimeType]
	if !ok {
		return "", fmt.Errorf("unsupported image MIME type: %s", mimeType)
	}
	return ext, nil
}

// ValidateImageFile checks the file header for valid size and content type
func ValidateImageFile(header *multipart.FileHeader) error {
	// Optional: limit size (e.g., 5MB)
	const maxFileSize = 5 << 20 // 5MB
	if header.Size > maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", header.Size, maxFileSize)
	}

	// Optional: validate by file extension (more robust: detect MIME using file content)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" || (ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".webp") {
		return fmt.Errorf("invalid file extension: %s", ext)
	}

	return nil
}

// SaveUploadedImage saves the uploaded file to the uploads folder
func SaveUploadedImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := EnsureUploadDir(); err != nil {
		return "", fmt.Errorf("creating upload dir: %w", err)
	}

	// Use timestamped unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("post_%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(UploadDir, filename)

	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", fmt.Errorf("saving file: %w", err)
	}

	return filename, nil
}

// DeleteImage removes the image file from disk
func DeleteImage(filename string) error {
	path := filepath.Join(UploadDir, filename)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error deleting image: %w", err)
	}
	return nil
}
