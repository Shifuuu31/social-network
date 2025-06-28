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

const UploadDir = "./pkg/db/data/uploads"

func ImageUpload(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
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

	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif"} // TODO add all types

	mimeType := http.DetectContentType(buffer)
	if !slices.Contains(allowedTypes, mimeType) {
		return "", errors.New("file format is not valid")
	}
	err = os.MkdirAll(UploadDir, os.ModePerm)
	if err != nil {
		return "", errors.New("failed to create data dir")
	}

	path := filepath.Join(UploadDir, handler.Filename)
	// Save the file
	outFile, err := os.Create(path)
	if err != nil {
		return "", errors.New("failed to create file")
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", errors.New("failed to save file")
	}

	return path, nil
}

// DeleteImage removes the image file from disk
func DeleteImage(filename string) error {
	path := filepath.Join(UploadDir, filename)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error deleting image: %w", err)
	}
	return nil
}
