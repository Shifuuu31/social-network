package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type ImageModel struct {
	DB *sql.DB
}

const UploadDir = "./pkg/db/data/uploads"

// GetImageIfAuthorized returns the image path from a given table (users, groups, posts, comments)
// if the requesting user has permission to view the image.
// Access is enforced per-table based on privacy rules.
func (imgm *ImageModel) GetImageIfAuthorized(uuid, table string, requesterID int) (path string, err error) {
	switch table {

	case "users":
		// Allow viewing user avatar if:
		// - Profile is public
		// - Requester is the user themselves
		// - Requester is a follower with accepted request
		query := `
			SELECT u.avatar_path
			FROM users u
			WHERE u.image_uuid = ?
			  AND (
				u.is_public = 1
				OR u.id = ?
				OR EXISTS (
					SELECT 1 FROM follow_requests fr
					WHERE fr.from_user_id = ? AND fr.to_user_id = u.id AND fr.status = 'accepted'
				)
			)`
		err = imgm.DB.QueryRow(query, uuid, requesterID, requesterID).Scan(&path)
		if err != nil {
			return "", fmt.Errorf("unauthorized or not found: %w", err)
		}

	case "groups":
		// Allow viewing group image if requester is a confirmed group member
		query := `
			SELECT g.image_path
			FROM groups g
			JOIN group_members gm ON gm.group_id = g.id
			WHERE g.image_uuid = ? AND gm.user_id = ? AND gm.status = 'member'`
		err = imgm.DB.QueryRow(query, uuid, requesterID).Scan(&path)
		if err != nil {
			return "", fmt.Errorf("unauthorized or not found: %w", err)
		}

	case "posts":
		// Allow viewing post image if:
		// - Public post
		// - Follower of the author (and accepted)
		// - Post shared with requester explicitly (selected)
		// - Group post & requester is a group member
		// - Post owner is requester
		query := `
			SELECT p.image_path
			FROM posts p
			WHERE p.image_uuid = ? AND (
				p.privacy = 'public'
				OR (p.privacy = 'followers' AND EXISTS (
					SELECT 1 FROM follow_requests fr
					WHERE fr.from_user_id = ? AND fr.to_user_id = p.user_id AND fr.status = 'accepted'
				))
				OR (p.privacy = 'selected' AND EXISTS (
					SELECT 1 FROM post_privacy_selected s
					WHERE s.post_id = p.id AND s.user_id = ?
				))
				OR (p.privacy = 'group' AND EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = p.group_id AND gm.user_id = ? AND gm.status = 'member'
				))
				OR p.user_id = ?
			)`
		err = imgm.DB.QueryRow(query, uuid, requesterID, requesterID, requesterID, requesterID).Scan(&path)
		if err != nil {
			return "", fmt.Errorf("unauthorized or not found: %w", err)
		}

	case "comments":
		// Allow viewing comment image if requester is allowed to view the parent post
		// Similar logic to post visibility
		query := `
			SELECT c.image_path
			FROM comments c
			JOIN posts p ON p.id = c.post_id
			WHERE c.image_uuid = ? AND (
				p.privacy = 'public'
				OR (p.privacy = 'followers' AND EXISTS (
					SELECT 1 FROM follow_requests fr
					WHERE fr.from_user_id = ? AND fr.to_user_id = p.user_id AND fr.status = 'accepted'
				))
				OR (p.privacy = 'selected' AND EXISTS (
					SELECT 1 FROM post_privacy_selected s
					WHERE s.post_id = p.id AND s.user_id = ?
				))
				OR (p.privacy = 'group' AND EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = p.group_id AND gm.user_id = ? AND gm.status = 'member'
				))
				OR p.user_id = ?
			)`
		err = imgm.DB.QueryRow(query, uuid, requesterID, requesterID, requesterID, requesterID).Scan(&path)
		if err != nil {
			return "", fmt.Errorf("unauthorized or not found: %w", err)
		}

	default:
		// Unsupported table
		return "", errors.New("unsupported image table")
	}

	// Return relative path under the uploads folder
	return filepath.Join(UploadDir, path), nil
}

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

func ImageServe(w http.ResponseWriter, imagePath string) (err error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return err
	}
	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer[:n])
	w.Header().Set("Content-Type", mimeType)
	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// DeleteImage removes the image file from disk
func DeleteImage(filename string) error {
	path := filepath.Join(UploadDir, filename)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error deleting image: %w", err)
	}
	return nil
}
